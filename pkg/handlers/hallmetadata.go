package handlers

import (
	"context"
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)


func HallMetadata(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//all the hall metadata will be retrive that has a relationship with that admin
	hallMetaData, err := database.GetHallMetadata(ctx, hallId)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}


	//encode hallMetaData struct to json data
	jsonHallMetaData, err := utils.EncodeJson(&hallMetaData) 
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	//finally send that hall metadata to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonHallMetaData)
}