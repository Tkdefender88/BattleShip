package solver

import (
	"math/rand"
	"sort"
	"time"
)

type (
	shipOrientation string

	// Ship represents a game ship
	Ship struct {
		Alive       bool
		Size        int
		Position    int
		Orientation shipOrientation
	}

	// Position marks a position in the battle ship grid as well as
	// useful information about ship placement.
	Position struct {
		Index       int
		Probability int
		Occupied    bool
		Fired       bool
		Hit         bool
		Confirmed   bool
		Row         int
		Col         int
		Directions  map[string]*Position
	}

	Positions []*Position

	// Strategy is a struct that is responsible for calculating the next moves for
	// a battle ship game based on grid and ship information
	Strategy struct {
		Grid     Positions
		Ships    map[string]*Ship
		GameOver bool
	}
)

const (
	vertical   shipOrientation = "vertical"
	horizontal                 = "horizontal"
)

var (
	shotCount  = 0
	randomness = 3
)

// NewStrategy creates a new strategy instance that with a grid
// initialised with all references to the tiles
func NewStrategy() *Strategy {

	positions := make([]*Position, 100)

	for i := 0; i < 100; i++ {
		row := i / 10
		col := i % 10
		positions[i] = &Position{
			Index:       i,
			Probability: 0,
			Occupied:    false,
			Fired:       false,
			Hit:         false,
			Confirmed:   false,
			Row:         row,
			Col:         col,
			Directions: map[string]*Position{
				"W": nil,
				"E": nil,
				"S": nil,
				"N": nil,
			},
		}
	}

	// Set references to adjacent tiles
	for i := 0; i < 100; i++ {
		if i != 0 && positions[i-1].Row == positions[i].Row {
			positions[i].Directions["W"] = positions[i-1]
		} else {
			positions[i].Directions["W"] = nil
		}
		if i < 99 && positions[i+1].Row == positions[i].Row {
			positions[i].Directions["E"] = positions[i+1]
		} else {
			positions[i].Directions["E"] = nil
		}
		if i-10 > 0 {
			positions[i].Directions["N"] = positions[i-10]
		} else {
			positions[i].Directions["N"] = nil
		}
		if i+10 < 99 {
			positions[i].Directions["S"] = positions[i+10]
		} else {
			positions[i].Directions["S"] = nil
		}
	}

	return &Strategy{
		Grid: positions,
		Ships: map[string]*Ship{
			"battleship": &Ship{
				Alive:       true,
				Size:        4,
				Position:    0,
				Orientation: vertical,
			},
			"destroyer": &Ship{
				Alive:       true,
				Size:        2,
				Position:    0,
				Orientation: vertical,
			},
			"carrier": &Ship{
				Alive:       true,
				Size:        5,
				Position:    0,
				Orientation: vertical,
			},
			"submarine": &Ship{
				Alive:       true,
				Size:        3,
				Position:    0,
				Orientation: vertical,
			},
			"cruiser": &Ship{
				Alive:       true,
				Size:        3,
				Position:    0,
				Orientation: vertical,
			},
		},
		GameOver: false,
	}
}

func (s *Strategy) Step() int {
	s.updateProbabilities()
	return s.fireNext()
}

func (s *Strategy) zeroBoard() {
	for _, pos := range s.Grid {
		pos.Probability = 0
	}
}

func (s *Strategy) updateProbabilities() {
	var lastPosition *Position
	var hitStreak int

	directions := []string{"W", "N", "E", "S"}

	rand.Seed(time.Now().Unix())

	s.zeroBoard()

	for _, ship := range s.Ships {
		if ship.Alive {
			for i := 0; i < 100; i++ {
				if s.tryShipAtPosition(ship, i, "S") {
					lastPosition = s.Grid[i]
					for j := 0; j < ship.Size; j++ {
						lastPosition.Probability++
						lastPosition = lastPosition.Directions["S"]
					}
				}
				if s.tryShipAtPosition(ship, i, "E") {
					lastPosition = s.Grid[i]
					for j := 0; j < ship.Size; j++ {
						lastPosition.Probability++
						lastPosition = lastPosition.Directions["E"]
					}
				}
			}
		}
	}

	for _, position := range s.Grid {
		if position.Probability > 0 {
			position.Probability += rand.Intn(randomness)
		}
		if position.Fired {
			if position.Hit && !position.Confirmed {
				for _, dir := range directions {
					lastPosition = position
					hitStreak = 1
					for lastPosition.Directions[dir] != nil && lastPosition.Directions[dir].Hit && !lastPosition.Directions[dir].Confirmed {
						hitStreak++
						lastPosition = lastPosition.Directions[dir]
					}
					lastPosition = lastPosition.Directions[dir]
					if lastPosition != nil && !lastPosition.Fired {
						lastPosition.Probability += hitStreak * 10
					}
				}
			}
			position.Probability = -1
		}
	}
}

func (s *Strategy) tryShipAtPosition(ship *Ship, index int, orientation string) bool {
	lastPosition := s.Grid[index]
	fit := true
	for i := 0; i < ship.Size; i++ {
		if lastPosition == nil || lastPosition.Confirmed {
			fit = false
			break
		}
		lastPosition = lastPosition.Directions[orientation]
	}

	return fit
}

func (s *Strategy) RemoveShip(shipName string) {
	allDead := true
	s.Ships[shipName].Alive = false

	for _, ship := range s.Ships {
		if ship.Alive {
			allDead = false
		}
	}

	if allDead {
		s.GameOver = true
	}
	s.updateProbabilities()
}

func (s *Strategy) fireNext() (index int) {
	shotCount++
	// Sort the positions based on probability
	sort.Slice(s.Grid, func(i, j int) bool {
		if s.Grid[i].Fired {
			return false
		}
		return s.Grid[i].Probability > s.Grid[j].Probability
	})

	index = s.Grid[0].Index

	s.Grid[0].Fired = true
	s.Grid[0].Probability = -1

	// Return the slice to it's original order
	sort.Slice(s.Grid, func(i, j int) bool {
		return s.Grid[i].Index < s.Grid[j].Index
	})

	return
}

func indexFromString(target []rune) int {
	row := int(target[0]) - 65
	col := int(target[1]) - 48
	return row*10 + col
}

func (s *Strategy) ConfirmShot(target string, hit bool) {
	index := indexFromString([]rune(target))
	s.Grid[index].Hit = hit
	if hit {
		s.Grid[index].Confirmed = false
	} else {
		s.Grid[index].Confirmed = true
	}
}
