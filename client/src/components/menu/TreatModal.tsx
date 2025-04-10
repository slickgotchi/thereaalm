// src/components/TreatModal.tsx
import React, { useState, useEffect } from 'react';
import TreatItemCard from './TreatItemCard';
import TreatDetails from './TreatDetails';
import './TreatModal.css';
import { GotchiEntity } from '../../phaser/entities/GotchiEntity';
import { eventBus } from '../../utils/EventBus';

interface TreatModalProps {
  entity: { type: string; [key: string]: any } | null;
}

interface TreatItem {
  name: string;
  description: string;
  cost: number;
  imageSrc: string;
}

interface GotchiData {
  treatTotal: number;
  stakedGhst: number;
}

const API_BASE_URL = 'http://localhost:8080'; // Adjust as needed

const TreatModal: React.FC<TreatModalProps> = ({ entity }) => {
  const items: TreatItem[] = [
    { name: 'Sushi Roll', description: 'Restores a little Ecto', cost: 500, imageSrc: '/assets/images/82_SushiRoll.png' },
    { name: 'Coconut', description: 'Restores a little Spark', cost: 500, imageSrc: '/assets/images/116_Coconut.png' },
    { name: 'Candy', description: 'Restores a little Pulse', cost: 500, imageSrc: '/assets/images/251_CoinGeckoCandies.png' },
  ];

  const [selectedItem, setSelectedItem] = useState<TreatItem | null>(items[0]);
  const [stakeAmount, setStakeAmount] = useState<string>('');
  const [unstakeAmount, setUnstakeAmount] = useState<string>('');
  const [treatTotal, setTreatTotal] = useState<number>(entity?.treatTotal || 0);
  const [stakedGhst, setStakedGhst] = useState<number>(entity?.stakedGhst || 0);

  // Fetch Gotchi data on mount and periodically
  useEffect(() => {
    const fetchGotchiData = async () => {
      const gotchiId = entity?.data?.gotchiId;
      const gotchiEntity = GotchiEntity.activeGotchis.get(gotchiId);

      setTreatTotal(gotchiEntity?.data?.stats?.treatTotal || 0);
      setStakedGhst(gotchiEntity?.data?.stats?.stakedGhst || 0);
    };

    fetchGotchiData();
    const interval = setInterval(fetchGotchiData, 5000); // Poll every 5 seconds
    return () => clearInterval(interval);
  }, []);

  // Update treatTotal and stakedGhst when entity changes
  useEffect(() => {
    if (entity) {
      setTreatTotal(entity?.data?.stats?.treatTotal || 0);
      setStakedGhst(entity?.data?.stats?.stakedGhst || 0);
    }
  }, [entity]);

  const handleItemClick = (item: TreatItem) => {
    setSelectedItem(item);
  };

  const handleStake = async () => {
    const amount = parseFloat(stakeAmount);
    if (isNaN(amount) || amount <= 0) {
      alert('Please enter a valid positive number for staking.');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/gotchi/stake`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          ghstAmount: amount,  
          uuid: entity?.data.uuid, 
          zoneId: entity?.data.zoneId
        }),
      });
      if (!response.ok) throw new Error('Failed to stake GHST');
      const data = await response.json();
      setTreatTotal(data.treatTotal);
      setStakedGhst(data.stakedGhst);
      eventBus.emit("treatEaten", {detail: {ecto:data.ecto, spark:data.spark, pulse:data.pulse}});
      setStakeAmount(''); // Clear input
    } catch (error) {
      console.error('[TreatModal] Failed to stake GHST:', error);
      alert('Failed to stake GHST. Please try again.');
    }
  };

  const handleUnstake = async () => {
    const amount = parseFloat(unstakeAmount);
    if (isNaN(amount) || amount <= 0) {
      alert('Please enter a valid positive number for unstaking.');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/gotchi/unstake`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          amount: amount,
          uuid: entity?.data.uuid, 
          zoneId: entity?.data.zoneId
         }),
      });
      if (!response.ok) throw new Error('Failed to unstake GHST');
      const data = await response.json();
      setTreatTotal(data.treatTotal);
      setStakedGhst(data.stakedGhst);
      eventBus.emit("treatEaten", {detail: {ecto:data.ecto, spark:data.spark, pulse:data.pulse}});
      setUnstakeAmount(''); // Clear input
    } catch (error) {
      console.error('[TreatModal] Failed to unstake GHST:', error);
      alert('Failed to unstake GHST. Please try again.');
    }
  };

  const handleEat = async (item: TreatItem) => {
    if (treatTotal < item.cost) {
      alert('Not enough TREAT to eat this treat.');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/gotchi/eat`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          treatName: item.name,
          uuid: entity?.data.uuid, 
          zoneId: entity?.data.zoneId
        }),
      });
      if (!response.ok) throw new Error('Failed to eat treat');
      const data = await response.json();
      setTreatTotal(data.treatTotal);
      setStakedGhst(data.stakedGhst);
      eventBus.emit("treatEaten", {detail: {ecto:data.ecto, spark:data.spark, pulse:data.pulse}});
    } catch (error) {
      console.error('[TreatModal] Failed to eat treat:', error);
      alert('Failed to eat treat. Please try again.');
    }
  };

  return (
    <div className="treat-modal">
      <div className="header-bar">
        <div className="header-content">Treats</div>
      </div>

      <div className="sub-header-content">
        <div className="treat-details">
          <div className="treat-details-content">
            <ul>
              <li>TREAT Rank: Aadept</li>
              <li>TREAT per Day: {stakedGhst}</li>
              <li>TREAT Total: {treatTotal.toFixed(2)}</li>
            </ul>
          </div>
        </div>

        <div className="ghst-staking">
          <div className="ghst-staking-content">
            <ul>
              <li className="stake-row">
                <input
                  type="number"
                  placeholder="Enter GHST Amount"
                  value={stakeAmount}
                  onChange={(e) => setStakeAmount(e.target.value)}
                />
                <button onClick={handleStake}>Stake</button>
              </li>
              <li className="stake-row">
                <input
                  type="number"
                  placeholder="Enter GHST Amount"
                  value={unstakeAmount}
                  onChange={(e) => setUnstakeAmount(e.target.value)}
                />
                <button onClick={handleUnstake}>Unstake</button>
              </li>
              <li style={{ marginTop: '0.1rem' }}>Staked GHST: {stakedGhst}</li>
            </ul>
          </div>
        </div>
      </div>

      <div className="main-content">
        <div className="item-details">
          {selectedItem && (
            <TreatDetails
              name={selectedItem.name}
              cost={selectedItem.cost}
              imageSrc={selectedItem.imageSrc}
              description={selectedItem.description}
              onEat={() => handleEat(selectedItem)}
              canEat={treatTotal >= selectedItem.cost}
            />
          )}
        </div>

        <div className="inventory-items"> 
          <div className="inventory-grid">
            {items.map((item) => (
              <TreatItemCard
                key={item.name}
                name={item.name}
                cost={item.cost}
                imageSrc={item.imageSrc}
                onClick={() => handleItemClick(item)}
              />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default TreatModal;
