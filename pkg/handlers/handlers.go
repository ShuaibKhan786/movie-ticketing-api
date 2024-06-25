package handlers

import (
	"net/http"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
)

func RegisterNormalRouter(router *http.ServeMux) {
	router.HandleFunc("GET /movies/ongoing",OnGoing)
	router.HandleFunc("GET /movies/commingsonn",CommingSoon)
}

func RegisterAdminRouter(router *http.ServeMux) {
	router.HandleFunc("POST /user/register",UserRegister)
	router.HandleFunc("POST /admin/register",AdminRegister)
	router.HandleFunc("POST /login",Login)
}

func RegisterVersion(router *http.ServeMux,versionRouter *http.ServeMux) {
	pattern := config.APIversion+"/"
	versionRouter.Handle(pattern,http.StripPrefix(config.APIversion,router))
}