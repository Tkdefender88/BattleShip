package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/battlestate"
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
		if rs.bsState.GameLost() {
			resp.Disposition = "WIN"
		} else {
			resp.Disposition = "INPROGRESS"
		}
	}

	//krs.UpdateClient()
	go rs.Target()

	ok(w)
	json.NewEncoder(w).Encode(resp)
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

	if resp.Status != battlestate.Miss {
		rs.strategy.ConfirmShot(resp.Tile, true)
	} else {
		rs.strategy.ConfirmShot(resp.Tile, false)
	}

	if resp.Disposition == "WIN" {
		go rs.Delete()
	}
}
func tileFromIndex(index int) string {
	row := rune((index / 10) + 65)
	col := rune((index % 10) + 48)
	return string([]rune{row, col})
}
