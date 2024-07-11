package database

import (
	"fmt"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

type UserDetails struct {
	models.Profile
	Email string `json:"email"`
}

func GetUserDetails(whichTable string, id int64) (UserDetails, error) {
	const query = `SELECT p.name, a.email, p.poster_url FROM %s as a INNER JOIN profile as p ON a.profile_id=p.id WHERE a.id = ?;`
	processedQuery := fmt.Sprintf(query, whichTable)

	stmt, err := db.Prepare(processedQuery)
	if err != nil {
		return UserDetails{}, fmt.Errorf("prepare statment : %w",err)
	}
	defer stmt.Close()

	var userDetails UserDetails
	if err := stmt.QueryRow(id).Scan(&userDetails.Name,
		&userDetails.Email,
		&userDetails.PosterUrl);
	err != nil {
		return UserDetails{}, fmt.Errorf("query row : %w",err)
	}

	return userDetails, nil
}