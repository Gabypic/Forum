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
	Guest     bool
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

// Function that hashes a password using bcrypt
func Password(password string) (string, error) {
	hach, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hach), err
}

// Function that checks a password against a hash using bcrypt
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Function that creates a user in the database
func CreateUser(user *User) error {
	hashedPassword, err := Password(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()

	query := `INSERT INTO users (username, email, password, created_at, admin, modo, guest) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, user.Username, user.Email, user.Password, user.CreatedAt, 0, 0, 0)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	return err
}

// Function that retrieves a user by email
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

// Function that retrieves a user by username
func GetUserByName(name string) (*User, error) {
	var user User
	query := `SELECT id, username, email, password, created_at, admin, modo, guest FROM users WHERE username = ?`
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
		errs := value.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.Admin, &user.Modo, &user.Guest)
		if errs != nil {
			fmt.Print("errs")
			fmt.Println(errs)
			return nil, errs
		}
	}
	return &user, nil
}

// Function that updates a user's information
func UpdateUser(ancientName string, name string, email string, password string) error {
	hashedPassword, err := Password(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	query := `UPDATE users SET username = ?, email = ?, password = ? WHERE username = ?`
	_, err = DB.Exec(query, name, email, hashedPassword, ancientName)
	if err != nil {
		log.Printf("Error updating user: %v", err)
	}
	return err
}

// Function that deletes a user from the database
func DeleteUser(username string) error {
	query := `DELETE FROM users WHERE username = ?`
	_, err := DB.Exec(query, username)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
	}
	return err
}

// Function that creates a category in the database
func CreateCategory(category *Category) error {
	category.CreatedAt = time.Now()

	query := `INSERT INTO categories (name, description, created_by, created_at) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, category.Name, category.Description, category.CreatedBy, category.CreatedAt)
	if err != nil {
		log.Printf("Error creating category: %v", err)
	}
	return err
}

// Function that retrieves a category by ID
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

// Function that updates a category's information
func UpdateCategory(category *Category) error {
	query := `UPDATE categories SET name = ?, description = ?, created_by = ? WHERE id = ?`
	_, err := DB.Exec(query, category.Name, category.Description, category.CreatedBy, category.ID)
	if err != nil {
		log.Printf("Error updating category: %v", err)
	}
	return err
}

// Function that deletes a category from the database
func DeleteCategory(id int) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
	}
	return err
}

// Function that creates a post in the database
func CreatePost(post *Post) error {
	post.CreatedAt = time.Now()

	query := `INSERT INTO posts (title, content, image_url, created_by, category_id, created_at, approved) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, post.Title, post.Content, post.ImageURL, post.CreatedBy, post.CategoryID, post.CreatedAt, post.Approved)
	if err != nil {
		log.Printf("Error creating post: %v", err)
	}
	return err
}

// Function that retrieves a post by ID
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

// Function that updates a post's information
func UpdatePost(post *Post) error {
	query := `UPDATE posts SET title = ?, content = ?, image_url = ?, created_by = ?, category_id = ?, approved = ? WHERE id = ?`
	_, err := DB.Exec(query, post.Title, post.Content, post.ImageURL, post.CreatedBy, post.CategoryID, post.Approved, post.ID)
	if err != nil {
		log.Printf("Error updating post: %v", err)
	}
	return err
}

// Function that deletes a post from the database
func DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
	}
	return err
}

// Function that creates a comment in the database
func CreateComment(comment *Comment) error {
	comment.CreatedAt = time.Now()

	query := `INSERT INTO comments (content, created_by, post_id, created_at, approved) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, comment.Content, comment.CreatedBy, comment.PostID, comment.CreatedAt, comment.Approved)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
	}
	return err
}

// Function that retrieves a comment by ID
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

// Function that updates a comment's information
func UpdateComment(comment *Comment) error {
	query := `UPDATE comments SET content = ? WHERE id = ?`
	_, err := DB.Exec(query, comment.Content, comment.ID)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
	}
	return err
}

// Function that deletes a comment from the database
func DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
	}
	return err
}

