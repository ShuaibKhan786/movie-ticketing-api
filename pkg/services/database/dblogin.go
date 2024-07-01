package database

import (
	"fmt"

	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

const QueryLogin Query = `SELECT id,password FROM %s WHERE email=?`

func (query *Query) DBLogin(model models.UserAdminCredentials, id *int64, password *string) error {
	processedQuery := fmt.Sprintf(string(*query), model.Role)

	stmt, err := db.Prepare(processedQuery)
	if err != nil {
		return dbError(ErrDBPrepareStmt, err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(model.Email).Scan(id, password); err != nil {
		return dbError(ErrDBNoRows, err)
	}

	return nil
}
