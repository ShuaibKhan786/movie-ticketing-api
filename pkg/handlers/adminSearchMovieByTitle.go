package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// url schema: http://localhost:3090/api/v1/auth/admin/movie?search_title="movieName"
// this route is specific for searching the registered show/movie by admin
func AdminSearchMovieByTitle(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	hallID, err := isHallRegistered(claims)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
	}

	title, err := getMovieTitleFromQuery(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movies, err := database.AdminSearchMoviesByTitle(ctx, hallID, title)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(movies) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonMovies, err := utils.EncodeJson(&movies)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(jsonMovies)
}
