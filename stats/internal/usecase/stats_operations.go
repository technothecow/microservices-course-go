package usecase

import (
	"database/sql"
	"errors"
	"time"

	"sn/libraries/clickhouse"
	"sn/libraries/proto/stats"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrPostNotFound = errors.New("post not found")

func GetPostStats(body *stats.GetPostStatsRequest) (*stats.PostStatsResponse, error) {
	// duplicated views by design
	const query = `SELECT
    views,
    likes,
    comments
FROM
    (
        SELECT post_id, count(*) AS views
        FROM stats.views
        WHERE post_id = ?
        GROUP BY post_id
    ) AS v
FULL OUTER JOIN
    (
        SELECT post_id, countIf(likes_count % 2 = 1) AS likes
        FROM (
            SELECT post_id, user_id, count() AS likes_count
            FROM stats.likes
            WHERE post_id = ?
            GROUP BY post_id, user_id
        )
        GROUP BY post_id
    ) AS l
    ON v.post_id = l.post_id
FULL OUTER JOIN
    (
        SELECT post_id, count(*) AS comments
        FROM stats.comments
        WHERE post_id = ?
        GROUP BY post_id
    ) AS c
    ON COALESCE(v.post_id, l.post_id) = c.post_id;`

	ch := clickhouse.GetClickhouseConnection()
	var response stats.PostStatsResponse
	row := ch.QueryRow(query, body.PostId, body.PostId, body.PostId)
	err := row.Scan(&response.Views, &response.Likes, &response.Comments)
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrPostNotFound
        }
		return nil, err
	}
	return &response, nil
}

