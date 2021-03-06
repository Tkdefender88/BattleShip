package routes

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Tkdefender88/BattleShip/battlestate"
	"github.com/Tkdefender88/BattleShip/repository"
	"github.com/Tkdefender88/BattleShip/solver"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	//SessionResource manages the session resource for responses to the /session endpoint
	SessionResource struct {
		Session string   `json:"session"` // A unique session ID to identify this game
		Roll    int      `json:"roll"`    // A random roll to determine who goes first
		Names   []string `json:"names"`   // The system name and player name of the opponent
		Epoch   int64    `json:"epoc"`    // The time the game started
		Latency int      `json:"latency"` // The delay between shots

		//repository to access battle models
		repo repository.ModelRepository

		// value to keep track if a session is active
		activeSesh bool
		// the strategy used to target enemy ships
		strategy *solver.Strategy
		// the url to send target requests to
		opponentURL string
		// the current state of the game
		bsState *battlestate.BsState
		// determines how the server responds to certain requests
		battlePhase bool
	}

	//SessionRequest is used for unmarshalling the post request body to the /session endpoint
	SessionRequest struct {
		// OpponentURL is the URL of the opponent that is requesting a match
		OpponentURL string `json:"opponentURL"`
		// Latency is the time to wait between sending requests to the opponents /target endpoint
		Latency int `json:"latency"`
	}
)

const (
	playerURL = "csdept16.cs.mtech.edu:30124"
)

var (
	EventBroker *Broker
)

func init() {
	EventBroker = NewServer()

	/*
		go func() {
			tiles := []int{34, 56, 42, 55, 33, 97}
			for i := 0; i < 6; i++ {
				time.Sleep(3 * time.Second)
				fe := FireEvent{
					Player: "player",
					Tile:   tiles[i],
					Hit:    tiles[i]%2 == 0,
				}
				b, _ := json.Marshal(&fe)
				EventBroker.Notifier <- b
			}
			for i := 0; i < 6; i++ {
				time.Sleep(3 * time.Second)
				fe := FireEvent{
					Player: "opponent",
					Tile:   tiles[i],
					Hit:    tiles[i]%2 == 0,
				}
				b, _ := json.Marshal(&fe)
				EventBroker.Notifier <- b
			}
		}()
	*/
}

// NewSession creates a new SessionResource object.
func NewSession(repo repository.ModelRepository) *SessionResource {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "csdept16"
	}
	return &SessionResource{
		Names:      []string{hostName, "Justin"},
		activeSesh: false,
		strategy:   solver.NewStrategy(),
		repo:       repo,
	}
}

