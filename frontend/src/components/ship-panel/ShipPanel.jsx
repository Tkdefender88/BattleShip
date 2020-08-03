import React from 'react';
import ShipDisplay from '../ship-display/ShipDisplay';

class ShipPanel extends React.Component {
    render() {

        const ships = ['carrier', 'battleship', 'cruiser', 'submarine', 'destroyer'];

        const shipDisplays = ships.map((ship, i) => 
            <ShipDisplay key={i} ship={ship}></ShipDisplay>
        );

        return (
            <div className="ship-panel">
               {shipDisplays} 
            </div>
        );
    }
}

export default ShipPanel;