package main

// http://localhost:8080/login

import (
	"Forum/server"
	"log"
)

func main() {
	log.Println("Starting server...")
	server.Start()
}

// Function that initializes the server.
