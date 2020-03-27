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

var (
	battlePhase = false
)

type (
	// BattleProtocol is the struct responsible for managing the /battle endpoint
	BattleProtocol struct{}

	//Session manages the session resource for responses to the /session endpoint
	Session struct {
		Session    string   `json:"session"`
		Roll       int      `json:"roll"`
		Names      []string `json:"names"`
		Epoch      int64    `json:"epoc"`
		Latency    int      `json:"latency"`
		activeSesh bool
	}

	//SessionRequest is used for unmarshalling the post request body to the /session endpoint
	SessionRequest struct {
		// OpponentURL is the URL of the opponent that is requesting a match
		OpponentURL string `json:"opponentURL"`
		// Latency is the time to wait between sending requests to the opponents /target endpoint
		Latency int `json:"latency"`
	}
)

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

func (rs *Session) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(BattlePhase)
	r.Use(rs.ActiveSessionCheck)

	r.Post("/", rs.Post)

	return r
}

func BattlePhase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !battlePhase {
			PRECONDITIONFAIL(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rs *Session) ActiveSessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rs.activeSesh {
			body, _ := json.Marshal(struct {
				Opponent []string `json:"opponent"`
			}{
				Opponent: rs.Names,
			})
			FORBIDDEN(w, body)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rs *Session) Post(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		INTERNALERROR(w)
		return
	}

	sessionReq := SessionRequest{}

	if err := json.Unmarshal(reqBody, &sessionReq); err != nil {
		log.Println(err)
		BADREQUEST(w, reqBody)
		return
	}
	local := "https://csdept16.mtech.edu:30124"

	rs.Epoch = time.Now().UnixNano() / int64(time.Millisecond)
	rs.Session = GetMD5Hash(local + r.RemoteAddr + strconv.Itoa(int(rs.Epoch)))
	rs.Names = []string{"Justin's Server", "Opponent Server"}

	if sessionReq.Latency <= 10000 && sessionReq.Latency >= 2000 {
		rs.Latency = sessionReq.Latency
	} else {
		rs.Latency = int(5000 * time.Millisecond)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rs.Roll = rand.Intn(2)

	respBody, err := json.Marshal(rs)
	if err != nil {
		INTERNALERROR(w)
		return
	}

	rs.activeSesh = true

	_, _ = w.Write(respBody)
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

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
