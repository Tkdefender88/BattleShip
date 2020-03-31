package BattleState

import "testing"

func TestCalculatePositions(t *testing.T) {
	cases := []struct {
		ship       Ship
		resultList [][]int
	}{
		{
			Ship{
				Size:      5,
				Placement: []int{0, 0, 0},
			},
			[][]int{
				{0, 0},
				{1, 0},
				{2, 0},
				{3, 0},
				{4, 0},
			},
		},
		{
			Ship{
				Size:      4,
				Placement: []int{5, 5, 1},
			},
			[][]int{
				{5, 5},
				{5, 6},
				{5, 7},
				{5, 8},
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
