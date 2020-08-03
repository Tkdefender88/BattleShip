import React from 'react';
import { NavBar, NavItem } from './components/navbar/navbar';
import DropdownMenu from './components/dropdown/Dropdown.jsx';
import ShipPanel from './components/ship-panel/ShipPanel.jsx';
import GameBoard from './components/game-board/game-board.jsx';

import {ReactComponent as CaretIcon} from './icons/caret.svg';

import './styles/App.scss';

function App() {
  return (
    <div className="battle-ship">
      <NavBar>
        <div className="title">
          BattleShip
        </div>
        <NavItem icon={<CaretIcon />}>
          <DropdownMenu/>
        </NavItem>
      </NavBar>
      <ShipPanel/>
      <div className="grid-container">
        <GameBoard/>
        <GameBoard/>
      </div>

    </div>
  );
}

export default App;
