import React, { Component } from 'react';
import ToggleButton from 'react-bootstrap/ToggleButton';
import ToggleButtonGroup from 'react-bootstrap/ToggleButtonGroup';
import style from './ship-display.module.css';
import ship from '../../images/carrier.png';

class ShipDisplay extends Component {
    render() {
        return (
            <div className={style.ship_display}>
                <div>
                    <img src={ship} alt="ship"></img>
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