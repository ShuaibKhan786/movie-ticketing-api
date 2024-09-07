package redisdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type SeatRegsSchema struct {
	TimingID     int64
	KeyExpiry    int64 //in seconds
	PreExpiryKey int64 //in seconds
}

type BookedSeatSchema struct {
	TimingID int64
	Seat     []string
}

const BOOKED_SEAT_KEY = `BOOKED:SEATS:%d`
const GARBAGE_MEMBER = "@"

// Lua script for atomic SADD and EXPIRE
// Syntax
//   - SADD   key member
//   - EXPIRY key tls
const luaScriptS = `
redis.call("SADD", KEYS[1], ARGV[1])
redis.call("EXPIRE", KEYS[1], ARGV[2])
return true
`

// Used Set i,e SAdd
// Used at routes
//   - /admin/hall/show/register
//   - /admin/hall/show/ticket/release/{timing_id}
func (schema *SeatRegsSchema) InitialBookedSeatsRegs(ctx context.Context) error {
	bookedKey := fmt.Sprintf(
		BOOKED_SEAT_KEY,
		schema.TimingID,
	)

	err := rdb.Eval(ctx, luaScriptS,
		[]string{
			bookedKey,
		},
		GARBAGE_MEMBER,         //ARGV[1]
		schema.KeyExpiry).Err() //ARGV[2]
	if err != nil {
		return fmt.Errorf("initial booked seat: %w", err)
	}

	return nil
}

// if seats are not booked it will return nil
func (schema *BookedSeatSchema) IsSeatAvilableBS(ctx context.Context, role string) ([]string, error) {
	var bookedSeats []string

	if role == "user" {
		userRBKey := fmt.Sprintf(
			USER_RB_SEAT_KEY,
			schema.TimingID,
		)

		_, err := rdb.Get(ctx, userRBKey).Result()
		if err == redis.Nil {
			return bookedSeats, redis.Nil
		}
		if err != nil {
			return bookedSeats, fmt.Errorf("checking pre expiry Booked Seat: %w", err)
		}
	}

	bookedKey := fmt.Sprintf(
		BOOKED_SEAT_KEY,
		schema.TimingID,
	)

	status, err := rdb.SMIsMember(ctx, bookedKey, schema.Seat).Result()
	if err != nil {
		return bookedSeats, fmt.Errorf("checking Reserved Seat: %w", err)
	}

	for index, s := range status {
		if s {
			bookedSeats = append(bookedSeats, schema.Seat[index])
		}
	}

	return bookedSeats, nil
}

// Call this when
//   - Permanent update of a seat
//   - On payment success
//   - no revokation required
func (schema *BookedSeatSchema) BookedSeatsRegs(ctx context.Context) error {
	bookedKey := fmt.Sprintf(
		BOOKED_SEAT_KEY,
		schema.TimingID,
	)

	err := rdb.SAdd(ctx, bookedKey, schema.Seat).Err()
	if err != nil {
		return fmt.Errorf("permanent booked seat: SAdd %w", err)
	}

	return nil
}

// retrieves all booked seats that are permanent
func (timingID TimingID) GetAllBookedSeats(ctx context.Context) ([]string, error) {
	bookedKey := fmt.Sprintf(
		BOOKED_SEAT_KEY,
		timingID,
	)

	seats, err := rdb.SMembers(ctx, bookedKey).Result()
	if err != nil {
		return nil, fmt.Errorf("get all Booked Seats: %w", err)
	}

	return seats, nil
}

func (timingID TimingID) CheckBookedSeatExists(ctx context.Context) (bool, error) {
	bookedKey := fmt.Sprintf(
		BOOKED_SEAT_KEY,
		timingID,
	)

	ok, err := rdb.Exists(ctx, bookedKey).Result()
	if err != nil {
		return false, fmt.Errorf("checking booked seats key exists: %w", err)
	}

	return ok == 1, nil
}