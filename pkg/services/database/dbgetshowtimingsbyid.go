package database

import (
	"context"
	"fmt"
)

type DBShowDate struct {
	Date string
	Timing []string
}

func GetShowTimingsByID(ctx context.Context, hallId, movieId int64) ([]DBShowDate, error) {
	var showDates []DBShowDate

	const query = `
		SELECT msd.id, msd.show_date, mst.show_timing
		FROM movie_show_dates msd
		INNER JOIN movie_show ms ON msd.movie_show_id = ms.id
		INNER JOIN movie_show_timings mst ON mst.movie_show_dates_id = msd.id
		WHERE ms.hall_id = ? AND ms.movie_id = ? AND mst.ticket_status = true;
	`

	rows, err := db.QueryContext(ctx, query, hallId, movieId)
	if err != nil {
		return nil, fmt.Errorf("error in query execution: %w", err)
	}
	defer rows.Close()

	dateMap := make(map[int64]*DBShowDate)
	for rows.Next() {
		var dateId int64 //this is just for using as a key in hmap
		var showDate string
		var showTiming string

		if err := rows.Scan(&dateId, &showDate, &showTiming); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := dateMap[dateId]; !exists {
			dateMap[dateId] = &DBShowDate{
				Date:   showDate,
				Timing: []string{},
			}
		}
		dateMap[dateId].Timing = append(dateMap[dateId].Timing, showTiming)
	}

	for _, showDate := range dateMap {
		showDates = append(showDates, *showDate)
	}

	return showDates, nil
}
