package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

//url schema: http://localhost:3090/api/v1/hall/{id} 
// here id means movieID

func GetHallByMovieID(w http.ResponseWriter, r *http.Request) {
	movieId, err := getIDFromPathParameter(r)
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
		utils.JSONResponse(&w, "error encoding hall details to JSON", http.StatusInternalServerError)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonHallDetails)
}