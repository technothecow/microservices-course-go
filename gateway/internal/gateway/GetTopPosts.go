package gateway

import (
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) GetTopPosts(ctx *gin.Context) {
	_, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.TopPostsRequest{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	comments, err := usecase.GetTopPosts(&body)
	if err != nil {
		log.Printf("Failed to get top posts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to get top posts",
			Code:    "failed_to_get_top_posts",
		})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}