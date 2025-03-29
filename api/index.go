package api

import (
	"log"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {

	log.Printf("GET /")

	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("error: GET /: %v", err)
	}

}
