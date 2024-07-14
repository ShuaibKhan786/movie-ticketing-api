package handlers

import (
	"net/http"
)

func RegisterUnprotectedRouter(router *http.ServeMux) {
	// router.HandleFunc("GET /movies/ongoing",OnGoing)
	// router.HandleFunc("GET /movies/commingsonn",CommingSoon)
	router.HandleFunc("POST /oauth/provider/signin", SignIn)
	router.HandleFunc("GET /oauth/provider/callback", Callback)
}

func RegisterProtectedRouter(router *http.ServeMux) {
	router.HandleFunc("GET /refresh/token", RefreshToken)
	router.HandleFunc("DELETE /logout", Logout)
	router.HandleFunc("POST /admin/hall/register", HallRegister)
	router.HandleFunc("GET /admin/hall/metadata", HallMetadata)
	router.HandleFunc("GET /profile/details", UserDetails)
}

func RegisterVersion(router *http.ServeMux, versionRouter *http.ServeMux, apiversion string) {
	pattern := apiversion + "/"
	versionRouter.Handle(pattern, http.StripPrefix(apiversion, router))
}
