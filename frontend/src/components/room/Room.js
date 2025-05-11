import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './Room.css';

import VotingCard from './voting-card/VotingCard';
import NamePromptDialog from './name-prompt-dialog/NamePromptDialog';
import ResultsDialog from './results-dialog/ResultsDialog';
import QRCodeDialog from './qr-code-dialog/QRCodeDialog';

import Room from '../../models/Room';
import VoteOption from '../../models/VoteOption';
import SessionStorage from '../../utils/SessionStorage';
import WebSocketService from '../../utils/WebSocketService';
import api from "../../utils/AxiosInstance";

const RoomComponent = () => {
  const { roomId } = useParams();
  const navigate = useNavigate();
  const [roomData, setRoomData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedVote, setSelectedVote] = useState('');
  const [copied, setCopied] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState('Connecting...');
  const [showNamePrompt, setShowNamePrompt] = useState(false);
  const [showResultsDialog, setShowResultsDialog] = useState(false);
  const [showQRCodeDialog, setShowQRCodeDialog] = useState(false);

  const sessionData = SessionStorage.getSession();
  const userId = sessionData.userId;
  const userName = sessionData.userName;

  const wsServiceRef = useRef(null);
  const voteOptions = VoteOption.getAllValues();
  const isScrumMaster = roomData && roomData.isScrumMaster(userId);

  useEffect(() => {
    if (!userId || !userName) {
      setShowNamePrompt(true);
      return;
    }

    initializeRoom();
  }, [userId, userName]);

  useEffect(() => {
    if (roomData && roomData.votesRevealed) {
      setShowResultsDialog(true);
    } else {
      setShowResultsDialog(false);
    }
  }, [roomData]);

  const initializeRoom = async () => {
    SessionStorage.saveSession(userId, userName, roomId);

    try {
      const response = await api.get(`/api/rooms/${roomId}`);
      const room = Room.fromApiResponse(response.data);
      setRoomData(room);
      setLoading(false);

      if (room.hasVoted(userId)) {
        setSelectedVote(room.getUserVote(userId));
      }

      const handleWebSocketMessage = async (data) => {
        console.log('handleWebSocketMessage', data);
        if (data.type === 'room_update') {
          const room = Room.fromApiResponse(data.payload);
          if (room.votesRevealed) {
            setSelectedVote('');
          }
          setRoomData(room);
        }
      };

      wsServiceRef.current = new WebSocketService(
          roomId,
          userId,
          handleWebSocketMessage,
          setConnectionStatus
      );

      wsServiceRef.current.connect();
    } catch (err) {
      setError('Failed to load room. Please try again.');
      setLoading(false);
    }
  };

  const handleJoinRoom = async (enteredName) => {
    try {
      const response = await api.post(`/api/rooms/${roomId}/join`, {
        user_name: enteredName
      });

      const newUserId = response.data.user.id;

      SessionStorage.saveSession(newUserId, enteredName, roomId);

      window.location.reload();
    } catch (err) {
      console.error('Error joining room:', err);
      setError('Failed to join room. Please try again.');
    }
  };

  const handleCancelJoin = () => {
    navigate('/');
  };

  const handleLeaveRoom = async () => {
    try {
      await api.post(`/api/rooms/${roomId}/leave`, {
        user_id: userId
      });
    } catch (err) {
      setError('Failed to leave room. Please try again.');
    }

    SessionStorage.clearSession();
    navigate('/');
  };

  const handleVote = async (vote) => {
    if (selectedVote === vote) {
      setSelectedVote('');
      try {
        await api.post(`/api/rooms/${roomId}/vote`, {
          user_id: userId,
          vote: ''
        });
      } catch (err) {
        setError('Failed to deselect vote. Please try again.');
      }
    } else {
      setSelectedVote(vote);
      try {
        await api.post(`/api/rooms/${roomId}/vote`, {
          user_id: userId,
          vote: vote
        });
      } catch (err) {
        setError('Failed to submit vote. Please try again.');
      }
    }
  };

  const handleRevealVotes = async () => {
    try {
      await api.post(`/api/rooms/${roomId}/reveal`, {
        user_id: userId
      });
    } catch (err) {
      setError('Failed to reveal votes. Please try again.');
    }
  };

  const handleResetVotes = async () => {
    try {
      await api.post(`/api/rooms/${roomId}/reset`, {
        user_id: userId
      });
      setSelectedVote('');
      setShowResultsDialog(false);
    } catch (err) {
      setError('Failed to reset votes. Please try again.');
    }
  };

  const handleTransferRole = async (newScrumMasterId) => {
    try {
      await api.post(`/api/rooms/${roomId}/transfer`, {
        user_id: userId,
        new_scrum_master_id: newScrumMasterId
      });
    } catch (err) {
      setError('Failed to transfer Scrum Master role. Please try again.');
    }
  };

  const getRoomLink = () => {
    return `${window.location.origin}/room/${roomId}`;
  };

  const copyRoomLink = () => {
    const roomLink = getRoomLink();
    navigator.clipboard.writeText(roomLink);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const showQRCode = () => {
    setShowQRCodeDialog(true);
  };

  const hideQRCode = () => {
    setShowQRCodeDialog(false);
  };

  if (showNamePrompt) {
    return <NamePromptDialog onSubmit={handleJoinRoom} onCancel={handleCancelJoin} />;
  }

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (!roomData) {
    return <div className="error">Room not found</div>;
  }

  return (
      <div className="room-container">
        {showResultsDialog && (
            <ResultsDialog
                votes={roomData.votes}
                voteOptions={voteOptions}
                isScrumMaster={isScrumMaster}
                onReset={handleResetVotes}
            />
        )}

        {showQRCodeDialog && (
            <QRCodeDialog
                url={getRoomLink()}
                onClose={hideQRCode}
            />
        )}

        <div className="room-header card">
          <div className="room-info">
            <h2>{roomData.name}</h2>
            <div className="room-id">
              Room ID: {roomId}
              <button className="btn btn-sm" onClick={copyRoomLink}>
                {copied ? 'Copied!' : 'Copy Link'}
              </button>
              <button className="btn btn-sm" onClick={showQRCode} style={{ marginLeft: '0.5rem' }}>
                Show QR Code
              </button>
            </div>
            <div className="connection-status">
              Status: {connectionStatus}
            </div>
          </div>
          <div className="room-actions">
            <button
                className="btn btn-danger"
                onClick={handleLeaveRoom}
            >
              Leave Room
            </button>
            {isScrumMaster && (
                <button
                    className="btn btn-primary"
                    onClick={handleRevealVotes}
                    disabled={roomData.votesRevealed}
                    style={{ marginLeft: '0.5rem' }}
                >
                  Reveal Votes
                </button>
            )}
          </div>
        </div>

        <div className="participants-section card">
          <h3>Participants</h3>
          <div className="participants-list">
            {roomData.getParticipantsArray().map((participant) => (
                  <div key={participant.id} className="participant">
                    <span className="participant-name">
                      {participant.name} {participant.id === roomData.scrumMaster && '(Scrum Master)'}
                    </span>
                    {roomData.hasVoted(participant.id) && (
                        <span className="vote-status">
                          {roomData.votesRevealed ? `Voted: ${roomData.getUserVote(participant.id)}` : 'Voted'}
                        </span>
                    )}
                    {isScrumMaster && participant.id !== userId && (
                        <button
                            className="btn btn-sm"
                            onClick={() => handleTransferRole(participant.id)}
                        >
                          Make Scrum Master
                        </button>
                    )}
                  </div>
            ))}
          </div>
        </div>

        <div className="voting-section card">
          <h3>Your Vote</h3>
          <div className="voting-cards">
            {voteOptions.map((vote) => (
                <VotingCard
                    key={vote}
                    value={vote}
                    selected={selectedVote === vote}
                    disabled={roomData.votesRevealed}
                    onClick={() => handleVote(vote)}
                />
            ))}
          </div>
        </div>
      </div>
  );
};

export default RoomComponent;
