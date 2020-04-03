package routes

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/BattleState"
	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/solver"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi"
)

type (
	//SessionResource manages the session resource for responses to the /session endpoint
	SessionResource struct {
		Session string   `json:"session"` // A unique session ID to identify this game
		Roll    int      `json:"roll"`    // A random roll to determine who goes first
		Names   []string `json:"names"`   // The system name and player name of the opponent
		Epoch   int64    `json:"epoc"`    // The time the game started
		Latency int      `json:"latency"` // The delay between shots

		// value to keep track if a session is active
		activeSesh bool
		// the strategy used to target enemy ships
		strategy *solver.Strategy
		// the url to send target requests to
		opponentURL string
		// the current state of the game
		bsState BattleState.BsState
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

	// TargetResource represents the response to a request to the /target endpoint
	TargetResource struct {
		// Status contains information of either a miss
		// or the name of the ship that was hit
		// CARRIER BATTLESHIP CRUISER SUBMARINE DESTROYER
		Status string `json:"status"`
		// The tile that was hit, in Row Column format
		// Rows being letters from [A - J] and Columns
		// being numbers [0 - 9]
		Tile string `json:"tile"`
		// The Disposition of the game, either 'INPROGRESS' or 'WIN'
		Disposition string `json:"disposition"`
	}

	// TargetRequest represents the body of a request sent to the /target endpoint
	// used during the battle phase when players are firing at eachothers ships.
	TargetRequest struct {
		Session string `json:"session"`
		Tile    string `json:"tile"`
	}
)

var (
	eventServer *sse.Server
)

// NewSession createss a new SessionResource object.
func NewSession() *SessionResource {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "csdept16"
	}
	return &SessionResource{
		Names:      []string{hostName, "Justin"},
		activeSesh: false,
		strategy:   solver.NewStrategy(),
	}
}

// EventServer sets the eventServer object and returns it
func EventServer() *sse.Server {
	eventServer = sse.NewServer(nil)
	return eventServer
}

