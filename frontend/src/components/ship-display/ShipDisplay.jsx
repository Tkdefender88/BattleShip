import React, { Component } from 'react';
import ToggleButton from 'react-bootstrap/ToggleButton';
import ToggleButtonGroup from 'react-bootstrap/ToggleButtonGroup';

class ShipDisplay extends Component {
    render() {
        return (
            <div>
                <ToggleButtonGroup type="radio" name="options" defaultValue={1}>
                    <ToggleButton value={1}>Vertical</ToggleButton>
                    <ToggleButton value={2}>Horizontal</ToggleButton>
                </ToggleButtonGroup>
            </div>
        )
    }
}

export default ShipDisplay;