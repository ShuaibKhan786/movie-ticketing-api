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

//url schema: http://localhost:3090/api/v1/movie/{id}
// using the id which is a movie id 
// will send all the movie details

func GetMovieByID(w http.ResponseWriter, r *http.Request) {
	movieId, err := getMovieIDFromPathParameter(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	movieDetails, err := database.GetMovieDetailsByID(ctx, movieId)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonMovieDetails, err := utils.EncodeJson(&movieDetails)
	if err != nil {
		utils.JSONResponse(&w, "error encoding movies details to JSON", http.StatusNotFound)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMovieDetails)
}

func getMovieIDFromPathParameter(r *http.Request) (int64, error) {
	movieIDStr := r.PathValue("id")
	if movieIDStr == "" {
		return 0, errors.New("missing movie id path parameter")
	}

	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid movie id")
	}

	return movieID, nil
}