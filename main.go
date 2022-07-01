package main

import (
	"Go-Postgres/router"
	"log"
	"net/http"
)

func main() {

	r := router.Router()

	log.Fatal(http.ListenAndServe(":8001", r))
}
