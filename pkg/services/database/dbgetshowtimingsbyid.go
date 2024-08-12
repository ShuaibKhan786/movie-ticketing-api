package database

import (
	"context"
	"fmt"
)

type DBTimingID struct {
	Id     int64  `json:"timing_id"`
	Timing string `json:"timing"`
}

type DBShowTimings struct {
	Id      int64        `json:"date_id"`
	Date    string       `json:"show_date"`
	Timings []DBTimingID `json:"timings"`
}

func GetShowTimingsByID(ctx context.Context, hallId, movieId int64) ([]DBShowTimings, error) {
	var showTimings []DBShowTimings

	const query = `
		SELECT msd.id, msd.show_date, mst.id, mst.show_timing
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

	dateMap := make(map[int64]*DBShowTimings)
	for rows.Next() {
		var dateId int64
		var showDate string
		var showTiming DBTimingID

		if err := rows.Scan(&dateId, &showDate, &showTiming.Id, &showTiming.Timing); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := dateMap[dateId]; !exists {
			dateMap[dateId] = &DBShowTimings{
				Id:      dateId,
				Date:    showDate,
				Timings: []DBTimingID{},
			}
		}
		dateMap[dateId].Timings = append(dateMap[dateId].Timings, showTiming)
	}

	for _, showDate := range dateMap {
		showTimings = append(showTimings, *showDate)
	}

	return showTimings, nil
}
