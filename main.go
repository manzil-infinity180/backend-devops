package main

import (
	"fmt"
	"github.com/manzil-infinity180/backend-devops/pkg/api"
	"github.com/manzil-infinity180/backend-devops/pkg/db"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("server has started")
	pgdb, err := db.StartDb()
	if err != nil {
		log.Printf("error starting the database %v", err)
	}
	router := api.StartAPI(pgdb)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Printf("error from router %v\n", err)
	}
}
