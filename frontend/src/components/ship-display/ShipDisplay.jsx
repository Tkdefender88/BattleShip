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

    render() {
        return (
            <div className="ship-display" key={this.props.key}>
                <div draggable onDragStart={(e) => {}}>
                    <img src={this.shipImage(this.props.ship)} alt="ship"></img>
                </div>
                <Toggle/>
            </div>
        )
    }
}

export default ShipDisplay;