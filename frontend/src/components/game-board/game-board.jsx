import React from 'react';
import { Component } from 'react';
import GridTile from '../grid-tile/grid-tile';

class GameBoard extends Component {

	onDrop(e) {
		let ship = e.dataTransfer.getData("ship")
	}

	render() {

		var items = [];

		for (let i = 0; i < 100; i ++ ) {
			items.push(<GridTile key={i}/>);
		}

		return (
			<div
				className="game-board dropable"
				onDragOver={(e) => e.preventDefault()}
				onDrop={(e) => {this.onDrop(e)}}
			>
				{items}
			</div>
		);
	}
}

export default GameBoard;
