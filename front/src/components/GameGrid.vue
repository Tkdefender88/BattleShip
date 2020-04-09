<template>
  <div
    class="gameBoard"
    @dragover="dragover"
    @drop.prevent="drop"
    @dragleave="dragleave"
  >
    <grid-tile v-for="n in 100" :key="n - 1" :label="n - 1"></grid-tile>
  </div>
</template>

<script>
import GridTile from "./GridTile.vue";
export default {
  components: {
    GridTile
  },

  methods: {
    dragover: e => {
      let board = document.getElementById("oceanGrid");
      board.children.forEach(child => {
        child.classList.add("gridShow");
        child.children[0].classList.remove("hideBackground");
      });
      e.preventDefault();
    },

    dragleave: e => {
      let board = document.getElementById("oceanGrid");
      board.children.forEach(child => {
        child.classList.remove("gridShow");
        child.children[0].classList.add("hideBackground");
      });
      e.preventDefault();
    },

    drop(e) {
      let board = document.getElementById("oceanGrid");

      board.children.forEach(child => {
        child.classList.remove("gridShow");
        child.children[0].classList.add("hideBackground");
      });
      this.$emit("onPlace", e);
    }
  }
};
</script>

<style>
.gameBoard {
  --tileSz: 10;
  display: inline-grid;
  grid-template-columns: repeat(10, 45px);
  grid-template-rows: repeat(10, 45px);
}

.gridShow {
  border: solid 1px;
}

#destroyer0 {
  background: url("../../public/images/destroyer.png") 0px 0px;
}

#destroyer1 {
  background: url("../../public/images/destroyer.png") -46px 0px;
}
/*
 * Submarine Sprites
 */

#submarine0 {
  background: url("../../public/images/submarine.png") 0px 0px;
}

#submarine1 {
  background: url("../../public/images/submarine.png") -46px 0px;
}

#submarine2 {
  background: url("../../public/images/submarine.png") -91px 0px;
}
/*
 * Cruiser Sprites
 */

#cruiser0 {
  background: url("../../public/images/cruiser.png") 0px 0px;
}

#cruiser1 {
  background: url("../../public/images/cruiser.png") -46px 0px;
}

#cruiser2 {
  background: url("../../public/images/cruiser.png") -91px 0px;
}
/*
 * Battleship Sprites
 */

#battleship0 {
  background: url("../../public/images/battleship.png") 0px 0px;
}

#battleship1 {
  /* -46px means to the right 46px; 0px is from top */
  background: url("../../public/images/battleship.png") -46px 0px;
}

#battleship2 {
  background: url("../../public/images/battleship.png") -91px 0px;
}

#battleship3 {
  background: url("../../public/images/battleship.png") -136px 0px;
}
/*
 * Carrier Sprites
 */

#carrier0 {
  background: url("../../public/images/carrier.png") 0px 0px;
}

#carrier1 {
  background: url("../../public/images/carrier.png") -46px 0px;
}

#carrier2 {
  background: url("../../public/images/carrier.png") -91px 0px;
}

#carrier3 {
  background: url("../../public/images/carrier.png") -136px 0px;
}

#carrier4 {
  background: url("../../public/images/carrier.png") -181px 0px;
}

.hit {
  background: url("../../public/images/hit.png") 0px 0px;
}

.miss {
  background: url("../../public/images/miss.png");
}

.vertical {
  transform: rotate(270deg);
}
</style>
