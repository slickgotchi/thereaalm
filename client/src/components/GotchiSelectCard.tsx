import { Aavegotchi } from "../utils/gotchi-loader/FetchGotchis";
import { eventBus } from "../utils/EventBus"; // Adjust path
import "./GotchiSelectCard.css";

interface GotchiSelectCardProps {
    gotchi: Aavegotchi;
    onSelect: () => void;
    isActive: boolean;
}

export const GotchiSelectCard = ({ gotchi, onSelect, isActive }: GotchiSelectCardProps) => {
    const handleClick = () => {
        if (!isActive) return;

        // Emit event via EventBus
        eventBus.emit("panToGotchi", { gotchiId: String(gotchi.id) });

        console.log(`Selected Gotchi: ${gotchi.name}`);
        onSelect();
    };

    return (
        <div
            className={`gotchi-select-card ${isActive ? "" : "inactive"}`}
            onClick={handleClick}
        >
            <span className="gotchi-name">{gotchi.name}</span>
            <span className="gotchi-brs">BRS: {gotchi.withSetsRarityScore}</span>
        </div>
    );
};