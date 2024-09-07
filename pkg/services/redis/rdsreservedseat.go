package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	FIFTEEN_MINUTES_EXPIRY = 15 * 60 // 15 minutes in seconds
	RESERVED_SEAT_KEY      = `RESERVED:SEATS:%d`
	USER_RB_SEAT_KEY       = `USER:RB:SEATS:%d`
)

type TimingID int64

type ReservedSeatSchema struct {
	TimingID int64
	Seats    []string
}

// Lua script for atomic ZADD and EXPIRE
// Syntax
//   - ZADD   key score member
//   - EXPIRY key tls
const luaScriptZ = `
redis.call("ZADD", KEYS[1], ARGV[1], ARGV[2])
redis.call("EXPIRE", KEYS[1], ARGV[3])
redis.call("SET", KEYS[2], ARGV[4])
redis.call("EXPIRE", KEYS[2], ARGV[5])
return true
`

// Used sorted Set i.e., ZAdd with Lua script for atomicity
// Used at routes
//   - /admin/hall/show/register
//   - /admin/hall/show/ticket/release/{timing_id}
func (schema *SeatRegsSchema) InitialReservedSeatsRegs(ctx context.Context) error {
	reservedKey := fmt.Sprintf(
		RESERVED_SEAT_KEY,
		schema.TimingID,
	)

	userRBKey := fmt.Sprintf(
		USER_RB_SEAT_KEY,
		schema.TimingID,
	)

	garbageExpiry := time.Now().Unix() + schema.KeyExpiry

	err := rdb.Eval(ctx,
		luaScriptZ,
		[]string{reservedKey, userRBKey},
		garbageExpiry,             //ARGV[1]
		GARBAGE_MEMBER,            //ARGV[2]
		schema.KeyExpiry,          //ARGV[3]
		userRBKey,                 //ARGV[4]
		schema.PreExpiryKey).Err() //ARGV[5] // TTL in seconds
	if err != nil {
		return fmt.Errorf("initial Reserved Seat: %w", err)
	}

	return nil
}
// ReservedSeatsRegs registers a seat reservation
func (schema *ReservedSeatSchema) ReservedSeatsRegs(ctx context.Context, role string) ([]string, error) {
	if role == "user" {
		userRBKey := fmt.Sprintf(
			USER_RB_SEAT_KEY,
			schema.TimingID,
		)

		_, err := rdb.Get(ctx, userRBKey).Result()
		if err == redis.Nil {
			return nil, redis.Nil
		}
		if err != nil {
			return nil, fmt.Errorf("checking pre expiry Reserved Seat: %w", err)
		}
	}

	luaScript := `
		local seats = cjson.decode(ARGV[1])
		local reservedKey = KEYS[1]
		local expiry = tonumber(ARGV[2])
		local currentTime = tonumber(redis.call("TIME")[1])

		local alreadyReservedSeats = {}

		-- Check if seats are available or expired
		for i, seat in ipairs(seats) do
			local score = redis.call("ZSCORE", reservedKey, seat)
			if score and tonumber(score) >= currentTime then
				table.insert(alreadyReservedSeats, seat)
			end
		end

		-- If there are already reserved seats, return them and do not proceed with reservation
		if #alreadyReservedSeats > 0 then
			return cjson.encode(alreadyReservedSeats)
		end

		-- Reserve the seats
		for i, seat in ipairs(seats) do
			redis.call("ZADD", reservedKey, expiry, seat)
		end

		return "SEATS_RESERVED_SUCCESSFULLY"
	`

	reservedKey := fmt.Sprintf(
		RESERVED_SEAT_KEY,
		schema.TimingID,
	)

	fmt.Println(reservedKey)
	fmt.Println(schema.Seats)
	fmt.Println("TIME STAMP ", time.Unix(int64(generateTTL(FIFTEEN_MINUTES_EXPIRY)), 0))

	seatsJSON, err := json.Marshal(schema.Seats)
	if err != nil {
		return nil, fmt.Errorf("encoding seats to JSON, Reserved Seat: %w", err)
	}

	fmt.Println(string(seatsJSON))

	result, err := rdb.Eval(
		ctx,
		luaScript,
		[]string{reservedKey},
		string(seatsJSON),
		generateTTL(FIFTEEN_MINUTES_EXPIRY),
	).Result()
	if err != nil {
		return nil, fmt.Errorf("registration Reserved Seat: %w", err)
	}


	if res, ok := result.(string); ok && res == "SEATS_RESERVED_SUCCESSFULLY" {
		return nil, nil
	}
	
	var alreadyReservedSeats []string
	if err := json.Unmarshal([]byte(result.(string)), &alreadyReservedSeats); err == nil && len(alreadyReservedSeats) > 0 {
		return alreadyReservedSeats, nil
	}


	return nil, fmt.Errorf("unexpected response from Lua script")
}

