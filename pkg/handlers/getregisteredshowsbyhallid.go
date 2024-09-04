package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// url schema: http://localhost:3090/api/v1/auth/admin/hall/shows?status=released&page=1&size=3
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

	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	// Ensure 'page' is provided and is valid
	if pageStr == "" {
		utils.JSONResponse(&w, "page query parameter is required", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.JSONResponse(&w, "invalid page number", http.StatusBadRequest)
		return
	}

	// Set default size if not provided
	size := 5
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			utils.JSONResponse(&w, "invalid size number", http.StatusBadRequest)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shows, err := database.GetRegisteredShowsByID(ctx, hallID, status, page, size)
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
