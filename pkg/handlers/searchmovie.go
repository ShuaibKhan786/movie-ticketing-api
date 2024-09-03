package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// url schema: http://localhost:3090/api/v1/movie?search_title="movieName"
// this is common route
// will be used by admin and user
func SearchMovieByTitle(w http.ResponseWriter, r *http.Request) {
	title, err := getMovieTitleFromQuery(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	movies, err := database.SearchMoviesByTitle(ctx, title)
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

func getMovieTitleFromQuery(r *http.Request) (string, error) {
	queryValue := r.URL.Query().Get("search_title")
	if queryValue == "" {
		return "", errors.New("missing or empty 'title' query parameter")
	}

	queryValue = strings.Trim(queryValue,`"`)
	return queryValue, nil
}