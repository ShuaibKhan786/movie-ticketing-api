package database

import (
	"context"
	"fmt" 
)

type AdminMinMovieMD struct {
	Id *int64 `json:"id"`
	Title *string `json:"title"`
	PortraitUrl *string `json:"portrait_url"`
}

func AdminSearchMoviesByTitle(ctx context.Context, hallID int64, title string) ([]AdminMinMovieMD, error) {
	const query = `
		SELECT 
			m.id,
			m.title,
			pup.url AS protrait_url
		FROM movie m
		INNER JOIN
			poster_urls pup ON pup.id = m.portrait_poster_url_id
		INNER JOIN 
			movie_show ms ON ms.movie_id = m.id 
		WHERE ms.status=true AND ms.hall_id=? AND m.title LIKE ?`
	
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	formattedTitle := fmt.Sprintf("%%%s%%",title)

	rows, err := stmt.QueryContext(ctx, hallID, formattedTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var movies []AdminMinMovieMD
	for rows.Next() {
		var movie AdminMinMovieMD
		if err := rows.Scan(
				&movie.Id,
				&movie.Title,
				&movie.PortraitUrl,
			); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return movies, nil
}