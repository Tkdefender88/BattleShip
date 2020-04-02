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

func (rs *SessionResource) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if err := rs.readBattleState(w, filename); err != nil {
		return
	}

	url := chi.URLParam(r, "url")
	if url != "" {
		rs.opponentURL = url
		go rs.StartSession()
	}

	rs.battlePhase = true
	OKReader(w, rs.bsState)
}

func (rs *SessionResource) GetURL(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if err := rs.readBattleState(w, filename); err != nil {
		return
	}
	rs.battlePhase = true
	OKReader(w, rs.bsState)
}

func (rs *SessionResource) readBattleState(w http.ResponseWriter, filename string) error {
	target := filepath.Join("./models", filename)

	if _, err := os.Stat(target); os.IsNotExist(err) {
		resp := struct {
			Filename string `json:"filename"`
		}{
			Filename: filename,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return err
	}

	reader, err := os.Open(target)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return err
	}

	if err := json.NewDecoder(reader).Decode(&rs.bsState); err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return err
	}
	if !rs.bsState.Valid() {
		BADREQUEST(w, "Invalid game state selected")
		return err
	}
	return nil
}
