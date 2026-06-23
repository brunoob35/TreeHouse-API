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
	PIIKey    []byte
	AppLocation = time.UTC
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

	appLocation, locationErr := time.LoadLocation("America/Sao_Paulo")
	if locationErr != nil {
		log.Println("falha ao carregar timezone America/Sao_Paulo, usando UTC:", locationErr)
		appLocation = time.UTC
	}
	AppLocation = appLocation

	Cfg = mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_ADDR"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Loc:       AppLocation,
		Params: map[string]string{
			"time_zone": "'-03:00'",
		},
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
	PIIKey = []byte(os.Getenv("PII_ENCRYPTION_KEY"))
	if len(PIIKey) == 0 {
		PIIKey = SecretKey
	}
}
