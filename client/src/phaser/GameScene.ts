import Phaser from "phaser";
import { WorldMap } from "./WorldMap";
import { CameraController } from "./CameraController";
import { TileMap } from "./TileMap";
import { HPBar } from "./HPBar";
import { BaseEntity, EntitySnapshot } from "./entities/BaseEntity";
import { EntityFactory } from "./entities/EntityFactory";
import { GotchiEntity } from "./entities/GotchiEntity";
import { NavigationGrid } from "./navigation/NavigationGrid";
import { SelectionManager } from "./SelectionManager";
import { EmoticonEmitter } from "./emoticons/EmoticonEmitter";
import { VFXManager } from "./VFXManager";

export const TILE_PIXELS = 64;
export const ZONE_TILES = 512;

export class GameScene extends Phaser.Scene {
    private entityMap: Map<string, BaseEntity> = new Map();
    private worldMap!: WorldMap;
    public cameraController!: CameraController;
    private worldWidth: number = 10 * ZONE_TILES * TILE_PIXELS;
    private worldHeight: number = 10 * ZONE_TILES * TILE_PIXELS;
    private tileMap!: TileMap;

    private navigationGrid!: NavigationGrid;

    private hpBars: HPBar[] = [];

    private currentZoneIndex = 0;

    // private selectedEntity: BaseEntity | null = null;
    selectionManager!: SelectionManager;
    vfxManager!: VFXManager;

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

        EmoticonEmitter.preload(this);
        VFXManager.preload(this);

        this.load.spritesheet(
            "static_entities",
            "assets/spritesheets/static_entities_spritesheet.png",
            {frameWidth: 64, frameHeight: 64}
        );
        // this.load.image("berrybush", "assets/images/berrybush.png");
        this.load.image("shop", "assets/images/shop.png");
        this.load.image("berry_icon", "assets/images/berry_icon.png");
        this.load.image("altar", "assets/images/golden_altar_l1.png");
        this.load.image("shadow", "assets/images/shadow.png");
        this.load.image("default_gotchi_svg", "assets/images/logo-gotchi-front.png");
        this.load.image("default_gotchi_left", "assets/images/logo-gotchi-left.png");
        this.load.image("default_gotchi_right", "assets/images/logo-gotchi-right.png");
        this.load.image("default_gotchi_back", "assets/images/logo-gotchi-back.png");
        this.load.image("fomoberry_bush", "assets/images/fomoberry_bush.png");
        this.load.image("kekwood_tree", "assets/images/kekwood_tree.png");
        this.load.image("alphaslate_boulders", "assets/images/alphaslate_boulders.png");
    }

    public getVFXManager(): VFXManager {
        return this.vfxManager;
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

        this.vfxManager = new VFXManager(this);

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

        this.selectionManager = new SelectionManager(this);
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
        // Get all entities in the current zone from entityMap to start
        // The list of entities to delete
        const deleteEntityIds = new Set<string>();

        // Iterate over entityMap entries to find entities in the given zone
        for (const [id, entity] of this.entityMap.entries()) {
            // Check if the entity's zoneId matches the given zoneId
            if (Number(entity.zoneId) === zoneId) {
                deleteEntityIds.add(id);
            } else {
                console.log(`Entity ID ${id} does NOT match zoneId ${zoneId} (entity.zoneId: ${entity.zoneId})`);
            }
        }

        // if entity is in entitySnapshots, remove it from deletion list
        entitySnapshots.forEach((entity) => deleteEntityIds.delete(entity.id));

        // DESTROY ENTITIES remaining ids as they are no longer in the zone
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