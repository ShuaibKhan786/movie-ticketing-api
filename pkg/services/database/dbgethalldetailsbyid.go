package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)


type HallDetails struct {
	HallID string `json:"hall_id"`
	HallName string `json:"hall_name"`
	models.Location
}

func GetHallDetailsByID(ctx context.Context, movieId int64) ([]HallDetails, error) {
	var hallDetails []HallDetails

	const query = `
		SELECT 
			h.id,
			h.name,
			hl.address,
			hl.city,
			hl.state,
			hl.postal_code,
			hl.latitude, 
			hl.longitude 
		FROM hall h
		INNER JOIN hall_location hl
			ON h.id = hl.hall_id
		INNER JOIN movie_show ms
			ON h.id = ms.hall_id
		WHERE ms.movie_id = ? ;
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return hallDetails,  fmt.Errorf("error preparing the query statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, movieId)
	if err != nil {
		return hallDetails, fmt.Errorf("error executing the query: %w", err)
	}

	for rows.Next() {
		var halldetial HallDetails

		if err := rows.Scan(
			&halldetial.HallID,
			&halldetial.HallName,
			&halldetial.Address,
			&halldetial.City,
			&halldetial.State,
			&halldetial.PostalCode,
			&halldetial.Latitude,
			&halldetial.Longitude,
		); err != nil {
			return hallDetails, fmt.Errorf("failed to scan the row: %w", err)
		}

		hallDetails = append(hallDetails, halldetial)
	}

	return hallDetails, nil
}