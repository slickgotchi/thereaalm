import Phaser from "phaser";
import { fetchBulkGotchiSVGs } from "./FetchGotchis"; // Adjust path as needed
import { resizeGame } from "./ResizeGame";

interface GotchiPosition {
    uuid: string; // Matches server's "uuid"
    gotchiId: string; // Matches server's "gotchiId"
    x: number;
    y: number;
}

const TILE_TO_PIXELS = 64;
const GAME_WIDTH = 1920;
const GAME_HEIGHT = 1200;

// Map to store Gotchi IDs and their state
interface GotchiState {
    svgState:
        | "ToBeFetched"
        | "Fetching"
        | "LoadingImage"
        | "ImageLoaded"
        | "ImageSet";
    sprite?: Phaser.GameObjects.Sprite; // Sprite reference
    position: { x: number; y: number }; // Store position for updates
    gotchiId: string; // Store gotchiId for SVG fetching
}

export class GameScene extends Phaser.Scene {
    private gotchiMap: Map<string, GotchiState> = new Map();

    constructor() {
        super("GameScene");
    }

    create() {
        this.gotchiMap.clear(); // Initialize the map
        console.log("Initialized gotchiMap:", this.gotchiMap.size); // Debug initial map state

        const rect = this.add
            .rectangle(100, 100, 100, 100, 0xffffff)
            .setOrigin(0, 0)
            .setAlpha(0.2);

        const circ = this.add.circle(0, 0, 5, 0xffffff).setOrigin(0.5, 0.5);

        // resizeGame(this);
        // window.addEventListener("resize", () => resizeGame(this));

        // Fetch initial snapshot for Zone 0
        this.fetchZoneSnapshot(0).then((data) => {
            this.addOrUpdateGotchis(data);
        });

        // Fetch snapshot every 10 seconds
        this.time.addEvent({
            delay: 10000, // 10 seconds
            callback: () => {
                this.fetchZoneSnapshot(0).then((data) => {
                    this.addOrUpdateGotchis(data);
                });
            },
            loop: true,
        });
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
            console.log("Fetched snapshot data:", data); // Debug snapshot data
            return data.gotchis || [];
        } catch (error) {
            console.error("Failed to fetch zone snapshot:", error);
            return [];
        }
    }

    private addOrUpdateGotchis(gotchisData: GotchiPosition[]) {
        console.log("Processing snapshot with Gotchis:", gotchisData); // Debug input data
        // console.log(
        //     "Current gotchiMap before update:",
        //     Array.from(this.gotchiMap.entries())
        // ); // Debug map state

        // get all current gotchi ids in this.gotchiMap
        const currentIds = new Set(this.gotchiMap.keys());

        // go through gotchisData received from latest snapshot
        // - iterating over each gotchisData datum, delete items from currentIds that match
        // - this leaves us with currentIds that are no longer in the snapshot
        gotchisData.forEach((gotchi) => currentIds.delete(gotchi.uuid));

        // delete items in gotchiMap that no longer exist according to the latest snapshot
        currentIds.forEach((uuid) => {
            const state = this.gotchiMap.get(uuid);
            if (state?.sprite) state.sprite.destroy();
            this.gotchiMap.delete(uuid);
            console.log(`Removed Gotchi UUID ${uuid} from map`); // Debug removal
        });

        // now process all gotchis in gotchisData
        gotchisData.forEach((gotchi) => {
            // get existing gotchi state
            const existingState = this.gotchiMap.get(gotchi.uuid);

            // NEW GOTCHI - lets make a new gotchi for the gotchiMap
            if (!existingState) {
                const newX = gotchi.x * TILE_TO_PIXELS * 0.1;
                const newY = gotchi.y * TILE_TO_PIXELS * 0.1;
                // New Gotchi: Initialize state and add to map
                this.gotchiMap.set(gotchi.uuid, {
                    svgState: "ToBeFetched",
                    position: {
                        x: newX,
                        y: newY,
                    },
                    gotchiId: gotchi.gotchiId, // Store gotchiId for SVG fetching
                });

                console.log(
                    `Added new Gotchi UUID ${gotchi.uuid} with gotchiId ${gotchi.gotchiId} to map`
                );
            } else {
                // Update position for existing Gotchi
                existingState.position = {
                    x: gotchi.x * TILE_TO_PIXELS,
                    y: gotchi.y * TILE_TO_PIXELS,
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
            console.log("gotchisToFetch array (gotchiIds):", gotchisToFetch); // Updated log
            this.fetchAndLoadSVGs(gotchisToFetch, gotchisData);
        } else {
            console.log("No Gotchis to fetch SVGs for.");
        }
    }

    private async fetchAndLoadSVGs(
        gotchiIDs: string[],
        gotchisData: GotchiPosition[]
    ) {
        console.log("Fetching SVGs for Gotchi IDs:", gotchiIDs); // Debug fetch
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
            console.log("Fetched SVG sets:", svgSets); // Debug SVG sets
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
        // finally {
        //     // Reset fetching state
        //     gotchiIDs.forEach((gotchiId) => {
        //         const gotchi = gotchisData.find((g) => g.gotchiId === gotchiId);
        //         if (gotchi) {
        //             const state = this.gotchiMap.get(gotchi.uuid);
        //             if (state) {
        //                 state.svgFetching = false;
        //                 console.log(
        //                     `Reset svgFetching to false for Gotchi UUID ${gotchi.uuid}`
        //                 ); // Debug reset
        //             }
        //         }
        //     });
        // }
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
                    console.log(`Loaded: ${key}`);

                    // If all views are loaded, mark as ImageLoaded
                    console.log(loadedViews.size, views.length);
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
        console.log("allImagesLoaded for: ", gotchiId);

        // Search through gotchiMap values to find the matching gotchiId
        const state = Array.from(this.gotchiMap.values()).find(
            (gotchi) => gotchi.gotchiId === gotchiId
        );

        if (!state) return; // Gotchi might have been removed

        const { x, y } = state.position;

        if (state.sprite) {
            // Update the existing sprite
            state.sprite.setTexture(`gotchi-${gotchiId}-svg`);
            console.log("Updated texture for: ", `gotchi-${gotchiId}-svg`);
        } else {
            // Create a new sprite
            state.sprite = this.add
                .sprite(x, y, `gotchi-${gotchiId}-svg`)
                .setDepth(1000)
                .setScale(0.5)
                .setName(gotchiId);
            console.log("Created new texture for: ", `gotchi-${gotchiId}-svg`);
        }

        // Mark as fully loaded
        state.svgState = "ImageSet";
        console.log(`Marked Gotchi UUID ${gotchiId} as svgFetched: true`);
    }
}
