package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

//url schema: http://localhost:3090/api/v1/hall/{hall_id}/seatlayout
func GetHallSeatLayoutByHallID(w http.ResponseWriter, r *http.Request) {
	hallID, err := getPathParameterValueInt64(r, "hall_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	seatLayout, err := database.GetHallSeatLayoutByHallID(ctx, hallID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonSeatLayout, err := utils.EncodeJson(&seatLayout)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonSeatLayout)
}