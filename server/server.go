package server

import (
	"Forum/application"
	"log"
	"net/http"
)

// Function that initializes and starts the server
func Start() {
	cfg := application.LoadConfig()
	application.Connect(cfg.DatabaseURL)
	application.SqlTable()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Handling static files

	http.HandleFunc("/login", handleLoginPage)
	http.HandleFunc("/create_account", handleRegisterPage)
	http.HandleFunc("/home", handleHomePage)
	http.HandleFunc("/register", handleHomePageRegister)
	http.HandleFunc("/create_category", handleCreateCategoryPage)
	http.HandleFunc("/view_category", handleGetCategoryPage)
	http.HandleFunc("/update_category", handleUpdateCategoryPage)
	http.HandleFunc("/delete_category", handleDeleteCategoryPage)
	http.HandleFunc("/profil", handleProfilPage)
	http.HandleFunc("/create_post", handleCreatePostPage)
	http.HandleFunc("/view_post", handleGetPostPage)
	http.HandleFunc("/update_post", handleUpdatePostPage)
	http.HandleFunc("/delete_post", handleDeletePostPage)
	http.HandleFunc("/create_comment", handleCreateCommentPage)
	http.HandleFunc("/view_comment", handleGetCommentPage)
	http.HandleFunc("/update_comment", handleUpdateCommentPage)
	http.HandleFunc("/delete_comment", handleDeleteCommentPage)
	http.HandleFunc("/view_category_posts", handleGetCategoryPostsPage)
	http.HandleFunc("/disconnect", disconnection)
	http.HandleFunc("/like", handleLike)
	http.HandleFunc("/unlike", handleUnlike)
	http.HandleFunc("/user_profil", handleUsersProfil)
	http.HandleFunc("/delete_user", delete_account)
	http.HandleFunc("/update_user", handle_modification_user)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	// Handling the various application routes

	log.Println("Listening on :8080...")
	// Starting the server on port 8080

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	// Starting the HTTP server.
}
