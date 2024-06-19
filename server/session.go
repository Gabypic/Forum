package server

import (
	"Forum/application"
	"errors"
	"fmt"
	"time"
)

// Structure containing user session information.
type Session struct {
	Id        string
	Mail      string
	Username  string
	joinDate  time.Time
	ExpiresAt time.Time
	Modo      bool
	Admin     bool
	Guest     bool
}

var sessions = map[string]*Session{}

func CreateSession(username string) *Session {
	sessionID := application.GenerateUUID()
	userDatas, _ := application.GetUserByName(username)
	// Function that creates a new session for a user.
	// Generates a new UUID for the session.
	// Retrieves user data by username.
	session := &Session{
		Id:        sessionID,
		Mail:      userDatas.Email,
		Username:  username,
		joinDate:  userDatas.CreatedAt,
		Modo:      userDatas.Modo,
		Admin:     userDatas.Admin,
		Guest:     userDatas.Guest,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	// Initializes the Session structure with user data and expiration date.
	sessions[sessionID] = session
	fmt.Println("error session")
	fmt.Println(sessions)
	return session
}

// Function that retrieves a session by its ID.
func GetSession(sessionID string) (*Session, error) {
	fmt.Println(sessions)
	session, exists := sessions[sessionID]
	if !exists || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session not found or expired")
	}
	fmt.Println(sessions)
	return session, nil
}

// DeleteSession function that deletes a session by its ID.
func DeleteSession(sessionID string) {
	delete(sessions, sessionID)
	fmt.Println(sessions)
}