// BattlePhase is middleware to block target and session requests if the server
// is not in phase 2, battle phase.
func (rs *SessionResource) BattlePhase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rs.battlePhase {
			PRECONDITIONFAIL(w)
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
			INTERNALERROR(w)
			return
		}
		r.Body.Close()

		err = json.Unmarshal(bodyBytes, s)
		if err != nil {
			BADREQUESTReader(w, ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if rs.activeSesh && s.Session != rs.Session {
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

// Routes sets up all the routes for the /session endpoint
func (rs *SessionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.With(rs.ActiveSessionCheck).Delete("/{session-id}", rs.DeleteSession)
	r.With(rs.BattlePhase, rs.ActiveSessionCheck).Post("/", rs.PostSession)

	return r
}

// BattleRoute sets up the endpoints for the /battle endpoint
func (rs *SessionResource) BattleRoute() chi.Router {
	r := chi.NewRouter()
	r.Get("/{filename}", rs.Get)
	r.Get("/{filename}/{url}", rs.GetURL)
	return r
}

// TargetRoute sets up the endpoints to the /target endpoint
func (rs *SessionResource) TargetRoute() chi.Router {
	r := chi.NewRouter()
	r.With(rs.BattlePhase, rs.ActiveSessionCheck).Post("/", rs.PostTarget)
	return r
}

// DeleteSession tears down a session after a game is completed.
func (rs *SessionResource) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session := chi.URLParam(r, "session-id")

	if session != rs.Session {
		BADREQUEST(w, session)
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

	rs = NewSession()
	OKReader(w, resp)
}

// PostTarget checks if the target the opponent just specified is a hit or a miss
// responds with which ship was hit and if the game has eneded.
func (rs *SessionResource) PostTarget(w http.ResponseWriter, r *http.Request) {
	req := &TargetRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		BADREQUESTReader(w, r.Body)
		return
	}

	if rs.Session != req.Session {
		fmt.Println(rs.Session, req.Session)
		UNAUTHORIZED(w)
		return
	}

	hit, ship := rs.bsState.Hit(req.Tile)

	resp := &TargetResource{}

	if !hit {
		resp.Status = BattleState.Miss
		resp.Tile = req.Tile
		resp.Disposition = "INPROGRESS"
		rs.bsState.Misses = append(rs.bsState.Misses, req.Tile)
	} else {
		resp.Status = ship
		resp.Tile = req.Tile
	}

	rs.UpdateClient()
	go rs.Target()

	OK(w)
	json.NewEncoder(w).Encode(resp)
}

// StartSession sends a new post request to the opponents /session endpoint to
// try and establish a new game session between the two servers
func (rs *SessionResource) StartSession() {
	client := http.Client{}

	body, _ := json.Marshal(SessionRequest{
		OpponentURL: "https://csdept16.mtech.edu:30124",
		Latency:     5000,
	})

	req, err := http.NewRequest(http.MethodPost, "https://"+rs.opponentURL+"/session", bytes.NewReader(body))
	if err != nil {
		log.Printf("Error %+v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error %+v\n", err)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(rs)
	if err != nil {
		log.Printf("Error %+v\n", err)
		return
	}
}

// UpdateClient will send and SSE message to the client with any state changes
// from the battle
func (rs *SessionResource) UpdateClient() {
	w := &strings.Builder{}
	err := json.NewEncoder(w).Encode(&rs.bsState)
	if err != nil {
		log.Printf("Error occured sending event message: %+v\n", err)
		return
	}
	eventServer.SendMessage("/events/updates", sse.SimpleMessage(w.String()))
}

// PostSession handles a POST request /session and builds a new game session
// between the requester and the server.
func (rs *SessionResource) PostSession(w http.ResponseWriter, r *http.Request) {
	req := &SessionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		log.Println(err)
		BADREQUESTReader(w, r.Body)
		return
	}

	local := "https://csdept16.mtech.edu:30124"

	rs.Epoch = milliSecondsTime(time.Now())
	rs.Session = getMD5hash(local + r.RemoteAddr + strconv.FormatInt(rs.Epoch, 10))

	if req.Latency <= 10000 && req.Latency >= 2000 {
		rs.Latency = req.Latency
	} else {
		rs.Latency = int(5000 * time.Millisecond)
	}
	rs.opponentURL = req.OpponentURL

	rand.Seed(time.Now().UTC().UnixNano())
	rs.Roll = rand.Intn(2)

	if rs.Roll == 0 {
		log.Println("Our turn first, firing shot")
		go rs.Target()
	}

	json.NewEncoder(w).Encode(rs)
	rs.activeSesh = true
}

// Target sends a target request out. Uses the strategy object to calculate the
// next shot and then confirmes if the shot was a hit or a miss from the response
func (rs *SessionResource) Target() {
	time.Sleep(time.Millisecond * time.Duration(rs.Latency))

	index := rs.strategy.FireNext()

	body := &TargetRequest{
		Session: rs.Session,
		Tile:    tileFromIndex(index),
	}

	b, _ := json.Marshal(body)
	r, err := http.Post("https://"+rs.opponentURL+"/target", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println("err", err)
		return
	}
	defer r.Body.Close()

	resp := &TargetResource{}

	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		fmt.Printf("Error: %+v", err)
	}

	if resp.Status != BattleState.Miss {
		rs.strategy.ConfirmShot(resp.Tile, true)
	} else {
		rs.strategy.ConfirmShot(resp.Tile, false)
	}

	if resp.Disposition == "WIN" {
		go rs.Delete()
	}
}

// Delete requests that the current session be terminated.
func (rs *SessionResource) Delete() {
	client := http.Client{}
	r, err := http.NewRequest(http.MethodDelete, rs.opponentURL+"/session/"+rs.Session, nil)
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

func tileFromIndex(index int) string {
	row := rune((index / 10) + 65)
	col := rune((index % 10) + 48)
	return string([]rune{row, col})
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
