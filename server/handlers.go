package server

import (
	"Forum/application"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseGlob("web/*.html"))

type SuggestionsData struct {
	Query       string
	Suggestions []string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", tmpl, err)
		http.Error(w, "Internal Server Error(212)", http.StatusInternalServerError)
		return
	}
}

func handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "account_creation", nil)
		return
	}
}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "login", nil)
		return
	}
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	categories, err := application.GetAllCategories()
	if err != nil {
		log.Printf("Error loading categories: %v", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}
	log.Printf("Loaded categories: %v", categories)

	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var showEditDeleteButtons bool
	var guestLogin bool
	var session *Session
	if userDatas == nil {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" && password == "" {
			email = "guest@guest"
			password = "guest"
		}

		user, err := application.GetUser(email)
		if err != nil || user == nil {
			fmt.Println("1")
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		fmt.Println(user.Username)
		session = CreateSession(user.Username)
		SetSessionCookie(w, session.Id)

		if !application.CheckPassword(password, user.Password) {
			fmt.Println("2")
			log.Println("Password does not match")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		fmt.Println("3")
		fmt.Print("yoyoyo")
		fmt.Println(user)
		userDatas, _ = GetSession(session.Id)
	}
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	if userDatas != nil && (userDatas.Admin == true || userDatas.Modo == true) {
		showEditDeleteButtons = true
	} else {
		showEditDeleteButtons = false
	}

	if userDatas != nil && (userDatas.Guest == true) {
		guestLogin = true
	} else {
		guestLogin = false
	}

	posts, err := application.GetUncategorizedPostsWithComments()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	users, err := application.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("query")
	var suggestions []string
	if query != "" {
		for _, user := range users {
			if strings.Contains(strings.ToLower(user), strings.ToLower(query)) {
				suggestions = append(suggestions, user)
			}
		}
	}

	data := map[string]interface{}{
		"Categories":            categories,
		"Posts":                 posts,
		"ShowEditDeleteButtons": showEditDeleteButtons,
		"Query":                 query,
		"Suggestions":           suggestions,
		"guestLogin":            guestLogin,
	}

	renderTemplate(w, "home", data)
}

func handleHomePageRegister(w http.ResponseWriter, r *http.Request) {
	RegisterUserHandler(w, r)
}

func handleProfilPage(w http.ResponseWriter, r *http.Request) {
	user, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(user)
	new_username := r.FormValue("new_username")
	new_email := r.FormValue("new_email")
	new_password := r.FormValue("new_password")

	if new_username != "" || new_email != "" || new_password != "" {
		err := application.UpdateUser(userDatas.Username, new_username, new_email, new_password)
		if err != nil {
			fmt.Println("erreur de la modification du compte ", err)
		}
	}

	var guestLogin bool
	if userDatas == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createdPosts, err := application.GetPostsByUser(userDatas.Username)
	if err != nil {
		log.Printf("Error fetching created posts: %v", err)
		http.Error(w, "Failed to load created posts", http.StatusInternalServerError)
		return
	}

	likedPosts, err := application.GetLikedPostsByUser(userDatas.Username)
	if err != nil {
		log.Printf("Error fetching liked posts: %v", err)
		http.Error(w, "Failed to load liked posts", http.StatusInternalServerError)
		return
	}

	if userDatas != nil && (userDatas.Guest == true) {
		guestLogin = true
	} else {
		guestLogin = false
	}

	log.Printf("Created posts: %v", createdPosts)
	log.Printf("Liked posts: %v", likedPosts)

	data := map[string]interface{}{
		"User":         userDatas.Username,
		"Mail":         userDatas.Mail,
		"JoinDate":     userDatas.joinDate,
		"CreatedPosts": createdPosts,
		"LikedPosts":   likedPosts,
		"guestLogin":   guestLogin,
	}

	renderTemplate(w, "profil", data)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	user := application.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user.ID = application.GenerateUUID()

	session := CreateSession(user.Username)
	SetSessionCookie(w, session.Id)

	err := application.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "home", nil)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := application.GetUser(email)
	if err != nil || user == nil {
		fmt.Println("1")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	fmt.Println(user.Username)
	session := CreateSession(user.Username)
	SetSessionCookie(w, session.Id)

	if !application.CheckPassword(password, user.Password) {
		fmt.Println("2")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	fmt.Println("3")
	fmt.Print("yoyoyo")
	fmt.Println(user)
	renderTemplate(w, "home", nil)
}

func handleCreateCategoryPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "create_category", nil)
		return
	}
	CreateCategoryHandler(w, r)
}

