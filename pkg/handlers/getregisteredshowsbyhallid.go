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

//TODO: update more details
func GetRegisteredShowsByHallID(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	shows, err := database.GetRegisteredShowsByID(ctx, hallID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(shows) == 0 {
		utils.JSONResponse(&w, "no show has been registered", http.StatusBadRequest)
		return
	}

	jsonShows, err := utils.EncodeJson(&shows)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonShows)
}