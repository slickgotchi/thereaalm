import Phaser from "phaser";
import { fetchBulkGotchiSVGs } from "./FetchGotchis"; // Adjust path as needed
import { resizeGame } from "./ResizeGame";
import { WorldMap } from "./WorldMap";
import { CameraController } from "./CameraController";
import { TileMap } from "./TileMap";

interface GotchiPosition {
    uuid: string; // Matches server's "uuid"
    gotchiId: string; // Matches server's "gotchiId"
    x: number;
    y: number;
}

export const TILE_PIXELS = 64;
export const ZONE_TILES = 512;
const GAME_WIDTH = 1920;
const GAME_HEIGHT = 1200;

// Map to store Gotchi IDs and their state
interface GotchiState {
    sprite?: Phaser.GameObjects.Sprite; // Sprite reference
    svgState:
        | "ToBeFetched"
        | "Fetching"
        | "LoadingImage"
        | "ImageLoaded"
        | "ImageSet";
    position: { x: number; y: number }; // Store position for updates
    gotchiId: string; // Store gotchiId for SVG fetching
    zoneId: number;
}

export class GameScene extends Phaser.Scene {
    private gotchiMap: Map<string, GotchiState> = new Map();
    private worldMap!: WorldMap;
    private cameraController!: CameraController;
    private worldWidth: number = 10 * ZONE_TILES * TILE_PIXELS;
    private worldHeight: number = 10 * ZONE_TILES * TILE_PIXELS;
    private newGotchiSVGCounter = 0;
    private tileMap!: TileMap;

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
    }

    async create() {
        // create tilemap
        this.tileMap.create();

        this.gotchiMap.clear(); // Initialize the map
        console.log("Initialized gotchiMap:", this.gotchiMap.size); // Debug initial map state

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
                this.addOrUpdateGotchis(this.currentZoneIndex, data);
            });
            this.currentZoneIndex++;
            if (this.currentZoneIndex >= numZones) {
                this.currentZoneIndex = 0;
            }
        };

        // Call immediately
        fetchAndProcessZone();

        // Continue every 5000ms
        setInterval(fetchAndProcessZone, 10000);
    }

    shutdown() {
        // window.removeEventListener("resize", () => resizeGame(this));
    }

    private async fetchZoneSnapshot(zoneId: number): Promise<GotchiPosition[]> {
        try {
            const response = await fetch(
                `http://localhost:8080/zones/${zoneId}/snapshot`
            );
            const data = await response.json();
            console.log(`Fetched zone ${zoneId} snapshot data:`, data); // Debug snapshot data
            return data.gotchis || [];
        } catch (error) {
            console.error("Failed to fetch zone snapshot:", error);
            return [];
        }
    }

    private addOrUpdateGotchis(zoneId: number, gotchisData: GotchiPosition[]) {
        // Get all gotchis in our current zone from gotchiMap
        const currentIds = new Set(
            Array.from(this.gotchiMap.entries())
                .filter(([_, gotchi]) => gotchi.zoneId === zoneId)
                .map(([id, _]) => id)
        );

        // Go through gotchis data and delete currentIds if we have a match
        gotchisData.forEach((gotchi) => currentIds.delete(gotchi.uuid));

        // Delete all remaining currentIds as they no longer exist in this zone
        currentIds.forEach((uuid) => {
            const state = this.gotchiMap.get(uuid);
            if (state?.sprite) state.sprite.destroy();
            this.gotchiMap.delete(uuid);
        });

        // process all gotchis in gotchisData
        gotchisData.forEach((gotchi) => {
            // get existing gotchi state
            const existingState = this.gotchiMap.get(gotchi.uuid);

            // NEW GOTCHI - lets make a new gotchi for the gotchiMap
            if (!existingState) {
                const newX = gotchi.x * TILE_PIXELS;
                const newY = gotchi.y * TILE_PIXELS;

                // create loading sprite
                const loadingSprite = this.add
                    .sprite(newX, newY, "loading_gotchi")
                    .setDepth(2000)
                    .setScale(1)
                    .setName(gotchi.gotchiId);
                loadingSprite.play("loading_gotchi_anim");

                // New Gotchi: Initialize state and add to map
                this.gotchiMap.set(gotchi.uuid, {
                    sprite: loadingSprite,
                    svgState: "ToBeFetched",
                    position: {
                        x: newX,
                        y: newY,
                    },
                    gotchiId: gotchi.gotchiId, // Store gotchiId for SVG fetching,
                    zoneId: zoneId,
                });

                // rect to check sizing
                // this.add
                //     .rectangle(newX, newY, 64, 64, 0x131313)
                //     .setDepth(2100)
                //     .setAlpha(0.5);

                // console.log(
                //     `Added new Gotchi UUID ${gotchi.uuid} with gotchiId ${gotchi.gotchiId} to map`
                // );
            } else {
                // Update position for existing Gotchi
                existingState.position = {
                    x: gotchi.x * TILE_PIXELS,
                    y: gotchi.y * TILE_PIXELS,
                };

                // set new position for sprite
                if (existingState.sprite) {
                    existingState.sprite.setPosition(
                        existingState.position.x,
                        existingState.position.y
                    );
                }
                // console.log(`Updated position for Gotchi UUID ${gotchi.uuid}`); // Debug update
            }
        });

        // Compile Gotchis that have not had svgFetched and are not svg fetching
        const gotchisToFetch = Array.from(this.gotchiMap.entries())
            .filter(([uuid, state]) => {
                return state.svgState === "ToBeFetched";
            })
            .map(([uuid, state]) => state.gotchiId); // Use gotchiId for SVG fetching

        if (gotchisToFetch.length > 0) {
            // console.log("gotchisToFetch array (gotchiIds):", gotchisToFetch); // Updated log
            this.fetchAndLoadSVGs(gotchisToFetch, gotchisData);
        } else {
            console.log("No Gotchis to fetch SVGs for.");
        }
    }

    private async fetchAndLoadSVGs(
        gotchiIDs: string[],
        gotchisData: GotchiPosition[]
    ) {
        // console.log("Fetching SVGs for Gotchi IDs:", gotchiIDs); // Debug fetch
        // Mark as fetching to avoid duplicate requests
        gotchiIDs.forEach((gotchiId) => {
            // Find the UUID corresponding to this gotchiId
            const gotchi = gotchisData.find((g) => g.gotchiId === gotchiId);
            if (gotchi) {
                const state = this.gotchiMap.get(gotchi.uuid);
                if (state) state.svgState = "Fetching";
            }
        });

        try {
            const svgSets = await fetchBulkGotchiSVGs(gotchiIDs);

            svgSets.forEach((svgSet: any, index: number) => {
                const gotchiId = gotchiIDs[index];
                // Find the UUID corresponding to this gotchiId
                const gotchi = gotchisData.find((g) => g.gotchiId === gotchiId);
                if (gotchi) {
                    const state = this.gotchiMap.get(gotchi.uuid);
                    if (state) state.svgState = "LoadingImage";
                    this.loadGotchiSVG(gotchi.gotchiId, svgSet);
                }
            });
        } catch (error) {
            console.error("Failed to fetch bulk SVGs:", error);
        }
    }

    private async loadGotchiSVG(
        gotchiId: string,
        svgSet: { svg: string; left: string; right: string; back: string }
    ) {
        try {
            const views: ("svg" | "left" | "right" | "back")[] = [
                "svg",
                "left",
                "right",
                "back",
            ];

            // Track loaded views
            const loadedViews = new Set<string>();

            // Add each image and listen for individual load completion
            views.forEach((view) => {
                const svgDataUrl = `data:image/svg+xml;base64,${btoa(
                    svgSet[view] || ""
                )}`;

                const key = `gotchi-${gotchiId}-${view}`;

                this.load.image(key, svgDataUrl);

                // Listen for individual image load
                this.load.once(`filecomplete-image-${key}`, () => {
                    loadedViews.add(view);
                    if (loadedViews.size === views.length) {
                        this.onAllImagesLoaded(gotchiId);
                    }
                });
            });

            this.load.start();
        } catch (err) {
            console.error(
                "Failed to load Gotchi SVG for UUID",
                gotchiId,
                ":",
                err
            );
        }
    }

    private onAllImagesLoaded(gotchiId: string) {
        // console.log("allImagesLoaded for: ", gotchiId);

        // Search through gotchiMap values to find the matching gotchiId
        const state = Array.from(this.gotchiMap.values()).find(
            (gotchi) => gotchi.gotchiId === gotchiId
        );

        if (!state) return; // Gotchi might have been removed

        const { x, y } = state.position;

        if (state.sprite) {
            // Update the existing sprite
            state.sprite.stop();
            state.sprite.setTexture(`gotchi-${gotchiId}-svg`);
            state.sprite.setPosition(x, y);
            console.log("set new texture");
            // console.log("Updated texture for: ", `gotchi-${gotchiId}-svg`);
        } else {
            // Create a new sprite
            // state.sprite = this.add
            //     .sprite(x, y, `gotchi-${gotchiId}-svg`)
            //     .setDepth(2000)
            //     .setScale(0.5)
            //     .setName(gotchiId);
            // console.log("Created new texture for: ", `gotchi-${gotchiId}-svg`);
        }

        this.newGotchiSVGCounter++;
        console.log("gotchi SVG count: ", this.newGotchiSVGCounter);

        // Mark as fully loaded
        state.svgState = "ImageSet";
        // console.log(`Marked Gotchi UUID ${gotchiId} as svgFetched: true`);
    }
}
