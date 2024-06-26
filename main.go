package main

import (
	"net/http"

	handlers "github.com/ShuaibKhan786/movie-ticketing-api/pkg/handlers"
	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	middlewares "github.com/ShuaibKhan786/movie-ticketing-api/pkg/middlewares"
)


func main() {
	unprotectedRouter := http.NewServeMux() //routes that need no auth
	handlers.RegisterUnprotectedRouter(unprotectedRouter)

	protectedRouter := http.NewServeMux() //routes that needs auth
	handlers.RegisterProtectedRouter(protectedRouter)

	unprotectedRouter.Handle("/",middlewares.EnsureAuth(protectedRouter)) //merging the routes

	versionRouter := http.NewServeMux() //adding a version to all the routes
	handlers.RegisterVersion(unprotectedRouter,versionRouter,config.APIversion)

	middlewareStack := middlewares.CreateStack(
		middlewares.Logging,
		middlewares.AllowCors,
	)

	server := http.Server{
		Addr: ":8090",
		Handler: middlewareStack(versionRouter),
	}

	server.ListenAndServe()
}