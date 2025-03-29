package api

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/frasmataz/go-chess/conf"
	"github.com/gorilla/mux"
)

var tmpl *template.Template

func StartServer(ctx context.Context) error {

	tmpl, _ = template.ParseGlob("templates/*.html")
	cfg := conf.DefaultConfig()

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: NewRouter(),
	}

	go func() {
		log.Println("Starting server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server")
	return server.Shutdown(ctx)

}

func NewRouter() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/game", GetGameState).Methods("GET")
	return router

}
