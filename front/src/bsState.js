const bsStateArrayHandler = function(routine) {
  return {
    set(property) {
      let retval = Reflect.set(...arguments);

      if (Number.isInteger(Number(property))) {
        routine();
      }

      return retval;
    },

    deleteProperty() {
      let retval = Reflect.deleteProperty(...arguments);
      routine();
      return retval;
    }
  };
};

// game ship class
class bsShip {
  constructor(name, size, redraw) {
    this.redraw = redraw;
    this._name = name;
    this._size = size;
    this._placed = false;
    this._placement = [];
    this.hitprofiles = [
      new Proxy(new Array(), bsStateArrayHandler(this.redraw)),
      new Proxy(new Array(), bsStateArrayHandler(this.redraw))
    ];
  }

  get name() {
    return this._name;
  }
  set name(value) {
    this._name = value;
  }

  get size() {
    return this._size;
  }
  set size(value) {
    this._size = value;
  }

  get placed() {
    return this._placed;
  }

  set placed(value) {
    this._placed = value;
  }

  get placement() {
    return this._placement;
  }
  set placement(value) {
    if (!this._placed) {
      this._placement = value;
      this._placed = true;
      this.redraw(this);
    }
  }
}

export default class {
  constructor(redraw) {
    this.redraw = redraw;
    this.destroyer = new bsShip("destroyer", 2, redraw);
    this.submarine = new bsShip("submarine", 3, redraw);
    this.cruiser = new bsShip("cruiser", 3, redraw);
    this.battleship = new bsShip("battleship", 4, redraw);
    this.carrier = new bsShip("carrier", 5, redraw);
    this.misses = new Proxy(new Array(), bsStateArrayHandler(this.redraw));
  } // end constructor

  place(coord, ship) {
    let minRow = 0;
    let maxCol = 9;
    let row = coord[0].charCodeAt(0) - 65;
    let col = Number(coord[1]);
    let shipLength = this[ship].size;
    let validPlace = true;

    let ships = ["carrier", "battleship", "cruiser", "submarine", "destroyer"];

    // check horizontal out of bounds
    if (col + shipLength - 1 > maxCol && coord[2] === 0) {
      col = maxCol - shipLength + 1;
    }
    coord[1] = col;

    // check vertical out of bounds
    if (row - shipLength + 1 < minRow && coord[2] === 1) {
      row = shipLength - 1;
    }
    coord[0] = row;

    // Check for ship collisions
    let newShipPlace = [];
    for (let i = 0; i < shipLength; i++) {
      let c = {
        y: coord[2] === 1 ? row - i : row,
        x: coord[2] === 0 ? col + i : col
      };
      newShipPlace.push(c);
    }

    ships.forEach(s => {
      let curShip = this[s];
      if (curShip.placement) {
        let y = curShip.placement[0],
          x = curShip.placement[1],
          rot = curShip.placement[2];

        let curShipPlace = [];
        for (let i = 0; i < curShip.size; i++) {
          let c = {
            x: rot === 0 ? x + i : x,
            y: rot === 1 ? y - i : y
          };
          curShipPlace.push(c);
        }

        for (let i = 0; i < curShipPlace.length; i++) {
          for (let j = 0; j < newShipPlace.length; j++) {
            if (
              newShipPlace[j].x === curShipPlace[i].x &&
              newShipPlace[j].y === curShipPlace[i].y
            ) {
              validPlace = false;
              return;
            }
          }
        }
      }
    });

    if (validPlace) {
      this[ship].placement = coord;
      console.log(this[ship].placement);
    }

    return true;
  }
}
