package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func RegisterShow(ctx context.Context, show models.Show, hallId int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
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

	tempMovieId, err := registerMovie(ctx, tx, show.Movie)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println(tempMovieId)

	tempShowId, err := registerActualShow(ctx, tx, tempMovieId, hallId, show.Status)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = registerShowTimings(ctx, tx, show.MovieShowTiming, tempShowId); err != nil {
		tx.Rollback()
		return err
	}

	if show.Movie.Id == nil {
		if err = registerAllTheCast(ctx, tx, show.Cast, tempMovieId); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func registerMovie(ctx context.Context, tx *sql.Tx, movie models.Movie) (int64, error) {
	var movieId int64

	if movie.Id == nil {
		portraitUrlId, err := registerPosterUrl(ctx, tx, movie.PortraitUrl)
		if err != nil {
			return movieId, fmt.Errorf("register movie: %w", err)
		}

		landscapeUrlId, err := registerPosterUrl(ctx, tx, movie.LandscapeUrl)
		if err != nil {
			return movieId, fmt.Errorf("register movie: %w", err)
		}

		const query = `INSERT INTO movie 
						(title, description, duration, genre, release_date, portrait_poster_url_id, landscape_poster_url_id)
						VALUES (?, ?, ?, ?, ?, ?, ?);`

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return movieId, fmt.Errorf("prepare statement: %w", err)
		}
		defer stmt.Close()

		res, err := stmt.ExecContext(ctx,
			movie.Title,
			movie.Description,
			movie.Duration,
			movie.Genre,
			movie.ReleaseDate,
			portraitUrlId,
			landscapeUrlId,
		)
		if err != nil {
			return movieId, fmt.Errorf("execution: %w", err)
		}
		movieId, err = res.LastInsertId()
		if err != nil {
			return movieId, fmt.Errorf("cannot get the id: %w", err)
		}

		return movieId, nil
	}

	movieId = *movie.Id
	return movieId, nil
}

func registerActualShow(ctx context.Context, tx *sql.Tx, movieId, hallId int64, status bool) (int64, error) {
	const query = `INSERT INTO movie_show (movie_id, hall_id, status) VALUES (?, ?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, movieId, hallId, status)
	if err != nil {
		return 0, fmt.Errorf("execution: %w", err)
	}

	showId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cannot get the id: %w", err)
	}

	return showId, nil
}

func registerShowTimings(ctx context.Context, tx *sql.Tx, dates []models.ShowDate, showId int64) error {
	const dateQuery = `
		INSERT INTO movie_show_dates 
		(show_date, movie_show_id) 
		VALUES (?, ?);
	`
	const timingQuery = `
		INSERT INTO movie_show_timings 
		(show_timing, ticket_status, movie_show_dates_id, pre_expiry_second, post_expiry_second) 
		VALUES (?, ?, ?, ?, ?);
	`

	stmtDate, err := tx.PrepareContext(ctx, dateQuery)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmtDate.Close()

	stmtTiming, err := tx.PrepareContext(ctx, timingQuery)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmtTiming.Close()

	for _, date := range dates {
		res, err := stmtDate.ExecContext(ctx, 
			date.Date,
			showId); 
		if err != nil {
			return fmt.Errorf("execution: %w", err)
		}

		dateId, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("error in getting the last inserted id: %w",err)
		}

		for _, timing := range date.Timing {
			res, err := stmtTiming.ExecContext(ctx, 
				timing.Time,
				timing.TicketStatus,
				dateId,
				timing.PreExpiry,
				timing.PostExpiry,
			)
			if err != nil {
				return fmt.Errorf("execution: %w", err)
			}

			timingID , err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("last id: hall_show_timings: %w", err)
			}

			// ----redis
			if timing.TicketStatus {
				showTimeSecs, err := utils.ConvertToSeconds(date.Date, timing.Time)
				if err != nil {
					return fmt.Errorf("timing conversion: hall_show_timings: %w", err)
				}

				preExpiry := showTimeSecs - int64(timing.PreExpiry)
				postExpiry := showTimeSecs + int64(timing.PostExpiry)
				schema := &redisdb.SeatRegsSchema{
					TimingID: timingID,
					KeyExpiry: postExpiry,
					PreExpiryKey: preExpiry,
				}
				//----regs
				err = schema.InitialReservedSeatsRegs(ctx)
				if err != nil {
					return fmt.Errorf("redis: %w", err)
				}
				err = schema.InitialBookedSeatsRegs(ctx)
				if err != nil {
					return fmt.Errorf("redis: %w", err)
				}
				//----regs
			}
			//----redis
		}
	}

	return nil
}

func registerAllTheCast(ctx context.Context, tx *sql.Tx, cast models.Cast, movieId int64) error {
	for _, actor := range cast.Actors {
		if err := registerActor(ctx, tx, actor, movieId); err != nil {
			return err
		}
	}

	for _, actress := range cast.Actress {
		if err := registerActress(ctx, tx, actress, movieId); err != nil {
			return err
		}
	}

	for _, director := range cast.Directors {
		if err := registerDirector(ctx, tx, director, movieId); err != nil {
			return err
		}
	}

	for _, producer := range cast.Producers {
		if err := registerProducer(ctx, tx, producer, movieId); err != nil {
			return err
		}
	}

	return nil
}

func registerActor(ctx context.Context, tx *sql.Tx, actor models.CastBlueprint, movieId int64) error {
	return registerCast(ctx, tx, actor, movieId, "actor", "movie_actor", "actor_id")
}

func registerActress(ctx context.Context, tx *sql.Tx, actress models.CastBlueprint, movieId int64) error {
	return registerCast(ctx, tx, actress, movieId, "actress", "movie_actress", "actress_id")
}

func registerDirector(ctx context.Context, tx *sql.Tx, director models.CastBlueprint, movieId int64) error {
	return registerCast(ctx, tx, director, movieId, "director", "movie_director", "director_id")
}

func registerProducer(ctx context.Context, tx *sql.Tx, producer models.CastBlueprint, movieId int64) error {
	fmt.Println("movieId ",movieId)
	return registerCast(ctx, tx, producer, movieId, "producer", "movie_producer", "producer_id")
}

func registerCast(ctx context.Context, tx *sql.Tx, castMember models.CastBlueprint, movieId int64, role, table, column string) error {
	if castMember.Id == nil {
		castId, err := registerNewCast(ctx, tx, castMember, role)
		if err != nil {
			return err
		}
		castMember.Id = new(int64)
		*castMember.Id = castId
	}
	return updateCastIdInCastList(ctx, tx, movieId, castMember.Id, castMember.Alias, table, column)
}

func updateCastIdInCastList(ctx context.Context, tx *sql.Tx, movieId int64, castId *int64, alias *string, table, column string) error {
	var query string

	switch table {
	case "movie_actor", "movie_actress":
		query = fmt.Sprintf(`INSERT INTO %s (movie_id, %s, alias) VALUES (?, ?, ?);`, table, column)

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("prepare statement: %w", err)
		}
		defer stmt.Close()

		if _, err := stmt.ExecContext(ctx, movieId, castId, alias); err != nil {
			return fmt.Errorf("execution: %w", err)
		}
	default:
		query = fmt.Sprintf(`INSERT INTO %s (movie_id, %s) VALUES (?, ?);`, table, column)

		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return fmt.Errorf("prepare statement: %w", err)
		}
		defer stmt.Close()

		if _, err := stmt.ExecContext(ctx, movieId, castId); err != nil {
			return fmt.Errorf("execution: %w", err)
		}
	}
	return nil
}

func registerNewCast(ctx context.Context, tx *sql.Tx, cast models.CastBlueprint, role string) (int64, error) {
	var castId int64

	var posterUrlId sql.NullInt64
	var err error
	if cast.PosterUrl != nil {
		posterUrlId.Int64, err = registerPosterUrl(ctx, tx, cast.PosterUrl)
		if err != nil {
			return castId, fmt.Errorf("register %s: %w", role, err)
		}
		posterUrlId.Valid = true	
	}

	query := fmt.Sprintf(`INSERT INTO %s (name, poster_url_id) VALUES (?, ?);`, role)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, cast.Name, posterUrlId)
	if err != nil {
		return 0, fmt.Errorf("execution: %w", err)
	} 
	
	castId, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cannot get the id: %w", err)
	}

	return castId, nil
}

func registerPosterUrl(ctx context.Context, tx *sql.Tx, url *string) (int64, error) {
	var posterId int64

	const query = `INSERT INTO poster_urls (url) VALUES (?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return posterId, fmt.Errorf("preparing poster_urls table query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, url)
	if err != nil {
		return posterId, fmt.Errorf("executing poster_urls table query: %w", err)
	}

	posterId, err = res.LastInsertId()
	if err != nil {
		return posterId, fmt.Errorf("last inserted id poster_urls table query: %w", err)
	}

	return posterId, nil
}