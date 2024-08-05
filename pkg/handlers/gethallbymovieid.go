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

//url schema: http://localhost:3090/api/v1/halls?movie_id=id
func GetHallByMovieID(w http.ResponseWriter, r *http.Request) {
	movieID, err := getQueryValueInt64(r, "movie_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}


	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	hallDetails, err := database.GetHallDetailsByID(ctx, movieID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonHallDetails, err := utils.EncodeJson(&hallDetails)
	if err != nil {
		utils.JSONResponse(&w, "error encoding hall details to JSON", http.StatusInternalServerError)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonHallDetails)
}

func getQueryValueInt64(r *http.Request, key string) (int64, error) {
	var int64ID int64

	stringID := r.URL.Query().Get(key)
	if stringID == "" {
		return int64ID, errors.New("no query parameter")
	}

	var err error
	int64ID, err = strconv.ParseInt(stringID, 10, 64)
	if err != nil {
		return int64ID, errors.New("invalid query")
	}

	return int64ID, nil
}