import { EmoticonEmitter } from "../emoticons/EmoticonEmitter";
import { fetchBulkGotchiSVGs, GotchiSVGSet } from "../FetchGotchis";
import { ZONE_TILES } from "../GameScene";
import { NavigationGrid } from "../navigation/NavigationGrid";
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
    lastEmoticonEmitTime_ms: number = 0;
    emoticonEmitInterval_ms: number = 3000;

    currentActionType: string = "";

    jumpY = 0;

    shadowSprite!: Phaser.GameObjects.Sprite;

    // static map of svg states
    static svgMap: Map<string, SVGState> = new Map();

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any, navigationGrid: NavigationGrid) {
        super(scene, id, zoneId, tileX, tileY, "gotchi", "loading_gotchi", data, navigationGrid);
        
        // Add animation
        this.sprite.play("loading_gotchi_anim");

        this.shadowSprite = scene.add.sprite(this.sprite.x+32, this.sprite.y+64, "shadow");
        this.shadowSprite.setDepth(this.sprite.depth-1);
        this.shadowSprite.setOrigin(0.5,0.5);
        this.shadowSprite.setAlpha(0.3);

        this.gotchiId = data.gotchiId;

        // add this gotchi to svg map if its not the default
        if (this.gotchiId !== "69420") {
            GotchiEntity.svgMap.set(this.gotchiId, {
                gotchiId: this.gotchiId,
                svgState: "ToBeFetched",
            });
        } else {
            // we set the default aavegotchi images to our texture set
            this.textureSet = {
                svg: "default_gotchi_svg",
                left: "default_gotchi_left",
                right: "default_gotchi_right",
                back: "default_gotchi_back",
            }
        }

        const delay_ms = Math.random() * 500 // Random delay between 0 and 500ms

        scene.tweens.add({
            targets: this,
            jumpY: -8,
            duration: 150,
            yoyo: true,
            repeat: -1,
            ease: `Quad.easeOut`,
            delay: delay_ms,
        });

        this.emoticonEmitter = new EmoticonEmitter(scene, 
            tileX*ZONE_TILES, tileY*ZONE_TILES);
    }

    protected frameUpdate(): void {
        super.frameUpdate();

        this.sprite.setPosition(this.currentPosition.x, this.currentPosition.y - this.jumpY);
        this.shadowSprite.setPosition(this.currentPosition.x+32, this.currentPosition.y+64);
        this.emoticonEmitter.setPosition(this.currentPosition.x+32, this.currentPosition.y+16);

        // if tweening we reset our timer
        if (this.tweenWorker.getIsTweening()) {
            this.lastEmoticonEmitTime_ms = 0;
        }

        var currTime_ms = Date.now();
        if (currTime_ms - this.lastEmoticonEmitTime_ms > this.emoticonEmitInterval_ms
            && !this.tweenWorker.getIsTweening() && this.currentActionType !== ""
        ) {
            this.lastEmoticonEmitTime_ms = currTime_ms;
            this.emoticonEmitter.emit(this.currentActionType, 750);
        }
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);

        this.updateSVG();


        var currentAction = snapshot.data.actionPlan.currentAction;
        if (currentAction) {
            this.currentActionType = currentAction.type;
        }
    }

    updateSVG() {
        // check svg is loaded
        if (!this.textureSet) {
            // try get svg map
            const svgMapItem = GotchiEntity.svgMap.get(this.gotchiId);
            if (!svgMapItem) return;
            if (svgMapItem.svgState === "ImageLoaded") {
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
                this.sprite.setTexture(this.textureSet.left);
                break;
            case "right":
                this.sprite.setTexture(this.textureSet.right);
                break;
            case "up":
                this.sprite.setTexture(this.textureSet.back);
                break;
            case "down":
            case "none":
            default:
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
