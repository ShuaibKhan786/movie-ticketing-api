package database

import (
	"fmt"
)


func IsValueExists(whichTable, whichColumn string, columnValue any) (bool, error) {
	const queryTemplate = `SELECT EXISTS(SELECT 1 FROM %s WHERE %s=?)`
	
	query := fmt.Sprintf(queryTemplate, whichTable, whichColumn)

	stmt, err := db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var exists bool
	if err := stmt.QueryRow(columnValue).Scan(&exists); err != nil {
		return false, fmt.Errorf("query execution failed: %w", err)
	}

	return exists, nil
}
