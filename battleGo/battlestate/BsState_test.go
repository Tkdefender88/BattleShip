package battlestate

import (
	"reflect"
	"testing"
)

func TestCalculatePositions(t *testing.T) {
	cases := []struct {
		ship       *Ship
		resultList [][]int
	}{
		{
			&Ship{
				Size:      5,
				Placement: []int{0, 0, 0},
			},
			[][]int{
				{0, 0},
				{0, 1},
				{0, 2},
				{0, 3},
				{0, 4},
			},
		},
		{
			&Ship{
				Size:      4,
				Placement: []int{5, 5, 1},
			},
			[][]int{
				{5, 5},
				{4, 5},
				{3, 5},
				{2, 5},
			},
		},
	}

	for _, c := range cases {
		p := calculatePositions(c.ship)
		if len(p) != len(c.resultList) {
			t.Errorf("Test Failed. Bad result length Got %d Want %d", len(p), len(c.resultList))
		}

		for i := 0; i < len(p); i++ {
			for j := 0; j < len(p[i]); j++ {
				if p[i][j] != c.resultList[i][j] {
					t.Errorf("Test failed. Bad result at (%d, %d) Got %d Want %d", i, j, p[i][j], c.resultList[i][j])
				}
			}
		}
	}
}

func Test_targetHitShip(t *testing.T) {
	type args struct {
		target    []int
		placement [][]int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "Good Case Hit",
			args: args{
				[]int{0, 0},
				[][]int{
					{0, 0},
					{1, 0},
					{2, 0},
					{3, 0},
					{4, 0},
				},
			},
			want: true,
		},
		{
			name: "Good Case Miss",
			args: args{
				[]int{5, 5},
				[][]int{
					{0, 0},
					{1, 0},
					{2, 0},
					{3, 0},
					{4, 0},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := targetHitShip(tt.args.target, tt.args.placement); got != tt.want {
				t.Errorf("targetHitShip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_placementFromPrettyString(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			"A0 to 0, 0",
			args{"A0"},
			[]int{0, 0},
		},
		{
			"B4 to 1, 4",
			args{"B4"},
			[]int{1, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := placementFromPrettyString(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("placementFromPrettyString() = %v, want %v", got, tt.want)
			}
		})
	}
}
