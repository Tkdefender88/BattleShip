package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/solver"
)

func TestStartBattleSession_SessionRequest_200(t *testing.T) {
	request := struct {
		OpponentURL string `json:"opponentURL"`
		Latency     int    `json:"latency"`
	}{
		OpponentURL: "https://csdept16.mtech.edu:30120",
		Latency:     2000,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
		t.Fatalf("Failed %+v\n", err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	s := NewSession()
	s.activeSesh = false
	s.strategy = solver.NewStrategy()
	s.battlePhase = true

	router := s.Routes()

	router.ServeHTTP(w, req)

	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Cannot read response body: %+v", err)
		return
	}

	expectedCode := 200
	if resp.StatusCode != expectedCode {
		t.Errorf("Test Failed: Got %d Want %d", resp.StatusCode, expectedCode)
	}

	respBody := struct {
		Session string   `json:"session"`
		Roll    int      `json:"roll"`
		Names   []string `json:"names"`
		Epoch   int64    `json:"epoc"`
		Latency int      `json:"latency"`
	}{}

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		t.Errorf("Test Failed, Bad response body: %s Error: %+v", string(body), err)
		return
	}

	if respBody.Session == "" {
		t.Errorf("Test Failed, SessionResource is empty %s", respBody.Session)
	}

	if respBody.Roll != 0 && respBody.Roll != 1 {
		t.Errorf("Test Failed, Roll invalid Got %d Want 0 or 1", respBody.Roll)
	}

	hostName, err := os.Hostname()
	if err != nil {
		t.Fatalf("Error %+v", err)
		return
	}

	if respBody.Names[0] != hostName {
		t.Errorf("Test Failed, Bad Host name Got %s", respBody.Names[0])
	}

	if respBody.Names[1] != "Justin" {
		t.Errorf("Test Failed, Bad player name Got %s Want %s", respBody.Names[1], "Justin")
	}

	if respBody.Epoch == 0 {
		t.Errorf("Test Failed, Expected nonZero epoc")
	}

	if respBody.Latency != 2000 {
		t.Errorf("Test Failed, Bad Latency Got %d Want %d", resp.Body, 2000)
	}

	if s.battlePhase != true {
		t.Errorf("Test Failed, Battle state not set. Got %t Want %t", s.battlePhase, true)
	}
}

func TestPostTarget_ValidTarget_OpponentAccept(t *testing.T) {
	b, _ := json.Marshal(
		struct {
			Session string `json:"session"`
			Tile    string `json:"tile"`
		}{
			Session: "validsession",
			Tile:    "A5",
		})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	w := httptest.NewRecorder()

	// Make sure the sessions match
	s := &SessionResource{
		Session:     "validsession",
		battlePhase: true,
	}

	router := s.TargetRoute()
	router.ServeHTTP(w, req)

	resp := w.Result()

	respBody := &TargetResource{}

	if exp := 200; resp.StatusCode != exp {
		t.Errorf("Test Failed, Got %d Want %d", resp.StatusCode, exp)
	}

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		t.Fatalf("Test failed due to error: %+v", err)
		return
	}

	if exp := "A5"; exp != respBody.Tile {
		t.Errorf("Test Failed, Bad Response Tile Got %s Want %s", respBody.Tile, exp)
	}

	if exp := "INPROGRESS"; exp != respBody.Disposition {
		t.Errorf("Test Failed, Got %s Want %s", respBody.Disposition, exp)
	}
}
