package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/google/uuid"
)

func UpdateBookingSeats(ctx context.Context, userID *int64, timingID int64, details models.BookedRequestPayload, seatsMD CheckoutMD, role string) (models.BookedResponsePayload, error) {
	var payload models.BookedResponsePayload

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return payload, fmt.Errorf("begin transaction: %w", err)
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

	movieShowID, err := getMovieShowID(ctx, tx, timingID)
	if err != nil {
		return payload, err
	}

	// Validate seats in Redis
	bookedSchema := &redisdb.BookedSeatSchema{
		TimingID: timingID,
		Seat:     *details.Seats,
	}
	reservedSchema := &redisdb.ReservedSeatSchema{
		TimingID: timingID,
		Seats:    *details.Seats,
	}

	_, err = bookedSchema.IsSeatAvilableBS(ctx, role)
	if err != nil {
		return payload, fmt.Errorf("failed to check booked seat in Redis: %w", err)
	}

	isNotExpired, err := reservedSchema.IsRSNotExpired(ctx, role)
	if err != nil {
		return payload, fmt.Errorf("failed to check reserved seat in Redis: %w", err)
	}
	if !isNotExpired {
		return payload, fmt.Errorf("one or more seat reservations have expired")
	}

	// update booking table
	bookingQuery := `
			INSERT INTO booking (
				user_id, movie_show_id, movie_show_timings_id, seat_type_id, booking_timing,
				role, payment_status, booking_status, transaction_id, discount_applied, amount, 
				mode_of_payment, created_at, updated_at, cash_amount
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := tx.ExecContext(
		ctx,
		bookingQuery,
		userID, movieShowID, timingID, details.ID, time.Now(),
		role, "completed", true, nil, nil, details.PayableAmount,
		details.PaymentMode, time.Now(), time.Now(), details.CashAmount,
	)
	if err != nil {
		return payload, fmt.Errorf("failed to insert booking table: %w", err)
	}
	bookingID, err := res.LastInsertId()
	if err != nil {
		return payload, fmt.Errorf("failed to get the last insertedID from booking table: %w", err)
	}

	// Generate tickets and update the booking table
	tickets := make([]models.Ticket, len(*details.Seats))
	for i, seat := range *details.Seats {
		ticketNumber := uuid.New().String()
		// Insert into ticket table
		_, err = tx.ExecContext(ctx, `
			INSERT INTO ticket (
				ticket_number, phone_number, booking_id, seat_number, ticket_issue_date
			) VALUES (?, ?, ?, ?, ?)`,
			ticketNumber, details.CustomerPhoneNo, bookingID, seat, time.Now(),
		)
		if err != nil {
			return payload, fmt.Errorf("failed to insert ticket table: %w", err)
		}

		tickets[i] = models.Ticket{
			TicketNumber: &ticketNumber,
			SeatNumber:   &seat,
		}
	}

	//finally update the seats to booked seats
	err = bookedSchema.BookedSeatsRegs(ctx)
	if err != nil {
		return payload, fmt.Errorf("redis booked seats update failed: %v: %w", details.Seats, err)
	}

	// Populate the response payload
	payload.CustomerPhoneNo = details.CustomerPhoneNo
	payload.HallName = &seatsMD.HallName
	payload.MovieName = &seatsMD.MovieName
	payload.ShowDate = &seatsMD.ShowDate
	payload.ShowTime = &seatsMD.ShowTiming
	payload.Tickets = tickets

	return payload, nil
}

func getMovieShowID(ctx context.Context, tx *sql.Tx, timingID int64) (int64, error) {
	query := `
		SELECT 
			msd.movie_show_id
		FROM movie_show_dates msd
		JOIN movie_show_timings mst
		ON msd.id = mst.movie_show_dates_id
		WHERE mst.id = ?
	`
	var movieShowID int64
	err := tx.QueryRowContext(ctx, query, timingID).Scan(&movieShowID)
	if err != nil {
		return movieShowID, fmt.Errorf("failed to retrieve movie_show_id: %w", err)
	}

	return movieShowID, nil
}
