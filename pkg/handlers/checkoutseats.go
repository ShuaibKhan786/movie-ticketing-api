package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type ActualSeats struct {
	Seats  []string `json:"seats"`
}

type Seats struct {
	ID     int64    `json:"id"`
	Counts int      `json:"counts"`
	ActualSeats
}

type SeatsCheckoutMD struct {
	Seats
	database.CheckoutMD
	TotalAmount int `json:"total_amount"`
}

// url schema: http://localhost:3090/api/v1/seats/checkout/{timing_id}
//	REQUEST payload:
//	{
//		"id": 1,
//		"counts": 2,
//	    "seats": ["d3"]
//	}
// RESPONSE payload:
// {
// 		"id": 1,
// 		"counts": 2,
// 	    "seats": ["d11", "d9"]
//		"name": "diamond"
//		"price": 250
//		"total_amount": 750
// 		"movie_name": "Inception",
// 		"hall_name": "IMAX Screen 1",
// 		"show_timing": "2024-08-22T19:00:00"
// }
func CheckoutSeats(w http.ResponseWriter, r *http.Request) {
	var role string
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if ok {
		if claims.Role == config.AdminRole {
			role = config.AdminRole
		} else {
			role = config.UserRole
		}
	} else {
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

	if len(seats.Seats) < 1 || len(seats.Seats) != seats.Counts {
		utils.JSONResponse(&w, "empty seats or mismatch seats", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	seatsMD, err := database.GetCheckoutMDbyID(ctx, seats.ID, timingID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	bookedSeatSchema := &redisdb.BookedSeatSchema{
		TimingID: timingID,
		Seat:     seats.Seats,
	}


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
		responseWithAvailableSeats(&w, bookedSeats, http.StatusBadRequest)
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
		responseWithAvailableSeats(&w, reservedSeats, http.StatusBadRequest)
		return
	}

	responseWithSeatsAndMetaData(&w, seats, seatsMD)
}

func responseWithAvailableSeats(w *http.ResponseWriter, seats []string, statusCode int) {
	s := &ActualSeats{
		Seats: seats,
	}

	jsonSeatRB, err := utils.EncodeJson(s)
	if err != nil {
		utils.JSONResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(statusCode)
	(*w).Write(jsonSeatRB)
}

func responseWithSeatsAndMetaData(w *http.ResponseWriter, seats Seats, seatsMD database.CheckoutMD) {
	// Calculate total amount based on seat metadata
	totalAmount := seats.Counts * seatsMD.Price

	seatsCheckoutMD := SeatsCheckoutMD{
		Seats:                 seats,
		CheckoutMD: seatsMD,
		TotalAmount:           totalAmount,
	}

	jsonSeatRB, err := utils.EncodeJson(seatsCheckoutMD)
	if err != nil {
		utils.JSONResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(jsonSeatRB)
}
