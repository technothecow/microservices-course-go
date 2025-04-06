package gateway

import (
	"encoding/json"
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/libraries/kafka"
	"time"

	"github.com/gin-gonic/gin"
)

type LikeEvent struct {
	UserId    string `json:"user_id"`
	PostId    string `json:"post_id"`
	Timestamp int64  `json:"timestamp"`
}

var PostsLikesTopic = "posts_likes"

func (*Server) LikePost(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PostId{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "bad_request",
		})
		return
	}

	event := LikeEvent{
		UserId:    userId,
		PostId:    body.Id.String(),
		Timestamp: time.Now().Unix(),
	}
	json_, err := json.Marshal(&event)
	if err != nil {
		log.Printf("json marshal error: %v", err)
	}

	kafka.SendMessageAsync(PostsLikesTopic, []byte(body.Id.String()), json_, true)
}