func handleGetCategoryPage(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var showEditDeleteButtons bool
	if userDatas != nil && (userDatas.Admin == true || userDatas.Modo == true) {
		showEditDeleteButtons = true
	} else {
		showEditDeleteButtons = false
	}

	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	category, err := application.GetCategory(id)
	if err != nil || category == nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	posts, err := application.GetPostsByCategoryID(id)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	for i, post := range posts {
		comments, err := application.GetCommentsByPostID(post.ID)
		if err != nil {
			http.Error(w, "Failed to load comments", http.StatusInternalServerError)
			return
		}
		posts[i].Comments = comments
	}

	data := map[string]interface{}{
		"Category":              category,
		"Posts":                 posts,
		"ShowEditDeleteButtons": showEditDeleteButtons,
	}
	renderTemplate(w, "view_category_posts", data)
}

func handleUpdateCategoryPage(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil || (!userDatas.Admin && !userDatas.Modo) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		request := r.URL.Query().Get("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		category, err := application.GetCategory(id)
		if err != nil || category == nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"Category": category,
		}
		renderTemplate(w, "update_category", data)
		return
	}
	UpdateCategoryHandler(w, r)
}

func handleDeleteCategoryPage(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil || (!userDatas.Admin && !userDatas.Modo) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		request := r.FormValue("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		DeleteCategoryHandler(w, r, id)
	}
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	category := application.Category{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		CreatedBy:   userDatas.Username,
	}

	err := application.CreateCategory(&category)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := application.GetCategory(id)
	if err != nil {
		http.Error(w, "Failed to get category", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "category", map[string]interface{}{"Category": category})
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	request := r.FormValue("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category := application.Category{
		ID:          id,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		CreatedBy:   r.FormValue("created_by"),
	}

	err = application.UpdateCategory(&category)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request, id int) {
	err := application.DeleteCategory(id)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func handleCreatePostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		categories, err := application.GetAllCategories()
		if err != nil {
			log.Printf("Error loading categories: %v", err)
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Categories": categories,
		}
		renderTemplate(w, "create_post", data)
		return
	}
	CreatePostHandler(w, r)
}

func handleGetPostPage(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var guestLogin bool
	var showEditDeleteButtons bool
	if userDatas != nil && (userDatas.Admin == true || userDatas.Modo == true) {
		showEditDeleteButtons = true
	} else {
		showEditDeleteButtons = false
	}
	if userDatas != nil && (userDatas.Guest == true) {
		guestLogin = true
	} else {
		guestLogin = false
	}
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := application.GetPost(id)
	if err != nil || post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, err := application.GetCommentsByPostID(id)
	if err != nil {
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	likeCount, err := application.GetLikeCount(&id, nil)
	if err != nil {
		http.Error(w, "Failed to get like count", http.StatusInternalServerError)
		return
	}
	unlikeCount, err := application.GetUnlikeCount(&id, nil)
	if err != nil {
		http.Error(w, "Failed to get unlike count", http.StatusInternalServerError)
		return
	}

	commentLikeCounts := make(map[int]int)
	commentUnlikeCounts := make(map[int]int)
	for _, comment := range comments {
		commentLikeCount, err := application.GetLikeCount(nil, &comment.ID)
		if err != nil {
			http.Error(w, "Failed to get comment like count", http.StatusInternalServerError)
			return
		}
		commentUnlikeCount, err := application.GetUnlikeCount(nil, &comment.ID)
		if err != nil {
			http.Error(w, "Failed to get comment unlike count", http.StatusInternalServerError)
			return
		}
		commentLikeCounts[comment.ID] = commentLikeCount
		commentUnlikeCounts[comment.ID] = commentUnlikeCount
	}

	data := map[string]interface{}{
		"Post":                  post,
		"Comments":              comments,
		"ShowEditDeleteButtons": showEditDeleteButtons,
		"LikeCount":             likeCount,
		"UnlikeCount":           unlikeCount,
		"CommentLikeCounts":     commentLikeCounts,
		"CommentUnlikeCounts":   commentUnlikeCounts,
		"guestLogin":            guestLogin,
	}
	renderTemplate(w, "view_post", data)
}

func handleUpdatePostPage(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil || (!userDatas.Admin && !userDatas.Modo) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		request := r.URL.Query().Get("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		post, err := application.GetPost(id)
		if err != nil || post == nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"Post": post,
		}
		renderTemplate(w, "update_post", data)
		return
	}
	UpdatePostHandler(w, r)
}

