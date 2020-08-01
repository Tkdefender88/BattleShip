import React from 'react';
import { Component } from 'react';

class GridTile extends Component {
    render() {
        return (
            <div className="grid-tile">
                <span class="bg-layer hide-background"></span>
                <span class="ship-layer"></span>
                <canvas class="fx-layer" width="45px" height="45px"></canvas>
            </div>
        )
    }
}

export default GridTile;