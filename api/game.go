package api

import (
	"net/http"

	"github.com/frasmataz/go-chess/model"
)

func GetGameState(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "game", model.Game{})
}
