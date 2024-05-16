package server

// http://localhost:8080/start
import (
	"log"
	"net/http"
)

// Start function initializes and starts the web server
func Start() {
	// Log a message indicating the server is starting
	log.Println("Starting server...")

	// Handle static files (e.g., CSS, images) by serving them from the /static/ path
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Define handlers for different HTTP routes
	http.HandleFunc("/start", handleStartPage)

	// Log a message indicating the server is listening on port :8080
	log.Println("Listening on :8080...")

	// Start the HTTP server on port :8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// Log a fatal error if the server fails to start
		log.Fatal(err)
	}
}