func GetPostDynamics(body *stats.GetPostDynamicsRequest) (*stats.PostDynamicsResponse, error) {
    query := `
    WITH days AS (
        SELECT toDate(now() - INTERVAL number DAY) AS date
        FROM numbers(7)
    ),
    views AS (
        SELECT toDate(created_at) AS date, count(*) AS cnt
        FROM stats.views
        WHERE post_id = ? AND created_at >= now() - INTERVAL 7 DAY
        GROUP BY date
    ),
    likes AS (
        SELECT date, countIf(like_count % 2 = 1) AS cnt
        FROM (
            SELECT toDate(created_at) AS date, user_id, count() AS like_count
            FROM stats.likes
            WHERE post_id = ? AND created_at >= now() - INTERVAL 7 DAY
            GROUP BY date, user_id
        )
        GROUP BY date
    ),
    comments AS (
        SELECT toDate(created_at) AS date, count(*) AS cnt
        FROM stats.comments
        WHERE post_id = ? AND created_at >= now() - INTERVAL 7 DAY
        GROUP BY date
    )
    SELECT
        d.date,
        COALESCE(v.cnt, 0) AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM days d
    LEFT JOIN views v ON d.date = v.date
    LEFT JOIN likes l ON d.date = l.date
    LEFT JOIN comments c ON d.date = c.date
    ORDER BY d.date;
    `

    ch := clickhouse.GetClickhouseConnection()

    rows, err := ch.Query(query, body.PostId, body.PostId, body.PostId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var response stats.PostDynamicsResponse

    for rows.Next() {
        var date time.Time
        var views, likes, comments int64

        if err := rows.Scan(&date, &views, &likes, &comments); err != nil {
            return nil, err
        }

        ts := timestamppb.New(date)

        response.Views = append(response.Views, &stats.PostDynamicsItem{
            Date:  ts,
            Value: views,
        })
        response.Likes = append(response.Likes, &stats.PostDynamicsItem{
            Date:  ts,
            Value: likes,
        })
        response.Comments = append(response.Comments, &stats.PostDynamicsItem{
            Date:  ts,
            Value: comments,
        })
    }

    return &response, nil
}

func GetTopPosts(body *stats.GetTopPostsRequest) (*stats.TopPostsResponse, error) {
    queryViews := `
    WITH
        views AS (
            SELECT post_id, count(*) AS cnt FROM stats.views GROUP BY post_id
        ),
        likes AS (
            SELECT post_id, countIf(like_count % 2 = 1) AS cnt
            FROM (
                SELECT post_id, user_id, count() AS like_count
                FROM stats.likes
                GROUP BY post_id, user_id
            )
            GROUP BY post_id
        ),
        comments AS (
            SELECT post_id, count(*) AS cnt FROM stats.comments GROUP BY post_id
        )

    SELECT
        v.post_id,
        v.cnt AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM views v
    LEFT JOIN likes l ON v.post_id = l.post_id
    LEFT JOIN comments c ON v.post_id = c.post_id
    ORDER BY views DESC
    LIMIT 10;`
    queryComments := `
    WITH
        views AS (
            SELECT post_id, count(*) AS cnt FROM stats.views GROUP BY post_id
        ),
        likes AS (
            SELECT post_id, countIf(like_count % 2 = 1) AS cnt
            FROM (
                SELECT post_id, user_id, count() AS like_count
                FROM stats.likes
                GROUP BY post_id, user_id
            )
            GROUP BY post_id
        ),
        comments AS (
            SELECT post_id, count(*) AS cnt FROM stats.comments GROUP BY post_id
        )

    SELECT
        v.post_id,
        v.cnt AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM views v
    LEFT JOIN likes l ON v.post_id = l.post_id
    LEFT JOIN comments c ON v.post_id = c.post_id
    ORDER BY comments DESC
    LIMIT 10;`
    queryLikes := `
    WITH
        views AS (
            SELECT post_id, count(*) AS cnt FROM stats.views GROUP BY post_id
        ),
        likes AS (
            SELECT post_id, countIf(like_count % 2 = 1) AS cnt
            FROM (
                SELECT post_id, user_id, count() AS like_count
                FROM stats.likes
                GROUP BY post_id, user_id
            )
            GROUP BY post_id
        ),
        comments AS (
            SELECT post_id, count(*) AS cnt FROM stats.comments GROUP BY post_id
        )

    SELECT
        v.post_id,
        v.cnt AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM views v
    LEFT JOIN likes l ON v.post_id = l.post_id
    LEFT JOIN comments c ON v.post_id = c.post_id
    ORDER BY likes DESC
    LIMIT 10;`

    var query string
    if body.Param == stats.TopRequestParam_VIEWS {
        query = queryViews
    } else if body.Param == stats.TopRequestParam_LIKES {
        query = queryLikes
    } else if body.Param == stats.TopRequestParam_COMMENTS {
        query = queryComments
    }

    ch := clickhouse.GetClickhouseConnection()

    rows, err := ch.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var response stats.TopPostsResponse

    for rows.Next() {
        var item stats.TopPostsItem

        if err := rows.Scan(&item.PostId, &item.Views, &item.Likes, &item.Comments); err != nil {
            return nil, err
        }

        response.TopPosts = append(response.TopPosts, &item)
    }

    return &response, nil
}

func GetTopUsers(body *stats.GetTopUsersRequest) (*stats.TopUsersResponse, error) {
    queryViews := `
    WITH
        views AS (
            SELECT user_id, count(*) AS cnt FROM stats.views GROUP BY user_id
        ),
        likes AS (
            SELECT user_id, countIf(like_count % 2 = 1) AS cnt
            FROM (
                SELECT user_id, post_id, count() AS like_count
                FROM stats.likes
                GROUP BY user_id, post_id
            )
            GROUP BY user_id
        ),
        comments AS (
            SELECT user_id, count(*) AS cnt FROM stats.comments GROUP BY user_id
        )

    SELECT
        COALESCE(v.user_id, l.user_id, c.user_id) AS user_id,
        COALESCE(v.cnt, 0) AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM views v
    FULL OUTER JOIN likes l ON v.user_id = l.user_id
    FULL OUTER JOIN comments c ON COALESCE(v.user_id, l.user_id) = c.user_id
    ORDER BY views DESC
    LIMIT 10
    `
    queryLikes := `
    WITH
        views AS (
            SELECT user_id, count(*) AS cnt FROM stats.views GROUP BY user_id
        ),
        likes AS (
            SELECT user_id, countIf(like_count % 2 = 1) AS cnt
            FROM (
                SELECT user_id, post_id, count() AS like_count
                FROM stats.likes
                GROUP BY user_id, post_id
            )
            GROUP BY user_id
        ),
        comments AS (
            SELECT user_id, count(*) AS cnt FROM stats.comments GROUP BY user_id
        )

    SELECT
        COALESCE(v.user_id, l.user_id, c.user_id) AS user_id,
        COALESCE(v.cnt, 0) AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM views v
    FULL OUTER JOIN likes l ON v.user_id = l.user_id
    FULL OUTER JOIN comments c ON COALESCE(v.user_id, l.user_id) = c.user_id
    ORDER BY likes DESC
    LIMIT 10
    `
    queryComments := `
    WITH
        views AS (
            SELECT user_id, count(*) AS cnt FROM stats.views GROUP BY user_id
        ),
        likes AS (
            SELECT user_id, countIf(like_count % 2 = 1) AS cnt
            FROM (
                SELECT user_id, post_id, count() AS like_count
                FROM stats.likes
                GROUP BY user_id, post_id
            )
            GROUP BY user_id
        ),
        comments AS (
            SELECT user_id, count(*) AS cnt FROM stats.comments GROUP BY user_id
        )

    SELECT
        COALESCE(v.user_id, l.user_id, c.user_id) AS user_id,
        COALESCE(v.cnt, 0) AS views,
        COALESCE(l.cnt, 0) AS likes,
        COALESCE(c.cnt, 0) AS comments
    FROM views v
    FULL OUTER JOIN likes l ON v.user_id = l.user_id
    FULL OUTER JOIN comments c ON COALESCE(v.user_id, l.user_id) = c.user_id
    ORDER BY comments DESC
    LIMIT 10
    `

    var query string
    if body.Param == stats.TopRequestParam_VIEWS {
        query = queryViews
    } else if body.Param == stats.TopRequestParam_LIKES {
        query = queryLikes
    } else if body.Param == stats.TopRequestParam_COMMENTS {
        query = queryComments
    }

    ch := clickhouse.GetClickhouseConnection()

    rows, err := ch.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var response stats.TopUsersResponse

    for rows.Next() {
        var item stats.TopUsersItem

        if err := rows.Scan(&item.UserId, &item.Views, &item.Likes, &item.Comments); err != nil {
            return nil, err
        }

        response.TopUsers = append(response.TopUsers, &item)
    }

    return &response, nil
}
