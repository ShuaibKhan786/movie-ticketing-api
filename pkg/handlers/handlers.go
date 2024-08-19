package handlers

import (
	"net/http"
)

func RegisterUnprotectedRouter(router *http.ServeMux) {
	router.HandleFunc("POST /oauth/provider/signin", SignIn)
	router.HandleFunc("GET /oauth/provider/callback", Callback)
	router.HandleFunc("GET /health", Health)

	router.HandleFunc("GET /movies", GetMovies)
	router.HandleFunc("GET /movie/{id}", GetMovieByID)
	router.HandleFunc("GET /halls", GetHallByMovieID)
	router.HandleFunc("GET /hall/{hall_id}/showtimes", GetShowTimingsByHallID)

	router.HandleFunc("GET /hall/{hall_id}/seatlayout", GetHallSeatLayoutByHallID)
	router.HandleFunc("POST /seats/checkout/{timing_id}", CheckoutSeats)
	//TODO: booking

	router.HandleFunc("GET /movie", SearchMovieByTitle)
}

func RegisterProtectedRouter(router *http.ServeMux) {
	router.HandleFunc("GET /refresh/token", RefreshToken)
	router.HandleFunc("DELETE /logout", Logout)

	router.HandleFunc("POST /admin/hall/register", HallRegister)
	router.HandleFunc("GET /admin/hall/metadata", HallMetadata)
	router.HandleFunc("PATCH /admin/hall/update", UpdateHallMD)

	router.HandleFunc("POST /admin/hall/seatlayout/register", SeatlayoutRegister)
	router.HandleFunc("GET /admin/hall/seatlayout/metdata", SeatLayoutMetadata)
	router.HandleFunc("PATCH /admin/hall/seatlayout/seattype/{seattype_id}", UpdateHallSeatType)
	router.HandleFunc("PATCH /admin/hall/seatlayout/seatrowname/{seatrowname_id}", UpdateHallSeatRowName)

	router.HandleFunc("POST /admin/hall/show/register", ShowRegister)
	//TODO: show patch update

	router.HandleFunc("GET /admin/hall/shows", GetRegisteredShowsByHallID)
	router.HandleFunc("GET /admin/hall/show/{show_id}/timings", GetShowTimingsByShowID)
	router.HandleFunc("POST /admin/hall/show/ticket/release/{timing_id}", SetTimingTicketStatusTrue)
	router.HandleFunc("POST /admin/seats/checkout/{timing_id}", CheckoutSeats)
	//TODO:
	//		booking think ?

	router.HandleFunc("POST /admin/hall/show/timings/avilability", CheckTimingsAvilablity)
	router.HandleFunc("GET /admin/cast", SearchCast)
	router.HandleFunc("POST /admin/image/upload", UploadImage)

	router.HandleFunc("GET /profile/details", UserDetails)
}

func RegisterVersion(router *http.ServeMux, versionRouter *http.ServeMux, apiversion string) {
	pattern := apiversion + "/"
	versionRouter.Handle(pattern, http.StripPrefix(apiversion, router))
}
