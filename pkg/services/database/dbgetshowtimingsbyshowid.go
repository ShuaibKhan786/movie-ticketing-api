package database

import (
	"context"
	"fmt"
)

type DBTimingStatus struct {
	DBTimingID
	PreExpiry  int `json:"pre_expiry_secs"`
	PostExpiry int `json:"post_expiry_secs"`
	TicketStatus bool   `json:"ticket_status"`
}

type DBShowTimingsAdmin struct {
	Date   string   `json:"show_date"`
	Timings []DBTimingStatus `json:"timings"`
}


func GetShowTimingsByShowID(ctx context.Context, showID int64) ([]DBShowTimingsAdmin, error) {
	var showTimings []DBShowTimingsAdmin

	const query = `
		SELECT 
			msd.id,
			msd.show_date,
			mst.id, 
			mst.show_timing,
			mst.ticket_status,
			mst.pre_expiry_second,
			mst.post_expiry_second
		FROM movie_show_dates msd
		INNER JOIN movie_show ms ON msd.movie_show_id = ms.id
		INNER JOIN movie_show_timings mst ON mst.movie_show_dates_id = msd.id
		WHERE ms.id = ?;
	`

	rows, err := db.QueryContext(ctx, query, showID)
	if err != nil {
		return nil, fmt.Errorf("error in query execution: %w", err)
	}
	defer rows.Close()

	dateMap := make(map[int64]*DBShowTimingsAdmin)
	for rows.Next() {
		var dateId int64
		var showDate string
		var showTiming DBTimingStatus

		if err := rows.Scan(
			&dateId,
			&showDate, 
			&showTiming.Id, 
			&showTiming.Timing, 
			&showTiming.TicketStatus,
			&showTiming.PreExpiry,
			&showTiming.PostExpiry); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := dateMap[dateId]; !exists {
			dateMap[dateId] = &DBShowTimingsAdmin{
				Date:    showDate,
				Timings: []DBTimingStatus{},
			}
		}
		dateMap[dateId].Timings = append(dateMap[dateId].Timings, showTiming)
	}

	for _, showDate := range dateMap {
		showTimings = append(showTimings, *showDate)
	}

	return showTimings, nil
}
