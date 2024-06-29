package database

import (
	"errors"
	"fmt"
)

var (
	ErrDBPrepareStmt 		= errors.New("error in preparing query statement")
	ErrDBExec        		= errors.New("error in executing the query")
	ErrDBNoRows      		= errors.New("error no records")
	ErrDBNoLastInsertedId	= errors.New("error no last inserted id")
)

func dbError(dbErr , originalErr error) error {
	return fmt.Errorf("%w : %w",dbErr, originalErr)
}