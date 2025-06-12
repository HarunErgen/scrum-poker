import React from 'react';
import {BrowserRouter as Router, Routes, Route, Link} from 'react-router-dom';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './App.css';
import logo from './assets/images/logo-s.png';

import Home from './components/home/Home';
import RoomComponent from './components/room/Room';
import NotFound from './components/not-found/NotFound';

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            <img src={logo} alt="Scrum Poker" className="header-logo" />
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
          <p>&copy; {new Date().getFullYear()} Scrum Poker. Created by Harun Reşid Ergen. All rights reserved.</p>
        </footer>
        <ToastContainer position="top-right" autoClose={2000} hideProgressBar={true} newestOnTop closeOnClick pauseOnFocusLoss draggable pauseOnHover />
      </div>
    </Router>
  );
}

export default App;
