package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

func GetMoviesByStatus(ctx context.Context, status, date string, limit, offset int) ([]models.Movie, error) {
	const query = `
		SELECT 
			m.id,
			m.title,
			m.description,
			m.duration,
			m.genre,
			m.release_date,
			pup.url AS portrait_url,
			pul.url AS landscape_url
		FROM 
			movie m
		INNER JOIN 
			poster_urls pup ON pup.id = m.portrait_poster_url_id
		INNER JOIN 
			poster_urls pul ON pul.id = m.landscape_poster_url_id
		INNER JOIN (
			SELECT 
				ms.movie_id
			FROM 
				movie_show ms
			INNER JOIN 
				movie_show_dates mst ON ms.id = mst.movie_show_id
			WHERE 
				ms.status = ? AND mst.show_date >= ?
			GROUP BY 
				ms.movie_id
			ORDER BY 
				MIN(mst.show_date) ASC
			LIMIT 
				? OFFSET ?
		) filtered_movies ON m.id = filtered_movies.movie_id;
	`

	rows, err := db.QueryContext(ctx, query, status, date, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

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
			&movie.PortraitUrl,
			&movie.LandscapeUrl,
		); err != nil {
			return nil, fmt.Errorf("rows scanning: %w", err)
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return movies, nil
}
