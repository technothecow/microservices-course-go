package gateway

import (
	"sn/gateway/internal/usecase"
	gen "sn/gateway/generated"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (*Server) GetPost(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PostId{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "bad_request",
		})
		return
	}

	post, err := usecase.GetPost(body.Id.String(), userId)
	if err != nil {
		if err == usecase.ErrPostNotFound {
			ctx.Status(http.StatusNotFound)
		} else {
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to get post",
				Code:    "failed_to_get_post",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, post)
}