package database

import "fmt"


const QueryGetId Query = `SELECT id from %s WHERE %s=?`

func (query *Query) DBGetId(whichTab , whichCol, colVal string,id *int64) error {
	processedQuery := fmt.Sprintf(string(*query), whichTab, whichCol)

	stmt, err := db.Prepare(processedQuery)
	if err != nil {
		return dbError(ErrDBPrepareStmt, err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(colVal).Scan(id) ; err != nil {
		return dbError(ErrDBNoRows, err)
	}

	return nil
}