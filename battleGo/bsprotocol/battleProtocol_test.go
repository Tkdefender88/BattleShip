package bsprotocol

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tkdefender88/BattleShip/battleGo/battlestate"
	"github.com/go-chi/chi"
)

var baseURL = "http://localhost:30124/battle/"

var testStacky = &battlestate.BsState{
	Destroyer: &battlestate.Ship{
		Name:        "destroyer",
		Size:        2,
		Placed:      true,
		Placement:   []int{3, 3, 1},
		HitProfiles: [][]string{},
	},
	Carrier: &battlestate.Ship{
		Name:        "carrier",
		Size:        5,
		Placed:      true,
		Placement:   []int{0, 0, 0},
		HitProfiles: [][]string{},
	},
	Battleship: &battlestate.Ship{
		Name:        "battleship",
		Size:        4,
		Placed:      true,
		Placement:   []int{1, 0, 0},
		HitProfiles: [][]string{},
	},
	Cruiser: &battlestate.Ship{
		Name:        "cruiser",
		Size:        3,
		Placed:      true,
		Placement:   []int{2, 0, 0},
		HitProfiles: [][]string{},
	},
	Submarine: &battlestate.Ship{
		Name:        "submarine",
		Size:        3,
		Placed:      true,
		Placement:   []int{3, 0, 0},
		HitProfiles: [][]string{},
	},
}

func TestStartBattleMode_NoURL_200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/stacky", nil)
	w := httptest.NewRecorder()

	session := NewSession()
	router := chi.NewRouter()
	router.Get("/{filename}", session.Get)

	router.ServeHTTP(w, req)

	stacky, err := json.Marshal(&testStacky)
	if err != nil {
		t.Fatalf("error occurred: %+v\n", err)
	}

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Cannot read response body: %+v", err)
		return
	}

	if resp.StatusCode != 200 {
		t.Errorf("Test Failed: Got %d Wanted %d", resp.StatusCode, 200)
	}

	// Compare the body of the response with the expected
	if bytes.Compare(stacky, body) != 0 {
		t.Errorf("Test Failed: Got %s Wanted %s", string(body), string(stacky))
	}
}

func TestStartBattleMode_NotFound_404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fooboi", nil)
	w := httptest.NewRecorder()

	s := NewSession()

	router := chi.NewRouter()
	router.Get("/{filename}", s.Get)
	router.ServeHTTP(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Cannot read response body: %+v", err)
		return
	}

	if resp.StatusCode != 404 {
		t.Errorf("Test Failed: Got %d Wanted %d", resp.StatusCode, 404)
	}

	expected := []byte{123, 34, 102, 105, 108, 101, 110, 97, 109, 101, 34, 58, 34, 102, 111, 111, 98, 111, 105, 34, 125, 10}
	if string(expected) != string(body) {
		t.Errorf("Test Failed: Got %s Wanted %s", string(body), string(expected))
	}
}
