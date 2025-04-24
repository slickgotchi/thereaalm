import Phaser from "phaser";
import { TILE_PIXELS, ZONE_TILES } from "./GameScene";
import { eventBus } from "../utils/EventBus"; // Adjust path
import { GotchiEntity } from "./entities/GotchiEntity";

export class CameraController {
    private scene: Phaser.Scene;
    private worldWidth: number;
    private worldHeight: number;
    private minZoom: number = 0.005;
    private maxZoom: number = 3.0;
    private isDragging: boolean = false;
    private dragStart: Phaser.Math.Vector2 = new Phaser.Math.Vector2();

    private newZoom!: number;
    private oldZoom!: number;

    constructor(scene: Phaser.Scene, worldWidth: number, worldHeight: number) {
        this.scene = scene;
        this.worldWidth = worldWidth;
        this.worldHeight = worldHeight;
        this.setupCamera();
        this.setupInput();

        // Subscribe to panToGotchi event
        eventBus.on("panToGotchi", this.handlePanToGotchi.bind(this));
    }

    private setupCamera() {
        const margin = ZONE_TILES * TILE_PIXELS;
        const width = this.worldWidth + 2 * margin;
        const height = this.worldHeight + 2 * margin;

        const initialZoom = Math.min(
            this.scene.scale.width / width,
            this.scene.scale.height / height
        );
        this.scene.cameras.main.setZoom(Math.max(initialZoom, this.minZoom));
        this.scene.cameras.main.setZoom(1);

        this.scene.cameras.main.centerOn(
            4.5 * ZONE_TILES * TILE_PIXELS,
            5.5 * ZONE_TILES * TILE_PIXELS
        );

        this.scene.cameras.main.setBackgroundColor(0x131313);
    }

    private setupInput() {
        this.scene.input.on("pointerdown", (pointer: Phaser.Input.Pointer) => {
            if (pointer.button === 0) {
                this.isDragging = true;
                this.dragStart.set(pointer.x, pointer.y);
            }
        });

        this.scene.input.on("pointerup", (pointer: Phaser.Input.Pointer) => {
            if (pointer.button === 0) {
                this.isDragging = false;
            }
        });

        this.scene.input.on("pointermove", (pointer: Phaser.Input.Pointer) => {
            if (this.isDragging) {
                const zoom = this.scene.cameras.main.zoom;
                const deltaX = (pointer.x - this.dragStart.x) / zoom;
                const deltaY = (pointer.y - this.dragStart.y) / zoom;
                this.scene.cameras.main.scrollX -= deltaX;
                this.scene.cameras.main.scrollY -= deltaY;
                this.dragStart.set(pointer.x, pointer.y);
            }
        });

        this.scene.input.on(
            "wheel",
            (
                pointer: Phaser.Input.Pointer,
                gameObjects: any[],
                deltaX: number,
                deltaY: number,
                deltaZ: number
            ) => {
                const camera = this.scene.cameras.main;
                this.oldZoom = camera.zoom;
                this.newZoom = this.oldZoom;

                const zoomFactor = this.oldZoom * 0.5;
                const zoomDelta = Math.sign(deltaY) * zoomFactor;
                this.newZoom = this.oldZoom - zoomDelta;

                this.newZoom = Phaser.Math.Clamp(
                    this.newZoom,
                    this.minZoom,
                    this.maxZoom
                );

                camera.zoomTo(this.newZoom, 100);
            }
        );
    }

    public panAndZoomTo(x: number, y: number, zoom: number = 1, duration: number = 500) {
        const camera = this.scene.cameras.main;
        this.scene.tweens.add({
            targets: camera,
            scrollX: x - (camera.width / 2) / zoom,
            scrollY: y - (camera.height / 2) / zoom,
            // zoom: zoom, // Add zoom to the tween
            duration: duration,
            ease: "Quad.easeInOut",
            onComplete: () => {
                camera.zoomTo(zoom, duration, "Quad.easeInOut");
            }
        });
    }

    private handlePanToGotchi(data: { gotchiId: string }) {
        const entity = GotchiEntity.activeGotchis.get(data.gotchiId);
        if (entity) {
            this.panAndZoomTo(
                entity.currentPosition.x + 32,
                entity.currentPosition.y +32+128,
                1,
                500
            );
        }
    }

    public destroy() {
        eventBus.off("panToGotchi", this.handlePanToGotchi.bind(this));
    }
}
