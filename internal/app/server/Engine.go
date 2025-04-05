package server

import (
	"baconi.co.uk/oauth/internal/app/token-exchange"
	"baconi.co.uk/oauth/internal/app/token-introspection"
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
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
	accessTokenRepository := token.NewInMemoryRepository[token.AccessToken]()
	accessTokenIssuer := token.NewAccessTokenIssuer(accessTokenRepository)

	refreshTokenRepository := token.NewInMemoryRepository[token.RefreshToken]()
	refreshTokenIssuer := token.NewRefreshTokenIssuer(refreshTokenRepository)

	scopeRepository := scope.NewInMemoryRepository()
	scopeService := scope.NewService(scopeRepository)

	userCredentialRepository := user.NewInMemoryCredentialRepository()
	userStatusRepository := user.NewInMemoryStatusRepository()
	userAuthenticationService := user.NewAuthenticationService(userCredentialRepository, userStatusRepository)

	passwordGrant := token_exchange.NewPasswordGrant(accessTokenIssuer, refreshTokenIssuer, userAuthenticationService)

	clientSecretRepository := client.NewInMemorySecretRepository()
	clientPrincipalRepository := client.NewInMemoryPrincipalRepository()
	clientAuthenticationService := client.NewAuthenticationService(clientSecretRepository, clientPrincipalRepository)

	tokenIntrospector := token_introspection.NewIntrospector(accessTokenRepository)

	//
	// Add Routes
	//
	engine.POST("/token",
		client.AuthenticateConfidentialClient(clientAuthenticationService),
		client.AuthenticatePublicClient(clientAuthenticationService),
		client.RequireClientAuthentication,
		// TODO - FormUrlEncodedRequired() ???
		token_exchange.Route(scopeService, passwordGrant),
	)

	token_introspection.Route(engine, clientAuthenticationService, tokenIntrospector)

	return engine, nil
}
