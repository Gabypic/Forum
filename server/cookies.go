package server

import (
	"fmt"
	"net/http"
	"time"
)

func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	fmt.Println(cookie)
	http.SetCookie(w, cookie)
}

func GetSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	fmt.Println(cookie)
	return cookie.Value, nil
}

func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	fmt.Println(cookie)
	http.SetCookie(w, cookie)
}
