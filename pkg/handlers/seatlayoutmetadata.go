package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func SeatLayoutMetadata(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	hallID, err := isHallRegistered(claims)
	if err != nil {
		utils.JSONResponse(&w, "hall not registered", http.StatusBadRequest)
		return
	}

	err = isSeatLayoutRegistered(claims)
	if err != nil {
		utils.JSONResponse(&w, "hall seatlayout not registered", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	seatLayouts, err := database.GetHallSeatLayoutAdminByHallID(ctx, hallID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	seatLayoutsJSON, err := utils.EncodeJson(&seatLayouts)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(seatLayoutsJSON)
}
