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

// FeedItem - тип для объединения постов и событий
type FeedItem struct {
	Type        string       `json:"type"` // "post" или "event"
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

// GetUserFeedHandler обрабатывает запрос на получение ленты пользователя
func GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем идентификатор сессии пользователя из куки
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	// Получаем глобальное соединение с базой данных из пакета db
	dbConnection := database.DB

	// Получаем идентификатор пользователя на основе текущей сессии
	userID, err := helpers.GetUserIDFromSession(dbConnection, cookie.Value)
	if err != nil {
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	// Получаем все посты и события для данного пользователя
	posts, events, groupPosts, err := GetUserFeed(dbConnection, userID)
	if err != nil {
		http.Error(w, "Ошибка при получении ленты пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразуем посты и события в массив FeedItem
	var feedItems []FeedItem
	for _, post := range posts {
		createdAt, err := time.Parse(time.RFC3339, post.CreatedAt)
		if err != nil {
			http.Error(w, "Ошибка при преобразовании времени: "+err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Ошибка при преобразовании времени: "+err.Error(), http.StatusInternalServerError)
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
			CreatedAt:  groupPost.CreatedAt, // Верное использование CreatedAt
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

	// Сортируем массив по дате создания
	sort.SliceStable(feedItems, func(i, j int) bool {
		return feedItems[i].CreatedAt.After(feedItems[j].CreatedAt)
	})

	// Преобразуем в JSON и отправляем в ответе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedItems)
}

// GetUserFeed возвращает все посты и события для данного пользователя
func GetUserFeed(dbConnection *sql.DB, userID string) ([]models.Post, []models.GroupEvent, []models.GroupPost, error) {
	// Получаем все посты пользователя
	posts, err := getPostsForUser(dbConnection, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Получаем все посты из групп, в которых состоит пользователь
	groupPosts, err := GetUserGroupPosts(dbConnection, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Получаем все события пользователя
	events, err := getEventsForUser(dbConnection, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	return posts, events, groupPosts, nil
}

// getPostsForUser получает все посты пользователей, на которых подписан текущий пользователь
func getPostsForUser(dbConnection *sql.DB, userID string) ([]models.Post, error) {
	// Выполнение запроса к базе данных для получения постов пользователей, на которых подписан текущий пользователь
	// Выполнение запроса к базе данных для получения постов, на которых подписан текущий пользователь
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

	// Создание списка для хранения постов пользователя
	posts := []models.Post{}

	// Итерация по результатам запроса
	for rows.Next() {
		// Создание переменных для хранения данных о посте
		var post models.Post

		// Сканирование значений из строки результата в переменные
		err := rows.Scan(&post.PostID, &post.AuthorID, &post.Content, &post.CreatedAt, &post.LikesCount, &post.Image, &post.Private, &post.AuthorFirstName, &post.AuthorLastName)
		if err != nil {
			return nil, err
		}

		// Добавление поста в список
		posts = append(posts, post)
	}

	// Проверка на наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращение списка постов пользователя
	return posts, nil
}

// getEventsForUser получает все события пользователя из базы данных
func getEventsForUser(dbConnection *sql.DB, userID string) ([]models.GroupEvent, error) {
	// Выполнение запроса к базе данных для получения событий пользователя
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

	// Создание списка для хранения событий пользователя
	events := []models.GroupEvent{}

	// Итерация по результатам запроса
	for rows.Next() {
		// Создание переменных для хранения данных о событии
		var event models.GroupEvent

		// Сканирование значений из строки результата в переменные
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

		// Добавление события в список
		events = append(events, event)
	}

	// Проверка на наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращение списка событий пользователя
	return events, nil
}

func GetUserGroupPosts(dbConnection *sql.DB, userID string) ([]models.GroupPost, error) {
	// Получаем список всех групп, в которых состоит пользователь
	groups, err := getUserGroups(dbConnection, userID)
	if err != nil {
		return nil, err
	}

	// Создаем список для хранения всех постов из групп пользователя
	allGroupPosts := []models.GroupPost{}

	// Для каждой группы пользователя получаем посты из этой группы и добавляем их в общий список
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
	// Выполнение запроса к базе данных для получения всех групп пользователя
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

	// Создание списка для хранения групп пользователя
	groups := []models.Group{}

	// Итерация по результатам запроса
	for rows.Next() {
		// Создание переменных для хранения данных о группе
		var group models.Group

		// Сканирование значений из строки результата в переменные
		err := rows.Scan(&group.GroupID, &group.GroupName, &group.GroupDescription)
		if err != nil {
			return nil, err
		}

		// Добавление группы в список
		groups = append(groups, group)
	}

	// Проверка на наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращение списка групп пользователя
	return groups, nil
}

func getGroupPosts(dbConnection *sql.DB, groupID string) ([]models.GroupPost, error) {
	// Выполнение запроса к базе данных для получения постов из указанной группы
	// Выполнение запроса к базе данных для получения постов из указанной группы
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

	// Создание списка для хранения постов из указанной группы
	posts := []models.GroupPost{}

	// Итерация по результатам запроса
	for rows.Next() {
		// Создание переменных для хранения данных о посте из группы
		var post models.GroupPost

		// Сканирование значений из строки результата в переменные
		err := rows.Scan(&post.PostID, &post.AuthorID, &post.Content, &post.CreatedAt, &post.Image, &post.AuthorFirstName, &post.AuthorLastName, &post.LikesCount)
		if err != nil {
			return nil, err
		}

		// Добавление поста из группы в список
		posts = append(posts, post)
	}

	// Проверка на наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращение списка постов из указанной группы
	return posts, nil
}
