package database

import "fmt"

const QueryIsExists Query = `SELECT EXISTS(SELECT 1 FROM %s WHERE %s=?)`

func (query *Query) DBIsExists(whichTab,whichCol,colVal string) (bool, error) {
	processedQuery := fmt.Sprintf(string(*query), whichTab, whichCol)

	stmt, err := db.Prepare(processedQuery)
	if err != nil {
		return false, dbError(ErrDBPrepareStmt, err)
	}
	defer stmt.Close()

	var exists bool
	if err := stmt.QueryRow(colVal).Scan(&exists) ; err != nil {
		return false, dbError(ErrDBNoRows, err)
	}

	return exists, nil
}