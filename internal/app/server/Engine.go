package server

import (
	"baconi.co.uk/oauth/internal/app/token"
	internalAuthentication "baconi.co.uk/oauth/internal/pkg/authentication"
	internalToken "baconi.co.uk/oauth/internal/pkg/token"
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

	//
	// Create stuff to be injected
	//
	accessTokenRepository := internalToken.NewInMemoryRepository[internalToken.AccessToken]()
	accessTokenService := internalToken.NewAccessTokenService(accessTokenRepository)

	refreshTokenRepository := internalToken.NewInMemoryRepository[internalToken.RefreshToken]()
	refreshTokenService := internalToken.NewRefreshTokenService(refreshTokenRepository)

	userAuthenticationService := internalAuthentication.UserAuthenticationService{}

	passwordGrant := token.NewPasswordGrant(accessTokenService, refreshTokenService, userAuthenticationService)

	//
	// Add Routes
	//
	engine.POST("/token", token.Route(passwordGrant))

	return engine, nil
}
