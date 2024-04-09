package groupEventHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	database "social/internal/db"
	"social/internal/models"
)

// GetGroupEvents retrieves all events for a specific group
func GetGroupEvents(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the group ID from the request
	groupID := r.FormValue("group_id")
	if groupID == "" {
		http.Error(w, "GroupID cannot be empty", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve group events with members
	eventsWithMembers, err := RetrieveGroupEventsWithMembers(dbConnection, groupID)
	if err != nil {
		log.Printf("Error retrieving group events: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with the retrieved events
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(eventsWithMembers)
}

// RetrieveGroupEventsWithMembers retrieves group events with members (going and not going) from the database
func RetrieveGroupEventsWithMembers(db *sql.DB, groupID string) ([]models.GroupEvent, error) {
	query := `
	SELECT
    ge.event_id,
    ge.group_id,
    ge.title,
    ge.description,
    ge.date_time,
    JSON_GROUP_ARRAY(egm.member_id) AS going_members,
    JSON_GROUP_ARRAY(engm.member_id) AS not_going_members
FROM
    group_events ge
LEFT JOIN
    event_going_members egm ON ge.event_id = egm.event_id
LEFT JOIN
    event_not_going_members engm ON ge.event_id = engm.event_id
WHERE
    ge.group_id = ?
GROUP BY
    ge.event_id, ge.group_id, ge.title, ge.description, ge.date_time

	`

	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventsWithMembers []models.GroupEvent

	for rows.Next() {
		var eventWithMembers models.GroupEvent
		var goingMembersJSON, notGoingMembersJSON []byte

		if err := rows.Scan(
			&eventWithMembers.EventID,
			&eventWithMembers.GroupID,
			&eventWithMembers.Title,
			&eventWithMembers.Description,
			&eventWithMembers.DateTime,
			&goingMembersJSON,
			&notGoingMembersJSON,
		); err != nil {
			return nil, err
		}

		// Unmarshal JSON data into slices of strings
		if err := json.Unmarshal(goingMembersJSON, &eventWithMembers.Options.Going); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(notGoingMembersJSON, &eventWithMembers.Options.NotGoing); err != nil {
			return nil, err
		}

		eventsWithMembers = append(eventsWithMembers, eventWithMembers)
	}

	return eventsWithMembers, nil
}
