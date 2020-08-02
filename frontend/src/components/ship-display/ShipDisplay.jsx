import React, { Component } from 'react';
import ship from '../../images/carrier.png';

class ShipDisplay extends Component {
    render() {
        return (
            <div className="ship-display">
                <div>
                    <img src={ship} alt="ship"></img>
                </div>
            </div>
        )
    }
}

export default ShipDisplay;