import React, { Component } from 'react';
import ToggleButton from 'react-bootstrap/ToggleButton';
import ToggleButtonGroup from 'react-bootstrap/ToggleButtonGroup';

class ShipDisplay extends Component {
    render() {
        return (
            <div className="ship-display">
                <div>
                    <img src="images/battleship.png"></img>
                </div>
                <ToggleButtonGroup type="radio" name="options" defaultValue={1}>
                    <ToggleButton value={1}>Horizontal</ToggleButton>
                    <ToggleButton value={2}>Vertical</ToggleButton>
                </ToggleButtonGroup>
            </div>
        )
    }
}

export default ShipDisplay;