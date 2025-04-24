import { useEffect } from "react";
import { Entity } from "../MenuSystem";

import './GotchiDetails.css';

export const GotchiDetails: React.FC<{ entity: Entity }> = ({ entity }) => {
    const { name, gotchiId, description, personality } = entity.data;
    const { nrg, agg, spk, brn, eys, eyc } = entity.data.stats;

    useEffect(() => {
        console.log(entity);
    }, [entity]);

    // Helper component to render a trait line with a positioned label
    const TraitLine: React.FC<{ label: string; value: number; personalityText: string }> = ({ label, value, personalityText }) => {
        // Clamp the value to the range 0-99 inclusive
        const clampedValue = Math.max(0, Math.min(99, value));
        // Calculate the percentage position of the label (value is 0-99, so map to 0-100%)
        const positionPercentage = (clampedValue / 99) * 100;

        // Determine rarity category based on trait value
        const getRarityClass = (value: number): string => {
            if ([0, 1, 98, 99].includes(value)) {
                return 'myth';
            } else if ((value >= 2 && value <= 9) || (value >= 91 && value <= 97)) {
                return 'legendary';
            } else if ((value >= 10 && value <= 24) || (value >= 75 && value <= 90)) {
                return 'uncommon';
            } else {
                // 25-74 inclusive
                return 'common';
            }
        };

        const rarityClass = getRarityClass(clampedValue);

        return (
            <div className="trait-line-container">
                <div className="trait-label">{label}</div>
                <div className="trait-line">
                    <div className="line">
                        <div className="middle-tick"></div>
                    </div>
                    <div
                        className={`personality-marker ${rarityClass}`}
                        style={{ left: `${positionPercentage}%` }}
                    >
                        {personalityText}
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className="gotchi-details">
            <h1>{name} ({gotchiId})</h1>
            <p><b>Entity Type:</b> {entity.type}</p>
            <p><b>Description:</b> {description}</p>
            <div className="personality-section">
                <div className="personality-title">Personality:</div>
                <TraitLine label="NRG" value={nrg} personalityText={personality[0]} />
                <TraitLine label="AGG" value={agg} personalityText={personality[1]} />
                <TraitLine label="SPK" value={spk} personalityText={personality[2]} />
                <TraitLine label="BRN" value={brn} personalityText={personality[3]} />
                <TraitLine label="EYS" value={eys} personalityText={personality[4]} />
                <TraitLine label="EYC" value={eyc} personalityText={personality[5]} />
            </div>
        </div>
    );
};


/*
import { useEffect } from "react";
import { Entity } from "../MenuSystem";

import './GotchiDetails.css';

export const GotchiDetails: React.FC<{ entity: Entity }> = ({ entity }) => {
    const { name, gotchiId, description, personality } = entity.data;
    const { nrg, agg, spk, brn, eys, eyc } = entity.data.stats;

    useEffect(() => {
        console.log(entity);
    }, [entity]);

    // Helper component to render a trait line with a positioned label
    const TraitLine: React.FC<{ label: string; value: number; personalityText: string }> = ({ label, value, personalityText }) => {
        // Clamp the value to the range 0-99 inclusive
        const clampedValue = Math.max(0, Math.min(99, value));
        // Calculate the percentage position of the label (value is 0-99, so map to 0-100%)
        const positionPercentage = (clampedValue / 99) * 100;

        return (
            <div className="trait-line-container">
                <div className="trait-label">{label}</div>
                <div className="trait-line">
                    <div className="line">
                        <div className="middle-tick"></div>
                    </div>
                    <div
                        className="personality-marker"
                        style={{ left: `${positionPercentage}%` }}
                    >
                        {personalityText}
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className="gotchi-details">
            <h1>{name} ({gotchiId})</h1>
            <p><b>Entity Type:</b> {entity.type}</p>
            <p><b>Description:</b> {description}</p>
            <div className="personality-section">
                <div className="personality-title">Personality:</div>
                <TraitLine label="NRG" value={nrg} personalityText={personality[0]} />
                <TraitLine label="AGG" value={agg} personalityText={personality[1]} />
                <TraitLine label="SPK" value={spk} personalityText={personality[2]} />
                <TraitLine label="BRN" value={brn} personalityText={personality[3]} />
                <TraitLine label="EYS" value={eys} personalityText={personality[4]} />
                <TraitLine label="EYC" value={eyc} personalityText={personality[5]} />
            </div>
        </div>
    );
};*/