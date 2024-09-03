package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
)

//url schema: http://localhost:3090/api/v1/hall/{hall_id}/showtimes?movie_id=2
func GetShowTimingsByHallID(w http.ResponseWriter, r *http.Request) {
	hallId , err := getPathParameterValueInt64(r, "hall_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	movieId, err := getQueryValueInt64(r, "movie_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()

	showTimings, err := database.GetShowTimingsByID(ctx, hallId, movieId)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError) 
		return
	}

	if len(showTimings) == 0 {
		w.WriteHeader(http.StatusNoContent)
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
