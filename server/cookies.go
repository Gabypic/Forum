package server

import (
	"net/http"
	"time"
)

// Function that sets a cookie.
// Creates a cookie with the name "session_id" and the provided sessionID value and expiration to 24 hours.
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

// Function that retrieves a cookie.
func GetSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Function that clears a cookie.
// Delete the cookies by setting expiration to now.
func ClearSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, cookie)
}
