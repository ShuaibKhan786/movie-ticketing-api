package handlers

import (
	"context"
	"net/http"
	"strings"

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

	adminId := claims.Id
	parentCtx := context.TODO()
	childCtx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	//all the hall metadata will be retrive that has a relationship with that admin
	hallMetaData, err := database.GetHallMetadata(childCtx, adminId)
	if err != nil {
		if strings.Contains(err.Error(),"no rows in result set") {
			utils.JSONResponse(&w, "no hall data found with that admin", http.StatusBadRequest)
			return
		}
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
	w.WriteHeader(http.StatusOK)
	w.Write(jsonHallMetaData)
}