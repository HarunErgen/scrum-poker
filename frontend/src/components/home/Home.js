import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './Home.css';

import api from "../../utils/AxiosInstance";

const Home = () => {
  const navigate = useNavigate();
  const [isCreating, setIsCreating] = useState(false);
  const [isJoining, setIsJoining] = useState(false);
  const [roomName, setRoomName] = useState('');
  const [userName, setUserName] = useState('');
  const [roomId, setRoomId] = useState('');
  const [error, setError] = useState('');

  const handleCreateRoom = async (e) => {
    e.preventDefault();
    setError('');

    if (!roomName || !userName) {
      setError('Room name and your name are required');
      return;
    }

    try {
      const response = await api.post('/rooms', {
        name: roomName,
        userName: userName
      });

      const userId = response.data.participants[response.data.scrumMaster].id;
      const newRoomId = response.data.id;

      try {
        await api.post(`/sessions/${userId}/${newRoomId}`);
      } catch (sessionErr) {
        console.error('Error creating session:', sessionErr);
      }

      navigate(`/room/${newRoomId}`, { state: { userId: userId, userName: userName } });
    } catch (err) {
      console.error('Error creating room:', err);
      setError('Failed to create room. Please try again.');
    }
  };

  const handleJoinRoom = async (e) => {
    e.preventDefault();
    setError('');

    if (!roomId || !userName) {
      setError('Room Id and your name are required');
      return;
    }

    try {
      const response = await api.post(`/rooms/${roomId}/join`, {
        userName: userName
      });

      const userId = response.data.user.id;

      try {
        await api.post(`/sessions/${userId}/${roomId}`);
      } catch (sessionErr) {
        console.error('Error creating session:', sessionErr);
      }

      navigate(`/room/${roomId}`, { state: { userId: userId, userName: userName } });
    } catch (err) {
      console.error('Error joining room:', err);
      setError('Failed to join room. Please check the Room Id and try again.');
    }
  };

  return (
    <div className="home-container">
      <div className="card">
        <h2 className="text-center">Welcome to Scrum Poker</h2>
        <p className="text-center">Estimate stories faster. Collaborate better. Stay Agile.</p>
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
              <label htmlFor="roomId">Room Id</label>
              <input
                type="text"
                id="roomId"
                className="form-control"
                value={roomId}
                onChange={(e) => setRoomId(e.target.value)}
                placeholder="Enter the Room Id"
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
