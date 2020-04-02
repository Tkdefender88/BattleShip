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

// This is a helper function that sets the 'content type' and the 'access
// control allow origin' headers.
func ContentHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func OKReader(w http.ResponseWriter, body interface{}) {
	OK(w)
	json.NewEncoder(w).Encode(body)
}

// OK sets the standard headers and writes status OK to the header.
func OK(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusOK)
}

// CREATED sets the standard headers and writes status CREATED to the header.
func CREATED(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusCreated)
}

// NOCONTENT sets the standard headers and writes status PUT to the header.
func NOCONTENT(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

func NOTFOUND(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusNotFound)
}

func INTERNALERROR(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusInternalServerError)
}

func BADREQUESTReader(w http.ResponseWriter, body io.Reader) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = io.Copy(w, body)
}

func BADREQUEST(w http.ResponseWriter, body string) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintln(w, body)
}

func FORBIDDEN(w http.ResponseWriter, body []byte) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusForbidden)
	w.Write(body)
}

func PRECONDITIONFAIL(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusPreconditionFailed)
	w.Write([]byte{})
}

func UNAUTHORIZED(w http.ResponseWriter) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)
}

func UNAUTHORIZEDPage(w http.ResponseWriter, msg ErrorMsg) {
	ContentHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)

	msg.Status = http.StatusUnauthorized

	tmpl := template.Must(template.ParseFiles("views/base.html", "views/error.html"))
	if err := tmpl.ExecuteTemplate(w, "base.html", msg); err != nil {
		log.Println(err)
	}
}
