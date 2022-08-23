package main

import (
	"fmt"
	"github.com/brunoob35/tree_house/src/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Rodando API na porta :5000")
	r := router.Generate()

	log.Fatal(http.ListenAndServe(":5000", r))
}
