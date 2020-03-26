<template>
  <div id="NavBar">
    <div id="titleBar">
      <span id="appTitle">
        <h1>Battleship</h1>
      </span>
      <span id="appVersion">
        <h4>Revision: 2020-02-03</h4>
      </span>
    </div>

    <div id="toolBar">
      <div class="dropdown">
        <button class="dropbtn">Game</button>
        <div class="dropdown-content">
          <a @click="newGameBtn" id="NewGameBtn" href="#">New Game</a>
          <a @click="loadGameBtn" id="LoadGameBtn" href="#">Load Board</a>
          <a @click="saveGameBtn" id="SaveGameBtn" href="#">Save Game</a>
          <a @click="exitGameBtn" id="ExitGameBtn" href="#">Exit Game</a>
        </div>
      </div>
      <div class="dropdown">
        <button class="dropbtn">Info</button>
        <div class="dropdown-content">
          <a @click="helpBtn" id="HelpBtn" href="#">Help</a>
          <a @click="aboutBtn" id="AboutBtn" href="#">About</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  props: {
    battleState: Object
  },
  methods: {
    ajaxPOST(uri, body, success) {
      var req = new XMLHttpRequest();

      req.onreadystatechange = function() {
        if (req.readyState === 3 && req.status === 201) {
          if (typeof callback === "function") {
            success(req.responseText);
          }
        }
      };

      req.open("POST", uri, true);
      req.send(body);
    },
    helpBtn() {
      alert("Help button pressed");
    },
    aboutBtn() {
      alert("About button pressed");
    },
    newGameBtn() {
      this.$emit("new-game");
    },
    loadGameBtn() {
      let model = prompt("Enter a file name to load:");
      if (model == null || model === "") {
        return;
      }
      let that = this;
      axios.get("bsState/" + model).then(function(resp) {
        that.$emit("update:battle-state", resp.data);
        alert("Loaded model!");
      });
    },
    saveGameBtn() {
      let modelName = prompt("Enter a file name to save this as:");
      if (modelName == null || modelName === "") {
        return;
      }
      axios
        .post("bsState/" + modelName, JSON.stringify(this.battleState))
        .then(() => {
          alert("Configuration Saved");
        })
        .catch(err => {
          alert("Error saving: " + err);
        });
      /*
      this.ajaxPOST(
        "bsState/" + modelName,
        JSON.stringify(this.battleState),
        resp => {
          if (resp.code === 201) {
            alert("Successful save");
          }
          if (resp.code === 400) {
            alert("There was a problem saving");
          }
        }
      );
      */
    },
    exitGameBtn() {
      alert("exit game button pressed");
    }
  }
};
</script>

<style>
:root {
  --header-bg: #93a8ac;
  --header-fg: #584b53;
  --header-fg-shade: #493e44;
}

h1 {
  font-size: 48px;
  font-weight: bold;
  margin: 0;
}

h4 {
  margin: 0;
}

#NavBar {
  background-color: var(--header-bg);

  /* grid properties */
  justify-self: stretch;

  display: grid;
  grid-template-rows: 1fr 0.5fr;
  grid-template-columns: 1fr;
  grid-gap: 5px;
  grid-template-areas:
    "titleBar"
    "toolBar";
}

#titleBar {
  padding-left: 5px;
}

#toolBar {
  grid-area: toolBar;
  width: 235px;
  display: grid;
  grid-template-rows: 1fr;
  grid-template-columns: 1fr 1fr;
  place-items: center;
}

.dropdown {
  display: inline-block;
  margin-bottom: 2px;
  padding-left: 5px;
}

.dropbtn {
  background-color: var(--header-fg);
  color: white;
  padding: 10px;
  font-size: 16px;
  min-width: 80px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

.dropdown-content {
  display: none;
  position: absolute;
  background-color: #f9f9f9;
  border-radius: 3px;
  box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
  z-index: 1;
}

.dropdown-content a {
  color: black;
  padding: 12px 16px;
  text-decoration: none;
  display: block;
  border-radius: 3px;
}

/* hover and active styles */

.dropdown-content a:hover {
  background-color: #c1c1c1;
}

.dropdown:hover .dropdown-content {
  display: block;
}

.dropdown:hover .dropbtn {
  background-color: var(--header-fg-shade);
}
</style>
