import React, { Component } from 'react';

import Toggle from '../toggle-switch/toggle-switch.jsx'
import carrier from '../../images/carrier.png';
import battleship from '../../images/battleship.png';
import cruiser from '../../images/cruiser.png';
import submarine from '../../images/submarine.png';
import destroyer from '../../images/destroyer.png';

class ShipDisplay extends Component {

    shipImage(shipName) {
        let shipImages = {
            "carrier" : carrier,
            "battleship" : battleship,
            "cruiser" : cruiser,
            "submarine" : submarine,
            "destroyer" : destroyer,
        }
        return shipImages[shipName];
    }

    onDragStart(e, ship) {
        console.log("dragstart: " + ship)
        e.dataTransfer.setData("ship", ship)
    }

    render() {
        return (
            <div className="ship-display" key={this.props.key}>
                <div key={this.props.ship} draggable onDragStart={(e) => {this.onDragStart(e, this.props.ship)}}>
                    <img src={this.shipImage(this.props.ship)} alt="ship"></img>
                </div>
                <Toggle/>
            </div>
        )
    }
}

export default ShipDisplay;