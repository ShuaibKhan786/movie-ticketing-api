package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


func GetHallSeatLayoutByHallID(ctx context.Context, hallID int64) (models.SeatLayout, error){
	var seatlayout models.SeatLayout

	query := `
		SELECT 
			st.name, 
			st.price, 
			st.seat_row, 
			st.seat_col, 
			st.seat_matrix, 
			st.order_from_screen, 
			GROUP_CONCAT(str.row_name) AS row_names
		FROM seat_type st
		INNER JOIN hall_seat_layout hsl ON hsl.id = st.hall_seat_layout_id
		INNER JOIN seat_type_row_name str ON st.id = str.seat_type_id
		WHERE hsl.hall_id = ?
		GROUP BY st.id;
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return seatlayout, fmt.Errorf("hall_seat_layout: prepare query %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, hallID)
	if err != nil {
		return seatlayout, fmt.Errorf("hall_seat_layout: executing query %v", err)
	}

	for rows.Next() {
		var seatType models.SeatType
		var rowNameString string
		if err := rows.Scan(
			&seatType.Name,
			&seatType.Price,
			&seatType.SeatRow,
			&seatType.SeatColumn,
			&seatType.SeatMatrix,
			&seatType.OrderFromScreen,
			&rowNameString,
		); err != nil {
			return seatlayout, fmt.Errorf("hall_seat_layout: scanning from the result %v", err)
		}

		rowNameSlice := strings.Split(rowNameString, ",")
		seatType.RowName = rowNameSlice

		seatlayout.SeatTypes = append(seatlayout.SeatTypes, seatType)
	}

	return seatlayout, nil
}


