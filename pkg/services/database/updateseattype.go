package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


func UpdateSeatType(ctx context.Context, seatTypeID int64, updates models.SeatType) error {
	query, values, err := constructQueryAndValues(seatTypeID, updates)
	if err != nil {
		return err
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to preapre the query seat_type: %d: %w", seatTypeID, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		return fmt.Errorf("failed to update seat_type: %d: %w", seatTypeID, err)
	}

	return nil
}
func constructQueryAndValues(seatTypeID int64, updates models.SeatType) (string, []interface{}, error) {
	query := `UPDATE seat_type SET %s WHERE id = ?`
	values := []interface{}{}
	columns := []string{}

	if updates.Name != nil {
		columns = append(columns, "name = ?")
		values = append(values, *updates.Name)
	}
	if updates.Price != nil {
		columns = append(columns, "price = ?")
		values = append(values, *updates.Price)
	}
	if updates.SeatRow != nil {
		columns = append(columns, "seat_row = ?")
		values = append(values, *updates.SeatRow)
	}
	if updates.SeatColumn != nil {
		columns = append(columns, "seat_col = ?")
		values = append(values, *updates.SeatColumn)
	}
	if updates.SeatMatrix != nil {
		columns = append(columns, "seat_matrix = ?")
		values = append(values, *updates.SeatMatrix)
	}
	if updates.OrderFromScreen != nil {
		columns = append(columns, "order_from_screen = ?")
		values = append(values, *updates.OrderFromScreen)
	}

	if len(columns) == 0 {
		return "", []interface{}{}, errors.New("no fields to update to the seat type")
	}

	values = append(values, seatTypeID)

	columnsString := strings.Join(columns, ", ")

	finalQuery := fmt.Sprintf(query, columnsString)

	return finalQuery, values, nil
}