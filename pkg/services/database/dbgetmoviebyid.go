package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


type MovieDetails struct {
	models.Movie
	models.Cast
}

func GetMovieDetailsByID(ctx context.Context, movieId int64) (MovieDetails, error) {
	var movieDetails MovieDetails

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return movieDetails, fmt.Errorf("failed to begin the transaction : %w",err)
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

	movieMetaData, err := getMovieMetaData(ctx, tx, movieId)
	if err != nil {
		tx.Rollback()
		return movieDetails, fmt.Errorf("failed to fetch movie metadata: %w", err)
	}

	movieCasts, err := getMovieCasts(ctx, tx, movieId)
	if err != nil {
		tx.Rollback()
		return movieDetails, fmt.Errorf("failed to fetch movie casts: %w", err)
	}


	movieDetails.Movie = movieMetaData
	movieDetails.Cast = movieCasts

	return movieDetails, nil
}

func getMovieMetaData(ctx context.Context, tx *sql.Tx, movieId int64) (models.Movie, error) {
	var movieMetaData models.Movie
	const query = `
		SELECT id, title, description, duration, genre, release_date
		FROM movie
		WHERE id = ?;
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return movieMetaData, fmt.Errorf("error preparing the query statement: %w",err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(movieId).Scan(
		&movieMetaData.Id,
		&movieMetaData.Title,
		&movieMetaData.Description,
		&movieMetaData.Duration,
		&movieMetaData.Genre,
		&movieMetaData.ReleaseDate,
	); err != nil {
		return movieMetaData, fmt.Errorf("error in executing the query: %w",err)
	}
	
	return movieMetaData, nil
}

func getMovieCasts(ctx context.Context, tx *sql.Tx, movieId int64) (models.Cast, error) {
	var casts models.Cast

	casts.Actors, err = getMovieCastsByRole(ctx, tx, movieId, "actor", "movie_actor", "actor_id")
	if err != nil {
		return casts, fmt.Errorf("failed to fetch actors: %w", err)
	}

	casts.Actress, err = getMovieCastsByRole(ctx, tx, movieId, "actress", "movie_actress", "actress_id")
	if err != nil {
		return casts, fmt.Errorf("failed to fetch actress: %w", err)
	}

	casts.Directors, err = getMovieCastsByRole(ctx, tx, movieId, "director", "movie_director", "director_id")
	if err != nil {
		return casts, fmt.Errorf("failed to fetch directors: %w", err)
	}

	casts.Producers, err = getMovieCastsByRole(ctx, tx, movieId, "producer", "movie_producer", "producer_id")
	if err != nil {
		return casts, fmt.Errorf("failed to fetch producers: %w", err)
	}

	return casts, nil
}

func getMovieCastsByRole(ctx context.Context, tx *sql.Tx, movieId int64, role, table, column string) ([]models.CastBlueprint, error) {
	var casts []models.CastBlueprint

	var query string

	switch role {
	case "actor", "actress":
		query = `
			SELECT mc.alias, c.id, c.name
			FROM %s as mc
			INNER JOIN %s as c 
			ON mc.%s = c.id
			WHERE mc.movie_id = ?
		`
	default:
		query = `
			SELECT c.id, c.name
			FROM %s as mc
			INNER JOIN %s as c 
			ON mc.%s = c.id
			WHERE mc.movie_id = ?
		`
	}

	processedQuery := fmt.Sprintf(query, table, role, column)

	stmt, err := tx.PrepareContext(ctx, processedQuery)
	if err != nil {
		return casts, fmt.Errorf("error preparing cast query for role %s: %w", role, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, movieId)
	if err != nil {
		return casts, fmt.Errorf("error executing cast query for role %s: %w", role, err)
	}

	for rows.Next() {
		var cast models.CastBlueprint 

		switch role {
		case "actor", "actress":
			if err := rows.Scan(
				&cast.Alias,
				&cast.Id,
				&cast.Name,
			); err != nil {
				return casts, fmt.Errorf("failed to scan the rows: %w", err)
			}
		default:
			if err := rows.Scan( 
				&cast.Id,
				&cast.Name,
			); err != nil {
				return casts, fmt.Errorf("error scanning cast rows for role %s: %w", role, err)
			}
		}

		casts = append(casts, cast)
	}

	return casts, nil
}