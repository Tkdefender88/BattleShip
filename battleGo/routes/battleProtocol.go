package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

type (
	// BattleProtocol is the struct responsible for managing the /battle endpoint
	BattleProtocol struct{}
)

func (rs *SessionResource) UrlParam(h http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, r *http.Request) {

		url := chi.URLParam(r, "url")
		if url != "" {
			rs.opponentURL = url
			go rs.StartSession()
		}

		h.ServeHTTP(response, r)
	}
}

// Get handles the /battle route where the optional parameter for the URL
// is not included
func (rs *SessionResource) Get(w http.ResponseWriter, r *http.Request) {
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

	reader, err := os.Open(target)
	if err != nil {
		log.Println(err)
		internalError(w)
		return
	}

	if err := json.NewDecoder(reader).Decode(&rs.bsState); err != nil {
		log.Println(err)
		internalError(w)
		return
	}

	log.Printf("%+v\n", rs.bsState.Carrier)

	if !rs.bsState.Valid() {
		badRequest(w, "Invalid game state selected")
		return
	}

	rs.battlePhase = true
	okReader(w, rs.bsState)
}
