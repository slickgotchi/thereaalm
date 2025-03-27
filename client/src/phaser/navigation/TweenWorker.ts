import Phaser from "phaser";

export type Direction = "up" | "down" | "left" | "right" | "none";
export type TweenUpdateCallback = (x: number, y: number, direction: Direction, isComplete?: boolean) => void;

export interface Waypoint {
    tileX: number;
    tileY: number;
    direction: Direction;
}

export class TweenWorker {
    private scene: Phaser.Scene;
    private tileSize: number = 64;
    private tweenDurationPerTile: number = 200; // Fixed 200ms per tile
    private totalTweenDuration: number = 0;
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

    public tweenToWaypoints(startTileX: number, startTileY: number, waypoints: Waypoint[]) {
        if (waypoints.length === 0) {
            this.updateCallback(startTileX * this.tileSize, startTileY * this.tileSize, "none");
            return;
        }

        if (this.isTweening && this.currentTween) {
            this.currentTween.stop();
            this.isTweening = false;
        }

        this.waypoints = waypoints;
        this.totalTweenDuration = this.waypoints.length * this.tweenDurationPerTile; // 200ms per tile
        console.log(`[TweenWorker.tweenToWaypoints] Starting tween from (${startTileX}, ${startTileY}) with ${waypoints.length} waypoints, total duration=${this.totalTweenDuration}ms`);
        this.startTweening(startTileX, startTileY);
    }

    private startTweening(startTileX: number, startTileY: number) {
        this.isTweening = true;
        const startX = startTileX * this.tileSize;
        const startY = startTileY * this.tileSize;
        const totalTiles = this.waypoints.length;

        this.updateCallback(startX, startY, "none");

        const tweenTarget = { progress: 0 };
        this.currentTween = this.scene.tweens.add({
            targets: tweenTarget,
            progress: 1,
            duration: this.totalTweenDuration,
            ease: "Linear",
            onUpdate: (tween) => {
                const t = tween.progress * totalTiles; // Progress scaled to number of waypoints
                const waypointIndex = Math.min(Math.floor(t), totalTiles - 1);
                const subProgress = t - waypointIndex;

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
                this.updateCallback(finalWaypoint.tileX * this.tileSize, finalWaypoint.tileY * this.tileSize, "none");
                console.log(`[TweenWorker.startTweening] Tween completed to (${finalWaypoint.tileX}, ${finalWaypoint.tileY})`);
                this.waypoints = [];
            },
        });
    }
}