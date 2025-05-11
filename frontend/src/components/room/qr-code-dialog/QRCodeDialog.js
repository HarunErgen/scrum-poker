import React from 'react';
import { QRCodeSVG } from 'qrcode.react';
import './QRCodeDialog.css';

const QRCodeDialog = ({ url, onClose }) => {
  return (
    <div className="qr-code-dialog-overlay">
      <div className="qr-code-dialog">
        <div className="qr-code-dialog-header">
          <h3>Scan QR Code to Join Room</h3>
          <button className="close-button" onClick={onClose}>×</button>
        </div>
        <div className="qr-code-container">
          <QRCodeSVG value={url} size={250} />
        </div>
        <div className="qr-code-url">
          <p>{url}</p>
        </div>
      </div>
    </div>
  );
};

export default QRCodeDialog;