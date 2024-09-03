package database

import (
	"context"
	"fmt"

)


type SearchCast struct {
	Id *int64 `json:"id"`
	Name *string `json:"name"`
	PosterUrl *string `json:"poster_url"`
}

//TODO: either send the whole details or only ID
func SearchCastByName(ctx context.Context, role, name string) ([]SearchCast, error) {
	const query = `
		SELECT 
			c.id,
			c.name,
			pu.url AS poster_url 
		FROM %s c
		LEFT JOIN poster_urls pu
		ON pu.id = c.poster_url_id
		WHERE c.name LIKE ?
	`
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

	var casts []SearchCast

	for rows.Next() {
		var cast SearchCast
		if err := rows.Scan(
				&cast.Id,
				&cast.Name,
				&cast.PosterUrl,
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