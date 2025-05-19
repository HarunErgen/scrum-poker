import React, { useState, useEffect, useRef } from 'react';
import {useParams, useNavigate, useLocation} from 'react-router-dom';
import './Room.css';

import VotingCard from './voting-card/VotingCard';
import NamePromptDialog from './name-prompt-dialog/NamePromptDialog';
import ResultsDialog from './results-dialog/ResultsDialog';
import QRCodeDialog from './qr-code-dialog/QRCodeDialog';

import Room from '../../models/Room';
import VoteOption from '../../models/VoteOption';
import WebSocketService from '../../utils/WebSocketService';
import api from "../../utils/AxiosInstance";
import ParticipantMenu from "./participant-menu/ParticipantMenu";
import User from "../../models/User";

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

  const location = useLocation();
  const [userId, setUserId] = useState(location.state && location.state.userId);
  const [userName, setUserName] = useState(location.state && location.state.userName);

  const wsServiceRef = useRef(null);
  const voteOptions = VoteOption.getAllValues();
  const isScrumMaster = roomData && roomData.isScrumMaster(userId);

  const [menuOpenFor, setMenuOpenFor] = useState(null);
  const [showRenameDialog, setShowRenameDialog] = useState(false);
  const [participantToRename, setParticipantToRename] = useState(null);

  useEffect(() => {
    const checkSession = async () => {
      try {
        const response = await api.get('/api/sessions', {
          params: {
            roomId: roomId,
          }
        });

        if (response.data && response.data.user) {
          setUserId(response.data.user.id);
          setUserName(response.data.user.name);
          setShowNamePrompt(false);
          initializeRoom();
          return;
        }
      } catch (err) {
        console.log('No active session found or session expired');
      }

      if (!userId || !userName) {
        setShowNamePrompt(true);
        return;
      }
      setShowNamePrompt(false);
      initializeRoom();
    };

    checkSession();
  }, [userId]);

  useEffect(() => {
    if (roomData && roomData.votesRevealed) {
      setShowResultsDialog(true);
    } else {
      setShowResultsDialog(false);
    }
  }, [roomData]);

  useEffect(() => {
    return () => {
      if (wsServiceRef.current) {
        wsServiceRef.current.disconnect();
      }
    };
  }, []);

  const initializeRoom = async () => {
    try {
      const response = await api.get(`/api/rooms/${roomId}`);
      const room = Room.fromApiResponse(response.data);
      setRoomData(room);
      setLoading(false);

      if (room.hasVoted(userId)) {
        setSelectedVote(room.getUserVote(userId));
      }

      const handleWebSocketMessage = async (data) => {
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
        userName: enteredName
      });

      const newUserId = response.data.user.id;

      setUserName(enteredName);
      setUserId(newUserId);
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
        userId: userId
      });
      await api.delete('/api/sessions');
    } catch (err) {
      setError('Failed to leave room. Please try again.');
    }
    navigate('/');
  };

  const handleVote = async (vote) => {
    if (selectedVote === vote) {
      setSelectedVote('');
      try {
        await api.post(`/api/rooms/${roomId}/vote`, {
          userId: userId,
          vote: ''
        });
      } catch (err) {
        setError('Failed to deselect vote. Please try again.');
      }
    } else {
      setSelectedVote(vote);
      try {
        await api.post(`/api/rooms/${roomId}/vote`, {
          userId: userId,
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
        userId: userId
      });
    } catch (err) {
      setError('Failed to reveal votes. Please try again.');
    }
  };

  const handleResetVotes = async () => {
    try {
      await api.post(`/api/rooms/${roomId}/reset`, {
        userId: userId
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
        userId: userId,
        newScrumMasterId: newScrumMasterId
      });
    } catch (err) {
      setError('Failed to transfer Scrum Master role. Please try again.');
    }
  };

  const handleRenameUser = (participantId) => {
    setParticipantToRename(participantId);
    setShowRenameDialog(true);
  };

  const handleRenameSubmit = async (newName) => {
    try {
      console.log('Renaming user:', participantToRename, newName);
      const user = new User(participantToRename, newName, true);
      await api.put(`/api/users`, user);

      if (participantToRename === userId) {
        setUserName(newName);
      }

      setShowRenameDialog(false);
      setParticipantToRename(null);
    } catch (err) {
      console.error('Error renaming user:', err);
      setError('Failed to rename user. Please try again.');
    }
  };

  const handleRenameCancel = () => {
    setShowRenameDialog(false);
    setParticipantToRename(null);
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

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (showNamePrompt) {
    return <NamePromptDialog onSubmit={handleJoinRoom} onCancel={handleCancelJoin} />;
  }

  if (loading) {
    return <div className="loading">Loading...</div>;
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
              Room Id: {roomId}
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
            {roomData.getParticipantsArray().map((participant) => {
              const isCurrentUser = participant.id === userId;
              const showMenu = isScrumMaster || isCurrentUser;

              const menuOptions = [];

              if (isCurrentUser) {
                menuOptions.push({
                  label: "Rename",
                  action: () => handleRenameUser(participant.id)
                });
              }

              if (isScrumMaster && !isCurrentUser) {
                menuOptions.push({
                  label: "Make Scrum Master",
                  action: () => handleTransferRole(participant.id)
                });
              }

              return (
                  <div
                      key={participant.id}
                      className={`participant ${
                          roomData.hasVoted(participant.id) && !roomData.votesRevealed ? 'voted' : ''
                      }`}
                  >
                    <span
                        className={`participant-name ${
                            !participant.isOnline ? 'offline-indicator' : ''
                        }`}
                    >
                      {participant.name}{' '}
                      {participant.id === roomData.scrumMaster && '(Scrum Master)'}
                      {!participant.isOnline && (
                          <span className="offline-indicator"> (Offline)</span>
                      )}
                    </span>

                    {roomData.hasVoted(participant.id) && roomData.votesRevealed && (
                        <span className="vote-status">
                          Voted: {roomData.getUserVote(participant.id)}
                        </span>
                    )}

                    {showMenu && (
                        <div className="menu-container">
                          <button
                              className="menu-button"
                              onClick={() => setMenuOpenFor(menuOpenFor === participant.id ? null : participant.id)}
                              aria-label="Options"
                          >
                            <span className="menu-dots">⋮</span>
                          </button>

                          <ParticipantMenu
                              isOpen={menuOpenFor === participant.id}
                              onClose={() => setMenuOpenFor(null)}
                              menuOptions={menuOptions}
                          />
                        </div>
                    )}
                  </div>
              );
            })}
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

        {showRenameDialog && (
            <NamePromptDialog
              onSubmit={handleRenameSubmit}
              onCancel={handleRenameCancel}
              initialValue={userName}
              buttonText="Rename"
              title="Change Username"
            />
        )}
      </div>
  );
};

export default RoomComponent;
