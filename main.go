package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/handlers"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/middlewares"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
)

func main() {
	//just for making sure that 
	//mysql server database is ready for connection
	time.Sleep(10 * time.Second)

	//loading the env into a global Env of ENV struct
	if !config.LoadConfig() {
		log.Fatal("Error in loading the configuration")
	}

	if err := database.InitDB(); err != nil {
		log.Fatal("mysql server", err)
	}
	
	if err := redisdb.InitRedis(); err != nil {
		log.Fatal("redis server", err)
	}


	// Unprotected routes (no token auth required)
	unprotectedRouter := http.NewServeMux()
	handlers.RegisterUnprotectedRouter(unprotectedRouter)

	// Protected routes (token auth required)
	protectedRouter := http.NewServeMux()
	handlers.RegisterProtectedRouter(protectedRouter)


	protectedRouterWithMiddleware := middlewares.EnsureTokenAuth(protectedRouter)

	// Combine unprotected and protected routers with authentication middleware
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
	addr := fmt.Sprintf(":%s",config.Env.PORT)
	server := http.Server{
		Addr:    addr,
		Handler: commonMiddlewareStack(versionRouter),
	}

	server.ListenAndServe()
}
