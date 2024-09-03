package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func HallRegister(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	//making sure an admin can registered only one hall
	_, err := isHallRegistered(claims) 
	if err == nil {
		utils.JSONResponse(&w, "hall already registered", http.StatusBadRequest)
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

	//making sure that only one hall be exists and does not conflict 
	exists, err := database.IsValueExists("hall", "name", hallMetadata.Name)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		utils.JSONResponse(&w, "a hall with that name already exists", http.StatusConflict)
		return
	}


	//registered the hall Metadata to the database
	databaseCtx, dbCancel := context.WithCancel(context.TODO())
	defer dbCancel()
	var hallId int64
	if err := database.RegisterHall(databaseCtx, hallMetadata, claims.Id, &hallId); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Register hall in Redis
	redisCtx, redisCancel := context.WithCancel(context.Background())
	defer redisCancel()
	redisKey := fmt.Sprintf("hall:registered:%s:%d", claims.Role, claims.Id)
	if err := redisdb.Set(redisCtx, redisKey, hallId, config.RedisZeroExpirationTime); err != nil {
		utils.JSONResponse(&w, "failed to update Redis", http.StatusInternalServerError)
		return
	}

	
	utils.JSONResponse(&w,"hall registered successfully", http.StatusCreated)
}