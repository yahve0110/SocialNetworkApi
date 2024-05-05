// api/api.go

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	database "social/internal/db"
	groupChat "social/internal/handlers/group/groupMessages"
	messageHandlers "social/internal/handlers/messages"
	"social/internal/handlers/notifications"
	"social/internal/middleware"
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

	// Create a middleware that will be applied to all routes
	allRoutesMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Perform pre-processing or checks here
			fmt.Println("Middleware applied to all routes")

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}

	// Wrap the entire router with the middleware
	routerWithMiddleware := allRoutesMiddleware(router)

	// Apply CORS middleware and pass the router with middleware
	http.Handle("/", middleware.CORSMiddleware(routerWithMiddleware))
	http.HandleFunc("/ws", messageHandlers.HandleConnections)
	http.HandleFunc("/wsGroupChat", groupChat.HandleGroupChatConnections)
	http.HandleFunc("/Wsnotifications", notifications.HandleConnectionsNotif)

	go messageHandlers.HandleMessages()
	go groupChat.HandleGroupMessages()
	go notifications.HandleNotifications()
	// Initialize database
	db, err := database.InitDB("./internal/db/database.db")
	if err != nil {
		log.Fatal("Error initializing database:", err)
		return err
	}
	defer db.Close()

	// Start the HTTP server with the correct router
	err = http.ListenAndServe(api.config.Port, nil) // Pass nil here to use DefaultServeMux
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
