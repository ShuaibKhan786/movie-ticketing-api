package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


//TODO: either send the whole details or only ID
func SearchCastByName(ctx context.Context, role, name string) ([]models.CastBlueprint, error) {
	const query = `SELECT id, name FROM %s WHERE name LIKE ?`
	processedQuery := fmt.Sprintf(query, role)

	stmt, err := db.PrepareContext(ctx, processedQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	formattedName := fmt.Sprintf("%%%s%%",name)

	rows, err := stmt.QueryContext(ctx, formattedName)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var casts []models.CastBlueprint

	for rows.Next() {
		var cast models.CastBlueprint
		if err := rows.Scan(
				&cast.Id,
				&cast.Name,
			); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
		casts = append(casts, cast)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return casts, nil
}