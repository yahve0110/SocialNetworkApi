package userFeedHandler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	database "social/internal/db"
	"social/internal/helpers"
	"social/internal/models"
	"sort"
	"time"
)

type FeedItem struct {
	Type        string       `json:"type"` 
	ID          string       `json:"id"`
	CreatedAt   time.Time    `json:"createdAt"`
	AuthorID    string       `json:"authorId"`
	Content     string       `json:"content,omitempty"`
	LikesCount  int          `json:"likesCount"`
	Image       string       `json:"image,omitempty"`
	Privacy     string       `json:"privacy,omitempty"`
	FirstName   string       `json:"firstName,omitempty"`
	LastName    string       `json:"lastName,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	DateTime    time.Time    `json:"dateTime,omitempty"`
	EventImg    string       `json:"eventImg,omitempty"`
	Options     EventOptions `json:"options,omitempty"` // Include options for events

}

type EventOptions struct {
	Going    []string `json:"going,omitempty"`
	NotGoing []string `json:"notGoing,omitempty"`
}

func GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Error", http.StatusUnauthorized)
		return
	}

	dbConnection := database.DB

	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "Error", http.StatusUnauthorized)
		return
	}

	posts, events, groupPosts, err := GetUserFeed(dbConnection, userID)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var feedItems []FeedItem
	for _, post := range posts {
		createdAt, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		feedItem := FeedItem{
			Type:       "post",
			ID:         post.PostID,
			CreatedAt:  createdAt,
			AuthorID:   post.AuthorID,
			Content:    post.Content,
			LikesCount: post.LikesCount,
			Image:      post.Image,
			Privacy:    post.Private,
			FirstName:  post.AuthorFirstName,
			LastName:   post.AuthorLastName,
		}
		feedItems = append(feedItems, feedItem)
	}
	for _, event := range events {
		createdAt, err := time.Parse(time.RFC3339, event.EventCreatedAt.Format(time.RFC3339))
		if err != nil {
			http.Error(w, "Error "+err.Error(), http.StatusInternalServerError)
			return
		}
		feedItem := FeedItem{
			Type:        "event",
			ID:          event.EventID,
			CreatedAt:   createdAt,
			Title:       event.Title,
			Description: event.Description,
			DateTime:    event.DateTime,
			EventImg:    event.EventImg,
			Options: EventOptions{ // Populate options for events
				Going:    event.Options.Going,
				NotGoing: event.Options.NotGoing,
			},
		}
		feedItems = append(feedItems, feedItem)
	}
	for _, groupPost := range groupPosts {
		feedItem := FeedItem{
			Type:       "groupPost",
			ID:         groupPost.PostID,
			CreatedAt:  groupPost.CreatedAt,
			AuthorID:   groupPost.AuthorID,
			Content:    groupPost.Content,
			LikesCount: groupPost.LikesCount,
			Image:      groupPost.Image,
			Privacy:    "",
			FirstName:  "",
			LastName:   "",
		}
		feedItems = append(feedItems, feedItem)
	}

	sort.SliceStable(feedItems, func(i, j int) bool {
		return feedItems[i].CreatedAt.After(feedItems[j].CreatedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedItems)
}

func GetUserFeed(dbConnection *sql.DB, userID string) ([]models.Post, []models.GroupEvent, []models.GroupPost, error) {
	posts, err := getPostsForUser(dbConnection, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	groupPosts, err := GetUserGroupPosts(dbConnection, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	events, err := getEventsForUser(dbConnection, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	return posts, events, groupPosts, nil
}

func getPostsForUser(dbConnection *sql.DB, userID string) ([]models.Post, error) {

	rows, err := dbConnection.Query(`
	SELECT
		p.post_id,
		p.author_id,
		p.content,
		p.post_created_at,
		COALESCE((SELECT COUNT(*) FROM postLikes pl WHERE pl.post_id = p.post_id), 0) AS likes_count,
		p.image,
		p.privacy,
		u.first_name,
		u.last_name
	FROM
		posts p
	JOIN
		users u ON p.author_id = u.user_id
	WHERE
		p.author_id IN (
			SELECT user_followed FROM followers WHERE user_following = ?
		)
	ORDER BY
		p.post_created_at DESC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.Post{}

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.PostID, &post.AuthorID, &post.Content, &post.CreatedAt, &post.LikesCount, &post.Image, &post.Private, &post.AuthorFirstName, &post.AuthorLastName)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func getEventsForUser(dbConnection *sql.DB, userID string) ([]models.GroupEvent, error) {
	rows, err := dbConnection.Query(`
		SELECT
			e.event_id,
			e.group_id,
			e.title,
			e.event_created_at,
			e.description,
			e.date_time,
			e.event_img
		FROM
			group_events e
		INNER JOIN
			group_members gm ON e.group_id = gm.group_id
		WHERE
			gm.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []models.GroupEvent{}

	for rows.Next() {
		var event models.GroupEvent

		err := rows.Scan(&event.EventID, &event.GroupID, &event.Title, &event.EventCreatedAt, &event.Description, &event.DateTime, &event.EventImg)
		if err != nil {
			return nil, err
		}

		// Query database to get users going to the event
		usersGoing, err := helpers.RetrieveUsersGoingToEvent(dbConnection, event.EventID)
		if err != nil {
			return nil, err
		}
		event.Options.Going = usersGoing

		// Query database to get users not going to the event
		usersNotGoing, err := helpers.RetrieveUsersNotGoingToEvent(dbConnection, event.EventID)
		if err != nil {
			return nil, err
		}
		event.Options.NotGoing = usersNotGoing

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetUserGroupPosts(dbConnection *sql.DB, userID string) ([]models.GroupPost, error) {
	groups, err := getUserGroups(dbConnection, userID)
	if err != nil {
		return nil, err
	}

	allGroupPosts := []models.GroupPost{}

	for _, group := range groups {
		groupPosts, err := getGroupPosts(dbConnection, group.GroupID)
		if err != nil {
			return nil, err
		}
		allGroupPosts = append(allGroupPosts, groupPosts...)
	}

	return allGroupPosts, nil
}

func getUserGroups(dbConnection *sql.DB, userID string) ([]models.Group, error) {
	rows, err := dbConnection.Query(`
		SELECT
			g.group_id,
			g.group_name,
			g.group_description
		FROM
			group_members gm
		INNER JOIN
			groups g ON gm.group_id = g.group_id
		WHERE
			gm.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []models.Group{}

	for rows.Next() {
		var group models.Group

		err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDescription)
		if err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func getGroupPosts(dbConnection *sql.DB, groupID string) ([]models.GroupPost, error) {

	rows, err := dbConnection.Query(`
	SELECT
		p.post_id,
		p.author_id,
		p.content,
		p.post_date,
		p.group_post_img,
		u.first_name,
		u.last_name,
		COALESCE(pl.likes_count, 0) AS likes_count
	FROM
		group_posts p
	INNER JOIN
		users u ON p.author_id = u.user_id
	LEFT JOIN
		(SELECT post_id, COUNT(*) AS likes_count FROM group_post_likes GROUP BY post_id) pl ON p.post_id = pl.post_id
	WHERE
		p.group_id = ?
`, groupID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []models.GroupPost{}

	for rows.Next() {
		var post models.GroupPost

		err := rows.Scan(&post.PostID, &post.AuthorID, &post.Content, &post.CreatedAt, &post.Image, &post.AuthorFirstName, &post.AuthorLastName, &post.LikesCount)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
