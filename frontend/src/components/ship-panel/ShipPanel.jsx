import React from 'react';
import ShipDisplay from '../ship-display/ShipDisplay';
import style from './ship-panel.module.css';


class ShipPanel extends React.Component {
    render() {
        return (
            <div className={style.shippanel}>
                <ShipDisplay></ShipDisplay>
            </div>
        );
    }
}

export default ShipPanel;