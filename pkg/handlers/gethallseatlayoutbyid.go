package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

//url schema: http://localhost:3090/api/v1/hall/{hall_id}/seatlayout?timing_id=id
func GetHallSeatLayoutByHallID(w http.ResponseWriter, r *http.Request) {
	hallID, err := getPathParameterValueInt64(r, "hall_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	timingID, err := getQueryValueInt64(r, "timing_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	seatLayout, err := database.GetHallSeatLayoutUserByHallID(ctx, hallID, timingID)
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

	if len(seatLayout.ReservedSeats) > 1 {
		go redisReservedSeatsCleanup(timingID)
	}
}

func redisReservedSeatsCleanup(timingID int64) {
	t := redisdb.TimingID(timingID)
	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Second)
	defer cancel()
	t.CleanupReservedSeats(ctx)
	//make sure the error in the log
	fmt.Println("It cleanup")
}