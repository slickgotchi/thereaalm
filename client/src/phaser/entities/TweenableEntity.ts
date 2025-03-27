import Phaser from "phaser";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { TweenWorker, Direction, TweenUpdateCallback } from "./TweenWorker";
import { TILE_PIXELS } from "../GameScene";

export class TweenableEntity extends BaseEntity {
    private tweenWorker: TweenWorker;
    protected direction: Direction = "none";
    private lastTileX: number;
    private lastTileY: number;
    private targetDirection: Direction = "none";

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, type: string, texture: string, data: any) {
        super(scene, id, zoneId, tileX, tileY, type, texture, data);

        this.lastTileX = tileX;
        this.lastTileY = tileY;

        const updateCallback: TweenUpdateCallback = (x: number, y: number, direction: Direction, isComplete?: boolean) => {
            this.sprite.setPosition(x, y);
            if (direction !== "none") {
                this.direction = direction;
            }
            this.updateDirection();
            if (isComplete) {
                this.lastTileX = x / TILE_PIXELS;
                this.lastTileY = y / TILE_PIXELS;
                if (!this.tweenWorker.getIsTweening()) {
                    this.direction = this.targetDirection;
                    this.updateDirection();
                }
            }
        };

        this.tweenWorker = new TweenWorker(scene, updateCallback);
    }

    update(snapshot: EntitySnapshot) {
        super.update(snapshot);
        this.targetDirection = snapshot.data.direction;

        if (this.lastTileX !== snapshot.tileX || this.lastTileY !== snapshot.tileY) {
            // Set sprite to last known position to prevent flicker
            this.sprite.setPosition(this.lastTileX * TILE_PIXELS, this.lastTileY * TILE_PIXELS);
            this.tweenWorker.setTargetTilePosition(this.lastTileX, this.lastTileY, snapshot.tileX, snapshot.tileY);
            this.tileX = snapshot.tileX; // Update tileX/tileY for state, but don't set sprite yet
            this.tileY = snapshot.tileY;
        }

        if (!this.tweenWorker.getIsTweening()) {
            this.direction = this.targetDirection;
            this.updateDirection();
        }
    }

    setDirection(direction: Direction) {
        this.targetDirection = direction;
        if (!this.tweenWorker.getIsTweening()) {
            this.direction = direction;
            this.updateDirection();
        }
    }

    protected updateDirection() {
        // Default does nothing; subclasses can override
    }
}