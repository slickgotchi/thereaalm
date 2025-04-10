// src/components/TreatModal.tsx
import React, { useState, useEffect } from 'react';
import './ESPModal.css';
import { GotchiEntity } from '../../phaser/entities/GotchiEntity';
import StatBar from './StatBar';
import { eventBus } from '../../utils/EventBus';

interface ESPModalProps {
  entity: { type: string; [key: string]: any } | null;
}

const ESPModal: React.FC<ESPModalProps> = ({ entity }) => {
  const [ecto, setEcto] = useState(0);
  const [spark, setSpark] = useState(0);
  const [pulse, setPulse] = useState(0);

  // Fetch Gotchi data on mount and periodically
  useEffect(() => {
    const fetchGotchiData = async () => {
      const gotchiId = entity?.data?.gotchiId;
      const gotchiEntity = GotchiEntity.activeGotchis.get(gotchiId);

      setEcto(gotchiEntity?.data?.stats?.ecto || 0);
      setSpark(gotchiEntity?.data?.stats?.spark || 0);
      setPulse(gotchiEntity?.data?.stats?.pulse || 0);
    };

    fetchGotchiData();
    const interval = setInterval(fetchGotchiData, 5000); // Poll every 5 seconds

    const handleTreatEaten = (data: { detail: any }) => {
      setEcto(data.detail.ecto);
      setSpark(data.detail.spark);
      setPulse(data.detail.pulse);
    }

    eventBus.on("treatEaten", handleTreatEaten);

    return () => {
      clearInterval(interval);
      eventBus.off("treatEaten", handleTreatEaten);
    }
  }, []);

  // Update treatTotal and stakedGhst when entity changes
  useEffect(() => {
    if (entity) {
      setEcto(entity?.data?.stats?.ecto || 0);
      setSpark(entity?.data?.stats?.spark || 0);
      setPulse(entity?.data?.stats?.pulse || 0);
    }
  }, [entity]);

  return (
    <div className="esp-modal">
      <div className="header-bar">
        <div className="header-content">ESP</div>
      </div>

      <div className="statbar-content">
        <StatBar label="Ecto" value={ecto} max={1000} color="#CA52C9" />
        <StatBar label="Spark" value={spark} max={1000} color="#0098DC" />
        <StatBar label="Pulse" value={pulse} max={1000} color="#F5555D" />
      </div>
    </div>
  );
};

export default ESPModal;