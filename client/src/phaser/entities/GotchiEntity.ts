import { fetchBulkGotchiSVGs } from "../FetchGotchis";
import { HPBar } from "../HPBar";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";

interface SVGState {
    gotchiId: string;
    svgState:
        | "ToBeFetched"
        | "Fetching"
        | "LoadingImage"
        | "ImageLoaded"
        | "ImageSet";
    entity: BaseEntity;
    // sprite: Phaser.GameObjects.Sprite;
}

export class GotchiEntity extends BaseEntity {
    // hpBar: HPBar;
    gotchiId: string = "";
    
    // static map of svg states
    static svgMap: Map<string, SVGState> = new Map();

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any) {
        super(scene, id, zoneId, tileX, tileY, "gotchi", "loading_gotchi", data);
        
        // Add animation
        this.sprite.play("loading_gotchi_anim");

        this.gotchiId = data.gotchiId;

        // add this gotchi to svg map
        GotchiEntity.svgMap.set(this.gotchiId, {
            gotchiId: this.gotchiId,
            svgState: "ToBeFetched",
            entity: this,
        });
    }

    update(snapshot: EntitySnapshot) {
        super.update(snapshot);
    }

    destroy(): void {
        super.destroy();
    }

    // static function to update all gotchi entity svgs
    static async fetchAndLoadSVGs(scene: Phaser.Scene) {
        // filter out all the gotchis that need svg fetched
        const gotchisToFetch = Array.from(GotchiEntity.svgMap.entries())
                .filter(([id, state]) => {
                return state.svgState === "ToBeFetched";
            })
            .map(([id, state]) => state.gotchiId); // Use gotchiId for SVG fetching

        if (gotchisToFetch.length <= 0) return;

        // set states to fetching
        gotchisToFetch.forEach(gotchiId => {
            const state = GotchiEntity.svgMap.get(gotchiId);
            if (state) state.svgState = "Fetching";
        });


        // try fetch
        try {
            const svgSets = await fetchBulkGotchiSVGs(gotchisToFetch);

            svgSets.forEach((svgSet: any, index: number) => {
                const gotchiId = gotchisToFetch[index];
                const state = GotchiEntity.svgMap.get(gotchiId);
                if (state) {
                    state.svgState = "LoadingImage";
                    this.loadGotchiSVG(scene, gotchiId, svgSet);
                } 
            });
        } catch (error) {
            console.error("Failed to fetch bulk SVGs:", error);
        }
    }

    static async loadGotchiSVG(
        scene: Phaser.Scene,
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

                scene.load.image(key, svgDataUrl);

                // Listen for individual image load
                scene.load.once(`filecomplete-image-${key}`, () => {
                    loadedViews.add(view);
                    if (loadedViews.size === views.length) {
                        this.onAllImagesLoaded(gotchiId);
                    }
                });
            });

            scene.load.start();
        } catch (err) {
            console.error(
                "Failed to load Gotchi SVG for UUID",
                gotchiId,
                ":",
                err
            );
        }
    }

    static onAllImagesLoaded(gotchiId: string) {
        // Search through gotchiMap values to find the matching gotchiId
        const state = Array.from(GotchiEntity.svgMap.values()).find(
            (gotchi) => gotchi.gotchiId === gotchiId
        );

        if (!state) return; // Gotchi might have been removed

        if (state.entity.sprite) {
            // Update the existing sprite
            state.entity.sprite.stop();
            state.entity.sprite.setTexture(`gotchi-${gotchiId}-svg`);

            // recalc outline
            state.entity.outlineEffect.rebuild();
            state.svgState = "ImageSet";
        } 
    }
}
