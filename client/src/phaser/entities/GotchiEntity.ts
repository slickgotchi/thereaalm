import { EmoticonEmitter } from "../emoticons/EmoticonEmitter";
import { fetchBulkGotchiSVGs, GotchiSVGSet } from "../FetchGotchis";
import { ZONE_TILES } from "../GameScene";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { Direction } from "../navigation/TweenWorker";
import { EntitySnapshot } from "./BaseEntity";
import { TweenableEntity } from "./TweenableEntity";

interface SVGState {
    gotchiId: string;
    svgState:
        | "ToBeFetched"
        | "Fetching"
        | "LoadingImage"
        | "ImageLoaded";
}

interface TextureSet {
    svg: string;
    left: string;
    right: string;
    back: string;
}

export class GotchiEntity extends TweenableEntity {
    // hpBar: HPBar;
    gotchiId: string = "";
    textureSet!: TextureSet;

    emoticonEmitter!: EmoticonEmitter;

    jumpY = 0;

    // static map of svg states
    static svgMap: Map<string, SVGState> = new Map();

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any, navigationGrid: NavigationGrid) {
        super(scene, id, zoneId, tileX, tileY, "gotchi", "loading_gotchi", data, navigationGrid);
        
        // Add animation
        this.sprite.play("loading_gotchi_anim");

        // var test = scene.add.sprite(this.sprite.x, this.sprite.y, "loading_gotchi");
        // test.setDepth(10000);
        // test.play("loading_gotchi_anim");

        // console.log("play loading anim");
        console.log(this.sprite.x, this.sprite.y);

        this.gotchiId = data.gotchiId;

        // add this gotchi to svg map
        GotchiEntity.svgMap.set(this.gotchiId, {
            gotchiId: this.gotchiId,
            svgState: "ToBeFetched",
        });

        this.emoticonEmitter = new EmoticonEmitter(scene, "emoticons", 
            tileX*ZONE_TILES, tileY*ZONE_TILES, 0);

        scene.tweens.add({
            targets: this,
            jumpY: -8,
            duration: 250,
            yoyo: true,
            repeat: -1,
            ease: `Quad.easeOut`
        });
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);

        this.updateSVG();
    }

    updateSVG() {
        // check svg is loaded
        if (!this.textureSet) {
            // try get svg map
            const svgMapItem = GotchiEntity.svgMap.get(this.gotchiId);
            if (!svgMapItem) return;
            if (svgMapItem.svgState === "ImageLoaded") {
                console.log("apply svg");
                console.log(this.sprite.x, this.sprite.y);
                this.textureSet = {
                    svg: `gotchi-${this.gotchiId}-svg`,
                    left: `gotchi-${this.gotchiId}-left`,
                    right: `gotchi-${this.gotchiId}-right`,
                    back: `gotchi-${this.gotchiId}-back`,
                }
                this.sprite.stop();
                this.sprite.setTexture(this.textureSet.svg);
                this.outlineEffect.rebuild();
            }
        }
    }

    protected updateDirection() {
        if (!this.textureSet) return;
        this.sprite.stop();
        switch (this.direction) {
            case "left":
                // console.log("LEFT");
                this.sprite.setTexture(this.textureSet.left);
                break;
            case "right":
                // console.log("RIGHT");
                this.sprite.setTexture(this.textureSet.right);
                break;
            case "up":
                // console.log("UP");
                this.sprite.setTexture(this.textureSet.back);
                break;
            case "down":
            case "none":
            default:
                // console.log("DOWN");
                this.sprite.setTexture(this.textureSet.svg);
                break;
        }
    }

    destroy(): void {
        super.destroy();
    }

    static async fetchAndLoadSVGs(scene: Phaser.Scene) {
        const gotchisToFetch = Array.from(GotchiEntity.svgMap.entries())
            .filter(([_, state]) => state.svgState === "ToBeFetched")
            .map(([_, state]) => state.gotchiId);
    
        if (gotchisToFetch.length <= 0) return;
    
        gotchisToFetch.forEach(gotchiId => {
            const state = GotchiEntity.svgMap.get(gotchiId);
            if (state) state.svgState = "Fetching";
        });
    
        try {
            const svgSets: GotchiSVGSet[] = await fetchBulkGotchiSVGs(gotchisToFetch);
            svgSets.forEach((svgSet) => {
                const gotchiId = svgSet.id; // Use 'id' from GotchiSVGSet
                const state = GotchiEntity.svgMap.get(gotchiId);
                if (state) {
                    state.svgState = "LoadingImage";
                    // Pass the full svgSet minus the id as the texture data
                    const { svg, left, right, back } = svgSet;
                    this.loadGotchiSVG(scene, gotchiId, { svg, left, right, back });
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
            (svgState) => svgState.gotchiId === gotchiId
        );

        if (!state) return; // Gotchi might have been removed
        state.svgState = "ImageLoaded";
    }
}
