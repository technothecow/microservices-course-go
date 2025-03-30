package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/oapi-codegen/gin-middleware"

	gen "sn/gateway/generated"
	"sn/gateway/internal/gateway"
)

func main() {
	logger := zap.NewExample()
	logger.Info("Starting server on port 50001")

	server := gateway.NewServer()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	validationMiddleware, err := ginmiddleware.OapiValidatorFromYamlFile("openapi.yaml")
	if err != nil {
		panic(err)
	}
	r.Use(validationMiddleware)
	gen.RegisterHandlers(r, server)
	s := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:50001",
	}

	log.Fatal(s.ListenAndServe())
}
