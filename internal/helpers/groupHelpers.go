package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"social/internal/models"
)

// isUserGroupMember checks if a user is a member of the specified group
func IsUserGroupMember(db *sql.DB, userID, groupID string) (bool, error) {
	// Query the group_members table to check if the user is a member
	query := "SELECT EXISTS(SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, groupID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		return false, err
	}

	return exists, nil
}

// isUserGroupCreator checks if a user is the creator of the specified group
func IsUserGroupCreator(db *sql.DB, userID, groupID string) (bool, error) {
	// Query the groups table to check if the user is the creator
	query := "SELECT EXISTS(SELECT 1 FROM groups WHERE creator_id = ? AND group_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, groupID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking group creator: %v", err)
		return false, err
	}

	return exists, nil
}

// saveInvitationToDatabase saves a group invitation to the database
func SaveInvitationToDatabase(invitationData models.GroupInvitation, db *sql.DB) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO group_invitations (group_id, sender_id, receiver_id, status) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(invitationData.GroupID, invitationData.InviterID, invitationData.ReceiverID, invitationData.Status)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}

// InvitationExists checks if an invitation with the specified parameters exists
func InvitationExists(db *sql.DB, groupID, receiverID string) (bool, error) {

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT COUNT(*) FROM group_invitations WHERE group_id = ? AND receiver_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return false, err
	}
	defer stmt.Close()
	fmt.Printf("Invitation")
	fmt.Printf("group id: %v", groupID)
	fmt.Printf("receiverID : %v", receiverID)

	// Execute the SQL statement
	var count int
	err = stmt.QueryRow(groupID, receiverID).Scan(&count)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return false, err
	}

	return count > 0, nil
}

// GetUserInvitations fetches group invitations for a specific user
func GetUserInvitations(dbConnection *sql.DB, userID string) ([]models.GroupInvitation, error) {
	// Query group invitations for the user
	rows, err := dbConnection.Query("SELECT group_id, sender_id, status FROM group_invitations WHERE receiver_id = ?", userID)
	if err != nil {
		log.Printf("Error querying group invitations from database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and create group invitation objects
	var invitations []models.GroupInvitation
	for rows.Next() {
		var invitation models.GroupInvitation
		if err := rows.Scan(&invitation.GroupID, &invitation.InviterID, &invitation.Status); err != nil {
			log.Printf("Error scanning group invitation rows: %v", err)
			return nil, err
		}
		invitation.ReceiverID = userID
		invitations = append(invitations, invitation)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over group invitation rows: %v", err)
		return nil, err
	}

	return invitations, nil
}

// DeleteInvitation deletes a group invitation from the database
func DeleteInvitation(db *sql.DB, GroupID, ReceiverID string) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM group_invitations WHERE group_id = ?  AND receiver_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(GroupID, ReceiverID)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	return nil
}

// GroupRequestExists checks if a group request exists in the database
func GroupRequestExists(db *sql.DB, groupID, userID string) (bool, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT COUNT(*) FROM group_requests WHERE group_id = ? AND user_id = ?")
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	var count int
	err = stmt.QueryRow(groupID, userID).Scan(&count)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return false, err
	}

	return count > 0, nil
}

// IsUserPostCreator checks if a user is the creator of the specified post
func IsUserPostCreator(db *sql.DB, userID, postID string) (bool, error) {
	// Query the posts table to check if the user is the creator
	query := "SELECT EXISTS(SELECT 1 FROM group_posts WHERE author_id = ? AND post_id = ?)"
	var exists bool
	err := db.QueryRow(query, userID, postID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking post creator: %v", err)
		return false, err
	}

	return exists, nil
}

// RetrieveUsersGoingToEvent retrieves user IDs of users who are going to the specified event
func RetrieveUsersGoingToEvent(db *sql.DB, eventID string) ([]string, error) {
	query := `
    SELECT
        member_id
    FROM
        event_going_members
    WHERE
        event_id = ?
    `

	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersGoing []string

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		usersGoing = append(usersGoing, userID)
	}

	return usersGoing, nil
}

// RetrieveUsersNotGoingToEvent retrieves user IDs of users who are not going to the specified event
func RetrieveUsersNotGoingToEvent(db *sql.DB, eventID string) ([]string, error) {
	query := `
    SELECT
        member_id
    FROM
        event_not_going_members
    WHERE
        event_id = ?
    `

	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersNotGoing []string

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		usersNotGoing = append(usersNotGoing, userID)
	}

	return usersNotGoing, nil
}
