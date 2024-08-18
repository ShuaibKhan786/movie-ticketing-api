package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type Seats struct {
	Seats []string `json:"seats"`
}

// url schema: http://localhost:3090/api/v1/seats/checkout/{timing_id}
//
//	{
//	    "seats": ["d3"]
//	}
func CheckoutSeats(w http.ResponseWriter, r *http.Request) {
	var role string
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if ok {
		if claims.Role == config.AdminRole {
			role = config.AdminRole
		}else {
			role = config.UserRole
		}
	}else {
		role = config.UserRole
	}

	timingID, err := getPathParameterValueInt64(r, "timing_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid payload", http.StatusBadRequest)
		return
	}

	var seats Seats
	err = utils.DecodeJson(body, &seats)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(seats.Seats) < 1 {
		utils.JSONResponse(&w, "empty seats", http.StatusBadRequest)
		return
	}

	bookedSeatSchema := &redisdb.BookedSeatSchema{
		TimingID: timingID,
		Seat:     seats.Seats,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	bookedSeats, err := bookedSeatSchema.IsSeatAvilableBS(ctx, role)
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			utils.JSONResponse(&w, "seats can no longer checkout", http.StatusBadRequest)
			return
		default:
			utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if len(bookedSeats) > 0 {
		responseWithSeats(&w, bookedSeats, http.StatusBadRequest)
		return
	}

	reservedSeatSchema := &redisdb.ReservedSeatSchema{
		TimingID: timingID,
		Seats:    seats.Seats,
	}

	reservedSeats, err := reservedSeatSchema.ReservedSeatsRegs(ctx, role)
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			utils.JSONResponse(&w, "seats can no longer checkout", http.StatusBadRequest)
			return
		default:
			utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	if reservedSeats != nil {
		responseWithSeats(&w, reservedSeats, http.StatusBadRequest)
		return
	}

	responseWithSeats(&w, seats.Seats, http.StatusOK)
}

func responseWithSeats(w *http.ResponseWriter, seats []string, statusCode int) {
	s := &Seats{
		Seats: seats,
	}

	jsonSeatRB, err := utils.EncodeJson(s)
	if err != nil {
		utils.JSONResponse(w, "internal server error", http.StatusInternalServerError)
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(statusCode)
	(*w).Write(jsonSeatRB)
}
