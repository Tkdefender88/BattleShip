package bsprotocol

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	errResp "github.com/Tkdefender88/BattleShip/battleGo/errorresponse"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type (
	// BattleProtocol is the struct responsible for managing the /battle endpoint
	BattleProtocol struct{}
)

// URLParam is a handler wrapper that parses out the optional url parameter
func (rs *SessionResource) URLParam(h http.HandlerFunc) http.HandlerFunc {
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
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	if err := json.NewDecoder(reader).Decode(&rs.bsState); err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	if !rs.bsState.Valid() {
		render.Render(w, r, errResp.ErrBadRequest(err, "Invalid game state selected"))
		return
	}

	rs.battlePhase = true
	render.Render(w, r, rs.bsState)
}
