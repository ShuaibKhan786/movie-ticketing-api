package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// url schema: http://localhost:3090/api/v1/auth/admin/hall/show/{show_id}/timings?status="released"
func GetShowTimingsByShowID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	showID, err := getPathParameterValueInt64(r, "show_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var jsonShowTimings []byte

	if strings.Trim(status, `"`) == "released" {
		movieID, err := database.GetId("movie_show", "id", showID)
		if err != nil {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}
		hallID, err := database.GetId("movie_show", "hall_id", showID)
		if err != nil {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		showTimings, err := database.GetShowTimingsByID(ctx, hallID, movieID)
		if err != nil {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonShowTimings, err = utils.EncodeJson(&showTimings)
		if err != nil {
			utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		showTimings, err := database.GetShowTimingsByShowID(ctx, showID)
		if err != nil {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonShowTimings, err = utils.EncodeJson(&showTimings)
		if err != nil {
			utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonShowTimings)
}
