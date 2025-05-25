package gateway

import (
	"errors"
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) GetPostDynamics(ctx *gin.Context) {
	_, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PostDynamicsRequest{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	comments, err := usecase.GetPostDynamics(&body)
	if err != nil {
		if errors.Is(err, usecase.ErrPostNotFound) {
			ctx.Status(http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrNotAuthorized) {
			ctx.Status(http.StatusForbidden)
		} else {
			log.Printf("Failed to get post dynamics: %v", err)
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to get post dynamics",
				Code:    "failed_to_get_post_dynamics",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, comments)
}