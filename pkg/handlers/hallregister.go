package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func HallRegister(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}
	
	body, err := io.ReadAll(r.Body) 
	if err != nil {
		utils.JSONResponse(&w, "failed to read the body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	
	//making sure data is send as JSON
	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid request payload", http.StatusBadRequest)
		return
	}

	//decoding the raw body which contains JSON into struct of models.Hall
	var hallMetadata models.Hall
	if err := utils.DecodeJson(body, &hallMetadata); err != nil {
		utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		return
	}

	//making sure the credeentials of hall is validate
	if err := utils.ValidateHall(&hallMetadata); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	//making sure that only one hall be exists
	exists, err := database.IsValueExists("hall", "hall_name", hallMetadata.Name)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		utils.JSONResponse(&w, "a hall with that name already exists", http.StatusConflict)
		return
	}


	parentCtx := context.TODO()
	childCtx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	if err := database.RegisterHall(childCtx, hallMetadata, claims.Id); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	utils.JSONResponse(&w,"hall registered successfully", http.StatusCreated)
}