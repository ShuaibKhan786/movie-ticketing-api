package handlers

import (
	"context"
	"database/sql"
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
	movieId, err := getPathParameterValueInt64(r, "id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	movieDetails, err := database.GetMovieDetailsByID(ctx, movieId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonMovieDetails, err := utils.EncodeJson(&movieDetails)
	if err != nil {
		utils.JSONResponse(&w, "error encoding movies details to JSON", http.StatusInternalServerError)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMovieDetails)
}

func getPathParameterValueInt64(r *http.Request, key string) (int64, error) {
	var int64ID int64

	stringID := r.PathValue(key)
	if stringID == "" {
		return int64ID, errors.New("no path parameter value")
	}

	var err error
	int64ID, err = strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid path parameter value")
	}

	return int64ID, nil
}