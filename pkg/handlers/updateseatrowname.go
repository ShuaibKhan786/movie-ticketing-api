package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)
// 	Method: PATCH
// 	url schema: http://localhost:3090/api/v1/auth/admin/hall/seatlayout/seatrowname/{seatrowname_id}
// payload:
// { 
//		"row_names": "F"
// }
func UpdateHallSeatRowName(w http.ResponseWriter, r *http.Request) {
	seatRowNameID, err := getPathParameterValueInt64(r, "seatrowname_id")
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid payload", http.StatusBadRequest)
		return
	}

	var updates models.SeatRowNameUpdate
	err = utils.DecodeJson(body, &updates)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	err = database.UpdateSeatRowName(ctx, seatRowNameID, updates)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}


	utils.JSONResponse(&w, "successfully updated seat type", http.StatusOK)
}