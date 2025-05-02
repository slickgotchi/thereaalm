import Phaser from "phaser";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { TILE_PIXELS } from "../GameScene";
import { Direction, TweenUpdateCallback, TweenWorker } from "../navigation/TweenWorker";
import { Pathfinder } from "../navigation/Pathfinder";
import { NavigationGrid } from "../navigation/NavigationGrid";

interface Props {
    scene: Phaser.Scene;
    id: string;
    zoneId: number;
    tileX: number;
    tileY: number;
    type: string;
    texture: string;
    data: any;
    navigationGrid: NavigationGrid;
}

export class TweenableEntity extends BaseEntity {
    protected tweenWorker: TweenWorker;
    private pathfinder: Pathfinder;
    protected direction: Direction = "none";

    private targetPosition = {x: 0, y: 0}
    private targetDirection: Direction = "none";

    public currentPosition = {x: 0, y: 0};
    public currentDirection: Direction = "none";

    constructor(props: Props) {
        const {scene, id, zoneId, tileX, tileY, type, texture, data, navigationGrid } = props;
        super({scene, id, zoneId, tileX, tileY, type, texture, data});

        this.currentPosition = {x: tileX * TILE_PIXELS, y: tileY * TILE_PIXELS}
        this.currentDirection = "down";

        const updateCallback: TweenUpdateCallback = (x: number, y: number, direction: Direction) => {
            // this.sprite.setPosition(x, y);
            this.currentPosition = {x, y}
            if (direction !== "none") {
                this.currentDirection = direction;
            }
        };

        this.tweenWorker = new TweenWorker(scene, updateCallback);
        this.pathfinder = new Pathfinder(navigationGrid);

        // listen for phaser updates
        scene.events.on("update", this.frameUpdate, this);
    }

    protected frameUpdate() {
        this.sprite.setPosition(this.currentPosition.x, this.currentPosition.y);

        this.outlineEffect.updatePosition();

        if (this.tweenWorker.getIsTweening()) {
            this.direction = this.currentDirection;
        } else {
            this.direction = this.targetDirection;
        }
        // this.updateDirection();
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        const lastTile = {x: this.tileX, y: this.tileY};

        super.snapshotUpdate(snapshot);
        this.targetDirection = snapshot.data.direction;
        this.targetPosition = {x: snapshot.tileX * TILE_PIXELS, y: snapshot.tileY * TILE_PIXELS}

        if (lastTile.x !== snapshot.tileX || lastTile.y !== snapshot.tileY && !this.tweenWorker.getIsTweening()) {
            // this.sprite.setPosition(lastTile.x * TILE_PIXELS, lastTile.y * TILE_PIXELS);
            this.currentPosition = {x: lastTile.x * TILE_PIXELS, y: lastTile.y * TILE_PIXELS}

            const waypoints = this.pathfinder.findPath(lastTile.x, lastTile.y, snapshot.tileX, snapshot.tileY, true);
            this.tweenWorker.tweenToWaypoints(lastTile.x, lastTile.y, waypoints);
        } else {
            this.currentPosition = this.targetPosition;
            this.currentDirection = this.targetDirection;
        }
    }

    // protected updateDirection() {
    //     // Default does nothing; subclasses can override
    //     // NOTE: GotchiEntity overrides this
    // }

    // Add destroy method to clean up
    destroy(): void {
        // Remove the update listener
        this.scene.events.off("update", this.frameUpdate, this);
        
        // Call parent destroy
        super.destroy();
    }
}
