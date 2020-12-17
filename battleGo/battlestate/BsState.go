package battlestate

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Miss       = "MISS"
	Carrier    = "CARRIER"
	BattleShip = "BATTLESHIP"
	Cruiser    = "CRUISER"
	Submarine  = "SUBMARINE"
	Destroyer  = "DESTROYER"
)

type Ship struct {
	// Name of the ship, Carrier, Battleship etc.
	Name string `json:"_name" bson:"_name,omitempty"`
	// Size is the length of the ship
	Size int `json:"_size" bson:"_size,omitempty"`
	// Placed tells if the ship has a placement on the board
	Placed bool `json:"_placed" bson:"_placed,omitempty"`
	// The placement of the player's ship, row column
	Placement []int `json:"_placement" bson:"_placement,omitempty"`
	// HitProfiles are organized as a 2xSize array
	// The first row in the array is player's ship hit profile
	// The second row is the enemies ship hit profile
	// When the hit profile is filled the ship is sunk
	HitProfiles [][]string `json:"hitprofiles" bson:"hitprofiles"`
}

// BsState represents the current battleship game state.
// sent to the client to update the view
type BsState struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
	Destroyer  *Ship              `json:"destroyer" bson:"destroyer"`
	Submarine  *Ship              `json:"submarine" bson:"submarine"`
	Cruiser    *Ship              `json:"cruiser" bson:"cruiser"`
	Battleship *Ship              `json:"battleship" bson:"battleship"`
	Carrier    *Ship              `json:"carrier" bson:"carrier"`
	// Misses tracks the missed shots on the board
	Misses []string `json:"misses" bson:"misses"`
}

// NewShip creates a new ship object of size, size
func NewShip() *Ship {
	hp := make([][]string, 2)
	return &Ship{
		HitProfiles: hp,
	}
}

// NewBsState creates a new BsState object with all the ships initialized
func NewBsState() *BsState {
	return &BsState{
		Destroyer:  NewShip(),
		Submarine:  NewShip(),
		Cruiser:    NewShip(),
		Battleship: NewShip(),
		Carrier:    NewShip(),
	}
}

// Valid ensures the user has selected a ship placement that is valid for a game
func (bs *BsState) Valid() bool {
	if !(bs.Destroyer != nil && bs.Submarine != nil && bs.Cruiser != nil &&
		bs.Battleship != nil && bs.Carrier != nil) {
		return false
	}
	return bs.Destroyer.Placed &&
		bs.Submarine.Placed &&
		bs.Cruiser.Placed &&
		bs.Battleship.Placed &&
		bs.Carrier.Placed
}

// Hit checks if a targeted shot hits any of ships on the players
// board. If there is a hit it returns true and the name of the hit
// ship. Otherwise false and the string "MISS"
func (bs *BsState) Hit(target string) (bool, string) {
	tar := placementFromPrettyString(target)
	carrier := calculatePositions(bs.Carrier)
	if targetHitShip(tar, carrier) {
		bs.Carrier.HitProfiles[0] = append(bs.Carrier.HitProfiles[0], target)
		return true, Carrier
	}
	battleship := calculatePositions(bs.Battleship)
	if targetHitShip(tar, battleship) {
		bs.Battleship.HitProfiles[0] = append(bs.Battleship.HitProfiles[0], target)
		return true, BattleShip
	}
	cruiser := calculatePositions(bs.Cruiser)
	if targetHitShip(tar, cruiser) {
		bs.Cruiser.HitProfiles[0] = append(bs.Cruiser.HitProfiles[0], target)
		return true, Cruiser
	}
	submarine := calculatePositions(bs.Submarine)
	if targetHitShip(tar, submarine) {
		bs.Submarine.HitProfiles[0] = append(bs.Submarine.HitProfiles[0], target)
		return true, Submarine
	}
	destroyer := calculatePositions(bs.Destroyer)
	if targetHitShip(tar, destroyer) {
		bs.Destroyer.HitProfiles[0] = append(bs.Destroyer.HitProfiles[0], target)
		return true, Destroyer
	}
	bs.Misses = append(bs.Misses, target)
	return false, Miss
}

// Sunk calculates if a particular ship is sunk
// returns two booleans the first for the player's ship and the second for the
// opponents ship.
func (s *Ship) Sunk() (bool, bool) {
	return len(s.HitProfiles[0]) == s.Size, len(s.HitProfiles[1]) == s.Size
}

// GameLost returns true if all the player's ships have been sunk
func (bs *BsState) GameLost() bool {
	carrier, _ := bs.Carrier.Sunk()
	battleship, _ := bs.Battleship.Sunk()
	cruiser, _ := bs.Cruiser.Sunk()
	submarine, _ := bs.Submarine.Sunk()
	destroyer, _ := bs.Destroyer.Sunk()

	return carrier && battleship && cruiser && submarine && destroyer
}

func (bs *BsState) ShipFromString(shipName string) *Ship {
	if strings.ToLower(shipName) == "carrier" {
		return bs.Carrier
	}
	if strings.ToLower(shipName) == "battleship" {
		return bs.Battleship
	}
	if strings.ToLower(shipName) == "cruiser" {
		return bs.Cruiser
	}
	if strings.ToLower(shipName) == "submarine" {
		return bs.Submarine
	}
	if strings.ToLower(shipName) == "destroyer" {
		return bs.Destroyer
	}
	return nil
}

// HitEnemy add a target position to the enemy ship hit profile
func (bs *BsState) HitEnemy(shipName string, target string) {
	if shipName == Carrier {
		bs.Carrier.HitProfiles[1] = append(bs.Carrier.HitProfiles[1], target)
	}
	if shipName == BattleShip {
		bs.Battleship.HitProfiles[1] = append(bs.Battleship.HitProfiles[1], target)
	}
	if shipName == Cruiser {
		bs.Cruiser.HitProfiles[1] = append(bs.Cruiser.HitProfiles[1], target)
	}
	if shipName == Submarine {
		bs.Submarine.HitProfiles[1] = append(bs.Submarine.HitProfiles[1], target)
	}
	if shipName == Destroyer {
		bs.Destroyer.HitProfiles[1] = append(bs.Destroyer.HitProfiles[1], target)
	}
}

func targetHitShip(target []int, placement [][]int) bool {
	for i := 0; i < len(placement); i++ {
		if placement[i][0] == target[0] && placement[i][1] == target[1] {
			return true
		}
	}
	return false
}

func calculatePositions(ship *Ship) [][]int {
	pos := make([][]int, ship.Size)
	for i := range pos {
		pos[i] = make([]int, 2)
	}
	for i := 0; i < ship.Size; i++ {
		if ship.Placement[2] == 1 { // vertical ship check
			pos[i] = []int{
				ship.Placement[0] - i,
				ship.Placement[1],
			}
		} else {
			pos[i] = []int{
				ship.Placement[0],
				ship.Placement[1] + i,
			}
		}
	}
	return pos
}

func placementFromPrettyString(target string) []int {
	t := []rune(target)
	row := int(t[0]) - 65
	col := int(t[1]) - 48
	return []int{row, col}
}
