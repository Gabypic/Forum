package server

import (
	"Forum/application"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates = template.Must(template.ParseGlob("web/*.html"))

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
	if userDatas == nil {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := application.GetUser(email)
		if err != nil || user == nil {
			fmt.Println("1")
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		fmt.Println(user.Username)
		session := CreateSession(user.Username)
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
	}
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	posts, err := application.GetRecentPosts()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	comments, err := application.GetAllComments()
	if err != nil {
		log.Printf("Error loading comments: %v", err)
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Categories": categories,
		"Posts":      posts,
		"Comments":   comments,
	}

	renderTemplate(w, "home", data)
}

func handleHomePageRegister(w http.ResponseWriter, r *http.Request) {
	RegisterUserHandler(w, r)
}

func handleProfilPage(w http.ResponseWriter, r *http.Request) {
	user, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(user)
	fmt.Println(user)
	if r.Method == http.MethodGet {
		if r.Method == http.MethodGet {
			data := map[string]interface{}{
				"User":     userDatas.Username,
				"Mail":     userDatas.Mail,
				"JoinDate": userDatas.joinDate,
			}
			renderTemplate(w, "profil", data)
			return
		}
	}
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
	renderTemplate(w, "view_category", map[string]interface{}{"Category": category})
}

func handleUpdateCategoryPage(w http.ResponseWriter, r *http.Request) {
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	category := application.Category{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		CreatedBy:   r.FormValue("created_by"),
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
		renderTemplate(w, "create_post", nil)
		return
	}
	CreatePostHandler(w, r)
}

func handleGetPostPage(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var authorButtons bool
	if userDatas.Admin == true || userDatas.Modo == true {
		authorButtons = true
	} else {
		authorButtons = false
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
	data := map[string]interface{}{
		"Post":                  post,
		"ShowEditDeleteButtons": authorButtons,
	}
	renderTemplate(w, "view_post", data)
}

func handleUpdatePostPage(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == http.MethodPost {
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
	user, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(user)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	post := application.Post{
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		ImageURL:   r.FormValue("image_url"),
		CreatedBy:  userDatas.Username,
		CategoryID: atoi(r.FormValue("category_id")),
		Approved:   r.FormValue("approved") == "true",
	}

	err := application.CreatePost(&post)
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

	post := application.Post{
		ID:         id,
		Title:      r.FormValue("title"),
		Content:    r.FormValue("content"),
		ImageURL:   r.FormValue("image_url"),
		CreatedBy:  r.FormValue("created_by"),
		CategoryID: atoi(r.FormValue("category_id")),
		Approved:   r.FormValue("approved") == "true",
	}

	err = application.UpdatePost(&post)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		renderTemplate(w, "create_comment", nil)
		return
	}
	CreateCommentHandler(w, r)
}

func handleGetCommentPage(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("id")
	userTest, _ := GetSessionCookie(r)
	userDatas, _ := GetSession(userTest)
	var authorButtons bool
	if userDatas.Admin == true || userDatas.Modo == true {
		authorButtons = true
	} else {
		authorButtons = false
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
		"ShowEditDeleteButtons": authorButtons,
	}
	renderTemplate(w, "view_comment", data)
}

func handleUpdateCommentPage(w http.ResponseWriter, r *http.Request) {
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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	comment := application.Comment{
		Content:   r.FormValue("content"),
		CreatedBy: r.FormValue("created_by"),
		PostID:    atoi(r.FormValue("post_id")),
		Approved:  r.FormValue("approved") == "true",
	}

	err := application.CreateComment(&comment)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}
