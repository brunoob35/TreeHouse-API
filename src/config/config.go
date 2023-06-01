package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	//ConncetionString = ""
	Port             = 0
	Cfg = mysql.Config{}
)

// LoadEnv loads the env variable
func LoadEnv() {
	var err error
	// old .env path "src/config/enviroment/enviroment.env"
	if err = godotenv.Load(".env"); err != nil {

		log.Fatal("entrou aqui, ", err)
	}
	log.Println("Carregou env")

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}


	Cfg = mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDR"),
		DBName: os.Getenv("DB_DATABASE"),
	}


	//ConncetionString = fmt.Sprintf("%s:%s@/(%s:%s)?charset=utf8&parseTime=True&loc=Local",
	//	os.Getenv("DB_USER"),
	//	os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_NAME"),
	//	os.Getenv("DB_PORT"),
	//)

}

//root@127.0.0.1:3306
//jdbc:mysql://127.0.0.1:3306/?user=root