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
		// event server
		events  *sse.Server
		eventID int
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
		Status string `json:"status"`
		// The tile that was hit, in Row Column format
		// Rows being letters from [A - J] and Columns
		// being numbers [0 - 9]
		Tile string `json:"tile"`
		// The Disposition of the game, either 'INPROGRESS' or 'WIN'
		Disposition string `json:"disposition"`
	}

	TargetRequest struct {
		Session string `json:"session"`
		Tile    string `json:"tile"`
	}
)

func NewSession() (*SessionResource, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &SessionResource{
		Names:      []string{hostName, "Justin"},
		activeSesh: false,
		strategy:   solver.NewStrategy(),
	}, nil
}

func (rs *SessionResource) RegisterEventSource(events *sse.Server) {
	rs.events = events
}

func (rs *SessionResource) BattlePhase(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rs.battlePhase {
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

	r.With(rs.ActiveSessionCheck).Delete("/{session-id}", rs.DeleteSession)
	r.With(rs.BattlePhase, rs.ActiveSessionCheck).Post("/", rs.PostSession)

	return r
}

func (rs *SessionResource) BattleRoute() chi.Router {
	r := chi.NewRouter()
	r.Get("/{filename}", rs.Get)
	r.Get("/{filename}/{url}", rs.GetURL)
	return r
}

func (rs *SessionResource) TargetRoute() chi.Router {
	r := chi.NewRouter()
	r.With(rs.BattlePhase, rs.ActiveSessionCheck).Post("/", rs.PostTarget)
	return r
}

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

	rs = &SessionResource{}
	OKReader(w, resp)
}

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

	OK(w)
	json.NewEncoder(w).Encode(resp)
}

<<<<<<< Updated upstream
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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error %+v\n", err)
		return
	}

	if err := json.Unmarshal(b, rs); err != nil {
		log.Printf("Error %+v\n", err)
		return
	}
=======
func (rs *SessionResource) UpdateClient() {
	w := &strings.Builder{}
	err := json.NewEncoder(w).Encode(&rs.bsState)
	if err != nil {
		log.Printf("Error occured sending event message: %+v\n", err)
		return
	}
	rs.events.SendMessage("/events/updates", sse.SimpleMessage(w.String()))
>>>>>>> Stashed changes
}

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
