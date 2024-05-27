package server

import (
	"Forum/application"
	"fmt"
	"html/template"
	"log"
	"net/http"
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

func handleHomePageConnection(w http.ResponseWriter, r *http.Request) {
	LoginUserHandler(w, r)
}

func handleHomePageRegister(w http.ResponseWriter, r *http.Request) {
	RegisterUserHandler(w, r)
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

	err := application.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "home", nil)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	//	return
	//}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := application.GetUser(email)
	fmt.Print("caca")
	fmt.Println(user)
	if err != nil || user == nil {
		fmt.Println("1")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !application.CheckPassword(password, user.Password) {
		fmt.Println("2")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	fmt.Println("3")
	renderTemplate(w, "home", nil)
}
