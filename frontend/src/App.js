import React from 'react';
import NavBar from './components/navbar/navbar';
import ShipPanel from './components/ship-panel/ShipPanel';


import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

function App() {
  return (
    <div className="App">
      <NavBar className="header"/>
      <ShipPanel className="ship-panel"/>
    </div>
  );
}

export default App;
