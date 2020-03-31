package BattleState

type Ship struct {
	// Name of the ship, Carrier, Battleship etc.
	Name string `json:"_name"`
	// Size is the length of the ship
	Size int `json:"_size"`
	// Placed tells if the ship has a placement on the board
	Placed bool `json:"_placed"`
	// The placement of the player's ship, row column
	Placement []int `json:"_placement"`
	// HitProfiles are organized as a 2xSize array
	// The first row in the array is player's ship hit profile
	// The second row is the enemies ship hit profile
	// When the hit profile is filled the ship is sunk
	HitProfiles [][]string `json:"hitprofiles"`
}

type BsState struct {
	Destroyer  Ship `json:"destroyer"`
	Submarine  Ship `json:"submarine"`
	Cruiser    Ship `json:"cruiser"`
	Battleship Ship `json:"battleship"`
	Carrier    Ship `json:"carrier"`
	// Misses tracks the missed shots on the board
	Misses []string `json:"misses"`
}

func (bs *BsState) Valid() bool {
	return bs.Destroyer.Placed &&
		bs.Submarine.Placed &&
		bs.Cruiser.Placed &&
		bs.Battleship.Placed &&
		bs.Carrier.Placed
}

// Hit checks if a targeted shot hits any of ships on the players
// board. If there is a hit it returns true and the name of the hit
// ship. Otherwise false and an empty string
func (bs *BsState) Hit(target string) (bool, string) {
	tar := placementFromPrettyString([]rune(target))
	carrier := calculatePositions(bs.Carrier)
	if targetHitShip(tar, carrier) {
		return true, "carrier"
	}
	battleship := calculatePositions(bs.Battleship)
	if targetHitShip(tar, battleship) {
		return true, "battleship"
	}
	cruiser := calculatePositions(bs.Cruiser)
	if targetHitShip(tar, cruiser) {
		return true, "cruiser"
	}
	submarine := calculatePositions(bs.Submarine)
	if targetHitShip(tar, submarine) {
		return true, "submarine"
	}
	destroyer := calculatePositions(bs.Destroyer)
	if targetHitShip(tar, destroyer) {
		return true, "destroyer"
	}
	return false, ""
}

func targetHitShip(target []int, placement [][]int) bool {
	for i := 0; i < len(placement); i++ {
		if placement[i][0] == target[0] && placement[i][1] == target[1] {
			return true
		}
	}
	return false
}

func calculatePositions(ship Ship) [][]int {
	pos := make([][]int, ship.Size)
	for i := range pos {
		pos[i] = make([]int, 2)
	}
	for i := 0; i < ship.Size; i++ {
		if ship.Placement[2] == 1 { // vertical ship check
			pos[i] = []int{
				ship.Placement[0],
				ship.Placement[1] + i,
			}
		} else {
			pos[i] = []int{
				ship.Placement[0] + i,
				ship.Placement[1],
			}
		}
	}
	return pos
}

func placementFromPrettyString(target []rune) []int {
	row := int(target[0]) - 65
	col := int(target[0]) - 48
	return []int{row, col}
}
