<template>
  <div id="app">
    <nav-bar
      class="NavBar"
      v-bind:battle-state="this.battleState"
      v-on:update:battle-state="stateUpdate"
      v-on:new-game="newGame"
    ></nav-bar>
    <ship-panel class="ShipPanel"></ship-panel>
    <div class="GameArea">
      <game-grid id="oceanGrid" @onPlace="place"></game-grid>
      <game-grid id="targetGrid"></game-grid>
    </div>
    <div class="Footer">
      <div id="noticeBar">
        <span id="appCopyright">
          &copy; 2020 Department of Computer Science, Montana Tech All rights
          reserved
        </span>
        <span id="author">Author: Justin Bak</span>
        <span id="version">version: 0.0.1</span>
      </div>
    </div>
  </div>
</template>

<script>
import NavBar from "./components/NavBar.vue";
import ShipPanel from "./components/ShipPanel.vue";
import GameGrid from "./components/GameGrid.vue";
import BsState from "./bsState";

export default {
  name: "App",
  data: () => ({
    loading: true,
    loadedAssets: [],
    assets: [
      "images/carrier.png",
      "images/battleship.png",
      "images/cruiser.png",
      "images/submarine.png",
      "images/destroyer.png"
    ]
  }),
  methods: {
    place(e) {
      let cell = e.target.parentNode;
      let coordinate = cell.querySelector(".bgLayer").innerHTML.split("-");
      let shipType = e.dataTransfer.getData("ship-type");
      let rot = e.dataTransfer.getData("rotation");

      coordinate.push(Number(rot));

      this.battleState.place(coordinate, shipType);
    },
    async preload() {
      const calls = [];
      this.assets.forEach(asset => calls.push(fetch(asset)));

      const results = await Promise.all(calls);

      results.forEach(result => this.loadedAssets.push(result.url));
      this.loading = false;
    },
    stateUpdate: function(newState) {
      console.log(newState);
      this.battleState.carrier.placement = newState.carrier._placement;
      this.battleState.battleship.placement = newState.battleship._placement;
      this.battleState.cruiser.placement = newState.cruiser._placement;
      this.battleState.submarine.placement = newState.submarine._placement;
      this.battleState.destroyer.placement = newState.destroyer._placement;
      this.battleState.misses = newState.misses;
    },
    newGame: function() {
      this.battleState = new BsState(this.redraw);
      let oceanGrid = document.querySelector("#oceanGrid");

      // clear the grid
      oceanGrid.children.forEach(cell => {
        let shipLayer = cell.querySelector(".shipLayer");
        shipLayer.removeAttribute("id");
      });
    },
    setupStream() {
      let es = new EventSource("/events/updates");

      es.onmessage = function(event) {
        console.log(event.data);
      };
    },
    redraw() {
      // get the ships position from 0-99
      var oceanGrid = document.querySelector("#oceanGrid");
      [
        this.battleState.carrier,
        this.battleState.battleship,
        this.battleState.cruiser,
        this.battleState.submarine,
        this.battleState.destroyer
      ].forEach(ship => {
        if (ship.placed) {
          let coordinate = ship.placement;
          let row = coordinate[0];
          let col = Number(coordinate[1]);
          let pos = row * 10 + col;
          let rotation = coordinate[2];

          for (let i = 0; i < ship.size; i++) {
            let shipSprite = oceanGrid.children[pos].querySelector(
              ".shipLayer"
            );
            shipSprite.setAttribute("id", `${ship.name}${i}`);
            if (rotation === 1) {
              pos -= 10;
              shipSprite.classList.add("vertical");
            } else {
              pos++;
            }
          }
        }
      });
    }
  },
  created() {
    this.preload();
    this.battleState = new BsState(this.redraw);
    this.setupStream();
  },
  components: {
    NavBar,
    ShipPanel,
    GameGrid
  }
};
</script>

<style>
/* animation definitions */
@keyframes ocean-motion {
  0% {
    background-position: -25px -25px;
  }
  25% {
    background-position: -50px -50px;
  }
  50% {
    background-position: -25px -100px;
  }
  75% {
    background-position: 0px -50px;
  }
  100% {
    background-position: -25px -25px;
  }
}

#oceanGrid {
  animation: ocean-motion 12s infinite linear;
  background: url("../public/images/ocean.jpg");
  background-size: 1000px;
}

body {
  margin: 0px;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;

  height: 100vh;
  width: 100vw;

  display: grid;
  grid-template-rows: 1fr 10fr 0.5fr;
  grid-template-columns: minmax(235px, 0.3fr) 2fr;
  grid-template-areas:
    "header    header"
    "shipPanel gameboard"
    "footer    footer";
}

.NavBar {
  grid-area: header;
}

.ShipPanel {
  grid-area: shipPanel;
}

.GameArea {
  grid-area: gameboard;
  width: 100%;
  height: 100%;
  background-color: #e6d5a6;

  display: grid;
  grid-template-columns: 451px 451px;
  grid-template-rows: 451px;
  grid-gap: 20px;
  place-content: center;
}

/* Footer styling */
.Footer {
  background-color: #4c5760;
  color: #c6c6c6;

  /* grid properties */
  justify-self: stretch;
  grid-area: footer;
}

#noticeBar {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: 1fr;
  grid-template-areas: "version author copyright";
}
</style>
