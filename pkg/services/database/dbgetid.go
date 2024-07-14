package database

import "fmt"

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
