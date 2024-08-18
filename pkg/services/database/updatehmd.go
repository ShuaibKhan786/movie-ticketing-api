package database

import (
	"context"
	"fmt"
	"strings"
)

func UpdateHallMetaData(ctx context.Context, hallID int64, metadata map[string]map[string]interface{}) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction : %w", err)
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

	for _, fields := range metadata {
		var updates []string
		var values []interface{}

		for key, value := range fields {
			if key == "tName" || key == "idName" {
				continue
			}

			updates = append(updates, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}

		updateStmt := strings.Join(updates, ", ")

		query := fmt.Sprintf(
			`UPDATE %s SET %s WHERE %s = ?`,
			fields["tName"], updateStmt, fields["idName"],
		)

		values = append(values, hallID)

		_, err := tx.ExecContext(ctx, query, values...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("hall metadata update: %s: %s: %w", fields["tName"], updateStmt, err)
		}
	}

	return nil
}
