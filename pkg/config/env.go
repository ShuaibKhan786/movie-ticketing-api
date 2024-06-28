package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type ENV struct {
	JWTSECRETKEY []byte 
	DSN          string 
}

var (
	Env  ENV
	once sync.Once
)

func LoadConfig() bool {
	if err := godotenv.Load(); err != nil {
		return false
	}
	return loadEnv()
}

func loadEnv() bool {
	var success = true
	once.Do(func() {
		Env.JWTSECRETKEY = []byte(os.Getenv("JWT_SECRET_KEY"))
		Env.DSN = os.Getenv("DSN")

		if len(Env.JWTSECRETKEY) == 0 {
			success = false
		}

		if Env.DSN == "" {
			success = false
		}
		
	})
	return success
}
