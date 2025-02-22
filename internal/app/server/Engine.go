package server

import (
	"baconi.co.uk/oauth/internal/app/token"
	"github.com/gin-gonic/gin"
)

func Engine(
	config *Config,
) *gin.Engine {

	// Engine setup
	// TODO - gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	_ = engine.SetTrustedProxies(nil) // TODO - or exit and return error

	// Add Routes
	engine.POST("/token", token.Route)

	return engine
}
