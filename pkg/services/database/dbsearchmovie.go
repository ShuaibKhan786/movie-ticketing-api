package database

import (
	"context"
	"fmt" 
)


type MinMovieMD struct {
	Id *int64 `json:"id"`
	Status *bool `json:"status"`
	Title *string `json:"title"`
	PortraitUrl *string `json:"portrait_url"`
}

func SearchMoviesByTitle(ctx context.Context, title string) ([]MinMovieMD, error) {
	const query = `
		SELECT 
			m.id,
			ms.status,
			m.title,
			pup.url AS protrait_url
		FROM movie m
		INNER JOIN
			poster_urls pup ON pup.id = m.portrait_poster_url_id
		INNER JOIN 
			movie_show ms ON ms.movie_id = m.id 
		WHERE m.title LIKE ?`
	
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	formattedTitle := fmt.Sprintf("%%%s%%",title)

	rows, err := stmt.QueryContext(ctx, formattedTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var movies []MinMovieMD
	for rows.Next() {
		var movie MinMovieMD
		if err := rows.Scan(
				&movie.Id,
				&movie.Status,
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