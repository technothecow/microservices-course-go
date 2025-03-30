package gateway

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"
)

func (*Server) CreatePost(ctx *gin.Context) {
	body := gen.CreatePostRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("failed to bind body: %v", err)
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "bad_request",
		})
		return
	}

	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	post, err := usecase.CreatePost(userId, &body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to create post",
			Code:    "internal_server_error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}
