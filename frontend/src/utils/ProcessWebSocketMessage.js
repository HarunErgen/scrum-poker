import User from "../models/User";
import ActionTypes from "../models/ActionTypes";
import Room from "../models/Room";
import { toast } from 'react-toastify';

export function processWebSocketMessage(message, roomData, setRoomData, setSelectedVote) {
    const { action, payload } = message;

    const actionHandlers = {
        [ActionTypes.JOIN]: () => {
            roomData.participants[payload.id] = new User(payload.id, payload.name, payload.isOnline);
            toast.info(`${payload.name} joined the room`);
        },

        [ActionTypes.OFFLINE]: () => {
            roomData.participants[payload.userId].isOnline = false;
        },

        [ActionTypes.ONLINE]: () => {
            roomData.participants[payload.userId].isOnline = true;
        },

        [ActionTypes.LEAVE]: () => {
            const userName = roomData.participants[payload.userId]?.name || 'Someone';
            delete roomData.participants[payload.userId];
            toast.info(`${userName} left the room`);
        },

        [ActionTypes.RENAME]: () => {
            roomData.participants[payload.userId].name = payload.name;
        },

        [ActionTypes.SUBMIT]: () => {
            roomData.votes[payload.userId] = payload.vote;
        },

        [ActionTypes.REVEAL]: () => {
            roomData.votes = payload.votes;
            roomData.votesRevealed = true;
        },

        [ActionTypes.RESET]: () => {
            roomData.votes = {};
            roomData.votesRevealed = false;
            setSelectedVote('');
        },

        [ActionTypes.TRANSFER]: () => {
            roomData.scrumMaster = payload.newScrumMasterId;
        }
    };

    if (actionHandlers[action]) {
        actionHandlers[action]();
        const newRoomData = new Room(
            roomData.id,
            roomData.name,
            roomData.scrumMaster,
            roomData.participants,
            roomData.votes,
            roomData.votesRevealed
        );
        setRoomData(newRoomData);
    } else {
        console.warn(`Unhandled WebSocket action: ${action}`);
    }
}
