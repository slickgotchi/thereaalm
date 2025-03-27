// TweenableEntity.ts
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

        const updateCallback: TweenUpdateCallback = (x: number, y: number, direction: Direction) => {
            this.sprite.setPosition(x, y);
            // Only update direction from tween if it's not "none"
            if (direction !== "none") {
                this.direction = direction;
            }
            this.updateDirection();
            if (x % TILE_PIXELS === 0 && y % TILE_PIXELS === 0) {
                this.lastTileX = x / TILE_PIXELS;
                this.lastTileY = y / TILE_PIXELS;
                // Apply target direction when tween completes
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

        if (!this.tweenWorker.getIsTweening() && (this.lastTileX !== this.tileX || this.lastTileY !== this.tileY)) {
            this.tweenWorker.setTargetTilePosition(this.lastTileX, this.lastTileY, this.tileX, this.tileY);
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

/*
import Phaser from "phaser";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { TweenWorker, Direction, TweenUpdateCallback } from "./TweenWorker";
import { TILE_PIXELS } from "../GameScene";

export class TweenableEntity extends BaseEntity {
    private tweenWorker: TweenWorker;
    protected direction: Direction = "none";
    private lastTileX: number;
    private lastTileY: number;

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, type: string, texture: string, data: any) {
        super(scene, id, zoneId, tileX, tileY, type, texture, data);

        this.lastTileX = tileX;
        this.lastTileY = tileY;

        // Callback now receives pixel coordinates
        const updateCallback: TweenUpdateCallback = (x: number, y: number, direction: Direction) => {
            this.sprite.setPosition(x, y);
            this.direction = direction;
            this.updateDirection();
            // Update last tile position when reaching a new tile
            if (x % TILE_PIXELS === 0 && y % TILE_PIXELS === 0) {
                this.lastTileX = x / TILE_PIXELS;
                this.lastTileY = y / TILE_PIXELS;
            }
        };

        this.tweenWorker = new TweenWorker(scene, updateCallback);
    }

    update(snapshot: EntitySnapshot) {
        super.update(snapshot);
        console.log("direction: ", snapshot.data.direction)

        if (!this.tweenWorker.getIsTweening() && (this.lastTileX !== this.tileX || this.lastTileY !== this.tileY)) {
            this.tweenWorker.setTargetTilePosition(this.lastTileX, this.lastTileY, this.tileX, this.tileY);
        }

        if (!this.tweenWorker.getIsTweening()) {
            this.setDirection(snapshot.data.direction)
        }
    }

    setDirection(direction: Direction) {
        this.direction = direction;
        this.updateDirection();
    }

    protected updateDirection() {
        // Default does nothing; subclasses can override
    }
}*/