package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

func RegisterSeatLayout(ctx context.Context, hallId, adminId int64, seatLayout models.SeatLayout) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
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

	seatlayoutID, err := registerHallSeatLayout(ctx, tx, hallId)
	if err != nil {
		return fmt.Errorf("seat register: %w", err)
	}

	err = registerHallSeatLayoutTypes(ctx, tx, seatlayoutID, seatLayout.SeatTypes)
	if err != nil {
		return fmt.Errorf("seat layout types register: %w", err)
	}

	err = updateAdminHallXColumnToTrue(ctx, tx, adminId, "hall_seat_layout_registered")
	if err != nil {
		return fmt.Errorf("failed to update the hall seatLayout boolean in admin: %w", err)
	}

	return nil
}

func registerHallSeatLayout(ctx context.Context, tx *sql.Tx, hallId int64) (int64, error) {
	var seatlayoutID int64
	query := `
		INSERT INTO hall_seat_layout (hall_id)
		VALUES (?);
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return seatlayoutID, fmt.Errorf("preparing hall_seat_layout query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, hallId)
	if err != nil {
		return seatlayoutID, fmt.Errorf("executing hall_seat_layout query: %w", err)
	}

	seatlayoutID, err = res.LastInsertId()
	if err != nil {
		return seatlayoutID, fmt.Errorf("last insert id for hall_seat_layout query: %w", err)
	}

	return seatlayoutID, nil
}

func registerHallSeatLayoutTypes(ctx context.Context, tx *sql.Tx, seatlayoutID int64, seatTypes []models.SeatType) error {
	query := `
		INSERT INTO seat_type
		(name, price, seat_row, seat_col, seat_matrix, order_from_screen, hall_seat_layout_id)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("preparing seat_type query: %w", err)
	}
	defer stmt.Close()

	for _, seatType := range seatTypes {
		res, err := stmt.ExecContext(
			ctx,
			seatType.Name,
			seatType.Price,
			seatType.SeatRow,
			seatType.SeatColumn,
			seatType.SeatMatrix,
			seatType.OrderFromScreen,
			seatlayoutID,
		)
		if err != nil {
			return fmt.Errorf("executing seat_type query: %w", err)
		}

		seatTypeID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("last insert id for seat_type query: %w", err)
		}

		if err := registerSeatTypeRowNames(ctx, tx, seatTypeID, seatType.RowName); err != nil {
			return fmt.Errorf("registering seat type row names: %w", err)
		}
	}

	return nil
}

func registerSeatTypeRowNames(ctx context.Context, tx *sql.Tx, seatTypeID int64, rowNames []string) error {
	query := `
		INSERT INTO seat_type_row_name (row_name, seat_type_id)
		VALUES (?, ?);
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("preparing seat_type_row_name query: %w", err)
	}
	defer stmt.Close()

	for _, rowName := range rowNames {
		if _, err := stmt.ExecContext(ctx, rowName, seatTypeID); err != nil {
			return fmt.Errorf("executing seat_type_row_name query: %w", err)
		}
	}

	return nil
}
