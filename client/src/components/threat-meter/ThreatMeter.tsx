import React from 'react';
import './ThreatMeter.css';

// Assuming the shield image is in your public folder as 'shield.png'
const SHIELD_IMAGE = '/assets/images/shield.png';

interface ThreatMeterProps {
  ThreatLevel: number; // Value between 0 and 100
}

const ThreatMeter: React.FC<ThreatMeterProps> = ({ ThreatLevel }) => {
  // Clamp ThreatLevel between 0 and 100
  const clampedThreatLevel = Math.min(100, Math.max(0, ThreatLevel));
  
  // Calculate the percentage of the bar that should be blue (Gotchis)
  const blueWidth = 100 - clampedThreatLevel;
  
  // The shield position will be based on the ThreatLevel (0% to 100%)
  const shieldPosition = clampedThreatLevel;

  return (
    <div className="threat-meter-container">
      <div className="threat-meter-labels">
        <span className="label-gotchis">Gotchis</span>
        <span className="label-lickquidators">Lickquidators</span>
      </div>
      <div className="threat-meter">
        <div className="threat-meter-inner">
          <div
            className="threat-meter-blue"
            style={{ width: `${blueWidth}%` }}
          />
          <div
            className="threat-meter-red"
            style={{ width: `${clampedThreatLevel}%` }}
          />
        </div>
        <img
          src={SHIELD_IMAGE}
          alt="Shield"
          className="threat-meter-shield"
          style={{ left: `${100-shieldPosition}%` }}
        />
      </div>
    </div>
  );
};

export default ThreatMeter;