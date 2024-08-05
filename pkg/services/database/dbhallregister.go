package database

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


func RegisterHall(ctx context.Context, hall models.Hall, adminId int64, hallId *int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction : %w",err)
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

	actualHall := models.ActualHall{
		Name: hall.Name,
		Manager: hall.Manager,
		Contact: hall.Contact,
		AdminId: adminId,
	}

	
	tempHallId, err := registerActualHall(ctx, tx, actualHall)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := registerLocation(ctx, tx, hall.Location, tempHallId); err != nil {
		tx.Rollback()
		return err
	}

	if err := registerOperation(ctx, tx, hall.OperationTime, tempHallId); err != nil {
		tx.Rollback()
		return err
	}

	if err := updateAdminHallXColumnToTrue(ctx, tx, adminId, "hall_registered"); err != nil {
		tx.Rollback()
		return err
	}
	*hallId = tempHallId
	return nil
}


func registerActualHall(ctx context.Context,tx *sql.Tx, actualHall models.ActualHall) (int64, error) {
	var hallId int64
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO hall (hall_name, hall_manager, hall_contact, admin_id) VALUES (?, ?, ?, ?);`)
	if err != nil {
		return hallId, fmt.Errorf("prepare context : %w",err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
		actualHall.Name,
		actualHall.Manager,
		actualHall.Contact,
		actualHall.AdminId)
	if err != nil { 
		return hallId, fmt.Errorf("query execution : %w",err)
	}

	hallId, err = res.LastInsertId()
	if err != nil { 
		return hallId, fmt.Errorf("last inserted : %w",err)
	}
	return hallId, nil
}



func registerLocation(ctx context.Context, tx *sql.Tx, location models.Location,hallId int64) error {
	stmt, err := tx.PrepareContext(ctx,`INSERT INTO hall_location (address, city, state, postal_code, latitude, longitude, hall_id) VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return fmt.Errorf("prepare context : %w",err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		location.Address,
		location.City,
		location.State,
		location.PostalCode,
		location.Latitude,
		location.Longitude,
		hallId)
	if err != nil { 
		return fmt.Errorf("query execution : %w",err)
	}

	return nil
}

func registerOperation(ctx context.Context, tx *sql.Tx, operation models.OperationTime, hallId int64) error {
	stmt, err := tx.PrepareContext(ctx,`INSERT INTO hall_operation_time (open_time, closed_time, hall_id) VALUES (?, ?, ?);`)
	if err != nil {
		return fmt.Errorf("prepare context : %w",err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		operation.OpenTime,
		operation.CloseTime,
		hallId)
	if err != nil { 
		return fmt.Errorf("query execution : %w",err)
	}

	return nil
}

func updateAdminHallXColumnToTrue(ctx context.Context, tx *sql.Tx, adminId int64, column string) error {
	query := fmt.Sprintf(`UPDATE admin SET %s=true WHERE id=?`, column)
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare context : %w",err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, adminId); err != nil {
		return fmt.Errorf("query execution : %w",err)
	}
	return nil
}