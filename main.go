package main

import (
	"fmt"
	"github.com/brunoob35/TreeHouse-API/src/config"
	"github.com/brunoob35/TreeHouse-API/src/router"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()
	fmt.Println("API running at door ", config.Port)
	r := router.Generate()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
