// src/components/ActionPlanModal.tsx
import React from 'react';
import './ActionPlanModal.css';

interface ActionPlanModalProps {
  entity: { type: string; [key: string]: any };
}

const ActionPlanModal: React.FC<ActionPlanModalProps> = ({ entity }) => {
  return (
    <div className="action-plan-modal">
      <div className="header-bar">
        <div className="header-content">Action Plan</div>
      </div>
      <div className="main-content">
        <p>Action Plan for {entity.type}</p>
        {/* Add action plan details here */}
      </div>
    </div>
  );
};

export default ActionPlanModal;