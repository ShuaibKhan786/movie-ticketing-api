package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

type SeatRowName struct {
	Id      *int64 `json:"id"`
	RowName string `json:"row_name"`
}

type SeatLayoutMD struct {
	Id           *int64        `json:"id"`
	models.SeatType
	SeatRowNameMD []SeatRowName `json:"seat_row_name_md"`
}

func GetHallSeatLayoutAdminByHallID(ctx context.Context, hallID int64) ([]SeatLayoutMD, error) {
	var seatLayouts []SeatLayoutMD
	var currentSeatLayout *SeatLayoutMD

	query := `
		SELECT 
			st.id,
			st.name, 
			st.price, 
			st.seat_row, 
			st.seat_col, 
			st.seat_matrix, 
			st.order_from_screen, 
			str.id,
			str.row_name
		FROM seat_type st
		INNER JOIN hall_seat_layout hsl ON hsl.id = st.hall_seat_layout_id
		INNER JOIN seat_type_row_name str ON st.id = str.seat_type_id
		WHERE hsl.hall_id = ?
		ORDER BY st.id, str.id;
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("hall_seat_layout: prepare query %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, hallID)
	if err != nil {
		return nil, fmt.Errorf("hall_seat_layout: executing query %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var seatRowName SeatRowName
		var seatLayout SeatLayoutMD

		if err := rows.Scan(
			&seatLayout.Id,
			&seatLayout.Name,
			&seatLayout.Price,
			&seatLayout.SeatRow,
			&seatLayout.SeatColumn,
			&seatLayout.SeatMatrix,
			&seatLayout.OrderFromScreen,
			&seatRowName.Id,
			&seatRowName.RowName,
		); err != nil {
			return nil, fmt.Errorf("hall_seat_layout: scanning from the result %v", err)
		}

		// If it's a new SeatLayoutMD, finalize the previous one and start a new one
		if currentSeatLayout != nil && *currentSeatLayout.Id != *seatLayout.Id {
			seatLayouts = append(seatLayouts, *currentSeatLayout)
			currentSeatLayout = &seatLayout
		}

		// Append the row name to the current SeatLayoutMD
		if currentSeatLayout == nil {
			currentSeatLayout = &seatLayout
		}
		currentSeatLayout.SeatRowNameMD = append(currentSeatLayout.SeatRowNameMD, seatRowName)
	}

	// Add the last processed seat layout to the list
	if currentSeatLayout != nil {
		seatLayouts = append(seatLayouts, *currentSeatLayout)
	}

	return seatLayouts, nil
}
