package routes_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tkdefender88/BattleShip/battlestate"
	"github.com/Tkdefender88/BattleShip/repository"
	"github.com/Tkdefender88/BattleShip/routes"
	"github.com/go-chi/chi"
)

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
	Misses: []string{},
}

func Test_GetBattle(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type fields struct {
		repo repository.ModelRepository
	}
	type expect struct {
		status int
		body   []byte
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		expect expect
	}{
		{
			name: "happy path",
			fields: fields{
				repo: &mockRepo{},
			},
			args: args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, "/stacky", nil),
			},
			expect: expect{
				status: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sesh := routes.NewSession(tt.fields.repo)

			r := chi.NewRouter().Route("/battle/{filename}", func(r chi.Router) {
				r.Get("/", sesh.Get)
				r.Get("/{url}", sesh.BattleURL(sesh.Get))
			})

			r.ServeHTTP(tt.args.w, tt.args.r)

			resp := tt.args.w.Result()

			if resp.StatusCode != tt.expect.status {
				t.Errorf("Incorrect response code expected %d got %d", tt.expect.status, resp.StatusCode)
			}
		})
	}
}

func Test_PostSession(t *testing.T) {

	sessionReq := routes.SessionRequest{
		OpponentURL: "https://csdept16.cs.mtech.edu/:30214",
		Latency:     1000,
	}
	happyPath := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(happyPath).Encode(sessionReq)
	if err != nil {
		t.Fatal(err)
	}

	badReq := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(badReq).Encode(struct{ foo float64}{ 420.69 })

	tests := []testCase{
		{
			name: "happy_path",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/", happyPath),
				w: httptest.NewRecorder(),
			},
			fields: fields{
				repo: &mockRepo{},
			},
			expect: expect{
				status: http.StatusOK,
				size: 126,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			sesh := routes.NewSession(tt.fields.repo)
			r := chi.NewRouter().Route("/session", func(r chi.Router) {
				r.Post("/", sesh.PostSession)
			})

			// act
			r.ServeHTTP(tt.args.w, tt.args.r)

			// assert
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.expect.status {
				t.Errorf("Bad StatusCode. Expected %d Got %d", tt.expect.status, resp.StatusCode)
			}

			rBytes, _ := ioutil.ReadAll(resp.Body)
			if len(rBytes) != tt.expect.size {
				t.Errorf("Bad Response Length. Expected %d Got %d", tt.expect.size, len(rBytes))
			}
		})
	}
}