func handleDeletePostPage(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil || (!userDatas.Admin && !userDatas.Modo) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		request := r.URL.Query().Get("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		post, err := application.GetPost(id)
		if err != nil || post == nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"Post": post,
		}
		renderTemplate(w, "delete_post", data)
	} else if r.Method == http.MethodPost {
		request := r.FormValue("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		DeletePostHandler(w, r, id)
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseMultipartForm(20 << 20)

	file, header, err := r.FormFile("image")
	var filePath string
	if err == nil {
		defer file.Close()
		fileType := header.Header.Get("Content-Type")
		if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		if header.Size > 20*1024*1024 {
			http.Error(w, "File is too large", http.StatusBadRequest)
			return
		}

		filePath = fmt.Sprintf("./images/%s", header.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Error copying the file", http.StatusInternalServerError)
			return
		}
	}

	post := application.Post{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		ImageURL:   filePath,
		CreatedBy:  userDatas.Username,
		CategoryID: atoi(r.FormValue("category_id")),
		Approved:   r.FormValue("approved") == "true",
	}

	err = application.CreatePost(&post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := application.GetPost(id)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "post", map[string]interface{}{"Post": post})
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	request := r.FormValue("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	existingPost, err := application.GetPost(id)
	if err != nil || existingPost == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(20 << 20)

	var filePath string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		fileType := header.Header.Get("Content-Type")
		if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		if header.Size > 20*1024*1024 {
			http.Error(w, "File is too large", http.StatusBadRequest)
			return
		}

		filePath = fmt.Sprintf("./images/%s", header.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to create the file for writing", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Error copying the file", http.StatusInternalServerError)
			return
		}
	} else {
		filePath = existingPost.ImageURL
	}

	post := application.Post{
		ID:         id,
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		ImageURL:   filePath,
		CreatedBy:  existingPost.CreatedBy,
		CategoryID: existingPost.CategoryID,
		Approved:   r.FormValue("approved") == "true",
	}

	err = application.UpdatePost(&post)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view_category?id=%d", existingPost.CategoryID), http.StatusSeeOther)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request, id int) {
	err := application.DeletePost(id)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func handleCreateCommentPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		postID := r.URL.Query().Get("post_id")
		categoryID := r.URL.Query().Get("category_id")
		data := map[string]interface{}{
			"PostID":     postID,
			"CategoryID": categoryID,
		}
		renderTemplate(w, "create_comment", data)
		return
	}
	CreateCommentHandler(w, r)
}

func handleGetCommentPage(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var showEditDeleteButtons bool
	if userDatas != nil && (userDatas.Admin == true || userDatas.Modo == true) {
		showEditDeleteButtons = true
	} else {
		showEditDeleteButtons = false
	}
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	comment, err := application.GetComment(id)
	if err != nil || comment == nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	data := map[string]interface{}{
		"Comment":               comment,
		"ShowEditDeleteButtons": showEditDeleteButtons,
	}
	renderTemplate(w, "view_comment", data)
}

func handleUpdateCommentPage(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil || (!userDatas.Admin && !userDatas.Modo) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		request := r.URL.Query().Get("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		comment, err := application.GetComment(id)
		if err != nil || comment == nil {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"Comment": comment,
		}
		renderTemplate(w, "update_comment", data)
		return
	}
	UpdateCommentHandler(w, r)
}

