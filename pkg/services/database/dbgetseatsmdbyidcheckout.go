package database

import (
	"context"
	"fmt"
)

type CheckoutMD struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	MovieName   string `json:"movie_name"`
	HallName    string `json:"hall_name"`
	ShowDate    string `json:"show_date"`
	ShowTiming  string `json:"show_timing"`
}

func GetCheckoutMDbyID(ctx context.Context, seat_type_id int64, timing_id int64) (CheckoutMD, error) {
	query := `
		SELECT 
			st.name,
			st.price,
			m.title AS movie_name,
			h.name AS hall_name,
			msd.show_date,
			mst.show_timing
		FROM seat_type st
		JOIN hall_seat_layout hsl ON st.hall_seat_layout_id = hsl.id
		JOIN hall h ON hsl.hall_id = h.id
		JOIN movie_show ms ON h.id = ms.hall_id
		JOIN movie m ON ms.movie_id = m.id
		JOIN movie_show_dates msd ON ms.id = msd.movie_show_id
		JOIN movie_show_timings mst ON msd.id = mst.movie_show_dates_id
		WHERE st.id = ? AND mst.id = ?
	`

	var seatMetadata CheckoutMD

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return seatMetadata, fmt.Errorf("failed to prepare query: seat_type: %d, timing_id: %d: %w", seat_type_id, timing_id, err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, seat_type_id, timing_id).Scan(
		&seatMetadata.Name,
		&seatMetadata.Price,
		&seatMetadata.MovieName,
		&seatMetadata.HallName,
		&seatMetadata.ShowDate,
		&seatMetadata.ShowTiming,
	)
	if err != nil {
		return seatMetadata, fmt.Errorf("failed to query result: seat_type: %d, timing_id: %d: %w", seat_type_id, timing_id, err)
	}

	return seatMetadata, nil
}
