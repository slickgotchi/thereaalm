// src/components/StatBar.tsx
import React from 'react';
import './StatBar.css';

interface StatBarProps {
  label: string;
  value: number;
  max: number;
  color: string;
}

const StatBar: React.FC<StatBarProps> = ({ label, value, max, color }) => {
  const percent = Math.min(100, (value / max) * 100);

//   const borderColor = percent < 10 ? "#F5555D" : "transparent";
  const isLow = percent < 10;

  return (
    <div
        className={`stat-bar ${isLow ? 'flash' : ''}`}
    >
      <div
        className="stat-bar-fill"
        style={{ width: `${percent}%`, backgroundColor: color }}
      />
      <div className="stat-bar-label">
        {label.toUpperCase()}: {value.toFixed(0)}/{max.toFixed(0)}
      </div>
    </div>
  );
};

export default StatBar;
