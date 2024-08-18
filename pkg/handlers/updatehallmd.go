package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func UpdateHallMD(w http.ResponseWriter, r *http.Request) {
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
	body, err := io.ReadAll(r.Body) 
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid payload", http.StatusBadRequest)
        return
	}

	var updateMetaData map[string]map[string]interface{}
	err = utils.DecodeJson(body, &updateMetaData)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.ValidateHallMDUpd(updateMetaData)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	err = database.UpdateHallMetaData(ctx, hallID, updateMetaData)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(&w, "hall updated successfully", http.StatusOK)
}
