package server

import (
	"baconi.co.uk/oauth/internal/app/token"

	"github.com/gin-gonic/gin"
)

func Engine(
	config *Config,
) *gin.Engine {

	// Engine setup
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	// Add Routes
	engine.POST("/token", token.Route)

	return engine
}
