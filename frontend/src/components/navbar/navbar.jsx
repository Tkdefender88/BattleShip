import React from 'react';
import Dropdown from 'react-bootstrap/Dropdown';
import './navbar.css';


const NavBar = () => {
    return (
        <div className="header">
            <div className="title">BattleShip</div>
            <Dropdown>
                <Dropdown.Toggle variant="primary">
                    Game
                </Dropdown.Toggle>
                <Dropdown.Menu>
                    <Dropdown.Item>Load</Dropdown.Item>
                    <Dropdown.Item>Save</Dropdown.Item>
                </Dropdown.Menu>
            </Dropdown>

            <Dropdown>
                <Dropdown.Toggle variant="secondary">
                    Help
                </Dropdown.Toggle>
                <Dropdown.Menu>
                    <Dropdown.Item>About</Dropdown.Item>
                    <Dropdown.Item>Rules</Dropdown.Item>
                </Dropdown.Menu>
            </Dropdown>

        </div>
    )
}

export default NavBar;