package database

import (
	"fmt"

	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

const QuerySignup Query = `INSERT INTO %s (email, password ) VALUES (?, ?)`

func (query *Query) DBsignup(model models.UserAdminCredentials, id *int64) error {
	processedQuery := fmt.Sprintf(string(*query), model.Role)

	stmt, err := db.Prepare(processedQuery)
	if err != nil {
		return dbError(ErrDBPrepareStmt, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(model.Email, model.Password)
	if err != nil { 
		return dbError(ErrDBExec, err)
	}

	insertedID, err := res.LastInsertId()
	if err != nil { 
		return dbError(ErrDBNoLastInsertedId, err)
	}

	*id = insertedID

	return nil
}