func handleDeleteCommentPage(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if userDatas == nil || (!userDatas.Admin && !userDatas.Modo) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		request := r.FormValue("id")
		id, err := strconv.Atoi(request)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		DeleteCommentHandler(w, r, id)
	}
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := atoi(r.FormValue("post_id"))
	categoryID := atoi(r.FormValue("category_id"))

	comment := application.Comment{
		Content:   r.FormValue("content"),
		CreatedBy: userDatas.Username,
		PostID:    postID,
		Approved:  r.FormValue("approved") == "true",
	}

	err := application.CreateComment(&comment)
	if err != nil {
		log.Printf("Failed to create comment: %v", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view_category_posts?id=%d", categoryID), http.StatusSeeOther)
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := application.GetComment(id)
	if err != nil {
		http.Error(w, "Failed to get comment", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "comment", map[string]interface{}{"Comment": comment})
}

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	request := r.FormValue("id")
	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment := application.Comment{
		ID:        id,
		Content:   r.FormValue("content"),
		CreatedBy: r.FormValue("created_by"),
		PostID:    atoi(r.FormValue("post_id")),
		Approved:  r.FormValue("approved") == "true",
	}

	err = application.UpdateComment(&comment)
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view_post?id=%d", comment.PostID), http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request, id int) {
	comment, err := application.GetComment(id)
	if err != nil || comment == nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	err = application.DeleteComment(id)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/view_post?id=%d", comment.PostID), http.StatusSeeOther)
}

func handleGetCategoryPostsPage(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var showEditDeleteButtons bool
	if userDatas != nil && (userDatas.Admin == true || userDatas.Modo == true) {
		showEditDeleteButtons = true
	} else {
		showEditDeleteButtons = false
	}

	id, err := strconv.Atoi(request)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	category, err := application.GetCategory(id)
	if err != nil || category == nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	posts, err := application.GetPostsByCategoryID(id)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	for i, post := range posts {
		comments, err := application.GetCommentsByPostID(post.ID)
		if err != nil {
			http.Error(w, "Failed to load comments", http.StatusInternalServerError)
			return
		}
		posts[i].Comments = comments
	}

	data := map[string]interface{}{
		"Category":              category,
		"Posts":                 posts,
		"ShowEditDeleteButtons": showEditDeleteButtons,
	}
	renderTemplate(w, "view_category_posts", data)
}

func handleLike(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionCookie(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userDatas, _ := GetSession(user)

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))

	reaction := &application.Reaction{
		Type:      "like",
		CreatedBy: userDatas.Username,
		PostID:    &postID,
		CommentID: &commentID,
	}

	err = application.ToggleReaction(reaction)
	if err != nil {
		http.Error(w, "Failed to like", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func handleUnlike(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionCookie(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userDatas, _ := GetSession(user)

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))

	reaction := &application.Reaction{
		Type:      "unlike",
		CreatedBy: userDatas.Username,
		PostID:    &postID,
		CommentID: &commentID,
	}

	err = application.ToggleReaction(reaction)
	if err != nil {
		http.Error(w, "Failed to unlike", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func disconnection(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	ClearSessionCookie(w, userDatas.Id)
	DeleteSession(userDatas.Id)
	renderTemplate(w, "login", nil)
}

func handleUsersProfil(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	var isAdmin bool

	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	userDatas, err := application.GetUserByName(username)
	if err != nil {
		http.Error(w, "Error retrieving user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	createdPost, err := application.GetPostsByUser(username)
	if err != nil {
		http.Error(w, "Error retrieving user posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	likedPosts, err := application.GetLikedPostsByUser(username)
	if err != nil {
		http.Error(w, "Error retrieving user liked posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(likedPosts, "liked")

	if userDatas != nil && (userDatas.Admin == true) {
		isAdmin = true
	} else {
		isAdmin = false
	}

	data := map[string]interface{}{
		"User":         userDatas.Username,
		"Mail":         userDatas.Email,
		"JoinDate":     userDatas.CreatedAt,
		"CreatedPosts": createdPost,
		"LikedPosts":   likedPosts,
		"isAdmin":      isAdmin,
	}

	renderTemplate(w, "user_profil", data)
}

func delete_account(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	application.DeleteUser(userDatas.Username)
	ClearSessionCookie(w, userDatas.Id)
	DeleteSession(userDatas.Id)
	renderTemplate(w, "login", nil)
}

func handle_modification_user(w http.ResponseWriter, r *http.Request) {
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)

	data := map[string]interface{}{
		"User":  userDatas.Username,
		"Email": userDatas.Mail,
	}

	renderTemplate(w, "update_user", data)
}

func atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}
