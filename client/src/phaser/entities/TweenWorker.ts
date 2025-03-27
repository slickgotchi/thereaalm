import Phaser from "phaser";

export type Direction = "up" | "down" | "left" | "right" | "none";
export type TweenUpdateCallback = (x: number, y: number, direction: Direction, isComplete?: boolean) => void;

interface Waypoint {
    tileX: number;
    tileY: number;
    direction: Direction;
}

export class TweenWorker {
    private scene: Phaser.Scene;
    private tileSize: number = 64;
    private tweenDurationPerTile: number = 200; // 200ms per waypoint
    private totalTweenDuration: number = 0; // Calculated based on waypoints
    private updateCallback: TweenUpdateCallback;
    private isTweening: boolean = false;
    private waypoints: Waypoint[] = [];
    private currentTween: Phaser.Tweens.Tween | null = null;

    constructor(scene: Phaser.Scene, updateCallback: TweenUpdateCallback) {
        this.scene = scene;
        this.updateCallback = updateCallback;
    }

    public getIsTweening(): boolean {
        return this.isTweening;
    }

    public setTargetTilePosition(startTileX: number, startTileY: number, targetTileX: number, targetTileY: number) {
        if (startTileX === targetTileX && startTileY === targetTileY) {
            this.updateCallback(startTileX * this.tileSize, startTileY * this.tileSize, "none");
            return;
        }

        // Stop any existing tween
        if (this.isTweening && this.currentTween) {
            this.currentTween.stop();
            this.isTweening = false;
        }

        // Generate waypoints
        this.waypoints = this.generateWaypoints(startTileX, startTileY, targetTileX, targetTileY);

        // Calculate total duration based on waypoints (200ms per tile)
        this.totalTweenDuration = this.waypoints.length * this.tweenDurationPerTile;

        // Start the tween
        this.startTweening(startTileX, startTileY);
    }

    private generateWaypoints(startTileX: number, startTileY: number, targetTileX: number, targetTileY: number): Waypoint[] {
        const waypoints: Waypoint[] = [];
        let currentX = startTileX;
        let currentY = startTileY;

        while (currentX !== targetTileX || currentY !== targetTileY) {
            const dx = targetTileX - currentX;
            const dy = targetTileY - currentY;
            const moveX = Math.random() < 0.5 && dx !== 0; // Randomly choose X or Y if both available
            let direction: Direction = "none";

            if (moveX) {
                currentX += Math.sign(dx);
                direction = dx > 0 ? "right" : "left";
            } else if (dy !== 0) {
                currentY += Math.sign(dy);
                direction = dy > 0 ? "down" : "up";
            }

            waypoints.push({ tileX: currentX, tileY: currentY, direction });
        }

        return waypoints;
    }

    private startTweening(startTileX: number, startTileY: number) {
        if (this.waypoints.length === 0) {
            this.isTweening = false;
            return;
        }

        this.isTweening = true;
        const startX = startTileX * this.tileSize;
        const startY = startTileY * this.tileSize;
        const totalTiles = this.waypoints.length;

        // Set initial position explicitly to avoid flicker
        this.updateCallback(startX, startY, "none");

        const tweenTarget = { progress: 0 };
        this.currentTween = this.scene.tweens.add({
            targets: tweenTarget,
            progress: 1, // Progress from 0 to 1 over the total duration
            duration: this.totalTweenDuration,
            ease: "Linear",
            onUpdate: (tween) => {
                const t = tween.progress * totalTiles; // Scale progress to waypoint count (0 to totalTiles)
                const waypointIndex = Math.min(Math.floor(t), totalTiles - 1); // Clamp to last waypoint
                const subProgress = t - waypointIndex; // Progress between current and next waypoint

                const prevWaypoint = waypointIndex === 0 ? { tileX: startTileX, tileY: startTileY, direction: "none" as Direction } : this.waypoints[waypointIndex - 1];
                const currentWaypoint = this.waypoints[waypointIndex];

                const prevX = prevWaypoint.tileX * this.tileSize;
                const prevY = prevWaypoint.tileY * this.tileSize;
                const nextX = currentWaypoint.tileX * this.tileSize;
                const nextY = currentWaypoint.tileY * this.tileSize;

                const currentX = prevX + (nextX - prevX) * subProgress;
                const currentY = prevY + (nextY - prevY) * subProgress;

                this.updateCallback(currentX, currentY, currentWaypoint.direction);
            },
            onComplete: () => {
                this.isTweening = false;
                this.currentTween = null;
                const finalWaypoint = this.waypoints[this.waypoints.length - 1];
                this.updateCallback(finalWaypoint.tileX * this.tileSize, finalWaypoint.tileY * this.tileSize, "none", true);
                this.waypoints = [];
                
            },
        });
    }
}