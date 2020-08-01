import React from 'react';
import { Component } from 'react';
import GridTile from '../grid-tile/grid-tile';

class GameBoard extends Component {
    render() {

        var items = [];

        for (let i = 0; i < 100; i ++ ) {
            items.push(<GridTile key={i}/>);
        }

        return (
            <div className="game-board">
                {items}
            </div>
        );
    }
}

export default GameBoard;