// src/components/InformationModal.tsx
import React from 'react';
import './InformationModal.css';

interface InformationModalProps {
  entity: { type: string; [key: string]: any };
}

const InformationModal: React.FC<InformationModalProps> = ({ entity }) => {
  return (
    <div className="information-modal">
      <div className="header-bar">
        <div className="header-content">Information</div>
      </div>
      <div className="main-content">
        <p>Entity Type: {entity.type}</p>
        {/* Add more entity details here */}
      </div>
    </div>
  );
};

export default InformationModal;