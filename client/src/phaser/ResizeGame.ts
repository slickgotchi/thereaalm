import Phaser from "phaser";

// Constants for game dimensions and aspect ratio
export const GAME_WIDTH = 1920; // Default game width
export const GAME_HEIGHT = 1200; // Default game height
const ASPECT_RATIO = 16 / 10; // 1.6 (16:10 aspect ratio)

// Function to resize the game based on window size
export function resizeGame(scene: Phaser.Scene) {
    if (!scene.cameras.main || !scene.scale) return;

    const availableWidth = window.innerWidth;
    const availableHeight = window.innerHeight;
    const widthBasedOnAvailableHeight = availableHeight * ASPECT_RATIO;
    const newWidth = Math.max(availableWidth, widthBasedOnAvailableHeight);
    const newHeight = newWidth / ASPECT_RATIO;

    // scene.scale.setGameSize(newWidth, newHeight);

    const zoom = Math.min(newWidth / GAME_WIDTH, newHeight / GAME_HEIGHT);
    // scene.cameras.main.setZoom(zoom);

    // console.log("resize zoom:", zoom, availableWidth, availableHeight);
}
