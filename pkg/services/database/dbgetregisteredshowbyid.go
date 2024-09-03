package database

import (
	"context"
	"fmt"
)

type ShowDetails struct {
	ShowID     *int64  `json:"show_id"`
	MovieID    *int64  `json:"movie_id"`
	Status     *bool   `json:"status"`
	MovieTitle *string `json:"movie_title"`
	PosterUrl  *string `json:"movie_poster_url"`
}

func GetRegisteredShowsByID(ctx context.Context, hallId int64) ([]ShowDetails, error) {
	const query = `
		SELECT 
			ms.id, 
			ms.movie_id, 
			ms.status, 
			m.title,
			pup.url 
		FROM movie_show ms
		INNER JOIN movie m ON m.id = ms.movie_id
		INNER JOIN poster_urls pup ON pup.id = m.portrait_poster_url_id
		WHERE hall_id = ?;
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
			&show.ShowID,
			&show.MovieID,
			&show.Status,
			&show.MovieTitle,
			&show.PosterUrl); err != nil {
			return shows, fmt.Errorf("error in scanning the row by hall id %d: %w", hallId, err)
		}

		shows = append(shows, show)
	}

	return shows, nil
}
