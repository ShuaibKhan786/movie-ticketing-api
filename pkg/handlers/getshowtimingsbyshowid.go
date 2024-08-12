package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
)

//url schema: http://localhost:3090/api/v1/auth/admin/hall/show/{show_id}/timings
func GetShowTimingsByShowID(w http.ResponseWriter, r *http.Request) {
	showID , err := getPathParameterValueInt64(r, "show_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()

	showTimings, err := database.GetShowTimingsByShowID(ctx, showID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError) 
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