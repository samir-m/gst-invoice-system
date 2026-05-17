package main

import (
	"net/http"
)

func (app *App) HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]interface{}{
		"Title": "Home",
		"Page":  "home",
	}

	err := app.Tmpl.ExecuteTemplate(w, "home", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
