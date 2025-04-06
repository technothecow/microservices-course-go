package gateway

import (
	"errors"
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) GetCommentsList(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PaginatedCommentsRequest{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	comments, err := usecase.GetCommentsList(userId, &body)
	if err != nil {
		if errors.Is(err, usecase.ErrPostNotFound) {
			ctx.Status(http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrPostNotAuthorized) {
			ctx.Status(http.StatusForbidden)
		} else {
			log.Printf("Failed to get comments list: %v", err)
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to get comments list",
				Code:    "failed_to_get_comments_list",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
