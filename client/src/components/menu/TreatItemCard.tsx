// src/components/TreatItemCard.tsx
import React from 'react';
import './TreatItemCard.css';

interface TreatItemCardProps {
  name: string;
  cost: number;
  imageSrc: string;
  onClick: () => void;
}

const TreatItemCard: React.FC<TreatItemCardProps> = ({ name, cost, imageSrc, onClick }) => {
  return (
    <div className="treat-item-card" onClick={onClick}>
      <div className="card-cost">{cost}</div>
      <img src={imageSrc} alt={name} className="card-image" />
      {/* <div className="card-name">{name}</div> */}
    </div>
  );
};

export default TreatItemCard;
