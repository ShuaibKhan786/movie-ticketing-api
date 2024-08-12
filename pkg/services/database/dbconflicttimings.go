package database

import (
	"context"
	"fmt"

)

type DBShowDate struct {
	Date   string   `json:"date"`
	Timing []string `json:"timings"`
}

func GetConflictTimings(ctx context.Context, hallId int64, providedTimings []DBShowDate) ([]DBShowDate, error) {
	var showDates []DBShowDate

	const query = `
		SELECT msd.id, msd.show_date, mst.show_timing
		FROM movie_show_dates msd
		INNER JOIN movie_show ms ON msd.movie_show_id = ms.id
		INNER JOIN movie_show_timings mst ON mst.movie_show_dates_id = msd.id
		WHERE ms.hall_id = ? AND msd.show_date >= ?;
	`

	rows, err := db.QueryContext(ctx, query, hallId, providedTimings[0].Date)
	if err != nil {
		return nil, fmt.Errorf("error in query execution: %w", err)
	}
	defer rows.Close()

	dateMap := make(map[int64]*DBShowDate)
	for rows.Next() {
		var dateId int64
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

	conflicts := checkForConflicts(showDates, providedTimings)

	return conflicts, nil
}

func checkForConflicts(existing []DBShowDate, provided []DBShowDate) []DBShowDate {
	var conflicts []DBShowDate

	existingMap := make(map[string][]string)
	for _, showDate := range existing {
		existingMap[showDate.Date] = showDate.Timing
	}

	for _, pDate := range provided {
		if eTimings, exists := existingMap[pDate.Date]; exists {
			conflictingTimings := make([]string, 0)
			for _, pTiming := range pDate.Timing {
				for _, eTiming := range eTimings {
					if pTiming == eTiming {
						conflictingTimings = append(conflictingTimings, pTiming)
					}
				}
			}
			if len(conflictingTimings) > 0 {
				conflicts = append(conflicts, DBShowDate{
					Date:   pDate.Date,
					Timing: conflictingTimings,
				})
			}
		}
	}

	return conflicts
}
