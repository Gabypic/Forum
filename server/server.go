package server

// http://localhost:8080/start
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

	http.HandleFunc("/start", handleLoginPage)

	http.HandleFunc("/create_account", handleRegisterPage)

	http.HandleFunc("/home", handleHomePage)

	log.Println("Listening on :8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
