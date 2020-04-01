package routes

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var baseURL = "http://localhost:30124/battle/"

func TestStartBattleMode_NoURL_200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/stacky", nil)
	w := httptest.NewRecorder()

	router := (&SessionResource{}).Routes()
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

	router := (&SessionResource{}).Routes()
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
