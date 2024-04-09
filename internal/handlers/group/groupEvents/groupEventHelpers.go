package groupEventHandlers

import "database/sql"

// IsEventExist checks if an event with the given eventID exists
func IsEventExist(db *sql.DB, eventID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM group_events WHERE event_id = ?)"
	var exists bool
	err := db.QueryRow(query, eventID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsUserJoinedEvent checks if a user with the given userID has joined the event with the given eventID
func IsUserJoinedEvent(db *sql.DB, userID, eventID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM event_going_members WHERE member_id = ? AND event_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, eventID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsUserNotJoinedEvent checks if a user with the given userID has not joined the event with the given eventID
func IsUserNotJoinedEvent(db *sql.DB, userID, eventID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM event_not_going_members WHERE member_id = ? AND event_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, eventID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// AddUserToGoingMembers adds a user with the given userID to the going members of the event with the given eventID
func AddUserToGoingMembers(db *sql.DB, userID, eventID string) error {
	_, err := db.Exec("INSERT INTO event_going_members (event_id, member_id) VALUES (?, ?)", eventID, userID)
	return err
}

// AddUserToNotGoingMembers adds a user with the given userID to the not going members of the event with the given eventID
func AddUserToNotGoingMembers(db *sql.DB, userID, eventID string) error {
	_, err := db.Exec("INSERT INTO event_not_going_members (event_id, member_id) VALUES (?, ?)", eventID, userID)
	return err
}
