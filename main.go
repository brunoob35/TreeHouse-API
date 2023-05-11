package main

import (
	"fmt"
	"github.com/brunoob35/TreeHouse-API/src/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("API running at door :5000")
	r := router.Generate()

	log.Fatal(http.ListenAndServe(":5000", r))
}
