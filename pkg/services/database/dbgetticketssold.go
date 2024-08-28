package database

import (
	"context"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

func GetTicketsSold(ctx context.Context, ID int64) ([]models.TicketSold, error) {
	query := `
		SELECT
			m.title,
			COUNT(t.id)
		FROM movie m
		JOIN movie_show ms ON ms.movie_id = m.id
		JOIN hall h ON h.id = ms.hall_id
		JOIN booking b ON b.movie_show_id = ms.id
		JOIN ticket t ON t.booking_id = b.id
		WHERE h.admin_id = ?
		GROUP BY m.title;
	`

	ticketsSold := make([]models.TicketSold, 0)

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return ticketsSold, fmt.Errorf("failed to prepare query statement: id=%d: %w", ID, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, ID)
	if err != nil {
		return ticketsSold, fmt.Errorf("failed to execute query: id=%d: %w", ID, err)
	}

	for rows.Next() {
		ticketSold := new(models.TicketSold)
		err = rows.Scan(&ticketSold.MovieName, &ticketSold.NoOfTicketsSold)
		if err != nil {
			return ticketsSold, fmt.Errorf("failed to scan the result: id=%d: %w", ID, err)
		}

		ticketsSold = append(ticketsSold, *ticketSold)
	}

	return ticketsSold, nil
}