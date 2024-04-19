package groupEventHandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	database "social/internal/db"
	"social/internal/helpers"
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
	newEvent.EventCreatedAt = time.Now().UTC()

	if newEvent.EventImg != "" {
		eventImageBase64 := newEvent.EventImg
		if eventImageBase64 != "" {
			cloudinaryURL, err := helpers.ImageToCloud(eventImageBase64, w)
			if err != nil {
				// Handle error
				return
			}
			newEvent.EventImg = cloudinaryURL
		}
	}

	// Initialize options
	newEvent.Options.Going = []string{}
	newEvent.Options.NotGoing = []string{}

	// Insert the new event into the database
	if err := InsertGroupEvent(dbConnection, newEvent); err != nil {
		http.Error(w, "Error inserting event into database", http.StatusInternalServerError)
		return
	}

	// Fetch the inserted event from the database
	createdEvent, err := GetGroupEventByID(dbConnection, newEvent.EventID)
	if err != nil {
		http.Error(w, "Error fetching created event from database", http.StatusInternalServerError)
		return
	}

	// Respond with the created event
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

// GetGroupEventByID retrieves a group event by its ID from the database
func GetGroupEventByID(db *sql.DB, eventID string) (*models.GroupEvent, error) {
	// Query the database to retrieve the event by its ID
	row := db.QueryRow(`
		SELECT event_id, group_id, title, description, date_time, event_created_at
		FROM group_events
		WHERE event_id = ?
	`, eventID)

	var event models.GroupEvent
	err := row.Scan(&event.EventID, &event.GroupID, &event.Title, &event.Description, &event.DateTime, &event.EventCreatedAt)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// InsertGroupEvent inserts a new event into the database
func InsertGroupEvent(db *sql.DB, event models.GroupEvent) error {
	// Insert the new event into the database
	_, err := db.Exec(`
		INSERT INTO group_events (event_id, group_id, title, description, date_time,event_created_at, event_img)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, event.EventID, event.GroupID, event.Title, event.Description, event.DateTime, event.EventCreatedAt, event.EventImg)
	if err != nil {
		log.Printf("Error inserting event into database: %v", err)
		return fmt.Errorf("error inserting event into database: %v", err)
	}

	return nil
}
