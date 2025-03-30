package gateway

import (
	"net/http"
	"sn/gateway/internal/usecase"
	gen "sn/gateway/generated"
	"github.com/gin-gonic/gin"
)

func (*Server) DeletePost(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PostId{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
	}

	err = usecase.DeletePost(userId, body.Id.String())
	if err != nil {
		if err == usecase.ErrPostNotFound {
			ctx.Status(http.StatusNotFound)
		} else if err == usecase.ErrPostNotAuthorized {
			ctx.Status(http.StatusForbidden)
		} else {
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to delete post",
				Code:    "failed_to_delete_post",
			})
		}
		return
	}

	ctx.Status(http.StatusOK)
}
