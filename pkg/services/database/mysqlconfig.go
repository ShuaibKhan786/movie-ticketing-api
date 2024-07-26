package database

import (
	"database/sql"
	"sync"
	"time"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

type Query string

var (
	db *sql.DB
	once sync.Once
	err error
)

func InitDB() error {
	once.Do(func() {
		driverName := "mysql"
		dsn := config.Env.DSN

		var errOpen error
		db, errOpen = sql.Open(driverName,dsn)
		if errOpen != nil {
			err = errOpen
			return
		} 

		if errPing := db.Ping(); errPing != nil {
			err = errPing
			return 
		}

		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(10)
		db.SetMaxOpenConns(10)

		err = nil
	})

	return err
}


func CloseDB() {
	if err := db.Ping(); err == nil {
		db.Close()
	}
}