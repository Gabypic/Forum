package server

import (
	"Forum/application"
	"errors"
	"fmt"
	"time"
)

type Session struct {
	ID        string
	Username  string
	ExpiresAt time.Time
}

var sessions = map[string]*Session{}

func CreateSession(username string) *Session {
	sessionID := application.GenerateUUID()
	session := &Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	sessions[sessionID] = session
	fmt.Println(sessions)
	return session
}

func GetSession(sessionID string) (*Session, error) {
	session, exists := sessions[sessionID]
	if !exists || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session not found or expired")
	}
	fmt.Println(sessions)
	return session, nil
}

func DeleteSession(sessionID string) {
	delete(sessions, sessionID)
	fmt.Println(sessions)
}
