package handlers

import (
	"context"
	"net/http"
	"time"
	"io"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)
//url Schema: POST http://localhost:3090/api/v1/admin/hall/show/timings/avilability
//payload: 
// [
//     {
//         "show_date": "2024-07-23",
//         "show_timing": [
//             "09:00:00"
//         ]
//     },
//		and many more
// ]
// dates must be send in ascending order e.g. : 2024-07-23, 2024-07-24, 2024-07-25


func CheckTimingsAvilablity(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	hallId, err := isHallRegistered(claims)
	if err != nil {
		utils.JSONResponse(&w, "hall not registered", http.StatusBadRequest)
		return
	}

	var timings []models.ShowDate
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := utils.DecodeJson(body, &timings); err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	conflictTimings, err := database.GetConflictTimings(ctx, hallId, timings)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonConflictTimings, err := utils.EncodeJson(&conflictTimings)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonConflictTimings)
}