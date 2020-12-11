package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/battlestate"
	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/repository"
	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/routes"
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
			sesh.Get(tt.args.w, tt.args.r)

			resp := tt.args.w.Result()

			if resp.StatusCode != tt.expect.status {
				t.Errorf("Incorrect response code expected %d got %d", tt.expect.status, resp.StatusCode)
			}
		})
	}
}

func Test_DeleteSession(t *testing.T) {
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
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sesh := routes.NewSession(tt.fields.repo)

			chi.NewRouter().Route("/{session_id}", func(r chi.Router) {
				r.Delete("/", sesh.Delete)
			}).ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
