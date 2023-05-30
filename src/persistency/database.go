package persistency

import (
	"database/sql"
	"github.com/brunoob35/TreeHouse-API/src/config"
	_ "github.com/go-sql-driver/mysql" // Drivers MySQL
	"log"
)

// Connect Opens persistency connection and returns it
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.ConncetionString)
	if err != nil {
		log.Println("DEBUG: Erro no sql Open")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
