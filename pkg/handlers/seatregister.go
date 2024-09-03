package handlers

import (
	"context"
	"io"
	"net/http"
	"time"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
)

func SeatlayoutRegister(w http.ResponseWriter, r *http.Request) {
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

	err = isSeatLayoutRegistered(claims)
	if err == nil {
		utils.JSONResponse(&w, "seat layout already registered", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body) 
	if err != nil {
		utils.JSONResponse(&w, "failed to read the body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var seatLayout models.SeatLayout
	if err := utils.DecodeJson(body, &seatLayout); err != nil {
		utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := database.RegisterSeatLayout(ctx, hallId, claims.Id, seatLayout); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	redisCtx, redisCancel := context.WithCancel(context.Background())
	defer redisCancel()
	redisKey := fmt.Sprintf("hall:seatlayout:registered:%s:%d", claims.Role, claims.Id)
	if err := redisdb.Set(redisCtx, redisKey, hallId, config.RedisZeroExpirationTime); err != nil {
		utils.JSONResponse(&w, "failed to update Redis", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(&w,"hall seatlayout registered successfully", http.StatusCreated)
}