package gateway

import (
	"errors"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) PostComment(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.CreateCommentRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	err = usecase.CreateComment(userId, &body)
	if err != nil {
		if errors.Is(err, usecase.ErrPostNotFound) {
			ctx.Status(http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrPostNotAuthorized) {
			ctx.Status(http.StatusForbidden)
		} else {
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to add comment",
				Code:    "failed_to_add_comment",
			})
		}
		return
	}

	ctx.Status(http.StatusOK)
}
