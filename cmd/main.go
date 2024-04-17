package main

import (
	"github.com/vladimirbannikov/BIP-backend/internal/app"
	"log"
)

func main() {
	log.Print("app started")
	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("app finished")
}
