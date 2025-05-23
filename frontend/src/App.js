import React from 'react';
import {BrowserRouter as Router, Routes, Route, Link} from 'react-router-dom';
import './App.css';

import Home from './components/home/Home';
import RoomComponent from './components/room/Room';
import NotFound from './components/not-found/NotFound';

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            <h1>Scrum Poker</h1>
          </Link>
        </header>
        <main className="App-main">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/room/:roomId" element={<RoomComponent />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </main>
        <footer className="App-footer">
          <p>&copy; {new Date().getFullYear()} Scrum Poker</p>
        </footer>
      </div>
    </Router>
  );
}

export default App;
