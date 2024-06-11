package server

import (
	"Forum/application"
	"log"
	"net/http"
)

func Start() {
	cfg := application.LoadConfig()
	application.Connect(cfg.DatabaseURL)
	application.SqlTable()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	log.Println("Listening on :8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
