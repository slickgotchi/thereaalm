// src/components/TreatItemCard.tsx
import React from 'react';
import './TreatDetails.css';

interface TreatDetailsProps {
  imageSrc: string;
  name: string;
  description: string;
  cost: number;
  onEat: () => void;
  canEat: boolean;
}

const TreatDetails: React.FC<TreatDetailsProps> = ({ name, cost, imageSrc, description, onEat, canEat }) => {
  return (
    <div className="item-content" onClick={onEat}>
      <img src={imageSrc} alt={name} className="item-image" />
      <div className="item-name">{name}</div>
      <div className="item-description">{description}</div>
      <div className="item-cost">Cost: {cost} TREAT</div>
      <button className='eat-button' disabled={!canEat}>EAT</button>
    </div>
  );
};

export default TreatDetails;
