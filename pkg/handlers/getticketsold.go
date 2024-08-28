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

func GetTicketsSold(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()

	responsePaylod, err := database.GetTicketsSold(ctx, claims.Id) 
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	responsePayloadJSON, err := utils.EncodeJson(&responsePaylod)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responsePayloadJSON)
}
