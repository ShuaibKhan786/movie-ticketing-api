package database

import (
	"context"
	"fmt"
)

type ShowDetails struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	ReleaseDate string `json:"release_date"`
}

func GetRegisteredShowsByID(ctx context.Context, hallId int64) ([]ShowDetails, error) {
	const query = `
		SELECT m.id, m.title, m.release_date
		FROM movie m
		INNER JOIN movie_show ms
		ON m.id = ms.movie_id
		WHERE ms.hall_id = ?;
	`

	var shows []ShowDetails

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return shows, fmt.Errorf("error in preapring query statement by hall id %d: %w", hallId, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, hallId)
	if err != nil {
		return shows, fmt.Errorf("error in executing query by hall id %d: %w", hallId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var show ShowDetails

		if err := rows.Scan(
			&show.Id,
			&show.Title,
			&show.ReleaseDate); err != nil {
				return shows, fmt.Errorf("error in scanning the row by hall id %d: %w", hallId, err)
			}
		
		shows = append(shows, show)
	}

	return shows, nil
}