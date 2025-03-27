// TweenWorker.ts
import Phaser from "phaser";

export type Direction = "up" | "down" | "left" | "right" | "none";
export type TweenUpdateCallback = (x: number, y: number, direction: Direction) => void;

export class TweenWorker {
    private scene: Phaser.Scene;
    private startTileX: number;
    private startTileY: number;
    private targetTileX: number;
    private targetTileY: number;
    private currentTileX: number;
    private currentTileY: number;
    private isTweening: boolean = false;
    private tileSize: number = 64;
    private tweenDuration: number = 200;
    private updateCallback: TweenUpdateCallback;

    constructor(scene: Phaser.Scene, updateCallback: TweenUpdateCallback) {
        this.scene = scene;
        this.updateCallback = updateCallback;
        this.startTileX = 0;
        this.startTileY = 0;
        this.targetTileX = 0;
        this.targetTileY = 0;
        this.currentTileX = 0;
        this.currentTileY = 0;
    }

    public getIsTweening(): boolean {
        return this.isTweening;
    }

    public setTargetTilePosition(startTileX: number, startTileY: number, targetTileX: number, targetTileY: number) {
        this.startTileX = startTileX;
        this.startTileY = startTileY;
        this.targetTileX = targetTileX;
        this.targetTileY = targetTileY;
        this.currentTileX = startTileX;
        this.currentTileY = startTileY;

        this.updateCallback(this.currentTileX * this.tileSize, this.currentTileY * this.tileSize, "none");

        if (!this.isTweening) {
            this.startTweening();
        }
    }

    private startTweening() {
        if (this.currentTileX === this.targetTileX && this.currentTileY === this.targetTileY) {
            return;
        }

        this.isTweening = true;
        this.tweenToNextTile();
    }

    private tweenToNextTile() {
        if (this.currentTileX === this.targetTileX && this.currentTileY === this.targetTileY) {
            this.isTweening = false;
            this.updateCallback(this.currentTileX * this.tileSize, this.currentTileY * this.tileSize, "none");
            return;
        }

        const dx = this.targetTileX - this.currentTileX;
        const dy = this.targetTileY - this.currentTileY;

        const moveX = Math.random() < 0.5 && dx !== 0;
        let nextTileX = this.currentTileX;
        let nextTileY = this.currentTileY;
        let direction: Direction = "none";

        if (moveX) {
            nextTileX = this.currentTileX + Math.sign(dx);
            direction = dx > 0 ? "right" : "left";
        } else if (dy !== 0) {
            nextTileY = this.currentTileY + Math.sign(dy);
            direction = dy > 0 ? "down" : "up";
        }

        const startX = this.currentTileX * this.tileSize;
        const startY = this.currentTileY * this.tileSize;
        const nextX = nextTileX * this.tileSize;
        const nextY = nextTileY * this.tileSize;

        const tweenTarget = { progress: 0 };
        this.scene.tweens.add({
            targets: tweenTarget,
            progress: 1,
            duration: this.tweenDuration,
            ease: "Linear",
            onUpdate: (tween) => {
                const t = tween.progress;
                const currentX = startX + (nextX - startX) * t;
                const currentY = startY + (nextY - startY) * t;
                this.updateCallback(currentX, currentY, direction);
            },
            onComplete: () => {
                this.currentTileX = nextTileX;
                this.currentTileY = nextTileY;
                this.updateCallback(nextX, nextY, direction);
                this.tweenToNextTile();
            },
        });
    }
}

/*
import Phaser from "phaser";

export type Direction = "up" | "down" | "left" | "right" | "none";

// Updated callback to use pixel coordinates instead of tile coordinates
export type TweenUpdateCallback = (x: number, y: number, direction: Direction) => void;

export class TweenWorker {
    private scene: Phaser.Scene;
    private startTileX: number;
    private startTileY: number;
    private targetTileX: number;
    private targetTileY: number;
    private currentTileX: number;
    private currentTileY: number;
    private isTweening: boolean = false;
    private tileSize: number = 64; // Game units per tile
    private tweenDuration: number = 200; // ms per tile
    private updateCallback: TweenUpdateCallback;
    private direction: Direction = "none";

    constructor(scene: Phaser.Scene, updateCallback: TweenUpdateCallback) {
        this.scene = scene;
        this.updateCallback = updateCallback;
        this.startTileX = 0;
        this.startTileY = 0;
        this.targetTileX = 0;
        this.targetTileY = 0;
        this.currentTileX = 0;
        this.currentTileY = 0;
    }

    public getIsTweening(): boolean {
        return this.isTweening;
    }

    public setDirection(direction: Direction): void {
        this.direction = direction;
    }

    public setTargetTilePosition(startTileX: number, startTileY: number, targetTileX: number, targetTileY: number) {
        this.startTileX = startTileX;
        this.startTileY = startTileY;
        this.targetTileX = targetTileX;
        this.targetTileY = targetTileY;
        this.currentTileX = startTileX;
        this.currentTileY = startTileY;

        // Emit initial pixel position
        this.updateCallback(this.currentTileX * this.tileSize, this.currentTileY * this.tileSize, this.direction);

        if (!this.isTweening) {
            this.startTweening();
        }
    }

    private startTweening() {
        if (this.currentTileX === this.targetTileX && this.currentTileY === this.targetTileY) {
            return;
        }

        this.isTweening = true;
        this.tweenToNextTile();
    }

    private tweenToNextTile() {
        if (this.currentTileX === this.targetTileX && this.currentTileY === this.targetTileY) {
            this.isTweening = false;
            this.updateCallback(this.currentTileX * this.tileSize, this.currentTileY * this.tileSize, "none");
            return;
        }

        const dx = this.targetTileX - this.currentTileX;
        const dy = this.targetTileY - this.currentTileY;

        // Randomly choose x or y movement (50% chance)
        const moveX = Math.random() < 0.5 && dx !== 0;
        let nextTileX = this.currentTileX;
        let nextTileY = this.currentTileY;
        let direction: Direction = "none";

        if (moveX) {
            nextTileX = this.currentTileX + Math.sign(dx);
            direction = dx > 0 ? "right" : "left";
        } else if (dy !== 0) {
            nextTileY = this.currentTileY + Math.sign(dy);
            direction = dy > 0 ? "down" : "up";
        }

        const startX = this.currentTileX * this.tileSize;
        const startY = this.currentTileY * this.tileSize;
        const nextX = nextTileX * this.tileSize;
        const nextY = nextTileY * this.tileSize;

        const tweenTarget = { progress: 0 };
        this.scene.tweens.add({
            targets: tweenTarget,
            progress: 1,
            duration: this.tweenDuration,
            ease: "Linear",
            onUpdate: (tween) => {
                // Interpolate pixel position between current and next tile
                const t = tween.progress;
                const currentX = startX + (nextX - startX) * t;
                const currentY = startY + (nextY - startY) * t;
                this.updateCallback(currentX, currentY, direction);
            },
            onComplete: () => {
                this.currentTileX = nextTileX;
                this.currentTileY = nextTileY;
                this.updateCallback(nextX, nextY, direction);
                this.tweenToNextTile();
            },
        });
    }
}*/