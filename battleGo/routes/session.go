package routes

import (
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
	"time"

	"github.com/go-chi/chi"
)

type (
	TargetStatus string
)

const (
	Miss       TargetStatus = "MISS"
	Carrier                 = "CARRIER"
	BattleShip              = "BATTLESHIP"
	Cruiser                 = "CRUISER"
	Submarine               = "SUBMARINE"
	Destroyer               = "DESTROYER"
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
	}
	//SessionRequest is used for unmarshalling the post request body to the /session endpoint
	SessionRequest struct {
		// OpponentURL is the URL of the opponent that is requesting a match
		OpponentURL string `json:"opponentURL"`
		// Latency is the time to wait between sending requests to the opponents /target endpoint
		Latency int `json:"latency"`
	}

	TargetResource struct {
		// Status contains information of either a miss
		// or the name of the ship that was hit
		// CARRIER BATTLESHIP CRUISER SUBMARINE DESTROYER
		Status TargetStatus `json:"status"`
		// The tile that was hit, in Row Column format
		// Rows being letters from [A - J] and Columns
		// being numbers [0 - 9]
		Tile string `json:"tile"`
		// The Disposition of the game, either 'INPROGRESS' or 'WIN'
		Disposition string `json:"disposition"`
	}
)

func BattlePhase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !battlePhase {
			PRECONDITIONFAIL(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rs *SessionResource) ActiveSessionCheck(next http.Handler) http.Handler {
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

func (rs *SessionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(BattlePhase)
	r.Use(rs.ActiveSessionCheck)

	r.Delete("/session/{session-id}", rs.DeleteSession)
	r.Post("/session", rs.PostSession)
	r.Post("/target", rs.PostTarget)

	return r
}

func (rs *SessionResource) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session := chi.URLParam(r, "session-id")

	if session != rs.Session {
		BADREQUEST(w, []byte(session))
		return
	}

	battlePhase = false

	resp := &struct {
		Session  string        `json:"session-id"`
		Duration time.Duration `json:"duration"`
	}{
		Session:  rs.Session,
		Duration: time.Since(int64ToTime(rs.Epoch)),
	}

	rs = &SessionResource{}
	OKReader(w, resp)
}

func (rs *SessionResource) PostTarget(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Session string `json:"session-id"`
		Tile    string `json:"tile"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		BADREQUESTReader(w, r.Body)
		return
	}

	// For now, until I have a functional battleship algorithm
	// this will be the only check to return a 200
	if rs.Session != req.Session {
		fmt.Println(rs.Session, req.Session)
		UNAUTHORIZED(w)
		return
	}

	resp := TargetResource{
		Status:      Miss,
		Tile:        req.Tile,
		Disposition: "INPROGRESS",
	}

	OK(w)
	json.NewEncoder(w).Encode(resp)
}

func (rs *SessionResource) PostSession(w http.ResponseWriter, r *http.Request) {
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

	rs.Epoch = milliSecondsTime(time.Now())
	rs.Session = getMD5hash(local + r.RemoteAddr + strconv.FormatInt(rs.Epoch, 10))
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Error"
	}
	rs.Names = []string{hostName, "Justin"}

	if sessionReq.Latency <= 10000 && sessionReq.Latency >= 2000 {
		rs.Latency = sessionReq.Latency
	} else {
		rs.Latency = int(5000 * time.Millisecond)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rs.Roll = rand.Intn(2)

	json.NewEncoder(w).Encode(rs)
	rs.activeSesh = true
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
