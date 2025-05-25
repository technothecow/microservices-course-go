package usecase

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"sn/libraries/kafka"
	"time"

	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sn/libraries/postgres"
	posts "sn/libraries/proto/posts"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

func CreatePost(req *posts.CreatePostRequest) (*posts.Post, error) {
	query := `
	INSERT INTO posts (title, description, creator_id, is_private, tags)
	VALUES ($1, $2, $3, $4, $5::TEXT[])
	RETURNING id, title, description, creator_id, created_at, updated_at, is_private, tags`

	post := &posts.Post{}
	var createdAt time.Time
	var updatedAt time.Time
	var tags []string

	db := postgres.GetPostgresConnection()

	requestTags := req.Tags.GetValues()
	if requestTags == nil {
		requestTags = []string{}
	}

	err := db.QueryRow(query, req.Title, req.Description, req.UserId, req.IsPrivate, pq.StringArray(requestTags)).Scan(
		&post.Id, &post.Title, &post.Description, &post.UserId, &createdAt, &updatedAt, &post.IsPrivate, pq.Array(&tags))
	if err != nil {
		return nil, err
	}

	post.CreatedAt = timestamppb.New(createdAt)
	post.UpdatedAt = timestamppb.New(updatedAt)
	if tags != nil {
		post.Tags = &posts.Tags{Values: tags}
	} else {
		post.Tags = nil
	}

	return post, nil
}

func GetPost(req *posts.GetPostRequest) (*posts.Post, error) {
	query := `
	SELECT id, title, description, creator_id, created_at, updated_at, is_private, tags
	FROM posts
	WHERE id = $1`

	post := &posts.Post{}
	var createdAt time.Time
	var updatedAt time.Time
	var tags []string

	db := postgres.GetPostgresConnection()

	log.Printf("query: %s", req.Id)
	err := db.QueryRow(query, req.Id).Scan(
		&post.Id, &post.Title, &post.Description, &post.UserId, &createdAt, &updatedAt, &post.IsPrivate, pq.Array(&tags))
	if err != nil {
		log.Printf("error: %v", err)
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	post.CreatedAt = timestamppb.New(createdAt)
	post.UpdatedAt = timestamppb.New(updatedAt)
	if tags != nil {
		post.Tags = &posts.Tags{Values: tags}
	} else {
		post.Tags = nil
	}

	return post, nil
}

func UpdatePost(req *posts.UpdatePostRequest) (*posts.Post, error) {
	query := `
	UPDATE posts
	SET title = $1, description = $2, is_private = $3, tags = $4
	WHERE id = $5
	RETURNING id, title, description, creator_id, created_at, updated_at, is_private, tags`

	post := &posts.Post{}
	var createdAt time.Time
	var updatedAt time.Time
	var tags []string

	db := postgres.GetPostgresConnection()

	reqTags := req.Tags.GetValues()
	if reqTags == nil {
		reqTags = []string{}
	}

	err := db.QueryRow(query, req.Title, req.Description, req.IsPrivate, pq.StringArray(reqTags), req.Id).Scan(
		&post.Id, &post.Title, &post.Description, &post.UserId, &createdAt, &updatedAt, &post.IsPrivate, pq.Array(&tags))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	post.CreatedAt = timestamppb.New(createdAt)
	post.UpdatedAt = timestamppb.New(updatedAt)
	if tags != nil {
		post.Tags = &posts.Tags{Values: tags}
	} else {
		post.Tags = nil
	}

	return post, nil
}

func DeletePost(req *posts.DeletePostRequest) error {
	query := `
	DELETE FROM posts
	WHERE id = $1`

	db := postgres.GetPostgresConnection()

	result, err := db.Exec(query, req.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}

func ListPosts(req *posts.ListPostsRequest) (*posts.ListPostsResponse, error) {
	query := `
	SELECT id, title, description, creator_id, created_at, updated_at, is_private, tags
	FROM posts
	WHERE (creator_id = $1 OR is_private = FALSE) AND (tags @> $2::TEXT[] OR tags = '{}')
	ORDER BY created_at DESC
	LIMIT $3 OFFSET $4`

	db := postgres.GetPostgresConnection()

	rows, err := db.Query(query, req.GetRequesterId(), pq.StringArray(req.Tags.GetValues()), req.GetPageSize(), req.GetPageNumber()*req.GetPageSize())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postsList := make([]*posts.Post, 0, req.GetPageSize())

	for rows.Next() {
		post := &posts.Post{}
		var createdAt time.Time
		var updatedAt time.Time
		var tags []string

		err := rows.Scan(
			&post.Id, &post.Title, &post.Description, &post.UserId, &createdAt, &updatedAt, &post.IsPrivate, pq.Array(&tags))
		if err != nil {
			return nil, err
		}

		post.CreatedAt = timestamppb.New(createdAt)
		post.UpdatedAt = timestamppb.New(updatedAt)
		if tags != nil {
			post.Tags = &posts.Tags{Values: tags}
		} else {
			post.Tags = nil
		}

		postsList = append(postsList, post)
	}

	totalCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE (creator_id = $1 OR is_private = FALSE) AND (tags @> $2::text[] OR tags = '{}')", req.GetRequesterId(), pq.Array(req.GetTags().GetValues())).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	return &posts.ListPostsResponse{
		Posts:      postsList,
		TotalCount: int32(totalCount),
	}, nil
}

var commentTopic = "posts_comments"

type commentEvent struct {
	PostId    string `json:"post_id"`
	Username  string `json:"username"`
	Timestamp int64  `json:"timestamp"`
}

func CreateComment(req *posts.CommentPostRequest, db *sql.DB) error {
	query := `INSERT INTO comments (content, creator_id, post_id)
	VALUES ($1, $2, $3)`

	_, err := db.Exec(query, req.GetText(), req.GetUserId(), req.GetPostId())
	if err != nil {
		return err
	}

	event := commentEvent{
		PostId:    req.GetPostId(),
		Username:  req.GetUserId(),
		Timestamp: time.Now().Unix(),
	}
	json_, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal comment event: %v", err)
		return err
	}
	_, _, err = kafka.SendMessageSync(commentTopic, []byte{}, json_, true)
	if err != nil {
		log.Printf("failed to send comment event to kafka: %v", err)
		return err
	}

	return nil
}

func ListComments(req *posts.ListCommentsRequest, db *sql.DB) ([]*posts.Comment, error) {
	query := `SELECT c.id, c.content, c.creator_id, c.created_at, c.post_id
	FROM comments c
	JOIN posts p ON c.post_id = p.id
	WHERE (p.is_private = FALSE OR p.creator_id = $1)
	AND c.post_id = $2
	ORDER BY c.created_at
	LIMIT $3 OFFSET $4;`

	rows, err := db.Query(query, req.GetUserId(), req.GetPostId(), req.GetPageSize(), req.GetPageSize()*req.GetPageNumber())
	if err != nil {
		return nil, err
	}

	var comments []*posts.Comment

	for rows.Next() {
		comment := &posts.Comment{}
		var createdAt time.Time

		err = rows.Scan(
			&comment.Id, &comment.Text, &comment.UserId, &createdAt, &comment.PostId)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = timestamppb.New(createdAt)

		comments = append(comments, comment)
	}

	return comments, nil
}
