import React from 'react';
import './InformationModal.css';
import { Entity } from './MenuSystem';
import { GotchiDetails } from './information-details/GotchiDetails';



interface InformationModalProps {
  entity: Entity;
}



const ItemDetails: React.FC<{ entity: Entity }> = ({ entity }) => (
  <div>
    <p>Entity Type: {entity.type}</p>
    <p>Item Name: {entity.name || 'Unknown Item'}</p>
    <p>Rarity: {entity.rarity || 'Common'}</p>
  </div>
);

const LocationDetails: React.FC<{ entity: Entity }> = ({ entity }) => (
  <div>
    <p>Entity Type: {entity.type}</p>
    <p>Location: {entity.name || 'Unknown Location'}</p>
    <p>Coordinates: ({entity.coords?.x || 0}, {entity.coords?.y || 0})</p>
  </div>
);

const DefaultDetails: React.FC<{ entity: Entity }> = ({ entity }) => (
  <div>
    <p>Entity Type: {entity.type}</p>
    <p>No specific details available for this entity.</p>
  </div>
);

// Component registry mapping entity types to components
const componentRegistry: Record<string, React.FC<{ entity: Entity }>> = {
  gotchi: GotchiDetails,
  item: ItemDetails,
  location: LocationDetails,
  unknown: DefaultDetails,
};

const InformationModal: React.FC<InformationModalProps> = ({ entity }) => {
  // Determine the component to render, default to 'unknown' if type is not recognized
  const EntityComponent = componentRegistry[entity.type] || componentRegistry.unknown;

  return (
    <div className="information-modal">
      <div className="header-bar">
        <div className="header-content">Information</div>
      </div>
      <div className="main-content">
        <EntityComponent entity={entity} />
      </div>
    </div>
  );
};

export default InformationModal;