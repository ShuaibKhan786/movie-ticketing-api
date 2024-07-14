package handlers

import (
	"net/http"
)

func RegisterUnprotectedRouter(router *http.ServeMux) {
	router.HandleFunc("GET /movies/ongoing", OnGoing)
	router.HandleFunc("GET /movies/commingsoon", CommingSoon)
}

func RegisterAccountRouter(router *http.ServeMux) {
	router.HandleFunc("POST /user/signup", Signup)
	router.HandleFunc("POST /admin/signup", Signup)
	router.HandleFunc("POST /login", Login)
}

func RegisterProtectedRouter(router *http.ServeMux) {
	router.HandleFunc("POST /user/details", UserDetails)
}

func RegisterVersion(router *http.ServeMux, versionRouter *http.ServeMux, apiversion string) {
	pattern := apiversion + "/"
	versionRouter.Handle(pattern, http.StripPrefix(apiversion, router))
}
