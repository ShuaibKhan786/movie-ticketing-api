package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

// RegisterUserDetails registers both admin and user details based on the role
func RegisterUserDetails(ctx context.Context, userDetails models.UserDetails, role string) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
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

	profileId, err := registerProfile(ctx, tx, userDetails.Profile)
	if err != nil {
		return fmt.Errorf("register profile: %w", err)
	}

	if err := registerAdminOrUser(ctx, tx, userDetails, role, profileId); err != nil {
		return fmt.Errorf("register admin/user: %w", err)
	}

	return nil
}

func registerProfile(ctx context.Context, tx *sql.Tx, profile models.Profile) (int64, error) {
	var profileId int64

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO profile (name, poster_url) VALUES (?, ?);`)
	if err != nil {
		return profileId, fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, 
		profile.Name,
		profile.PosterUrl)
	if err != nil {
		return profileId, fmt.Errorf("query execution: %w", err)
	}

	profileId, err = res.LastInsertId()
	if err != nil {
		return profileId, fmt.Errorf("last inserted: %w", err)
	}
	return profileId, nil
}


func registerAdminOrUser(ctx context.Context, tx *sql.Tx, userDetails models.UserDetails, role string, profileId int64) error {
	query := fmt.Sprintf(`INSERT INTO %s (email, email_verified, profile_id) VALUES (?, ?, ?);`, role)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, 
		userDetails.Email, 
		userDetails.EmailVerified, 
		profileId); 
	err != nil {
		return fmt.Errorf("query execution: %w", err)
	}
	return nil
}
