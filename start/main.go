package main

// http://localhost:8080/start

import (
	"Forum/server"
	"log"
)

func main() {
	log.Println("Starting server...")
	server.Start()
}
