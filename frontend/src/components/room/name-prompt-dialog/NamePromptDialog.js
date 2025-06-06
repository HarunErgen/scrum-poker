import React, { useState, useEffect } from 'react';
import './NamePromptDialog.css';

const NamePromptDialog = ({
                              onSubmit,
                              onCancel,
                              initialValue = '',
                              title = 'Join Room',
                              buttonText = 'Join Room'
                          }) => {
    const [userName, setUserName] = useState(initialValue);
    const [error, setError] = useState('');

    useEffect(() => {
        setUserName(initialValue);
    }, [initialValue]);

    const handleSubmit = (e) => {
        e.preventDefault();

        if (!userName.trim()) {
            setError('Please enter a name');
            return;
        }

        onSubmit(userName);
    };

    return (
        <div className="name-dialog-overlay">
            <div className="name-dialog card">
                <h3>{title}</h3>

                {error && <div className="error-message">{error}</div>}

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="userName">Name</label>
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
                            {buttonText}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default NamePromptDialog;
