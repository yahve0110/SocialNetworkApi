package groupEventHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	database "social/internal/db"
	"social/internal/models"

	"github.com/google/uuid"
)

// CreateGroupEventHandler handles the creation of a group event
func CreateGroupEventHandler(w http.ResponseWriter, r *http.Request) {
	var newEvent models.GroupEvent

	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Access the global database connection from the db package
	dbConnection := database.DB

	// Generate a unique event ID
	newEvent.EventID = uuid.New().String()

	// Set current date and time
	newEvent.DateTime = time.Now().UTC()

	// Initialize options
	newEvent.Options.Going = []string{}
	newEvent.Options.NotGoing = []string{}

	// Insert the new event into the database
	if err := InsertGroupEvent(dbConnection, newEvent); err != nil {
		http.Error(w, "Error inserting event into database", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Event created successfully"))
}

// InsertGroupEvent inserts a new event into the database
func InsertGroupEvent(db *sql.DB, event models.GroupEvent) error {
	// Insert the new event into the database
	_, err := db.Exec(`
		INSERT INTO group_events (event_id, group_id, title, description, date_time)
		VALUES (?, ?, ?, ?, ?)
	`, event.EventID, event.GroupID, event.Title, event.Description, event.DateTime)
	if err != nil {
		log.Printf("Error inserting event into database: %v", err)
		return fmt.Errorf("error inserting event into database: %v", err)
	}

	return nil
}
