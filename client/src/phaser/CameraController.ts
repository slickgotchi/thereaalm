// src/phaser/CameraController.ts
import Phaser from "phaser";
import { TILE_PIXELS, ZONE_TILES } from "./GameScene";

export class CameraController {
    private scene: Phaser.Scene;
    private worldWidth: number;
    private worldHeight: number;
    private minZoom: number = 0.005; // Zoomed out to see the whole map
    private maxZoom: number = 1.0; // Zoomed in to 1:1 scale
    private isDragging: boolean = false;
    private dragStart: Phaser.Math.Vector2 = new Phaser.Math.Vector2();

    private newZoom!: number;
    private oldZoom!: number;

    constructor(scene: Phaser.Scene, worldWidth: number, worldHeight: number) {
        this.scene = scene;
        this.worldWidth = worldWidth;
        this.worldHeight = worldHeight;

        // Set up the camera
        this.setupCamera();

        // Set up input handlers
        this.setupInput();
    }

    private setupCamera() {
        const margin = ZONE_TILES * TILE_PIXELS;
        const width = this.worldWidth + 2 * margin;
        const height = this.worldHeight + 2 * margin;

        // Start the camera zoomed out to show the entire world
        const initialZoom = Math.min(
            this.scene.scale.width / width,
            this.scene.scale.height / height
        );
        this.scene.cameras.main.setZoom(Math.max(initialZoom, this.minZoom));
        this.scene.cameras.main.setZoom(0.01);

        // Center the camera on the world
        this.scene.cameras.main.centerOn(
            -margin + width / 2,
            -margin + height / 2
        );

        this.scene.cameras.main.setBackgroundColor(0x131313);
    }

    private setupInput() {
        // Panning with left click + drag
        this.scene.input.on("pointerdown", (pointer: Phaser.Input.Pointer) => {
            if (pointer.button === 0) {
                // Left mouse button
                this.isDragging = true;
                this.dragStart.set(pointer.x, pointer.y);
            }
        });

        this.scene.input.on("pointerup", (pointer: Phaser.Input.Pointer) => {
            if (pointer.button === 0) {
                // Left mouse button
                this.isDragging = false;
            }
        });

        this.scene.input.on("pointermove", (pointer: Phaser.Input.Pointer) => {
            if (this.isDragging) {
                const zoom = this.scene.cameras.main.zoom;
                const deltaX = (pointer.x - this.dragStart.x) / zoom;
                const deltaY = (pointer.y - this.dragStart.y) / zoom;

                // Pan the camera
                this.scene.cameras.main.scrollX -= deltaX;
                this.scene.cameras.main.scrollY -= deltaY;

                // Update drag start position
                this.dragStart.set(pointer.x, pointer.y);
            }
        });

        // Zooming with middle mouse wheel
        this.scene.input.on(
            "wheel",
            (
                pointer: Phaser.Input.Pointer,
                gameObjects: any[],
                deltaX: number,
                deltaY: number,
                deltaZ: number
            ) => {
                // 1. store current zoom in oldZoom
                const camera = this.scene.cameras.main;
                this.oldZoom = camera.zoom;
                this.newZoom = this.oldZoom;

                // 3. Calculate the new zoom level
                const zoomFactor = this.oldZoom * 0.5;
                const zoomDelta = Math.sign(deltaY) * zoomFactor;
                this.newZoom = this.oldZoom - zoomDelta;

                this.newZoom = Phaser.Math.Clamp(
                    this.newZoom,
                    this.minZoom,
                    this.maxZoom
                );

                // 4. set the new zoom
                camera.zoomTo(this.newZoom, 100);
            }
        );
    }
}
