export default class {
    hit;
    player;
    tile;
    constructor(data) {
        this.hit = data.hit;
        this.tile= data.tile;
        this.player = data.player;
    }

    updateGrid() {
        let gridId = this.player === "player" ? "oceanGrid" : "targetGrid";
        let grid = document.getElementById(gridId);
        if (this.hit) {
            grid.children[this.tile].querySelector(".fxLayer").classList.add("hit");
        } else {
            grid.children[this.tile].querySelector(".fxLayer").classList.add("miss");
        }
    }
}