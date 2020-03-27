package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseURL = "http://localhost:30124/battle/"

func TestStartBattleMode_NoURL_200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/stacky", nil)
	w := httptest.NewRecorder()

	router := BattleProtocol{}.Routes()
	router.ServeHTTP(w, req)

	filename := "./models/stacky"
	stacky, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Precondition failed, %s doesn't exist", filename)
		return
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

	router := BattleProtocol{}.Routes()
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

func TestStartBattleSession_SessionRequest_200(t *testing.T) {
	request := struct {
		OpponentURL string `json:"opponentURL"`
		Latency     int    `json:"latency"`
	}{
		OpponentURL: "https://csdept16.mtech.edu:30120",
		Latency:     2000,
	}

	reqBody, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	router := (&Session{}).Routes()

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
		t.Errorf("Test Failed, Bad response body: %+v Error: %+v", respBody, err)
		return
	}

	if respBody.Session == "" {
		t.Errorf("Test Failed, Session is empty %s", respBody.Session)
	}

	if respBody.Roll != 0 && respBody.Roll != 1 {
		t.Errorf("Test Failed, Roll invalid Got %d Want 0 or 1", respBody.Roll)
	}

	if respBody.Epoch == 0 {
		t.Errorf("Test Failed, Expected nonZero epoc")
	}

	if respBody.Latency != 2000 {
		t.Errorf("Test Failed, Bad Latency Got %d Want %d", resp.Body, 2000)
	}

	if battlePhase != true {
		t.Errorf("Test Failed, Battle state not set. Got %t Want %t", battlePhase, true)
	}
}
