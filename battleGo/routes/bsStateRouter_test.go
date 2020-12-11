package routes_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/battlestate"
	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/repository"
	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/routes"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	stackModel = &battlestate.BsState{
		ID: primitive.NewObjectID(),
		Destroyer: &battlestate.Ship{
			Name:        "destroyer",
			Size:        2,
			Placed:      true,
			Placement:   []int{3, 3, 1},
			HitProfiles: [][]string{},
		},
		Submarine: &battlestate.Ship{
			Name:        "submarine",
			Size:        3,
			Placed:      true,
			Placement:   []int{3, 0, 0},
			HitProfiles: [][]string{},
		},
		Cruiser: &battlestate.Ship{
			Name:        "cruiser",
			Size:        3,
			Placed:      true,
			Placement:   []int{2, 0, 0},
			HitProfiles: [][]string{},
		},
		Battleship: &battlestate.Ship{
			Name:        "battleship",
			Size:        4,
			Placed:      true,
			Placement:   []int{1, 0, 0},
			HitProfiles: [][]string{},
		},
		Carrier: &battlestate.Ship{
			Name:        "carrier",
			Size:        5,
			Placed:      true,
			Placement:   []int{0, 0, 0},
			HitProfiles: [][]string{},
		},
		Misses: []string{},
	}
	badModel = &battlestate.BsState{
		Destroyer: &battlestate.Ship{
			Name:        "destroyer",
			Size:        2,
			Placed:      true,
			Placement:   []int{3, 3, 1},
			HitProfiles: [][]string{},
		},
		Cruiser: &battlestate.Ship{
			Name:        "cruiser",
			Size:        3,
			Placed:      true,
			Placement:   []int{2, 0, 0},
			HitProfiles: [][]string{},
		},
		Battleship: &battlestate.Ship{
			Name:        "battleship",
			Size:        4,
			Placed:      true,
			Placement:   []int{1, 0, 0},
			HitProfiles: [][]string{},
		},
		Carrier: &battlestate.Ship{
			Name:        "carrier",
			Size:        5,
			Placed:      true,
			Placement:   []int{0, 0, 0},
			HitProfiles: [][]string{},
		},
		Misses: []string{},
	}
)

type mockRepo struct {
}

func (mr *mockRepo) FindModel(name string) (*battlestate.BsState, error) {
	if name == "foo" {
		return nil, mongo.ErrNoDocuments
	}
	return stackModel, nil
}

func (mr *mockRepo) ListModels() ([]string, error) {
	return []string{"a", "b", "c"}, nil
}

func (mr *mockRepo) DeleteModel(name string) error {
	return nil
}

func (mr *mockRepo) CreateModel(name string, model *battlestate.BsState) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}

//////////////////////////////////// Unit Test Funcs //////////////////////////////////////////////

type fields struct {
	repo repository.ModelRepository
}
type args struct {
	w *httptest.ResponseRecorder
	r *http.Request
}
type expect struct {
	status int
	body   []byte
}

func TestBsStateResource_Get(t *testing.T) {
	stacky, _ := json.Marshal(stackModel)

	tests := []struct {
		name   string
		fields fields
		args   args
		expect expect
	}{
		{
			name: "happyPath",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/stacky", nil),
			},
			fields: fields{
				repo: &mockRepo{},
			},
			expect: expect{
				status: http.StatusOK,
				body:   append(stacky, byte('\n')),
			},
		},
		{
			name: "filenotfound",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/foo", nil),
			},
			fields: fields{
				repo: &mockRepo{},
			},
			expect: expect{
				status: http.StatusNotFound,
				body:   []byte(`{"message":"Document foo Not Found"}` + "\n"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			r := chi.NewRouter()
			rs := routes.NewBsStateController(tt.fields.repo)
			r.Route("/{filename}", func(r chi.Router) {
				r.Get("/", rs.Get)
			})

			// act
			r.ServeHTTP(tt.args.w, tt.args.r)

			// assert
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.expect.status {
				t.Errorf("Bad status code: Expected %d Got %d", tt.expect.status, resp.StatusCode)
			}

			respBody, _ := ioutil.ReadAll(resp.Body)
			if bytes.Compare(respBody, tt.expect.body) != 0 {
				t.Errorf("Error in response body: Expected %s Got %s", string(tt.expect.body), string(respBody))
			}
		})
	}
}

func TestBsStateResource_Post(t *testing.T) {
	stacky, _ := json.Marshal(stackModel)
	badM, _ := json.Marshal(badModel)

	tests := []struct {
		name   string
		fields fields
		args   args
		expect expect
	}{
		{
			name: "happyPath",
			fields: fields{
				repo: &mockRepo{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/stacky", bytes.NewBuffer(stacky)),
			},
			expect: expect{
				status: 201,
				body:   []byte("{\"id\":\"ObjectID(\\\"000000000000000000000000\\\")\"}\n"),
			},
		},
		{
			name: "NoBody",
			fields: fields{
				repo: &mockRepo{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/stacky", nil),
			},
			expect: expect{
				status: http.StatusInternalServerError,
				body:   []byte("{\"message\":\"\"}\n"),
			},
		},
		{
			name: "InvalidModel",
			fields: fields{
				repo: &mockRepo{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/stacky", bytes.NewBuffer(badM)),
			},
			expect: expect{
				status: http.StatusBadRequest,
				body:   []byte("{\"message\":\"Invalid model\"}\n"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			rs := routes.NewBsStateController(tt.fields.repo)
			r := chi.NewRouter()
			r.Route("/{filename}", func(r chi.Router) {
				r.Post("/", rs.Post)
			})

			// act
			r.ServeHTTP(tt.args.w, tt.args.r)

			// assert
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.expect.status {
				t.Errorf("expected %d got %d", tt.expect.status, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}
			defer resp.Body.Close()

			if bytes.Compare(body, tt.expect.body) != 0 {
				t.Errorf("Expected %+v got %+v", tt.expect.body, body)
				t.Logf("%s", string(body))
				t.Logf("%s", string(tt.expect.body))
			}
		})
	}
}
