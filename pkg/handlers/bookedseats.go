package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// URL_SCHEMA:
//
//	admin: http://localhost:3090/api/v1/seats/book/{timing_id}
//	user: http://localhost:3090/api/v1/auth/admin/seats/book/{timing_id}
//
// REQUEST
// payload
//
//	{
//		"id": 2,
//		"counts": 2,
//		"seats": ["d9", "d11"],
//		"payable_amount": 350,
//		"payment_mode": "upi",
//		"customer_phone_no": "7301271044",
//		"cash_amount": 400
//	}
//
// RESPONSE
// payload
//
//	{
//		"customer_phone_no": "7301271044",
//		"movie_name": "Example Movie",
//		"hall_name": "Example Hall",
//		"show_date": "2024-08-22",
//		"show_time": "19:00",
//		"tickets": [
//			{
//				"ticket_number": "UUID-12345",
//				"seat_number": "d9",
//			},
//			{
//				"ticket_number": "UUID-67891",
//				"seat_number": "d11",
//			}
//		]
//	}
func BookedSeats(w http.ResponseWriter, r *http.Request) {
	var role string
	var userID *int64
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if ok {
		if claims.Role == config.AdminRole {
			role = config.AdminRole
		} else {
			role = config.UserRole
		}
		userID = &claims.Id
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

	var details models.BookedRequestPayload
	err = utils.DecodeJson(body, &details)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !utils.ValidateBookedRequestPayload(details) {
		utils.JSONResponse(&w, "missing payload", http.StatusBadRequest)
		return
	}

	//there can be no zero seats to booked
	if len(*details.Seats) < 1 || len(*details.Seats) != *details.Counts {
		utils.JSONResponse(&w, "empty seats or mismatch seats", http.StatusBadRequest)
		return
	}

	//only admin can proceed with cash payment
	if role == config.UserRole && *details.PaymentMode == config.Cash {
		utils.JSONResponse(&w, "invalid payment mode", http.StatusBadRequest)
		return
	}

	//if role is admin and he choose cash and there is no amount then it wont work
	if role == config.AdminRole {
		if *details.PaymentMode == config.Cash && details.CashAmount == nil {
			utils.JSONResponse(&w, "invalid payment mode", http.StatusBadRequest)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	seatsMD, err := database.GetCheckoutMDbyID(ctx, *details.ID, timingID)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	//payable amount must same after calculated
	payableAmount := seatsMD.Price * (*details.Counts)
	if payableAmount != *details.PayableAmount {
		utils.JSONResponse(&w, "invalid payment amount", http.StatusBadRequest)
		return
	}

	tickets, err := database.UpdateBookingSeats(ctx, userID, timingID, details, seatsMD, role)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	ticketJson, err := utils.EncodeJson(tickets)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(ticketJson)
}
