package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	ConncetionString = ""
	Port = 0
)

//LoadEnv loads the env variable
func LoadEnv() {
	var err error
	if err = godotenv.Load("src/config/enviroment/enviroment.env"); err != nil {

		log.Fatal("entrou aqui, ", err)
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}


	ConncetionString = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)


}