func (schema *ReservedSeatSchema) IsRSNotExpired(ctx context.Context, role string) (bool, error) {
	if role == "user" {
		userRBKey := fmt.Sprintf(
			USER_RB_SEAT_KEY,
			schema.TimingID,
		)

		_, err := rdb.Get(ctx, userRBKey).Result()
		if err == redis.Nil {
			return false, redis.Nil
		}
		if err != nil {
			return false, fmt.Errorf("checking pre expiry Reserved Seat: %w", err)
		}
	}

	reservedKey := fmt.Sprintf(
		RESERVED_SEAT_KEY,
		schema.TimingID,
	)

	currentTime := time.Now().Unix()

	// Lua script to check if any seat is expired
	luaScript := `
		local seats = cjson.decode(ARGV[1])
		local reservedKey = KEYS[1]
		local currentTime = tonumber(ARGV[2])

		-- Check each seat's expiration
		for i, seat in ipairs(seats) do
			local score = redis.call("ZSCORE", reservedKey, seat)
			if score and tonumber(score) < currentTime then
				return 0
			end
		end

		return 1
	`

	seatsJSON, err := json.Marshal(schema.Seats)
	if err != nil {
		return false, fmt.Errorf("encoding seats to JSON: %w", err)
	}

	// Execute the Lua script
	result, err := rdb.Eval(
		ctx,
		luaScript,
		[]string{reservedKey},
		string(seatsJSON),
		currentTime,
	).Result()

	if err != nil {
		return false, fmt.Errorf("checking Reserved Seat expiration: %w", err)
	}

	if result == int64(0) {
		return false, nil
	}

	return true, nil
} 


// GetAllReservedSeats retrieves all reserved seats that are not expired
func (timingID TimingID) GetAllReservedSeats(ctx context.Context) ([]string, error) {
	reservedKey := fmt.Sprintf(
		RESERVED_SEAT_KEY,
		timingID,
	)

	min := fmt.Sprintf("%f", float64(time.Now().Unix()))
	max := "+inf"

	seats, err := rdb.ZRangeByScore(ctx, reservedKey, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("get all Reserved Seats: %w", err)
	}

	return seats, nil
}

// CleanupReservedSeats removes expired seats from the reservation set
func (timingID TimingID) CleanupReservedSeats(ctx context.Context) error {
	reservedKey := fmt.Sprintf(
		RESERVED_SEAT_KEY,
		timingID,
	)

	min := "-inf"
	max := fmt.Sprintf("%f", float64(time.Now().Unix()))

	_, err := rdb.ZRemRangeByScore(ctx, reservedKey, min, max).Result()
	if err != nil {
		return fmt.Errorf("removed expired seats Reserved Seat: %w", err)
	}

	return nil
}

func (timingID TimingID)CheckReservedSeatExists(ctx context.Context) (bool, error) {
	reservedKey := fmt.Sprintf(
		RESERVED_SEAT_KEY,
		timingID,
	)

	ok, err := rdb.Exists(ctx, reservedKey).Result()
	if err != nil {
		return false, fmt.Errorf("checking booked seats key exists: %w", err)
	}

	return ok == 1, nil
}

// Utility functions

func generateTTL(duration int64) float64 {
	// Generate a TTL value based on the current time plus the given duration
	return float64(time.Now().Unix() + duration)
}
