package server

import (
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

func handleStartPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering start template...")
	renderTemplate(w, "templates", nil)
}
