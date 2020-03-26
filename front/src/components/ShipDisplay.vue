<template>
  <div class="shipSelect">
    <toggle-switch @onToggle="updateToggle">Vertical</toggle-switch>
    <div
      :class="{
        displayCarrier: shipType === 'carrier',
        displayBattleship: shipType === 'battleship',
        displayCruiser: shipType === 'cruiser',
        displaySubmarine: shipType === 'submarine',
        displayDestroyer: shipType === 'destroyer'
      }"
      :draggable="draggable"
      @dragstart="dragStart"
      @dragover.stop
    ></div>
  </div>
</template>

<script>
import ToggleSwitch from "./ToggleSwitch.vue";
export default {
  props: ["draggable", "shipType"],
  components: {
    ToggleSwitch
  },
  methods: {
    updateToggle() {
      this.toggled = !this.toggled;
    },
    dragStart(e) {
      let imgSrc = new Image();
      e.dataTransfer.clearData();
      let ship = e.target.classList[0];
      let rotation = this.toggled ? 1 : 0;
      let shipName = ship.toLowerCase().slice(7, ship.length);

      imgSrc.src = `images/${shipName}.png`;

      if (rotation === 1) {
        imgSrc.classList.add("rotateImg90");
      }

      e.dataTransfer.setDragImage(imgSrc, 22, 22);
      e.dataTransfer.setData("ship-type", shipName);
      e.dataTransfer.setData("rotation", rotation);
    }
  }
};
</script>

<style>
.displayDestroyer {
  width: 90px;
  height: 45px;
  background: url("../../public/images/destroyer.png");
}

.displaySubmarine {
  width: 135px;
  height: 45px;
  background: url("../../public/images/submarine.png");
}

.displayCruiser {
  width: 135px;
  height: 45px;
  background: url("../../public/images/cruiser.png");
}

.displayBattleship {
  width: 180px;
  height: 45px;
  background: url("../../public/images/battleship.png");
}

.displayCarrier {
  width: 225px;
  height: 45px;
  background: url("../../public/images/carrier.png");
}

.rotateImg90 {
  transform: rotate(90deg);
}
</style>
