package routes

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/battlestate"
)

const (
	pem  = "/etc/certs/wildcard_cs_mtech_edu.key"
	cert = "/etc/certs/wildcard_cs_mtech_edu.cer"
)

type (
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

	// FireEvent is an event that is sent to the client during the battle phase via
	// SSE. This gives the client the information needed to update the board.
	FireEvent struct {
		Player string `json:"player"`
		Tile   int    `json:"tile"`
		Hit    bool   `json:"hit"`
	}
)

// PostTarget checks if the target the opponent just specified is a hit or a miss
// responds with which ship was hit and if the game has eneded.
func (rs *SessionResource) PostTarget(w http.ResponseWriter, r *http.Request) {
	req := &TargetRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		badRequestReader(w, r.Body)
		return
	}

	if rs.Session != req.Session {
		fmt.Println(rs.Session, req.Session)
		unauthorized(w)
		return
	}

	hit, ship := rs.bsState.Hit(req.Tile)

	resp := &TargetResource{}

	if !hit {
		resp.Status = battlestate.Miss
		resp.Tile = req.Tile
		resp.Disposition = "INPROGRESS"
	} else {
		resp.Status = ship
		resp.Tile = req.Tile
		log.Printf("%+v\n", rs.bsState.Cruiser.HitProfiles)
		if rs.bsState.GameLost() {
			resp.Disposition = "WIN"
		} else {
			resp.Disposition = "INPROGRESS"
		}
	}

	index := indexFromTile(req.Tile)
	event := FireEvent{
		Player: "player",
		Tile:   index,
		Hit:    hit,
	}

	rs.UpdateClient(event)

	go rs.Target()

	ok(w)
	json.NewEncoder(w).Encode(resp)
}

// Target sends a target request out. Uses the strategy object to calculate the
// next shot and then confirmes if the shot was a hit or a miss from the response
func (rs *SessionResource) Target() {
	time.Sleep(time.Millisecond * time.Duration(rs.Latency))

	index := rs.strategy.Step()
	tile := tileFromIndex(index)

	body := &TargetRequest{
		Session: rs.Session,
		Tile:    tile,
	}

	b, _ := json.Marshal(body)
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

	r, err := client.Post("https://"+rs.opponentURL+"/bsProtocol/target", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println("err", err)
		return
	}
	defer r.Body.Close()

	resp := &TargetResource{}

	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		fmt.Printf("Error: %+v", err)
	}

	fmt.Printf("Body: %+v\n", resp)

	if resp.Status != battlestate.Miss {
		rs.strategy.ConfirmShot(resp.Tile, true)
		rs.bsState.HitEnemy(resp.Status, resp.Tile)

		ship := rs.bsState.ShipFromString(resp.Status)
		_, sunk := ship.Sunk()

		if sunk {
			rs.strategy.RemoveShip(strings.ToLower(resp.Status))
		}
	} else {
		rs.strategy.ConfirmShot(resp.Tile, false)
	}

	event := FireEvent{
		Player: "opponent",
		Tile:   index,
		Hit:    resp.Status != battlestate.Miss,
	}

	rs.UpdateClient(event)
	if resp.Disposition == "WIN" {
		go rs.Delete()
	}
}

// UpdateClient will send and SSE message to the client with any state changes
// from the battle
func (rs *SessionResource) UpdateClient(event FireEvent) {
	eventData, err := json.Marshal(&event)
	if err != nil {
		log.Printf("Error occured sending event message: %+v\n", err)
		return
	}
	EventBroker.Notifier <- eventData
}

func tileFromIndex(index int) string {
	row := rune((index / 10) + 65)
	col := rune((index % 10) + 48)
	return string([]rune{row, col})
}

func indexFromTile(tile string) (index int) {
	t := []rune(tile)
	row := int(t[0]) - 65
	col := int(t[1]) - 48
	index = (row * 10) + col
	return
}
