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

func (rs SessionResource) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	url := chi.URLParam(r, "url")
	if url != "" {

	}

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
		INTERNALERROR(w)
		return
	}

	_ = json.NewDecoder(reader).Decode(rs.bsState)
	rs.battlePhase = true
	OKReader(w, rs.bsState)
}
