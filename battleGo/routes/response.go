package routes

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type ErrorMsg struct {
	Message string
	Status  int
}

// ContentHeaders is a helper function that sets the 'content type' and the 'access
// control allow origin' headers.
func ContentHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func okReader(w http.ResponseWriter, body interface{}) {
	ok(w)
	json.NewEncoder(w).Encode(body)
}

// OK sets the standard headers and writes status OK to the header.
func ok(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusOK)
}

// CREATED sets the standard headers and writes status CREATED to the header.
func created(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusCreated)
}

// NOCONTENT sets the standard headers and writes status PUT to the header.
func noContent(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

func notFound(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusNotFound)
}

func internalError(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusInternalServerError)
}

func badRequestReader(w http.ResponseWriter, body io.Reader) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = io.Copy(w, body)
}

func badRequest(w http.ResponseWriter, body string) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintln(w, body)
}

func forbidden(w http.ResponseWriter, body []byte) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusForbidden)
	w.Write(body)
}

func preconditionFail(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusPreconditionFailed)
	w.Write([]byte{})
}

func unauthorized(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)
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
