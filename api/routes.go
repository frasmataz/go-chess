package api

import "github.com/gorilla/mux"

func SetupRoutes() {

	router := mux.NewRouter()

	router.HandleFunc("/", GetIndex).Methods("GET")
	router.HandleFunc("/game", GetGameState).Methods("GET")
}
