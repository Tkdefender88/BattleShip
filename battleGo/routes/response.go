package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type ErrorMsg struct {
	Message string
	Status  int
}

func respondJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respondJSON(w, code, map[string]string{"message": msg})
}

// ContentHeaders is a helper function that sets the 'content type' and the 'access
// control allow origin' headers.
func ContentHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func unauthorizedPage(w http.ResponseWriter, msg ErrorMsg) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)

	msg.Status = http.StatusUnauthorized

	tmpl := template.Must(template.ParseFiles("views/base.html", "views/error.html"))
	if err := tmpl.ExecuteTemplate(w, "base.html", msg); err != nil {
		log.Println(err)
	}
}
