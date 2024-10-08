package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	UserID      int64
	DiscordID   string
	DiscordName string
	CharUsage   int
	TTSCount    int
	SFXCount    int
	CreatedAt   string
}

type Usage struct {
	UserID        int64
	AudioType     string
	AudioService  string
	Prompt        string
	AudioFilename string
}

func get_db() *sql.DB {
	db, err := sql.Open("sqlite3", db_path)

	if err != nil {
		// TODO
	}

	return db
}

func get_user(discord_id string) *User {
	db := get_db()

	query := fmt.Sprintf("SELECT * FROM user WHERE discord_id = ?")

	row := db.QueryRow(query, discord_id)

	var user User
	err := row.Scan(&user.UserID, &user.DiscordID, &user.DiscordName, &user.CharUsage, &user.TTSCount, &user.SFXCount, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	return &user
}

func add_user(discord_id string, discord_name string) (*User, error) {
	db := get_db()

	query := `
	INSERT into user(discord_id, discord_name) values(?, ?)
	`

	_, err := db.Exec(query, discord_id, discord_name)
	if err != nil {
		return nil, err
	}

	user := get_user(discord_id)

	return user, nil
}

func add_usage(usage Usage) error {
	db := get_db()

	usage_query := `
	INSERT into audio(user, audio_type, audio_service, prompt, audio_filename) values (?, ?, ?, ?, ?)
	`

	// TODO: Update character usage and tts/sfx count

	_, err := db.Exec(usage_query, usage.UserID, usage.AudioType, usage.AudioService, usage.Prompt, usage.AudioFilename)

	if err != nil {
		return err
	}

	return nil
}

func AddUsage(discord_user DiscordUser, usage Usage) {
	user := get_user(discord_user.ID)
	if user == nil {
		log.Println("User doesn't exist")
		user, _ = add_user(discord_user.ID, discord_user.Name)
	} else {
		log.Println("User exists")
	}

	usage.UserID = user.UserID

	add_usage(usage)
}
