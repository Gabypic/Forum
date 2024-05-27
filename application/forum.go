package application

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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

func Password(password string) (string, error) {
	hach, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hach), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(user *User) error {
	hashedPassword, err := Password(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	query := `INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)`
	_, err = DB.Exec(query, user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	return err
}

func GetUser(email string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = ?`
	value, err := DB.Query(query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("err")
			log.Printf("Error login user: %v", err)
			return nil, err
		}
		log.Printf("Error login user: %v", err)
	}
	value.Next()
	errs := value.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if errs != nil {
		fmt.Print("errs")
		fmt.Println(errs)
		return nil, errs
	}
	return &user, nil
}

func UpdateUser(user *User) error {
	hashedPassword, err := Password(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.Password = hashedPassword

	query := `UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?`
	_, err = DB.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
	}
	return err
}

func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
	}
	return err
}
