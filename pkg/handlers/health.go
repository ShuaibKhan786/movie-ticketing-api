package handlers

import (
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(
		&w,
		"server is healthy",
		http.StatusOK,
	)
}