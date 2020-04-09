<template>
  <div id="app">
    <nav-bar
      class="NavBar"
      v-bind:battle-state="this.battleState"
      v-on:update:battle-state="stateUpdate"
      v-on:battle-modal="showModal"
      v-on:new-game="newGame"
    ></nav-bar>
    <ship-panel class="ShipPanel"></ship-panel>

    <modal v-show="isModalVisible" @close="closeModal" @done="submit">
      <div slot="header">
        <h2>Start a Battle!</h2>
      </div>
      <div id="modalBody" slot="body">
        <form>
          <input type="text" id="stateName" name="stateName" value="stacky" />
          <label for="stateName">Battle State</label>
          <br />
          <input type="text" id="oppUrl" name="oppUrl" value />
          <label for="oppUrl">Opponent's URL (optional)</label>
          <br />
          <br />
        </form>
      </div>
    </modal>
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
import fireEvent from "./fireEvent";
import modal from "./components/modal.vue";
import axios from "axios"

export default {
  name: "App",
  data() {
    return {
      loading: true,
      loadedAssets: [],
      isModalVisible: false,
      assets: [
        "images/carrier.png",
        "images/battleship.png",
        "images/cruiser.png",
        "images/submarine.png",
        "images/destroyer.png"
      ]
    };
  },
  methods: {
    submit() {
      let modelName = document.getElementById("stateName").value;
      let url = document.getElementById("oppUrl").value;
      console.log(url);
      console.log(modelName);
      let endpoint = url === "" ? "/battle/"+modelName : "/battle/"+modelName+"/"+url;
      let that = this;
      axios.get(endpoint ).then((resp) => {
        that.stateUpdate(resp.data);
      }).catch(reason => {
        alert("An error occured: " + reason);
      });
      document.getElementById("stateName").value = "stacky";
      document.getElementById("oppUrl").value = "";
      this.isModalVisible = false;
    },
    showModal() {
      this.isModalVisible = true;
    },
    closeModal() {
      this.isModalVisible = false;
    },
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
      let targetGrid = document.querySelector("#targetGrid");

      targetGrid.children.forEach(cell => {
        let fxLayer = cell.querySelector(".fxLayer");
        fxLayer.classList.remove("miss");
        fxLayer.classList.remove("hit");
      });

      // clear the grid
      oceanGrid.children.forEach(cell => {
        let shipLayer = cell.querySelector(".shipLayer");
        let fxLayer = cell.querySelector(".fxLayer");
        shipLayer.removeAttribute("id");
        shipLayer.classList.remove("vertical");
        fxLayer.classList.remove("miss");
        fxLayer.classList.remove("hit");
      });
    },
    setupStream() {
      let es = new EventSource("/events/updates");

      es.addEventListener(
        "message",
        e => {
            let evt = new fireEvent(JSON.parse(e.data));
            console.log(e.data);
            evt.updateGrid();
        },
        false
      );
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
    GameGrid,
    modal
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

#targetGrid {
  background: #ad9a5a;
}

body {
  margin: 0px;
}

input {
  padding: 10px;
  border: solid 1px #e5e5e5;
  outline: 0;
  width: 200px;
  background: #ffffff;
  margin: 5px;
  box-shadow: rgba(0,0,0,0.1) 0 0 8px;
}

input:hover,
input:focus {
  border-color: #c9c9c9;
}

form label {
  margin-left: 10px;
  color: #999999;
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
