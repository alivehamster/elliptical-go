package utils

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alivehamster/elliptical-go/types"
	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dsn string) (*sql.DB, error) {
	if err := ensureDatabase(dsn); err != nil {
		return nil, fmt.Errorf("failed to ensure database exists: %w", err)
	}

	dbDSN := addDatabaseToDSN(dsn, "elliptical")
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := ensureTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func ensureDatabase(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS elliptical")
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

func addDatabaseToDSN(dsn, database string) string {
	if strings.Contains(dsn, "?") {
		parts := strings.Split(dsn, "?")
		if strings.HasSuffix(parts[0], "/") {
			return parts[0] + database + "?" + parts[1]
		} else {
			return parts[0] + "/" + database + "?" + parts[1]
		}
	} else {
		if strings.HasSuffix(dsn, "/") {
			return dsn + database
		} else {
			return dsn + "/" + database
		}
	}
}

func ensureTables(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS room (
        roomid INT PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(255) NOT NULL
    )`)
	if err != nil {
		return fmt.Errorf("failed to create room table: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
        msgid INT PRIMARY KEY AUTO_INCREMENT,
        content TEXT NOT NULL,
        roomid INT NOT NULL,
        FOREIGN KEY (roomid) REFERENCES room(roomid) ON DELETE CASCADE
    )`)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}
	return nil
}

func GetRooms(db *sql.DB) ([]types.Room, error) {
	rows, err := db.Query("SELECT roomid, title FROM room")
	if err != nil {
		return nil, fmt.Errorf("failed to query rooms: %w", err)
	}
	defer rows.Close()

	var rooms []types.Room
	for rows.Next() {
		var room types.Room
		if err := rows.Scan(&room.RoomID, &room.Title); err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rooms: %w", err)
	}

	return rooms, nil
}

func CreateRoom(db *sql.DB, title string) (int64, error) {
	result, err := db.Exec("INSERT INTO room (title) VALUES (?)", title)
	if err != nil {
		return 0, fmt.Errorf("failed to create room: %w", err)
	}

	roomID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return roomID, nil
}

func GetChats(db *sql.DB, roomID string) ([]types.Chat, error) {
	rows, err := db.Query("SELECT msgid, content FROM messages WHERE roomid = ?", roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chats: %w", err)
	}
	defer rows.Close()

	var chats []types.Chat
	for rows.Next() {
		var chat types.Chat
		if err := rows.Scan(&chat.ID, &chat.Msg); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over chats: %w", err)
	}

	return chats, nil
}

func StoreChat(db *sql.DB, roomID string, msg string) (int64, error) {
	result, err := db.Exec("INSERT INTO messages (content, roomid) VALUES (?, ?)", msg, roomID)
	if err != nil {
		return 0, fmt.Errorf("failed to store chat: %w", err)
	}

	msgID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return msgID, nil
}
