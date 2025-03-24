import { Depth } from "../Depth";
import { TILE_PIXELS } from "../GameScene";
import { OutlineEffect } from "../OutlineEffect";

export interface EntitySnapshot {
    id: string; // Matches server's "uuid"
    zoneId: number;
    type: string; // Matches server's "gotchiId"
    tileX: number;
    tileY: number;
    data: any;
}

export class BaseEntity {
    id: string;
    zoneId: number;
    sprite: Phaser.GameObjects.Sprite;
    scene: Phaser.Scene;
    x: number;
    y: number;
    type: string;

    outlineEffect: OutlineEffect;
    data: any;
    private hoverUpdateTimer: Phaser.Time.TimerEvent | null = null; // Timer for hover updates


    constructor(scene: Phaser.Scene, id: string, zoneId: number, x: number, y: number, type: string, texture: string, data: any) {
        this.scene = scene;
        this.id = id;
        this.zoneId = zoneId;
        this.x = x * TILE_PIXELS;
        this.y = y * TILE_PIXELS;
        this.type = type;

        this.sprite = scene.add.sprite(this.x, this.y, texture)
            .setOrigin(0, 0)
            .setDepth(Depth.ENTITIES);

        this.outlineEffect = new OutlineEffect({
            scene: this.scene,
            target: this.sprite,
        });

        this.setupHoverEvents();

        this.data = data;
    }

    update(snapshot: EntitySnapshot) {
        const { tileX, tileY, data } = snapshot;

        this.x = tileX * TILE_PIXELS;
        this.y = tileY * TILE_PIXELS;
        this.data = data;

        this.sprite.setPosition(this.x, this.y);

        // If currently hovered, keep updating the hover info
        if (this.hoverUpdateTimer) {
            this.dispatchHoverData();
        }
    }

    destroy() {
        this.sprite.destroy();
    }

    private setupHoverEvents() {
        this.sprite.setInteractive({ useHandCursor: true });

        this.sprite.on("pointerover", () => {
            this.outlineEffect.showOutline();

            // Start updating hover data
            this.startHoverUpdates();
        });

        this.sprite.on("pointerout", () => {
            this.outlineEffect.hideOutline();

            // Stop updating hover data
            this.stopHoverUpdates();

            window.dispatchEvent(new CustomEvent("entityHover", { detail: null }));
        });
    }

    private startHoverUpdates() {
        // Dispatch immediately
        this.dispatchHoverData();

        // Set up a timer to update hover data periodically
        this.hoverUpdateTimer = this.scene.time.addEvent({
            delay: 100, // Update every 100ms (adjust as needed)
            callback: this.dispatchHoverData,
            callbackScope: this,
            loop: true,
        });
    }

    private stopHoverUpdates() {
        if (this.hoverUpdateTimer) {
            this.hoverUpdateTimer.remove(false); // Stop the timer
            this.hoverUpdateTimer = null;
        }
    }

    private dispatchHoverData() {
        window.dispatchEvent(new CustomEvent("entityHover", {
            detail: {
                type: this.type,
                data: this.data,
            },
        }));
    }
}
