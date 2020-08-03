import React, { Component } from 'react';
import ShipDisplay from '../ship-display/ShipDisplay';

class ShipPanel extends Component {
    render() {

        const ships = ['carrier', 'battleship', 'cruiser', 'submarine', 'destroyer'];

        const shipDisplays = ships.map((ship) => 
            <ShipDisplay ship={ship}></ShipDisplay>
        );

        return (
            <div className="ship-panel">
               {shipDisplays} 
            </div>
        );
    }
}

export default ShipPanel;