package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

var (
	battlePhase = false
)

type (
	// BattleProtocol is the struct responsible for managing the /battle endpoint
	BattleProtocol struct{}
)

func (rs BattleProtocol) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{filename}", func(r chi.Router) {
		r.Get("/", rs.Get)
		//r.Get("/{url}", rs.GetURL)
	})

	return r
}

func (rs BattleProtocol) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	target := filepath.Join("./models", filename)

	if _, err := os.Stat(target); os.IsNotExist(err) {
		resp := struct {
			Filename string `json:"filename"`
		}{
			Filename: filename,
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	body, err := ioutil.ReadFile(target)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}
	battlePhase = true
	OK(w)
	w.Write(body)
}
