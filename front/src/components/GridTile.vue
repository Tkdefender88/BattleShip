<template>
  <div class="gridTile">
    <span class="bgLayer hideBackground">{{ position }}</span>
    <span class="shipLayer"></span>
    <canvas class="fxLayer" width="45px" height="45px"></canvas>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: Number,
      required: true
    }
  },
  data(obj) {
    let row = Math.floor(obj.label / 10) + 65;
    let col = Number(obj.label) % 10;

    return {
      position: String.fromCharCode(row) + "-" + col
    };
  }
};
</script>

<style>
.bgLayer {
  position: absolute;
  width: 45px;
  height: 45px;
  z-index: 0;
  padding-top: 12px;
  text-align: center;
}

.hideBackground {
  display: none;
}

.shipLayer {
  position: absolute;
  width: 45px;
  height: 45px;
  z-index: 1;
}

.fxLayer {
  position: absolute;
  z-index: 2;
}

.gridTile {
  width: 45px;
  height: 45px;
  position: relative;
}
</style>
