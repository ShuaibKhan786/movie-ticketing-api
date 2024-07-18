package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


func SearchMoviesByTitle(ctx context.Context, title string) ([]models.Movie, error) {
	const query = `SELECT id, title, description, duration, genre, release_date
					FROM movie
					WHERE title LIKE ?`
	
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

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(
				&movie.Id,
				&movie.Title,
				&movie.Description,
				&movie.Duration,
				&movie.Genre,
				&movie.ReleaseDate,
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