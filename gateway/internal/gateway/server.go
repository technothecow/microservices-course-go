package gateway

import (
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/crypto"

	"github.com/gin-gonic/gin"
)

type Server struct {
	gen.ServerInterface
}

func NewServer() *Server {
	return &Server{}
}

func GetUserIdFromContext(ctx *gin.Context) (string, error) {
	authToken, err := ctx.Cookie("auth_token")
	if err != nil {
		log.Printf("failed to get auth token: %v", err)
		ctx.Status(http.StatusUnauthorized)
		return "", err
	}

	userId, err := crypto.GetUserIdFromToken(authToken)
	if err != nil {
		log.Printf("failed to get user id from token: %v", err)
		ctx.Status(http.StatusUnauthorized)
		return "", err
	}

	if userId == "" {
		log.Printf("failed to get user id from token: empty user id")
		ctx.Status(http.StatusUnauthorized)
		return "", err
	}

	return userId, nil
}
