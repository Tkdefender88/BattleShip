import React from 'react';
import NavBar from './components/navbar/navbar';
import ShipPanel from './components/ship-panel/ShipPanel';
import GameGrid from './components/game-board/game-board';


import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

function App() {
  return (
    <div className="App">
      <NavBar className="header"/>
      <ShipPanel className="ship-panel"/>
      <div className="grid-container">
        <GameGrid />
        <GameGrid />
      </div>
    </div>
  );
}

export default App;
