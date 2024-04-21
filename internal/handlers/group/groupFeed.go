package groupHandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
)

// GetGroupFeed retrieves both group events and group posts for a specific group
func GetGroupFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Access the global database connection from the db package
	dbConnection := database.DB

	// Get the group ID from the request
	groupID := r.FormValue("group_id")
	if groupID == "" {
		http.Error(w, "GroupID cannot be empty", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve both group events and group posts
	groupFeed, err := RetrieveGroupFeed(dbConnection, groupID)
	if err != nil {
		log.Printf("Error retrieving group feed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond with the retrieved group feed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groupFeed)
}

// RetrieveGroupFeed retrieves both group events and group posts for a specific group
// RetrieveGroupFeed retrieves both group events and group posts for a specific group
func RetrieveGroupFeed(db *sql.DB, groupID string) ([]interface{}, error) {
	// Retrieve group events
	events, err := RetrieveGroupEvents(db, groupID)
	if err != nil {
		return nil, err
	}

	// Retrieve group posts
	posts, err := RetrieveGroupPosts(db, groupID)
	if err != nil {
		return nil, err
	}

	// Combine events and posts into a single slice
	var groupFeed []interface{}

	// Append events to groupFeed
	for _, event := range events {
		groupFeed = append(groupFeed, event)
	}

	// Append posts to groupFeed
	for _, post := range posts {
		groupFeed = append(groupFeed, post)
	}

	// Sort groupFeed by creation date
	sort.Slice(groupFeed, func(i, j int) bool {
		switch feedItem := groupFeed[i].(type) {
		case models.GroupEvent:
			return feedItem.EventCreatedAt.After(getCreationDate(groupFeed[j]))
		case models.GroupPost:
			return feedItem.CreatedAt.After(getCreationDate(groupFeed[j]))
		default:
			return false
		}
	})

	return groupFeed, nil
}
func getCreationDate(item interface{}) time.Time {
	switch item := item.(type) {
	case models.GroupEvent:
		return item.EventCreatedAt
	case models.GroupPost:
		return item.CreatedAt
	default:
		// Return a default time indicating an error condition
		return time.Time{}
	}
}

// RetrieveGroupEvents retrieves group events for a specific group
func RetrieveGroupEvents(db *sql.DB, groupID string) ([]models.GroupEvent, error) {
	query := `
	SELECT
		event_id,
		group_id,
		title,
		description,
		date_time,
		event_created_at,
		event_img
	FROM
		group_events
	WHERE
		group_id = ?
	`

	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.GroupEvent

	for rows.Next() {
		var event models.GroupEvent
		if err := rows.Scan(
			&event.EventID,
			&event.GroupID,
			&event.Title,
			&event.Description,
			&event.DateTime,
			&event.EventCreatedAt,
			&event.EventImg,
		); err != nil {
			return nil, err
		}

		// Query database to get users going to the event
		usersGoing, err := helpers.RetrieveUsersGoingToEvent(db, event.EventID)
		if err != nil {
			return nil, err
		}
		event.Options.Going = usersGoing

		// Query database to get users not going to the event
		usersNotGoing, err := helpers.RetrieveUsersNotGoingToEvent(db, event.EventID)
		if err != nil {
			return nil, err
		}
		event.Options.NotGoing = usersNotGoing

		events = append(events, event)
	}

	return events, nil
}

// RetrieveGroupPosts retrieves group posts for a specific group with author's first name, last name, and group post image
func RetrieveGroupPosts(db *sql.DB, groupID string) ([]models.GroupPost, error) {
	query := `
	SELECT
		gp.post_id,
		gp.group_id,
		gp.author_id,
		gp.content,
		gp.post_date,
		u.first_name,
		u.last_name,
		gp.group_post_img
	FROM
		group_posts gp
	INNER JOIN
		users u ON gp.author_id = u.user_id
	WHERE
		gp.group_id = ?
	`

	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.GroupPost

	for rows.Next() {
		var post models.GroupPost
		if err := rows.Scan(
			&post.PostID,
			&post.GroupID,
			&post.AuthorID,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorFirstName,
			&post.AuthorLastName,
			&post.Image,
		); err != nil {
			return nil, err
		}

		// Fetch likes count for the post
		likesCount, err := getLikeCountForPost(db, post.PostID)
		if err != nil {
			return nil, err
		}
		post.LikesCount = likesCount

		posts = append(posts, post)
	}

	return posts, nil
}

func getLikeCountForPost(dbConnection *sql.DB, postID string) (int, error) {
	// Query the group_post_likes table to get the count of likes for the post
	query := "SELECT COUNT(*) FROM group_post_likes WHERE post_id = ?"
	var likesCount int
	err := dbConnection.QueryRow(query, postID).Scan(&likesCount)
	if err != nil {
		log.Printf("Error fetching like count for post %s: %v", postID, err)
		return 0, err
	}

	return likesCount, nil
}
