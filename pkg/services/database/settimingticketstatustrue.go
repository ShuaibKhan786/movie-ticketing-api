package database

import (
	"context"
	"database/sql"
	"fmt"

	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func SetTimingTicketStatusTrue(ctx context.Context, timingID int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction : %w", err)
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

	err = setStatusTrue(ctx, tx, timingID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = initialSeatRegsRedis(ctx, tx, timingID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func setStatusTrue(ctx context.Context, tx *sql.Tx, timingID int64) error {
	query := `
		UPDATE movie_show_timings 
		SET ticket_status = true
		WHERE id = ?;
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("movie_show_timings: preparing query: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, timingID)
	if err != nil {
		return fmt.Errorf("movie_show_timings: executing query: %v", err)
	}

	return nil
}

func initialSeatRegsRedis(ctx context.Context, tx *sql.Tx, timingID int64) error {
	query := `
		SELECT 
			msd.id,
			msd.show_date,
			mst.show_timing,
			mst.pre_expiry_second,
			mst.post_expiry_second
		FROM movie_show_dates msd
		INNER JOIN movie_show_timings mst
		ON mst.movie_show_dates_id = msd.id
		WHERE mst.id = ?;
	`

	var dateID int64
	var preExp, postExp int
	var date, time string

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("redis: ticket status: preparing query: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, timingID)
	if row.Err() != nil {
		return fmt.Errorf("redis: ticket status: executing row query: %v", err)
	}

	err = row.Scan(
		&dateID,
		&date,
		&time,
		&preExp,  //in seconds
		&postExp) //in seconds
	if err != nil {
		return fmt.Errorf("redis: ticket status: scannig row: %v", err)
	}

	//----redis
	showTimeSecs, err := utils.ConvertToSeconds(date, time)
	if err != nil {
		return fmt.Errorf("timing conversion: hall_show_timings: %w", err)
	}

	preExpiry := showTimeSecs - int64(preExp)
	postExpiry := showTimeSecs + int64(postExp)
	schema := &redisdb.SeatRegsSchema{
		TimingID:     timingID,
		KeyExpiry:    postExpiry,
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
	//----redis

	return nil
}
