package gateway

import (
	"errors"
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) GetPostStats(ctx *gin.Context) {
	_, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.PostStatsRequest{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	comments, err := usecase.GetPostStats(&body)
	if err != nil {
		if errors.Is(err, usecase.ErrPostNotFound) {
			ctx.Status(http.StatusNotFound)
		} else {
			log.Printf("Failed to get post stats: %v", err)
			ctx.JSON(http.StatusInternalServerError, gen.Error{
				Message: "Failed to get post stats",
				Code:    "failed_to_get_post_stats",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
