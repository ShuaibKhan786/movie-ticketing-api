package database

import (
	"fmt"
	"errors"
)

func GetId(whichTable, whichColumn string, columnValue any) (int64, error) {
	const queryTemplate = `SELECT id from %s WHERE %s=?`

	query := fmt.Sprintf(queryTemplate, whichTable, whichColumn)

	var id int64
	stmt, err := db.Prepare(query)
	if err != nil {
		return id, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(columnValue).Scan(&id); err != nil {
		return id, fmt.Errorf("query execution failed: %w", err)
	}

	return id, nil
}


func IsSeatLayoutRegistered(adminId int64) error {
	query := `
		SELECT hall_seat_layout_registered FROM admin WHERE id = ?;
	`

	var status bool

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare context admin query: %w", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(adminId).Scan(&status); err != nil {
		return fmt.Errorf("query execution admin query: %w", err)
	}

	if !status {
		return errors.New("seat layout not registered")
	}

	return nil
}