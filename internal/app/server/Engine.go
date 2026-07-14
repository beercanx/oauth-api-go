package server

import (
	"context"

	"baconi.co.uk/oauth/internal/app/token-exchange"
	"baconi.co.uk/oauth/internal/app/token-introspection"
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/gin-gonic/gin"
)

func Engine(
	ctx context.Context,
	config Config,
) (*gin.Engine, error) {

	// TODO - gin.SetMode(gin.ReleaseMode)

	// Engine setup
	engine := gin.New(func(engine *gin.Engine) {
		engine.HandleMethodNotAllowed = true
		engine.Use(gin.Logger(), gin.Recovery())
	})

	// Because GO likes to have errors returned.
	if proxyError := engine.SetTrustedProxies(nil); proxyError != nil {
		return nil, proxyError
	}

	//
	// Create stuff to be injected
	//
	database, databaseError := db.Connect(config.DatabaseSource)
	if databaseError != nil {
		return nil, databaseError
	}
	if databaseMigrationsError := db.RunMigrations(database, config.DatabaseMigrations); databaseMigrationsError != nil {
		return nil, databaseMigrationsError
	}

	accessTokenRepository := token.NewAccessTokenRepository(ctx, database)
	accessTokenIssuer := token.NewAccessTokenIssuer(accessTokenRepository)
	accessTokenAuthenticator := token.NewAccessTokenAuthenticator(accessTokenRepository)

	refreshTokenRepository := token.NewRefreshTokenRepository(ctx, database)
	refreshTokenIssuer := token.NewRefreshTokenIssuer(refreshTokenRepository)

	userCredentialRepository := user.NewInMemoryCredentialRepository()
	userStatusRepository := user.NewInMemoryStatusRepository()
	userAuthenticator := user.NewAuthenticator(userCredentialRepository, userStatusRepository)

	passwordGrant := token_exchange.NewPasswordGrant(accessTokenIssuer, refreshTokenIssuer, userAuthenticator)

	clientSecretRepository := client.NewSecretRepository(ctx, database)
	clientPrincipalRepository := client.NewPrincipalRepository(ctx, database)
	clientAuthenticator := client.NewAuthenticator(clientSecretRepository, clientPrincipalRepository)

	tokenIntrospector := token_introspection.NewIntrospector(accessTokenAuthenticator)

	//
	// Add Routes
	//
	token_exchange.Route(engine, clientAuthenticator, passwordGrant)
	token_introspection.Route(engine, clientAuthenticator, tokenIntrospector)

	return engine, nil
}
