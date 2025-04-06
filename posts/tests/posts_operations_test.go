package tests

import (
	"database/sql"
	"sn/libraries/kafka"
	"sn/libraries/postgres"
	"sn/libraries/proto/posts"
	"sn/posts/internal/usecase"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var validUUID = "123e4567-e89b-12d3-a456-426614174000"
var validUserID = "123e4567-e89b-12d3-a456-426614174000"
var columns = []string{"id", "title", "description", "creator_id", "created_at", "updated_at", "is_private", "tags"}

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows(columns).
		AddRow(validUUID, "Test Post", "Test Content", validUserID, time.Now(), time.Now(), false, "{}")
	mock.ExpectQuery(".*").WillReturnRows(mockRows)

	postgres.GetPostgresConnection = func() *sql.DB {
		return db
	}

	post := &posts.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Content",
		UserId:      validUserID,
	}

	result, err := usecase.CreatePost(post)

	assert.NoError(t, err)
	assert.Equal(t, result.Id, validUUID)
}

func TestGetPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows(columns).
		AddRow(validUUID, "Test Post", "Test Content", validUserID, time.Now(), time.Now(), false, "{}")
	mock.ExpectQuery(".*").WillReturnRows(mockRows)

	postgres.GetPostgresConnection = func() *sql.DB {
		return db
	}

	req := &posts.GetPostRequest{
		Id:          validUUID,
		RequesterId: validUserID,
	}

	post, err := usecase.GetPost(req)

	assert.NoError(t, err)
	assert.Equal(t, post.Id, validUUID)
}

func TestListPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows(columns).
		AddRow(validUUID, "Test Post", "Test Content", validUserID, time.Now(), time.Now(), false, "{}").
		AddRow(validUUID, "Test Post 2", "Test Content 2", validUserID, time.Now(), time.Now(), false, "{}")
	mock.ExpectQuery(".*").WillReturnRows(mockRows)
	mockRows = sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(".*").WillReturnRows(mockRows)

	postgres.GetPostgresConnection = func() *sql.DB {
		return db
	}

	req := &posts.ListPostsRequest{
		RequesterId: validUserID,
		PageNumber:  0,
		PageSize:    10,
		Tags:        &posts.Tags{Values: []string{}},
	}

	posts_, err := usecase.ListPosts(req)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts_.Posts))
	assert.Equal(t, validUUID, posts_.Posts[0].Id)
	assert.Equal(t, validUUID, posts_.Posts[1].Id)
}

func TestUpdatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	newTitle := "New Title"
	newDescription := "New Description"

	mockRows := sqlmock.NewRows(columns).
		AddRow(validUUID, newTitle, newDescription, validUserID, time.Now(), time.Now(), false, "{}")
	mock.ExpectQuery(".*").WillReturnRows(mockRows)

	postgres.GetPostgresConnection = func() *sql.DB {
		return db
	}

	req := &posts.UpdatePostRequest{
		Id:          validUUID,
		Title:       &newTitle,
		Description: &newDescription,
		RequesterId: validUserID,
	}

	post, err := usecase.UpdatePost(req)

	assert.NoError(t, err)
	assert.Equal(t, newTitle, post.Title)
	assert.Equal(t, newDescription, post.Description)
}

func TestDeletePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))

	postgres.GetPostgresConnection = func() *sql.DB {
		return db
	}

	req := &posts.DeletePostRequest{
		Id:          validUUID,
		RequesterId: validUserID,
	}

	err = usecase.DeletePost(req)
	assert.NoError(t, err)
}

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mock.ExpectExec(".*").WillReturnResult(sqlmock.NewResult(0, 1))

	kafka.SetTestingMode(true)

	req := &posts.CommentPostRequest{
		PostId: validUUID,
		Text:   "Comment Text",
		UserId: validUserID,
	}

	err = usecase.CreateComment(req, db)
	assert.NoError(t, err)
}

func TestListComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "text", "user_id", "created_at", "post_id"}).
		AddRow(validUUID, "Comment 1", validUserID, time.Now(), validUUID).
		AddRow(validUUID, "Comment 2", validUserID, time.Now(), validUUID)
	mock.ExpectQuery(".*").WillReturnRows(mockRows)

	req := &posts.ListCommentsRequest{
		PostId:     validUUID,
		UserId:     validUserID,
		PageNumber: 0,
		PageSize:   10,
	}

	result, err := usecase.ListComments(req, db)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
}
