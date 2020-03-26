package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

type BsStateResource struct{}

type BsState struct {
	Destroyer struct {
		Name        string     `json:"_name"`
		Size        int        `json:"_size"`
		Placed      bool       `json:"_placed"`
		Placement   []int      `json:"_placement"`
		HitProfiles [][]string `json:"hitprofiles"`
	} `json:"destroyer"`
	Submarine struct {
		Name        string     `json:"_name"`
		Size        int        `json:"_size"`
		Placed      bool       `json:"_placed"`
		Placement   []int      `json:"_placement"`
		HitProfiles [][]string `json:"hitprofiles"`
	} `json:"submarine"`
	Cruiser struct {
		Name        string     `json:"_name"`
		Size        int        `json:"_size"`
		Placed      bool       `json:"_placed"`
		Placement   []int      `json:"_placement"`
		HitProfiles [][]string `json:"hitprofiles"`
	} `json:"cruiser"`
	Battleship struct {
		Name        string     `json:"_name"`
		Size        int        `json:"_size"`
		Placed      bool       `json:"_placed"`
		Placement   []int      `json:"_placement"`
		HitProfiles [][]string `json:"hitprofiles"`
	} `json:"battleship"`
	Carrier struct {
		Name        string     `json:"_name"`
		Size        int        `json:"_size"`
		Placed      bool       `json:"_placed"`
		Placement   []int      `json:"_placement"`
		HitProfiles [][]string `json:"hitprofiles"`
	} `json:"carrier"`
	Misses []string `json:"misses"`
}

const (
	modelsDir = "./models/"
)

// BsStateRouter manages all the routes related to the bsState
func (rs BsStateResource) Routes() chi.Router {
	r := chi.NewRouter()

	//r.Use(Authenticated)

	r.Get("/", rs.List)

	r.Route("/{filename}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Delete("/", rs.Delete)
		r.Post("/", rs.Post)
	})

	return r
}

func (rs BsStateResource) Post(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	bs := &BsState{}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}
	defer r.Body.Close()

	// Unmarshal body to ensure it fits the structure of a bs state
	if err := json.Unmarshal(bytes, bs); err != nil {
		log.Println(err)
		BADREQUEST(w)
		return
	}

	b, err := json.Marshal(bs)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	if err := ioutil.WriteFile(filepath.Join(modelsDir, filename), b, 0666); err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	CREATED(w)
}

func (rs BsStateResource) List(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(filepath.Dir(modelsDir))
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	fileList := []string{}
	for _, f := range files {
		if f.Name() != "." && f.Name() != ".." {
			fileList = append(fileList, f.Name())
		}
	}

	res := struct {
		Files []string `json:"files"`
	}{
		fileList,
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	w.Write(b)
	OK(w)
}

func (rs BsStateResource) Get(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "filename")

	target := filepath.Join(modelsDir, val)

	if s, err := os.Stat(target); os.IsNotExist(err) {
		fmt.Println(target)
		fmt.Println(s)
		NOTFOUND(w)
		return
	}

	d, err := ioutil.ReadFile(target)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	ContentHeaders(w)
	w.Write(d)
}

func (rs BsStateResource) Delete(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	if _, err := os.Stat(filepath.Join(modelsDir, filename)); os.IsNotExist(err) {
		log.Println(err)
		NOTFOUND(w)
		return
	}

	if err := os.Remove(filepath.Join(modelsDir, filename)); err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	NOCONTENT(w)
}
