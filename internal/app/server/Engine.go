package server

import (
	"baconi.co.uk/oauth/internal/app/token"
	"github.com/gin-gonic/gin"
)

func Engine(
	config *Config,
) (*gin.Engine, error) {

	// TODO - gin.SetMode(gin.ReleaseMode)

	// Engine setup
	engine := gin.New(func(engine *gin.Engine) {
		engine.HandleMethodNotAllowed = true
		engine.Use(gin.Logger(), gin.Recovery())
	})

	// Because GO likes to have errors returned.
	if err := engine.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	// Add Routes
	engine.POST("/token", token.Route)

	return engine, nil
}
