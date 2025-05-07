import React, { useState } from 'react';
import './NamePromptDialog.css';

const NamePromptDialog = ({ onSubmit, onCancel }) => {
    const [userName, setUserName] = useState('');
    const [error, setError] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();

        if (!userName.trim()) {
            setError('Enter your name');
            return;
        }

        onSubmit(userName);
    };

    return (
        <div className="name-dialog-overlay">
            <div className="name-dialog card">
                <h3>Join Room</h3>

                {error && <div className="error-message">{error}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="userName">Your Name</label>
                        <input
                            type="text"
                            id="userName"
                            className="form-control"
                            value={userName}
                            onChange={(e) => setUserName(e.target.value)}
                            placeholder="e.g., John Doe"
                            autoFocus
                        />
                    </div>

                    <div className="dialog-actions">
                        <button type="button" className="btn" onClick={onCancel}>
                            Cancel
                        </button>
                        <button type="submit" className="btn btn-primary">
                            Join Room
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default NamePromptDialog;
