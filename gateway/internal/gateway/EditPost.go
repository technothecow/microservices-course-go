package gateway

import (
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) EditPost(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.EditPostRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	post, err := usecase.EditPost(userId, &body)
	if err != nil {
		if err == usecase.ErrPostNotFound {
			ctx.Status(http.StatusNotFound)
		} else if err == usecase.ErrPostNotAuthorized {
			ctx.Status(http.StatusForbidden)
		} else {
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to edit post",
				Code:    "failed_to_edit_post",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, post)
}
