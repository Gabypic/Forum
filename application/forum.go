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
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	Admin     bool
	Modo      bool
}

type Category struct {
	ID          int
	Name        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
}

type Post struct {
	ID         int
	Title      string
	Content    string
	ImageURL   string
	CreatedBy  string
	CategoryID int
	CreatedAt  time.Time
	Approved   bool
	Comments   []Comment
}

type Comment struct {
	ID        int
	Content   string
	CreatedBy string
	PostID    int
	CreatedAt time.Time
	Approved  bool
}

type Reaction struct {
	ID        int
	Type      string
	CreatedBy string
	PostID    *int
	CommentID *int
	CreatedAt time.Time
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

	query := `INSERT INTO users (username, email, password, created_at, admin, modo) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, user.Username, user.Email, user.Password, user.CreatedAt, 0, 0)
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
		return nil, err
	}
	defer value.Close()

	if value.Next() {
		errs := value.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
		if errs != nil {
			fmt.Print("errs")
			fmt.Println(errs)
			return nil, errs
		}
	}
	return &user, nil
}

func GetUserByName(name string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password, created_at, admin, modo FROM users WHERE username = ?`
	value, err := DB.Query(query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("err")
			log.Printf("Error login user: %v", err)
			return nil, err
		}
		log.Printf("Error login user: %v", err)
		return nil, err
	}
	defer value.Close()

	if value.Next() {
		errs := value.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.Admin, &user.Modo)
		if errs != nil {
			fmt.Print("errs")
			fmt.Println(errs)
			return nil, errs
		}
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

func CreateCategory(category *Category) error {
	category.CreatedAt = time.Now()

	query := `INSERT INTO categories (name, description, created_by, created_at) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, category.Name, category.Description, category.CreatedBy, category.CreatedAt)
	if err != nil {
		log.Printf("Error creating category: %v", err)
	}
	return err
}

func GetCategory(id int) (*Category, error) {
	var category Category
	query := `SELECT id, name, description, created_by, created_at FROM categories WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedBy, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error retrieving category: %v", err)
	}
	return &category, err
}

func UpdateCategory(category *Category) error {
	query := `UPDATE categories SET name = ?, description = ?, created_by = ? WHERE id = ?`
	_, err := DB.Exec(query, category.Name, category.Description, category.CreatedBy, category.ID)
	if err != nil {
		log.Printf("Error updating category: %v", err)
	}
	return err
}

func DeleteCategory(id int) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
	}
	return err
}

func CreatePost(post *Post) error {
	post.CreatedAt = time.Now()

	query := `INSERT INTO posts (title, content, image_url, created_by, category_id, created_at, approved) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, post.Title, post.Content, post.ImageURL, post.CreatedBy, post.CategoryID, post.CreatedAt, post.Approved)
	if err != nil {
		log.Printf("Error creating post: %v", err)
	}
	return err
}

func GetPost(id int) (*Post, error) {
	var post Post
	query := `SELECT id, title, content, image_url, created_by, category_id, created_at, approved FROM posts WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedBy, &post.CategoryID, &post.CreatedAt, &post.Approved)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error retrieving post: %v", err)
	}
	return &post, err
}

func UpdatePost(post *Post) error {
	query := `UPDATE posts SET title = ?, content = ?, image_url = ?, created_by = ?, category_id = ?, approved = ? WHERE id = ?`
	_, err := DB.Exec(query, post.Title, post.Content, post.ImageURL, post.CreatedBy, post.CategoryID, post.Approved, post.ID)
	if err != nil {
		log.Printf("Error updating post: %v", err)
	}
	return err
}

func DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
	}
	return err
}

func CreateComment(comment *Comment) error {
	comment.CreatedAt = time.Now()

	query := `INSERT INTO comments (content, created_by, post_id, created_at, approved) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, comment.Content, comment.CreatedBy, comment.PostID, comment.CreatedAt, comment.Approved)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
	}
	return err
}

func GetComment(id int) (*Comment, error) {
	var comment Comment
	query := `SELECT id, content, created_by, post_id, created_at, approved FROM comments WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&comment.ID, &comment.Content, &comment.CreatedBy, &comment.PostID, &comment.CreatedAt, &comment.Approved)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error retrieving comment: %v", err)
	}
	return &comment, err
}

func UpdateComment(comment *Comment) error {
	query := `UPDATE comments SET content = ? WHERE id = ?`
	_, err := DB.Exec(query, comment.Content, comment.ID)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
	}
	return err
}

func DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
	}
	return err
}

func GetAllCategories() ([]Category, error) {
	rows, err := DB.Query("SELECT id, name, description, created_by, created_at FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedBy, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetRecentPostsWithComments() ([]Post, error) {
	rows, err := DB.Query("SELECT id, title, content, image_url, created_by, category_id, created_at, approved FROM posts ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedBy, &post.CategoryID, &post.CreatedAt, &post.Approved); err != nil {
			return nil, err
		}
		comments, err := GetCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}

	return posts, nil
}

func GetAllComments() ([]Comment, error) {
	rows, err := DB.Query("SELECT id, content, created_by, post_id, created_at, approved FROM comments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedBy, &comment.PostID, &comment.CreatedAt, &comment.Approved); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func GetCommentsByPostID(postID int) ([]Comment, error) {
	rows, err := DB.Query("SELECT id, content, created_by, post_id, created_at, approved FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedBy, &comment.PostID, &comment.CreatedAt, &comment.Approved); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
