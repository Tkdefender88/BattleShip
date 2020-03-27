package routes

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type BattleProtocol struct{}

type Session struct {
	OpponentURL string `json:"opponentURL"`
	Latency     int    `json:"latency"`
}

type SessionResource struct {
	Session string   `json:"session"`
	Roll    int      `json:"roll"`
	Names   []string `json:"names"`
	Epoch   int64    `json:"epoc"`
	Latency int      `json:"latency"`
}

func (rs BattleProtocol) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/{filename}", func(r chi.Router) {
		r.Get("/", rs.Get)
		//r.Get("/{url}", rs.GetURL)
		//r.Delete("/", rs.Delete)
		//r.Post("/", rs.Post)
	})

	return r
}

func (rs Session) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Post)

	return r
}

func (rs Session) Post(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	if err := json.Unmarshal(reqBody, &rs); err != nil {
		log.Println(err)
		BADREQUEST(w, reqBody)
		return
	}
	local := "https://csdept16.mtech.edu:30124"

	sesh := SessionResource{}
	sesh.Epoch = time.Now().UnixNano() / int64(time.Millisecond)
	sesh.Session = GetMD5Hash(local + r.RemoteAddr + strconv.Itoa(int(sesh.Epoch)))
	sesh.Names = []string{"Justin's Server", "Opponent Server"}

	if rs.Latency <= 10000 && rs.Latency >= 2000 {
		sesh.Latency = rs.Latency
	} else {
		sesh.Latency = int(5000 * time.Millisecond)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	sesh.Roll = rand.Intn(2)

	respBody, err := json.Marshal(sesh)
	if err != nil {
		INTERNALERROR(w)
		return
	}

	w.Write(respBody)
}

func (rs BattleProtocol) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	target := filepath.Join("../models", filename)

	if _, err := os.Stat(target); os.IsNotExist(err) {
		resp := struct {
			Filename string `json:"filename"`
		}{
			Filename: filename,
		}

		body, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusNotFound)
		w.Write(body)
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

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
