// api/api.go

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	database "social/internal/db"
	myrouter "social/internal/router"
)

// API is the base API server description
type API struct {
	config *Config
}

// New creates a new instance of the base API
func New(config *Config) *API {
	return &API{
		config: config,
	}
}

// Start initializes server loggers, router, database, etc.
func (api *API) Start() error {
	fmt.Printf("Server is starting on port %s with logger level %s\n", api.config.Port, api.config.LoggerLevel)
	log.Println("Log message from Start function")
	// Add your server initialization logic here

	// Create a new router and define routes
	router := myrouter.DefineRoutes()

	// Initialize database
	db, err := database.InitDB("./internal/db/database.db")
	if err != nil {
		log.Fatal("Error initializing database:", err)
		return err
	}
	defer db.Close()

	// Start the HTTP server
	err = http.ListenAndServe(api.config.Port, router)
	if err != nil {
		log.Fatal("Error starting server:", err)
		return err
	}

	return nil
}

// ReadConfigFromFile reads the configuration from a JSON file
func ReadConfigFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

