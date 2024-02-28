package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// MigrateDB applies all pending migrations
func MigrateDB(dbPath, migrationsPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %v", err)
	}

	return nil
}

// InitDB opens a connection to the database, creates necessary tables, and applies migrations
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		db.Close()
		log.Printf("Error connecting to database: %v", err)
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Connected to the database")

	// Create all necessary tables
	if err := CreateAllTables(db); err != nil {
		log.Printf("Error initializing tables: %v", err)
		db.Close()
		return nil, fmt.Errorf("error initializing tables: %v", err)
	}

	// Apply migrations
	// migrationsPath := "./internal/db/migrations"
	// if err := MigrateDB(dbPath, migrationsPath); err != nil {

	// 	log.Printf("Error applying migrations: %v", err)
	// 	db.Close()
	// 	return nil, fmt.Errorf("error applying migrations: %v", err)
	// }
	DB = db
	return db, nil
}

// CreateAllTables creates all necessary tables
func CreateAllTables(db *sql.DB) error {

	createTables := `
	CREATE TABLE IF NOT EXISTS users (
		user_id TEXT PRIMARY KEY UNIQUE,
		username TEXT UNIQUE NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		gender TEXT NOT NULL,
		birth_date DATE NOT NULL,
		profile_picture TEXT ,
		role TEXT DEFAULT 'user',
		about TEXT,
		privacy TEXT
	);


	CREATE TABLE IF NOT EXISTS posts (
		post_id TEXT PRIMARY KEY NOT NULL UNIQUE,
		author_id TEXT NOT NULL ,
		content TEXT NOT NULL,
		post_created_at TIMESTAMP NOT NULL,
		likes_count INTEGER NOT NULL,
		privacy TEXT ,
		image TEXT NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS postLikes (
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL UNIQUE,
		FOREIGN KEY (post_id) REFERENCES posts(post_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);



	CREATE TABLE IF NOT EXISTS comments (
		comment_id TEXT UNIQUE NOT NULL,
		content TEXT NOT NULL,
		comment_created_at TIMESTAMP NOT NULL,
		author_id TEXT NOT NULL,
		post_id TEXT NOT NULL,
		author_nickname TEXT NOT NULL,
		image TEXT NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(user_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id)
	);

	CREATE TABLE IF NOT EXISTS CommentLikes (
		comment_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS profile (
		profile_id INTEGER NOT NULL PRIMARY KEY,
		follower_id INTEGER NOT NULL UNIQUE,
		profile_status TEXT NOT NULL,
		FOREIGN KEY (follower_id) REFERENCES users(user_id),
		FOREIGN KEY (profile_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS Followers (
		user_followed_status TEXT NOT NULL,
		user_followed TEXT NOT NULL,
		user_following TEXT NOT NULL,
		FOREIGN KEY (user_followed) REFERENCES users(user_id),
		FOREIGN KEY (user_following) REFERENCES users(user_id)
	);


	CREATE TABLE IF NOT EXISTS friendship (
		friendship_id INTEGER NOT NULL PRIMARY KEY,
		user1_id INTEGER NOT NULL,
		user2_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		action_user_id INTEGER NOT NULL,
		FOREIGN KEY (user1_id) REFERENCES users(user_id),
		FOREIGN KEY (user2_id) REFERENCES users(user_id)
	);
	CREATE TABLE IF NOT EXISTS groups (
		group_id TEXT NOT NULL PRIMARY KEY,
		group_name TEXT NOT NULL,
		group_description TEXT NOT NULL,
		group_image TEXT NOT NULL,
		creation_date TIMESTAMP NOT NULL,
		creator_id TEXT NOT NULL,
		FOREIGN KEY (creator_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_members (
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_invitations (
		group_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		receiver_id TEXT NOT NULL,
		status TEXT NOT NULL, -- 'Pending', 'Accepted', 'Declined'
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (sender_id) REFERENCES users(user_id),
		FOREIGN KEY (receiver_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_requests (
		request_id TEXT NOT NULL PRIMARY KEY,
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		status TEXT NOT NULL, -- 'Pending', 'Accepted', 'Declined'
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_posts (
		post_id TEXT NOT NULL PRIMARY KEY,
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		post_date TIMESTAMP NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_comments (
		comment_id TEXT NOT NULL PRIMARY KEY,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		comment_date TIMESTAMP NOT NULL,
		FOREIGN KEY (post_id) REFERENCES group_posts(post_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_events (
		event_id TEXT NOT NULL PRIMARY KEY,
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		day_time DATETIME NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(group_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS event_participants (
		event_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		response TEXT NOT NULL, -- 'Going', 'Not Going'
		PRIMARY KEY (event_id, user_id),
		FOREIGN KEY (event_id) REFERENCES group_events(event_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS chat (
		message_id INTEGER NOT NULL PRIMARY KEY,
		user1_id INTEGER NOT NULL,
		user2_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		chat_created_at TEXT NOT NULL,
		FOREIGN KEY (message_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_chat (
		chat_id INTEGER NOT NULL PRIMARY KEY,
		chat_name TEXT NOT NULL,
		creator_id INTEGER NOT NULL,
		creation_date TEXT NOT NULL,
		FOREIGN KEY (creator_id) REFERENCES users(user_id)
	);

	CREATE TABLE IF NOT EXISTS group_chat_members (
		member_id INTEGER NOT NULL PRIMARY KEY,
		chat_id INTEGER NOT NULL,
		FOREIGN KEY (member_id) REFERENCES users(user_id),
		FOREIGN KEY (chat_id) REFERENCES group_chat(chat_id)
	);

	CREATE TABLE IF NOT EXISTS group_chat_messages (
		message_id INTEGER NOT NULL PRIMARY KEY,
		content TEXT NOT NULL,
		sender_id INTEGER NOT NULL,
		chat_id INTEGER NOT NULL,
		group_chat_message_created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (sender_id) REFERENCES users(user_id),
		FOREIGN KEY (chat_id) REFERENCES group_chat(chat_id)

	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		session_id TEXT NOT NULL,
		expiration_time DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
		UNIQUE(session_id)
	);


    `

	_, err := db.Exec(createTables)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return fmt.Errorf("error creating users table: %v", err)
	}

	return nil
}
