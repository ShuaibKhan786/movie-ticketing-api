package main

import (
	"log"
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/handlers"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/middlewares"
)

func main() {
	//loading the env into a global Env of ENV struct
	if !config.LoadConfig() {
		log.Fatal("Error in loading the configuration")
	}

	// Unprotected routes (no token auth required)
	unprotectedRouter := http.NewServeMux()
	handlers.RegisterUnprotectedRouter(unprotectedRouter)

	// Protected routes (token auth required)
	protectedRouter := http.NewServeMux()
	handlers.RegisterProtectedRouter(protectedRouter)

	// account routes (signup, login)
	accountRouter := http.NewServeMux()
	handlers.RegisterAccountRouter(accountRouter)

	// Middleware stack for account routes
	accountMiddlewareStack := middlewares.CreateStack(
		middlewares.IsValidJSONCred,
		middlewares.IsValidCredentials,
		middlewares.IsEmailExists,
	)
	accountRouterWithMiddleware := accountMiddlewareStack(accountRouter)

	protectedRouterWithMiddleware := middlewares.EnsureTokenAuth(protectedRouter)

	// Combine unprotected and protected routers with authentication middleware
	unprotectedRouter.Handle("/account/", http.StripPrefix("/account", accountRouterWithMiddleware)) // merging with main router
	unprotectedRouter.Handle("/auth/", http.StripPrefix("/auth", protectedRouterWithMiddleware))

	// Add versioning to the routes
	versionRouter := http.NewServeMux()
	handlers.RegisterVersion(unprotectedRouter, versionRouter, config.APIversion)

	// Common middleware stack
	commonMiddlewareStack := middlewares.CreateStack(
		middlewares.Logging,
		middlewares.AllowCors,
	)

	// Final server setup with all middlewares and routers
	server := http.Server{
		Addr:    ":8090",
		Handler: commonMiddlewareStack(versionRouter),
	}

	server.ListenAndServe()
}
