// src/components/ActivityLogModal.tsx
import React from 'react';
import './JobModal.css';

interface JobModalProps {
  entity: { type: string; [key: string]: any };
}

const JobModal: React.FC<JobModalProps> = ({ entity }) => {

  var jobStr = entity.data.job;
  jobStr = jobStr[0].toUpperCase() + jobStr.slice(1);

  return (
    <div className="job-modal">
      <div className="header-bar">
        <div className="header-content">Job</div>
      </div>
      <div className="main-content">
        <p><b>Job:</b> {jobStr}</p>
        {/* Add activity log details here */}
      </div>
    </div>
  );
};

export default JobModal;