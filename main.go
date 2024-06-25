package main

import (
	"net/http"

	handlers "github.com/ShuaibKhan786/movie-ticketing-api/pkg/handlers"
)


func main() {
	normalRouter := http.NewServeMux() //routes that need no auth
	handlers.RegisterNormalRouter(normalRouter)

	adminRouter := http.NewServeMux() //routes that needs auth
	handlers.RegisterAdminRouter(adminRouter)

	normalRouter.Handle("/",adminRouter) //merging the routes

	versionRouter := http.NewServeMux() //adding a version to all the routes
	handlers.RegisterVersion(normalRouter,versionRouter)

	server := http.Server{
		Addr: ":8090",
		Handler: versionRouter,
	}

	server.ListenAndServe()
}