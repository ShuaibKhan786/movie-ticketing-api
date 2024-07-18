package handlers

import (
	"net/http"
	"context"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
)

func GetHallByMovieID(w http.ResponseWriter, r *http.Request) {
	movieId, err := getMovieIDFromPathParameter(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return		
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	hallDetails, err := database.GetHallDetailsByID(ctx, movieId)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonHallDetails, err := utils.EncodeJson(&hallDetails)
	if err != nil {
		utils.JSONResponse(&w, "error encoding hall details to JSON", http.StatusNotFound)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonHallDetails)
}