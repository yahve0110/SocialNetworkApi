package main

import (
	"log"
	"social/internal/app"
)

var ()

func init() {

}

func main() {
	config, err := api.ReadConfigFromFile("./configs/config.json")
	if err != nil {
		log.Fatal("Error reading configuration:", err)
	}

	// Create a new instance of the API with the configuration
	server := api.New(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
