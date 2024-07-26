package handlers

import (
	"net/http"
)

func RegisterUnprotectedRouter(router *http.ServeMux) {
	router.HandleFunc("POST /oauth/provider/signin", SignIn)
	router.HandleFunc("GET /oauth/provider/callback", Callback)
	router.HandleFunc("GET /health", Health)
	router.HandleFunc("GET /movie", SearchMovieByTitle)
	router.HandleFunc("GET /movies", GetMovies)
	router.HandleFunc("GET /movie/{id}", GetMovieByID)
	router.HandleFunc("GET /hall/{id}", GetHallByMovieID)
	router.HandleFunc("GET /hall/{id}/showtimes", GetShowTimingsByHallID)
}

func RegisterProtectedRouter(router *http.ServeMux) {
	router.HandleFunc("GET /refresh/token", RefreshToken)
	router.HandleFunc("DELETE /logout", Logout)
	router.HandleFunc("POST /admin/hall/register", HallRegister)
	router.HandleFunc("GET /admin/hall/metadata", HallMetadata)
	router.HandleFunc("POST /admin/hall/show/register", ShowRegister)
	router.HandleFunc("GET /admin/hall/shows", GetRegisteredShowsByHallID)
	router.HandleFunc("POST /admin/hall/show/timings/avilability", CheckTimingsAvilablity)
	router.HandleFunc("GET /admin/cast", SearchCast)
	router.HandleFunc("GET /profile/details", UserDetails)
}

func RegisterVersion(router *http.ServeMux, versionRouter *http.ServeMux, apiversion string) {
	pattern := apiversion + "/"
	versionRouter.Handle(pattern, http.StripPrefix(apiversion, router))
}
