package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


func UpdateSeatRowName(ctx context.Context, seatRowNameID int64, updates models.SeatRowNameUpdate) error {
	query := `
		UPDATE seat_type_row_name
		SET row_name = ?
		WHERE id = ?;
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to preapre the query seat_type_row_name: %d: %w", seatRowNameID, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, updates.RowName, seatRowNameID)
	if err != nil {
		return fmt.Errorf("failed to update seat_type_row_name: %d: %w", seatRowNameID, err)
	}

	return nil
}