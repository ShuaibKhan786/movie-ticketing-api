package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// url schema: http://localhost:3090/api/v1/auth/admin/hall/show/ticket/release/{timing_id}
func SetTimingTicketStatusTrue(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	_, err := isHallRegistered(claims)
	if err != nil {
		utils.JSONResponse(&w, "hall not registered", http.StatusBadRequest)
		return
	}

	timingID, err := getPathParameterValueInt64(r, "timing_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	redisTimingID := redisdb.TimingID(timingID)
	ok, err = redisTimingID.CheckBookedSeatExists(ctx)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok {
		utils.JSONResponse(&w, "timing already released for ticket", http.StatusBadRequest)
		return
	}

	ok, err = redisTimingID.CheckReservedSeatExists(ctx)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok {
		utils.JSONResponse(&w, "timing already released for ticket", http.StatusBadRequest)
		return
	}

	if err := database.SetTimingTicketStatusTrue(ctx, timingID); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(&w, "successfully release a ticket", http.StatusOK)
}
