package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

//url schema: http://localhost:3090/api/v1/movies?status=incinemas&page=1&size=3
// here status can be either "incinemas" or "upcoming"
// page and size is for pagination

type MoviesQuery struct{
	Status string
	Date string
	Limit int
	Offset int
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	moviesQuery, err := getMoviesQueryFromQuery(r)	
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()

	movies, err := database.GetMoviesByStatus(ctx,
		moviesQuery.Status,
		moviesQuery.Date,
		moviesQuery.Limit,
		moviesQuery.Offset)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	if (len(movies) == 0) {
		utils.JSONResponse(&w, "no more movies", http.StatusNotFound)
		return
	}

	jsonMovies, err := utils.EncodeJson(&movies)
	if err != nil {
		utils.JSONResponse(&w, "error encoding movies to JSON", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMovies)
}


func getMoviesQueryFromQuery(r *http.Request) (MoviesQuery, error){
	status := r.URL.Query().Get("status")
	if status == "" {
		return MoviesQuery{}, errors.New("missing or empty 'status' query parameter")
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		return MoviesQuery{}, errors.New("missing or empty 'page' query parameter")
	}
	size := r.URL.Query().Get("size")
	if size == "" {
		return MoviesQuery{}, errors.New("missing or empty 'size' query parameter")
	}

	if !isThisValidStatus(status) {
		return MoviesQuery{}, errors.New("invalid status")
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return MoviesQuery{}, errors.New("invalid page")
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return MoviesQuery{}, errors.New("invalid size")
	}

	offset := (pageInt - 1) * sizeInt
	date := time.Now().Format(time.DateOnly) // returns in this format yyyy-mm-dd

	return MoviesQuery{
		Status: status,
		Limit: sizeInt,
		Offset: offset,
		Date: date,
	}, nil
}

func isThisValidStatus(status string) bool {
	switch status {
	case "incinemas",
		"upcoming":
		return true
	default:
		return false
	}
}