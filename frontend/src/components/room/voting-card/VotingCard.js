import React from 'react';
import VoteOption from '../../../models/VoteOption';

const VotingCard = ({ value, selected, disabled, onClick }) => {
  const voteOption = VoteOption.getByValue(value);
  const color = voteOption.getColor();

  const cardStyle = {
    borderColor: color,
    ...(selected && { backgroundColor: color })
  };

  return (
    <div 
      className={`voting-card ${selected ? 'selected' : ''} ${disabled && !selected ? 'disabled' : ''}`}
      onClick={!disabled ? onClick : undefined}
      style={cardStyle}
    >
      <div className="card-content">
        <span className="card-value">{value}</span>
      </div>
    </div>
  );
};

export default VotingCard;
