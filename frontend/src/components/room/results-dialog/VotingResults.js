import React, {useMemo} from 'react';
import './VotingResults.css';
import VoteOption from '../../../models/VoteOption';

const nearestFibonacci = (value) => {
    return VoteOption.getAllNumericValues().reduce((prev, curr) =>
        Math.abs(curr - value) < Math.abs(prev - value) ? curr : prev
    );
};

const VotingResults = ({ votes, voteOptions }) => {
    const { average, roundedAverage } = useMemo(() => {
        const numericVotes = Object.values(votes)
            .map((v) => Number(v))
            .filter((v) => !Number.isNaN(v));

        if (numericVotes.length === 0) return { average: null, roundedAverage: null };

        const avg = numericVotes.reduce((sum, v) => sum + v, 0) / numericVotes.length;

        return { average: avg, roundedAverage: nearestFibonacci(avg) };
    }, [votes]);

    return (
        <div className="voting-results">
            {average !== null && (
                <div className="average-wrapper">
                    <span className="label">Average&nbsp;</span>
                    <span className="average-value">{average.toFixed(2)}</span>
                    <span className="average-rounded">(~&nbsp;{roundedAverage})</span>
                </div>
            )}

            <div className="results">
                {Object.keys(votes).length === 0 ? (
                    <p>No votes yet</p>
                ) : (
                    <div className="vote-summary">
                        {voteOptions.map((vote) => {
                            const count = Object.values(votes).filter(v => v === vote).length;
                            if (count === 0) return null;

                            const voteOption = VoteOption.getByValue(vote);
                            const color = voteOption.getColor();

                            return (
                                <div key={vote} className="vote-count" style={{ borderLeft: `4px solid ${color}` }}>
                                    <span className="vote" style={{ color }}>
                                      {vote}
                                    </span>
                                    <span className="count">{count}</span>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
};

export default VotingResults;
