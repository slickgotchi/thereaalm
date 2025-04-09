import { useState, useEffect, useRef } from "react";
import { fetchAavegotchis, Aavegotchi } from "../phaser/FetchGotchis";
import { GotchiSelectCard } from "./GotchiSelectCard";
import { GotchiEntity } from "../phaser/entities/GotchiEntity";
import "./GotchiSelectModal.css";

interface GotchiSelectModalProps {
    account: string;
    onClose: () => void;
}

export const GotchiSelectModal = ({ account, onClose }: GotchiSelectModalProps) => {
    const [gotchis, setGotchis] = useState<Aavegotchi[]>([]);
    const modalRef = useRef<HTMLDivElement>(null); // Ref to track the modal content

    useEffect(() => {
        const loadGotchis = async () => {
            try {
                const fetchedGotchis = await fetchAavegotchis(account);
                console.log({account, fetchedGotchis});
                const sortedGotchis = fetchedGotchis.sort(
                    (a, b) => b.withSetsRarityScore - a.withSetsRarityScore
                );
                setGotchis(sortedGotchis);
            } catch (error) {
                console.error("Failed to fetch Gotchis:", error);
            }
        };
        loadGotchis();
    }, [account]);

    // Handle clicks outside the modal
    const handleOverlayClick = (event: React.MouseEvent<HTMLDivElement>) => {
        if (modalRef.current && !modalRef.current.contains(event.target as Node)) {
            onClose(); // Close modal if click is outside modal content
        }
    };

    return (
        <div className="modal-overlay" onClick={handleOverlayClick}>
            <div className="gotchi-select-modal" ref={modalRef}>
                <button className="close-button" onClick={onClose}>
                    ×
                </button>
                <h2>Select Your Gotchi</h2>
                <div className="gotchi-list">
                    {gotchis.map((gotchi) => (
                        <GotchiSelectCard
                            key={gotchi.id}
                            gotchi={gotchi}
                            onSelect={onClose}
                            isActive={GotchiEntity.activeGotchis.has(String(gotchi.id))}
                        />
                    ))}
                </div>
            </div>
        </div>
    );
};

/*
import { useState, useEffect } from "react";
import { fetchAavegotchis, Aavegotchi } from "../phaser/FetchGotchis";
import { GotchiSelectCard } from "./GotchiSelectCard";
import { GotchiEntity } from "../phaser/entities/GotchiEntity"; // Adjust path
import "./GotchiSelectModal.css";

interface GotchiSelectModalProps {
    account: string;
    onClose: () => void;
}

export const GotchiSelectModal = ({ account, onClose }: GotchiSelectModalProps) => {
    const [gotchis, setGotchis] = useState<Aavegotchi[]>([]);

    useEffect(() => {
        const loadGotchis = async () => {
            try {
                const fetchedGotchis = await fetchAavegotchis(account);
                const sortedGotchis = fetchedGotchis.sort(
                    (a, b) => b.withSetsRarityScore - a.withSetsRarityScore
                );
                setGotchis(sortedGotchis);
            } catch (error) {
                console.error("Failed to fetch Gotchis:", error);
            }
        };
        loadGotchis();
    }, [account]);

    return (
        <div className="modal-overlay">
            <div className="gotchi-select-modal">
                <button className="close-button" onClick={onClose}>
                    ×
                </button>
                <h2>Select Your Gotchi</h2>
                <div className="gotchi-list">
                    {gotchis.map((gotchi) => (
                        <GotchiSelectCard
                            key={gotchi.id}
                            gotchi={gotchi}
                            onSelect={onClose}
                            isActive={GotchiEntity.activeGotchis.has(String(gotchi.id))}
                        />
                    ))}
                </div>
            </div>
        </div>
    );
};
*/