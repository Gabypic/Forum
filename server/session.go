package server

import (
	"Forum/application"
	"errors"
	"fmt"
	"time"
)

type Session struct {
	Id        string
	Mail      string
	Username  string
	joinDate  time.Time
	ExpiresAt time.Time
}

var sessions = map[string]*Session{}

func CreateSession(username string) *Session {
	sessionID := application.GenerateUUID()
	userDatas, _ := application.GetUserByName(username)
	session := &Session{
		Id:        sessionID,
		Mail:      userDatas.Email,
		Username:  username,
		joinDate:  userDatas.CreatedAt,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	sessions[sessionID] = session
	fmt.Println("error session")
	fmt.Println(sessions)
	return session
}

func GetSession(sessionID string) (*Session, error) {
	fmt.Println(sessions)
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