// BattlePhase is middleware to block target and session requests if the server
// is not in phase 2, battle phase.
func (rs *SessionResource) BattlePhase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rs.battlePhase {
			respondError(w, http.StatusPreconditionFailed, "")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ActiveSessionCheck is a middleware to ensure that requests have a valid session
// if there is an active game occuring.
func (rs *SessionResource) ActiveSessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := &struct {
			Session string `json:"session"`
		}{}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error: %+v\n", err)
			respondError(w, http.StatusInternalServerError, "")
			return
		}
		r.Body.Close()

		err = json.Unmarshal(bodyBytes, s)
		if err != nil {
			respondError(w, http.StatusBadRequest, string(bodyBytes))
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if rs.activeSesh && s.Session != rs.Session {
			body, _ := json.Marshal(struct {
				Opponent []string `json:"opponent"`
			}{
				Opponent: rs.Names,
			})
			respondError(w, http.StatusForbidden, string(body))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Routes sets up all the routes for the /session endpoint
func (rs *SessionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(rs.ActiveSessionCheck).Delete("/{session-id}", rs.DeleteSession)
	r.With(rs.BattlePhase, rs.ActiveSessionCheck).Post("/", rs.PostSession)

	return r
}

// URLParam is a handler wrapper that parses out the optional url parameter
func (rs *SessionResource) BattleURL(h http.HandlerFunc) http.HandlerFunc {
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
	modelName := chi.URLParam(r, "filename")

	m, err := rs.repo.FindModel(modelName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp := struct {
				Filename string `json:"filename"`
			}{
				Filename: modelName,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rs.bsState = m

	if !rs.bsState.Valid() {
		respondError(w, http.StatusBadRequest, "Invalid game state selected")
		return
	}

	rs.battlePhase = true
	respondJSON(w, http.StatusOK, rs.bsState)
}

// DeleteSession tears down a session after a game is completed.
func (rs *SessionResource) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session := chi.URLParam(r, "session-id")

	if session != rs.Session {
		respondError(w, http.StatusBadRequest, session)
		return
	}

	rs.battlePhase = false

	resp := &struct {
		Session  string        `json:"session-id"`
		Duration time.Duration `json:"duration"`
	}{
		Session:  rs.Session,
		Duration: time.Since(int64ToTime(rs.Epoch)),
	}

	rs = NewSession(rs.repo)
	respondJSON(w, http.StatusOK, resp)
}

// StartSession sends a new post request to the opponents /session endpoint to
// try and establish a new game session between the two servers
func (rs *SessionResource) StartSession() {
	body, _ := json.Marshal(SessionRequest{
		OpponentURL: playerURL,
		Latency:     5000,
	})

	certPool, err := x509.SystemCertPool()
	if err != nil || certPool == nil {
		certPool = x509.NewCertPool()
	}

	pemkey, err := ioutil.ReadFile(pem)
	if err != nil {
		log.Printf("Error occured reading cert key: %+v\n", err)
		return
	}
	certificate, err := tls.LoadX509KeyPair(cert, pem)
	certPool.AppendCertsFromPEM(pemkey)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				Certificates:       []tls.Certificate{certificate},
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(http.MethodPost, "https://"+rs.opponentURL+"/bsProtocol/session", bytes.NewReader(body))
	if err != nil {
		log.Printf("Error: %+v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error %+v\n", err)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		log.Printf("Error %+v\n", err)
		return
	}

	if rs.Roll == 0 {
		go rs.Target()
	}
}

// PostSession handles a POST request /session and builds a new game session
// between the requester and the server.
func (rs *SessionResource) PostSession(w http.ResponseWriter, r *http.Request) {
	req := &SessionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, "")
		return
	}

	rs.Epoch = milliSecondsTime(time.Now())
	rs.Session = getMD5hash(playerURL + r.RemoteAddr + strconv.FormatInt(rs.Epoch, 10))

	if req.Latency <= 10000 && req.Latency >= 2000 {
		rs.Latency = req.Latency
	} else {
		rs.Latency = int(5000 * time.Millisecond)
	}
	rs.opponentURL = req.OpponentURL

	rand.Seed(time.Now().UTC().UnixNano())
	rs.Roll = rand.Intn(2)

	if rs.Roll == 1 {
		log.Println("Our turn first, firing shot")
		go rs.Target()
	}

	json.NewEncoder(w).Encode(&rs)
	rs.activeSesh = true
}

// Delete requests that the current session be terminated.
func (rs *SessionResource) Delete() {
	client := http.Client{}
	r, err := http.NewRequest(http.MethodDelete, rs.opponentURL+"/bsProtocol/session/"+rs.Session, nil)
	if err != nil {
		log.Printf("err %+v\n", err)
		return
	}

	resp, err := client.Do(r)
	if err != nil {
		log.Printf("err %+v\n", err)
		return
	}

	body := &struct {
		Session  string `json:"session"`
		Duration int64  `json:"duration"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		log.Printf("err %+v\n", err)
		return
	}
	rs.battlePhase = false
	log.Printf("The game lasted for %d ms\n", body.Duration)
}

func getMD5hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func milliSecondsTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func int64ToTime(t int64) time.Time {
	t *= int64(time.Millisecond)
	return time.Unix(0, t)
}
