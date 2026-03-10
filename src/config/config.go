package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	Port      = 0
	Cfg       = mysql.Config{}
	SecretKey []byte
)

// LoadEnv loads the env variable
func LoadEnv() {
	var err error
	if err = godotenv.Load(".env"); err != nil {
		log.Fatal("entrou aqui, ", err)
	}
	log.Println("Carregou env")

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}

	Cfg = mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_ADDR"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Loc:       time.Local,
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
