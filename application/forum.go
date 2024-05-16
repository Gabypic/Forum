package application

import (
	"time"
)

// Modèles
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ImageURL   string    `json:"image_url"`
	CreatedBy  string    `json:"created_by"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	Approved   bool      `json:"approved"`
}

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"created_by"`
	PostID    int       `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	Approved  bool      `json:"approved"`
}

type Reaction struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	CreatedBy string    `json:"created_by"`
	PostID    *int      `json:"post_id"`
	CommentID *int      `json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
}
