import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Home.css';

import SessionStorage from '../../utils/SessionStorage';
import api from "../../utils/AxiosInstance";

const Home = () => {
  const navigate = useNavigate();
  const [isCreating, setIsCreating] = useState(false);
  const [isJoining, setIsJoining] = useState(false);
  const [roomName, setRoomName] = useState('');
  const [userName, setUserName] = useState('');
  const [roomId, setRoomId] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    const session = SessionStorage.getSession();

    if (SessionStorage.hasSession()) {
      setUserName(session.userName || '');
      setRoomId(session.roomId || '');
    }
  });

  const handleCreateRoom = async (e) => {
    e.preventDefault();
    setError('');

    if (!roomName || !userName) {
      setError('Room name and your name are required');
      return;
    }

    try {
      const response = await api.post('/api/rooms', {
        name: roomName,
        user_name: userName
      });

      const userId = response.data.participants[response.data.scrum_master].id;
      const newRoomId = response.data.id;

      SessionStorage.saveSession(userId, userName, newRoomId);

      navigate(`/room/${newRoomId}`);
    } catch (err) {
      console.error('Error creating room:', err);
      setError('Failed to create room. Please try again.');
    }
  };

  const handleJoinRoom = async (e) => {
    e.preventDefault();
    setError('');

    if (!roomId || !userName) {
      setError('Room ID and your name are required');
      return;
    }

    try {
      const response = await api.post(`/api/rooms/${roomId}/join`, {
        user_name: userName
      });

      const userId = response.data.user.id;

      SessionStorage.saveSession(userId, userName, roomId);

      navigate(`/room/${roomId}`);
    } catch (err) {
      console.error('Error joining room:', err);
      setError('Failed to join room. Please check the Room ID and try again.');
    }
  };

  return (
    <div className="home-container">
      <div className="card">
        <h2 className="text-center">Welcome to Scrum Poker</h2>
        <p className="text-center">A lightweight, real-time Scrum Poker application for Agile teams</p>

        <div className="actions mt-3">
          <button 
            className="btn btn-primary" 
            onClick={() => { setIsCreating(true); setIsJoining(false); }}
          >
            Create a Room
          </button>
          <button 
            className="btn" 
            onClick={() => { setIsJoining(true); setIsCreating(false); }}
            style={{ marginLeft: '1rem' }}
          >
            Join a Room
          </button>
        </div>

        {error && <div className="error-message mt-3">{error}</div>}

        {isCreating && (
          <form onSubmit={handleCreateRoom} className="mt-3">
            <div className="form-group">
              <label htmlFor="roomName">Room Name</label>
              <input
                type="text"
                id="roomName"
                className="form-control"
                value={roomName}
                onChange={(e) => setRoomName(e.target.value)}
                placeholder="e.g., Sprint Planning"
              />
            </div>
            <div className="form-group">
              <label htmlFor="userName">Your Name</label>
              <input
                type="text"
                id="userName"
                className="form-control"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                placeholder="e.g., John Doe"
              />
            </div>
            <button type="submit" className="btn btn-success">
              Create Room
            </button>
          </form>
        )}

        {isJoining && (
          <form onSubmit={handleJoinRoom} className="mt-3">
            <div className="form-group">
              <label htmlFor="roomId">Room ID</label>
              <input
                type="text"
                id="roomId"
                className="form-control"
                value={roomId}
                onChange={(e) => setRoomId(e.target.value)}
                placeholder="Enter the Room ID"
              />
            </div>
            <div className="form-group">
              <label htmlFor="joinUserName">Your Name</label>
              <input
                type="text"
                id="joinUserName"
                className="form-control"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                placeholder="e.g., Jane Smith"
              />
            </div>
            <button type="submit" className="btn btn-success">
              Join Room
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

export default Home;
