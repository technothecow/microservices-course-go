package gateway

import (
	"net/http"
	"log"
	"sn/gateway/internal/usecase"
	gen "sn/gateway/generated"
	"github.com/gin-gonic/gin"
)

func (*Server) GetPostsList(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PaginatedPostsRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	posts, err := usecase.ListPosts(userId, &body)
	if err != nil {
		log.Printf("failed to get posts list: %v", err)
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to get posts list",
			Code:    "failed_to_get_posts_list",
		})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}
