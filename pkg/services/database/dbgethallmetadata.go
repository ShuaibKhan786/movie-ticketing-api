package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

func GetHallMetadata(ctx context.Context,adminId int64) (models.Hall, error) {
	var hallMetaData models.Hall
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return hallMetaData, fmt.Errorf("begin transaction : %w",err)
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

	actualHallData, hallId, err := getActualHallData(ctx, tx, adminId)
	if err != nil {
		tx.Rollback()
		return hallMetaData, err
	}
	
	location, err := getHallLocationData(ctx, tx, hallId)
	if err != nil {
		tx.Rollback()
		return hallMetaData, err
	}

	seatlayout, err := getHallSeatlayoutData(ctx, tx, hallId)
	if err != nil {
		tx.Rollback()
		return hallMetaData, err
	}

	operationTime, err := getHallOperationTimeData(ctx, tx, hallId)
	if err != nil {
		tx.Rollback()
		return hallMetaData, err
	}

	hallMetaData.Name = actualHallData.Name
	hallMetaData.Manager = actualHallData.Manager
	hallMetaData.Contact = actualHallData.Contact
	hallMetaData.Location = location
	hallMetaData.SeatLayout = seatlayout
	hallMetaData.OperationTime = operationTime

	return hallMetaData, nil
}


func getActualHallData(ctx context.Context, tx *sql.Tx, adminId int64) (models.ActualHall, int64, error) {
	var actualHall models.ActualHall
	var hallId int64

	stmt, err := tx.PrepareContext(ctx, `SELECT id, hall_name, hall_manager, hall_contact, admin_id FROM hall WHERE admin_id=?`)
	if err != nil {
		return actualHall, hallId, fmt.Errorf("prepare context : %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, adminId)
	if err := row.Scan(&hallId, 
		&actualHall.Name, 
		&actualHall.Manager, 
		&actualHall.Contact, 
		&actualHall.AdminId); err != nil {
		return actualHall, hallId, fmt.Errorf("query row exec : %w", err)
	}

	return actualHall, hallId, nil
}


func getHallLocationData(ctx context.Context, tx *sql.Tx, hallId int64) (models.Location, error) {
	var location models.Location

	stmt, err := tx.PrepareContext(ctx, `SELECT address, city, state, postal_code, latitude, longitude FROM hall_location WHERE hall_id=?`)
	if err != nil {
		return location, fmt.Errorf("prepare context : %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, hallId)
	if err := row.Scan(&location.Address,
		&location.City,
		&location.State,
		&location.PostalCode,
		&location.Latitude,
		&location.Longitude); err != nil {
		return location, fmt.Errorf("query row exec : %w", err)
	}

	return location, nil
}

func getHallSeatlayoutData(ctx context.Context, tx *sql.Tx, hallId int64) (models.SeatLayout, error) {
	var seatlayout models.SeatLayout

	stmt, err := tx.PrepareContext(ctx, `SELECT max_capacity, h_rows, h_columns, types, layout FROM hall_seat_layout WHERE hall_id=?`)
	if err != nil {
		return seatlayout, fmt.Errorf("prepare context : %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, hallId)
	if err := row.Scan(&seatlayout.MaxCapacity,
		&seatlayout.Rows,
		&seatlayout.Columns,
		&seatlayout.Types,
		&seatlayout.Layout); err != nil {
		return seatlayout, fmt.Errorf("query row exec : %w", err)
	}

	return seatlayout, nil
}

func getHallOperationTimeData(ctx context.Context, tx *sql.Tx, hallId int64) (models.OperationTime, error) {
	var operation models.OperationTime

	stmt, err := tx.PrepareContext(ctx, `SELECT open_time, closed_time FROM hall_operation_time WHERE hall_id=?`)
	if err != nil {
		return operation, fmt.Errorf("prepare context : %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, hallId)
	if err := row.Scan(&operation.OpenTime,
		&operation.CloseTime); err != nil {
		return operation, fmt.Errorf("query row exec : %w", err)
	}

	return operation, nil
}