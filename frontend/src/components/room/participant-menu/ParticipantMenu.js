import React, { useEffect, useRef } from 'react';
import './ParticipantMenu.css';

const ParticipantMenu = ({ isOpen, onClose, menuOptions }) => {
    const menuRef = useRef(null);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (menuRef.current && !menuRef.current.contains(event.target)) {
                onClose();
            }
        };

        if (isOpen) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [isOpen, onClose]);

    if (!isOpen) return null;

    return (
        <div className="participant-menu" ref={menuRef}>
            {menuOptions.map((option) => (
                <button
                    key={option.label}
                    className="menu-item"
                    onClick={() => {
                        option.action();
                        onClose();
                    }}
                >
                    {option.label}
                </button>
            ))}
        </div>
    );
};

export default ParticipantMenu;