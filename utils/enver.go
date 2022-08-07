package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv() {
	stage, exists := os.LookupEnv("STAGE")
	if !exists {
		stage = "local"
	}

	err := godotenv.Load(".env." + stage)
	if err != nil {
		log.Printf("Env-file .env.%s couldn't be loaded. Is env set up?\nError: %s", stage, err.Error())
	}
}

func GetEnvVal(key string, defaultVal string) string {
	val, _ := os.LookupEnv(key)
	if val == "" {
		val = defaultVal
		os.Setenv(key, val)
	}
	return val
}
