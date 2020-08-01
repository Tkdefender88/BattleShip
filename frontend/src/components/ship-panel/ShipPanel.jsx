import React from 'react';
import ShipDisplay from '../ship-display/ShipDisplay';


class ShipPanel extends React.Component {
    render() {
        return (
            <div className="ship-panel">
                <ShipDisplay></ShipDisplay>
            </div>
        );
    }
}

export default ShipPanel;