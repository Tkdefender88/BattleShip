package battlestate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	errResp "github.com/Tkdefender88/BattleShip/battleGo/errorresponse"
	"github.com/Tkdefender88/BattleShip/battleGo/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	modelsDir = "./models/"
)

// BsStateResource is responsible for all the routes to /bsState
type BsStateResource struct{}

type stateListResource struct {
	Files []string `json:"files"`
}

func (s *stateListResource) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Routes manages all the routes related to the bsState
func (rs BsStateResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(routes.Refresh)
	r.Use(routes.Authenticated)

	r.Get("/", rs.ListStates)

	r.Route("/{filename}", func(r chi.Router) {
		r.Get("/", rs.GetState)
		r.Delete("/", rs.DeleteState)
		r.Post("/", rs.SaveState)
	})

	return r
}

// SaveState handles post reqests to the bsState endpoint. Accepts a bsState object
// as the body and will save it to the file system
func (rs BsStateResource) SaveState(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	bs := &BsState{}

	err := json.NewDecoder(r.Body).Decode(bs)
	if err != nil {
		render.Render(w, r, errResp.ErrBadRequest(err, nil))
		return
	}

	defer r.Body.Close()

	b, err := json.Marshal(bs)
	if err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	if err := ioutil.WriteFile(filepath.Join(modelsDir, filename), b, 0666); err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusCreated)
}

// ListStates will respond with a list of the battlestates currently stored on the
// filesystem
func (rs BsStateResource) ListStates(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(filepath.Dir(modelsDir))
	if err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	fileList := []string{}
	for _, f := range files {
		if f.Name() != "." && f.Name() != ".." {
			fileList = append(fileList, f.Name())
		}
	}

	res := &stateListResource{
		fileList,
	}

	render.Render(w, r, res)
}

// GetState will respond with the requested battlestate
func (rs BsStateResource) GetState(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "filename")

	target := filepath.Join(modelsDir, val)

	if s, err := os.Stat(target); os.IsNotExist(err) {
		fmt.Println(target)
		fmt.Println(s)
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	file, err := os.Open(target)
	if err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	resp := &BsState{}
	err = json.NewDecoder(file).Decode(resp)
	if err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}
	render.Render(w, r, resp)
}

// DeleteState will remove a battlestate from the filesystem
func (rs BsStateResource) DeleteState(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	if _, err := os.Stat(filepath.Join(modelsDir, filename)); os.IsNotExist(err) {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	if err := os.Remove(filepath.Join(modelsDir, filename)); err != nil {
		render.Render(w, r, errResp.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusNoContent)
	render.Render(w, r, nil)
}
