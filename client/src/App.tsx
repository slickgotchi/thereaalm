import { useEffect, useLayoutEffect, useRef } from "react";
import Phaser from "phaser";
import { GameScene } from "./phaser/GameScene";
import "./index.css";
import "./App.css";
import { HoverInfo } from "./components/HoverInfo";
import { ConnectWallet } from "./components/ConnectWallet";

function App() {
    const gameRef = useRef<Phaser.Game | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const renderCount = useRef(0);

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

        // Cleanup on unmount
        return () => {
            console.log("Cleaning up Phaser game");
            if (gameRef.current) {
                gameRef.current.destroy(true);
                gameRef.current = null;
            }
        };
    }, []); // Empty dependency array ensures it only runs once

    return (
    <div>
        <HoverInfo />
        <ConnectWallet />
        <div ref={containerRef} className="game-container" />
    </div>
    );
}

export default App;
