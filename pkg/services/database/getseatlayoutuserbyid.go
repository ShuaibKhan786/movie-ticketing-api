package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
)

type SeatLayout struct {
	models.SeatLayout
	ReservedSeats []string `json:"reserved_seats"`
	BookedSeats []string `json:"booked_seats"`
}

func GetHallSeatLayoutUserByHallID(ctx context.Context, hallID, timingID int64) (SeatLayout, error){
	var seatlayout SeatLayout

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return seatlayout, fmt.Errorf("seat layout: begin transaction : %w",err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = getDefaultSeatLayout(ctx, tx, &seatlayout, hallID)
	if err != nil {
		tx.Rollback()
		return seatlayout, err
	}

	err = getReservedBookedSeats(ctx, &seatlayout, timingID)
	if err != nil {
		tx.Rollback()
		return seatlayout, err
	}

	return seatlayout, nil
}

func getDefaultSeatLayout(ctx context.Context, tx *sql.Tx, seatLayout *SeatLayout, hallID int64) error {
	query := `
		SELECT 
			st.id,
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

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("hall_seat_layout: prepare query %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, hallID)
	if err != nil {
		return fmt.Errorf("hall_seat_layout: executing query %v", err)
	}

	for rows.Next() {
		var seatType models.SeatType
		var rowNameString string
		if err := rows.Scan(
			&seatType.ID,
			&seatType.Name,
			&seatType.Price,
			&seatType.SeatRow,
			&seatType.SeatColumn,
			&seatType.SeatMatrix,
			&seatType.OrderFromScreen,
			&rowNameString,
		); err != nil {
			return fmt.Errorf("hall_seat_layout: scanning from the result %v", err)
		}

		rowNameSlice := strings.Split(rowNameString, ",")
		seatType.RowName = rowNameSlice

		seatLayout.SeatTypes = append(seatLayout.SeatTypes, seatType)
	}

	return nil
}

func getReservedBookedSeats(ctx context.Context, seatLayout *SeatLayout, timingID int64) error {
	t := redisdb.TimingID(timingID)
	
	seatLayout.ReservedSeats, err = t.GetAllReservedSeats(ctx)
	if err != nil {
		return fmt.Errorf("redis: hall_seat_layout: %w", err)
	}

	seatLayout.BookedSeats, err = t.GetAllBookedSeats(ctx)
	if err != nil {
		return fmt.Errorf("redis: hall_seat_layout: %w", err)
	}

	return nil
}