// Function that creates a reaction in the database
func CreateReaction(reaction *Reaction) error {
	reaction.CreatedAt = time.Now()

	query := `INSERT INTO reactions (type, created_by, post_id, comment_id, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, reaction.Type, reaction.CreatedBy, reaction.PostID, reaction.CommentID, reaction.CreatedAt)
	if err != nil {
		log.Printf("Error creating reaction: %v", err)
		return err
	}
	return nil
}

// Function that deletes a reaction from the database
func DeleteReaction(reaction *Reaction) error {
	query := `DELETE FROM reactions WHERE type = ? AND created_by = ? AND post_id = ? AND comment_id = ?`
	_, err := DB.Exec(query, reaction.Type, reaction.CreatedBy, reaction.PostID, reaction.CommentID)
	if err != nil {
		log.Printf("Error deleting reaction: %v", err)
		return err
	}
	return nil
}

// Function that gets the reaction count for a post or comment
func GetReactionCount(postID *int, commentID *int, reactionType string) (int, error) {
	var count int
	var query string

	if postID != nil && commentID == nil {
		query = `SELECT COUNT(*) FROM reactions WHERE type = ? AND post_id = ?`
		err := DB.QueryRow(query, reactionType, *postID).Scan(&count)
		if err != nil {
			return 0, err
		}
	} else if postID == nil && commentID != nil {
		query = `SELECT COUNT(*) FROM reactions WHERE type = ? AND comment_id = ?`
		err := DB.QueryRow(query, reactionType, *commentID).Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// Function that toggles a reaction (like/unlike)
func ToggleReaction(reaction *Reaction) error {
	hasLiked, err := UserHasReacted(reaction.CreatedBy, reaction.PostID, reaction.CommentID, "like")
	if err != nil {
		return err
	}
	hasUnliked, err := UserHasReacted(reaction.CreatedBy, reaction.PostID, reaction.CommentID, "unlike")
	if err != nil {
		return err
	}

	if hasLiked || hasUnliked {
		err = DeleteReaction(&Reaction{
			Type:      "like",
			CreatedBy: reaction.CreatedBy,
			PostID:    reaction.PostID,
			CommentID: reaction.CommentID,
		})
		if err != nil {
			return err
		}
		err = DeleteReaction(&Reaction{
			Type:      "unlike",
			CreatedBy: reaction.CreatedBy,
			PostID:    reaction.PostID,
			CommentID: reaction.CommentID,
		})
		if err != nil {
			return err
		}
	}

	err = CreateReaction(reaction)
	if err != nil {
		return err
	}

	return nil
}

// Function that checks if a user has reacted to a post or comment
func UserHasReacted(createdBy string, postID *int, commentID *int, reactionType string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM reactions WHERE type = ? AND created_by = ? AND post_id = ? AND comment_id = ?`
	err := DB.QueryRow(query, reactionType, createdBy, postID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Function that gets the like count for a post or comment
func GetLikeCount(postID *int, commentID *int) (int, error) {
	return GetReactionCount(postID, commentID, "like")
}

// Function that gets the unlike count for a post or comment
func GetUnlikeCount(postID *int, commentID *int) (int, error) {
	return GetReactionCount(postID, commentID, "unlike")
}

// Function that gets posts created by a user
func GetPostsByUser(username string) ([]Post, error) {
	rows, err := DB.Query(`
        SELECT id, title, content, image_url, created_by, category_id, created_at, approved
        FROM posts
        WHERE created_by = ?`, username)
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
		posts = append(posts, post)
	}
	return posts, nil
}

// Function that gets posts liked by a user
func GetLikedPostsByUser(username string) ([]Post, error) {
	rows, err := DB.Query(`
        SELECT p.id, p.title, p.content, p.image_url, p.created_by, p.category_id, p.created_at, p.approved
        FROM posts p
        JOIN reactions r ON p.id = r.post_id
        WHERE r.type = 'like' AND r.created_by = ?`, username)
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
		posts = append(posts, post)
	}
	return posts, nil
}

// Function that gets the count of posts liked by a user
func GetLikedPostCountByUser(username string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM reactions WHERE type = 'like' AND created_by = ? AND post_id IS NOT NULL`
	err := DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Function that retrieves all categories
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

// Function that retrieves recent posts with their comments
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

// Function that retrieves all comments
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

// Function that retrieves comments by post ID
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

// Function that retrieves posts by category ID
func GetPostsByCategoryID(categoryID int) ([]Post, error) {
	rows, err := DB.Query("SELECT id, title, content, image_url, created_by, category_id, created_at, approved FROM posts WHERE category_id = ? ORDER BY created_at DESC", categoryID)
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
		posts = append(posts, post)
	}

	return posts, nil
}

// Function that retrieves uncategorized posts with their comments
func GetUncategorizedPostsWithComments() ([]Post, error) {
	rows, err := DB.Query("SELECT id, title, content, image_url, created_by, category_id, created_at, approved FROM posts WHERE category_id IS NULL ORDER BY created_at DESC LIMIT 10")
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
		posts = append(posts, post)
	}

	return posts, nil
}

// Function that retrieves all usernames
func GetAllUsers() ([]string, error) {
	rows, err := DB.Query("SELECT username FROM users")
	if err != nil {
		return nil, err
	}
	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	fmt.Print("users2")
	fmt.Println(users)
	return users, nil
}
