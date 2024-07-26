package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
)

//url schema: http://localhost:3090/api/v1/hall/{id}/showtimes?movieId=2
// here id is the id of the hall
// TODO: log the errors to a log file

func GetShowTimingsByHallID(w http.ResponseWriter, r *http.Request) {
	hallId , movieId, err := parseHallIDAndMovieID(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()

	showTimings, err := database.GetShowTimingsByID(ctx, hallId, movieId)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError) 
		return
	}

	jsonShowTimings, err := utils.EncodeJson(&showTimings)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError) 
		return
	}
    
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonShowTimings)
}

func parseHallIDAndMovieID(r *http.Request) (int64, int64, error) {
	var hallId, movieId int64
	tempHallId, err := getIDFromPathParameter(r)
	if err != nil {
		return hallId, movieId, err
	}

	hallId = tempHallId

	strMovieId := r.URL.Query().Get("movieId")
	if strMovieId == "" {
		return hallId, movieId, errors.New("missing or empty 'movieId' query parameter")
	}

	movieId, err = strconv.ParseInt(strMovieId, 10, 64)
	if err != nil {
		return hallId, movieId, errors.New("invalid movie id")
	}

	return hallId, movieId, nil
}