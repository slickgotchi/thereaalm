import { useEffect, useLayoutEffect, useRef, useState } from "react";
import Phaser from "phaser";
import { GameScene } from "./phaser/GameScene";
import "./index.css";
import "./App.css";
import { HoverInfo } from "./components/HoverInfo";
import { ConnectWallet } from "./components/ConnectWallet";
import TreatModal from "./components/menu/TreatModal";
import MenuSystem from "./components/menu/MenuSystem";
import { eventBus } from "./utils/EventBus";
import { GotchiHUD } from "./components/gotchi-hud/GotchiHUD";

function App() {
    const gameRef = useRef<Phaser.Game | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const renderCount = useRef(0);

    const [selectedEntity, setSelectedEntity] = useState<{ type: string; [key: string]: any } | null>(null);

    useEffect(() => {
        if (!containerRef.current) return;

        // increment render count (this prevents double renders in)
        if (
            import.meta.env.MODE === "development" &&
            renderCount.current == 0
        ) {
            renderCount.current++;
            return;
        }

        const config: Phaser.Types.Core.GameConfig = {
            type: Phaser.AUTO,
            parent: containerRef.current,
            // parent: "game-parent",
            scene: [GameScene],
            scale: {
                mode: Phaser.Scale.ENVELOP,
                width: 1920,
                height: 1200,
                autoCenter: Phaser.Scale.CENTER_BOTH,
            },
            pixelArt: true,
            roundPixels: true,
        };

        // Initialize Phaser game
        gameRef.current = new Phaser.Game(config);

        const handleEntitySelection = (data: { detail: any }) => {
            console.log("[App] Received entitySelection event:", data);
            setSelectedEntity(data.detail); // Update selectedEntity with the event data
        };

        eventBus.on("entitySelection", handleEntitySelection);


        // Cleanup on unmount
        return () => {
            console.log("Cleaning up Phaser game");
            if (gameRef.current) {
                gameRef.current.destroy(true);
                gameRef.current = null;
            }

            eventBus.off("entitySelection", handleEntitySelection);

        };
    }, []); // Empty dependency array ensures it only runs once

    return (
    <div>
        {/* <HoverInfo /> */}
        <ConnectWallet />
        <MenuSystem
            selectedEntity={selectedEntity}
            onClose={() => setSelectedEntity(null)}
        />
        {selectedEntity && selectedEntity.type === "gotchi" &&
        <GotchiHUD
            selectedGotchiEntity={selectedEntity}
        />
        }
        <div ref={containerRef} className="game-container" />
    </div>
    );
}

export default App;
