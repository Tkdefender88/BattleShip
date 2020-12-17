package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Tkdefender88/BattleShip/battlestate"
	"github.com/Tkdefender88/BattleShip/repository"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

// BsStateResource is responsible for all the routes to /bsState
type BsStateResource struct {
	repo repository.ModelRepository
}

type Controller interface {
	Routes() chi.Router
}

func NewBsStateController(repo repository.ModelRepository) *BsStateResource {
	return &BsStateResource{
		repo: repo,
	}
}

// Routes manages all the routes related to the bsState
func (rs BsStateResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)

	r.Route("/{filename}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Delete("/", rs.Delete)
		r.Post("/", rs.Post)
	})

	return r
}

// Post handles post reqests to the bsState endpoint. Accepts a bsState object
// as the body and will save it to the file system
func (rs BsStateResource) Post(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if len(filename) == 0 {
		respondError(w, http.StatusBadRequest, "No model name given")
		return
	}
	bs := &battlestate.BsState{}

	err := json.NewDecoder(r.Body).Decode(bs)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}
	defer r.Body.Close()

	if !bs.Valid() {
		respondError(w, http.StatusBadRequest, "Invalid model")
		return
	}

	id, err := rs.repo.CreateModel(filename, bs)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, "")
		return
	}
	respondJSON(w, http.StatusCreated, map[string]string{"id": id.String()})
}

// List will respond with a list of the battlestates currently stored on the
// filesystem
func (rs BsStateResource) List(w http.ResponseWriter, r *http.Request) {

	fileList, err := rs.repo.ListModels()
	if err != nil {
		return
	}

	res := struct {
		Files []string `json:"files"`
	}{
		fileList,
	}

	respondJSON(w, http.StatusOK, res)
}

// Get will respond with the requested battlestate
func (rs BsStateResource) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	// get the resource from the repository
	bsState, err := rs.repo.FindModel(filename)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(filename)
			respondError(w, http.StatusNotFound, "Document "+filename+" Not Found")
			return
		}
		respondError(w, http.StatusInternalServerError, "")
		return
	}

	// return the resouce to the requester
	respondJSON(w, http.StatusOK, bsState)
}

// Delete will remove a battlestate from the filesystem
func (rs BsStateResource) Delete(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "filename")
	w.WriteHeader(http.StatusNoContent)
}
