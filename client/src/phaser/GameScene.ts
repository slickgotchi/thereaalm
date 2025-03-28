import Phaser from "phaser";
import { WorldMap } from "./WorldMap";
import { CameraController } from "./CameraController";
import { TileMap } from "./TileMap";
import { HPBar } from "./HPBar";
import { BaseEntity, EntitySnapshot } from "./entities/BaseEntity";
import { EntityFactory } from "./entities/EntityFactory";
import { GotchiEntity } from "./entities/GotchiEntity";
import { NavigationGrid } from "./navigation/NavigationGrid";



export const TILE_PIXELS = 64;
export const ZONE_TILES = 512;
// const GAME_WIDTH = 1920;
// const GAME_HEIGHT = 1200;

export class GameScene extends Phaser.Scene {
    private entityMap: Map<string, BaseEntity> = new Map();
    private worldMap!: WorldMap;
    private cameraController!: CameraController;
    private worldWidth: number = 10 * ZONE_TILES * TILE_PIXELS;
    private worldHeight: number = 10 * ZONE_TILES * TILE_PIXELS;
    private tileMap!: TileMap;

    private navigationGrid!: NavigationGrid;

    private hpBars: HPBar[] = [];

    private currentZoneIndex = 0;

    constructor() {
        super("GameScene");
    }

    preload() {
        this.tileMap = new TileMap(
            this,
            "yield_fields_2",
            "terrain_tileset",
            "terrain.png"
        );
        this.tileMap.preload();
        this.load.spritesheet(
            "loading_gotchi",
            "assets/spritesheets/loading_gotchi_spritesheet.png",
            { frameWidth: 64, frameHeight: 64 }
        );
        this.load.spritesheet(
            "lickquidator",
            "assets/spritesheets/lickquidator_spritesheet.png",
            { frameWidth: 64, frameHeight: 64 }
        );
        this.load.spritesheet(
            "emoticons",
            "assets/emoticons/emoticons_48px.png",
            {frameWidth: 48, frameHeight: 48, margin: 2, spacing: 4}
        )
        this.load.image("berrybush", "assets/images/berrybush.png");
        this.load.image("shop", "assets/images/shop.png");
        this.load.image("berry_icon", "assets/images/berry_icon.png");
        this.load.image("altar", "assets/images/golden_altar_l1.png");
    }

    async create() {
        // create tilemap
        this.tileMap.create();

        this.navigationGrid = new NavigationGrid(10*ZONE_TILES,10*ZONE_TILES);

        // Set up the camera controller
        this.cameraController = new CameraController(
            this,
            this.worldWidth,
            this.worldHeight
        );

        // Create and draw the world map
        this.worldMap = new WorldMap(this);
        const numZones = await this.worldMap.draw();

        // create loading gtochi animation
        this.anims.create({
            key: "loading_gotchi_anim",
            frames: this.anims.generateFrameNumbers("loading_gotchi", {
                start: 0,
                end: 7,
            }),
            frameRate: 10,
            repeat: -1,
        });

        // Function to fetch and process zone snapshot
        const fetchAndProcessZone = () => {
            // lets just focus on zone 36 (yield fields) initially
            this.currentZoneIndex = 42;

            this.fetchZoneSnapshot(this.currentZoneIndex).then((data) => {
                // console.log(data);
                this.addOrUpdateEntities(this.currentZoneIndex, data);
            });
            /*
            this.currentZoneIndex++;
            if (this.currentZoneIndex >= numZones) {
                this.currentZoneIndex = 0;
            }
                */
        };

        // Call immediately
        fetchAndProcessZone();

        // Continue every 5000ms
        setInterval(fetchAndProcessZone, 3000);
    }

    update(time: number, delta: number): void {
        this.hpBars.forEach(hpBar => {
            hpBar.updateHP(100);
        })
    }

    postUpdate(time: number, delta: number) {
        this.events.emit("postUpdate");
    }

    shutdown() {
        // window.removeEventListener("resize", () => resizeGame(this));
    }

    private async fetchZoneSnapshot(zoneId: number): Promise<EntitySnapshot[]> {
        try {
            const response = await fetch(
                `http://localhost:8080/zones/${zoneId}/snapshot`
            );
            const data = await response.json();
            // console.log(`Fetched zone ${zoneId} snapshot data:`, data); // Debug snapshot data
            return data.entitySnapshots || [];
        } catch (error) {
            console.error("Failed to fetch zone snapshot:", error);
            return [];
        }
    }

    private addOrUpdateEntities(zoneId: number, entitySnapshots: EntitySnapshot[]) {
        // get all entities in current zone from entityMap to start
        // the list of entities to delete
        const deleteEntityIds = new Set(
            Array.from(this.entityMap.entries())
                .filter(([_, entity]) => entity.zoneId === zoneId)
                .map(([id, _]) => id)
        );

        // if entity is in entitySnapshots, remove it from deletion list
        entitySnapshots.forEach((entity) => deleteEntityIds.delete(entity.id));

        // delete remaining ids as they are no longer in the zone
        deleteEntityIds.forEach((id) => {
            const baseEntity = this.entityMap.get(id);
            baseEntity?.destroy();
            this.entityMap.delete(id);
        });

        // process all entities
        entitySnapshots.forEach(entitySnapshot => {
            const existingState = this.entityMap.get(entitySnapshot.id);

            // NEW ENTITY
            if (!existingState) {
                const newEntity = EntityFactory.create(this, entitySnapshot, this.navigationGrid);
                this.entityMap.set(entitySnapshot.id, newEntity);
            }
            // EXISTING ENTITY
            else {
                existingState.snapshotUpdate(entitySnapshot);
            }
        });

        GotchiEntity.fetchAndLoadSVGs(this);
    }
}
