<template>
  <span id="shipPanel">
    <ShipDisplay shipType="carrier" :draggable="true"></ShipDisplay>
    <ShipDisplay shipType="battleship" :draggable="true"></ShipDisplay>
    <ShipDisplay shipType="cruiser" :draggable="true"></ShipDisplay>
    <ShipDisplay shipType="submarine" :draggable="true"></ShipDisplay>
    <ShipDisplay shipType="destroyer" :draggable="true"></ShipDisplay>
  </span>
</template>

<script>
import ShipDisplay from "./ShipDisplay";

export default {
  name: "ShipPanel",
  components: {
    ShipDisplay
  }
};
</script>

<style>
#shipPanel {
  background-color: #b57e4d;

  /* grid properties */
  align-items: center;
  justify-items: left;
  padding-left: 5px;
  border-radius: 4px;
  display: grid;
  grid-template-rows: repeat(5, 1fr);
}
</style>
