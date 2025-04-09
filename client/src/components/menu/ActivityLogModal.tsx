// src/components/ActivityLogModal.tsx
import React from 'react';
import './ActivityLogModal.css';

interface ActivityLogModalProps {
  entity: { type: string; [key: string]: any };
}

const ActivityLogModal: React.FC<ActivityLogModalProps> = ({ entity }) => {
  return (
    <div className="activity-log-modal">
      <div className="header-bar">
        <div className="header-content">Activity Log</div>
      </div>
      <div className="main-content">
        <p>Activity Log for {entity.type}</p>
        {/* Add activity log details here */}
      </div>
    </div>
  );
};

export default ActivityLogModal;