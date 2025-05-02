import { Emoticons } from "../emoticons/Emoticons";
import { ESPBar } from "../ESPBar";
import { fetchBulkGotchiSVGs, GotchiSVGSet } from "../../utils/gotchi-loader/FetchGotchis";
import { TILE_PIXELS } from "../GameScene";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { EntitySnapshot } from "./BaseEntity";
import { TweenableEntity } from "./TweenableEntity";

interface SVGState {
    gotchiId: string;
    svgState: "ToBeFetched" | "Fetching" | "LoadingImage" | "ImageLoaded";
}

// this is used to access the specific texture we are after
interface TextureSet {
    front: string;
    left: string;
    right: string;
    back: string;
}

export class GotchiEntity extends TweenableEntity {
    gotchiId: string = "";
    textureSet!: TextureSet;
    currentActionType: string = "";
    jumpY = 0;

    actionSprite!: Phaser.GameObjects.Sprite;
    shadowSprite!: Phaser.GameObjects.Sprite;
    espBar!: ESPBar;

    entityState!: string;

    jumpTween!: Phaser.Tweens.Tween;

    isDeathTriggered = false;


    static svgMap: Map<string, SVGState> = new Map();
    static activeGotchis: Map<string, GotchiEntity> = new Map(); // Tracks active GotchiEntity instances by gotchiId

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "gotchi", texture: "loading_gotchi", data, navigationGrid});
        
        this.sprite.play("loading_gotchi_anim");
        this.shadowSprite = scene.add.sprite(this.sprite.x + 32, this.sprite.y + 64, "shadow");
        this.shadowSprite.setDepth(this.sprite.depth - 1);
        this.shadowSprite.setOrigin(0.5, 0.5);
        this.shadowSprite.setAlpha(0.3);
        this.shadowSprite.setScale(0.8);

        this.gotchiId = data.gotchiId;

        this.isDeathTriggered = false;

        this.actionSprite = scene.add.sprite(this.sprite.x, this.sprite.y, 
            "actionicons", 12);
        this.actionSprite.setDepth(this.sprite.depth + 1);
        this.actionSprite.setOrigin(0, 0.5);
        this.actionSprite.setAlpha(1);
        this.actionSprite.setScale(10/48);

        // Add to activeGotchis
        GotchiEntity.activeGotchis.set(this.gotchiId, this);

        if (this.gotchiId !== "69420") {
            GotchiEntity.svgMap.set(this.gotchiId, {
                gotchiId: this.gotchiId,
                svgState: "ToBeFetched",
            });
        } else {
            this.textureSet = {
                front: "default_gotchi_svg",
                left: "default_gotchi_left",
                right: "default_gotchi_right",
                back: "default_gotchi_back",
            };
        }

        this.jumpTween = this.createJumpTween(150);

        // this.emoticonEmitter = new EmoticonEmitter(scene, tileX * ZONE_TILES, tileY * ZONE_TILES);

        // this.hpBar = new HPBar({
        //     scene,
        //     x: tileX*TILE_PIXELS,
        //     y: tileY*TILE_PIXELS,
        //     currentHP: data.stats.pulse,
        //     maxHP: data.stats.maxpulse,
        //     trackingSprite: this.sprite,
        // });

        this.espBar = new ESPBar({
            scene,
            x: tileX*TILE_PIXELS,
            y: tileY*TILE_PIXELS,
            ecto: data.stats.ecto,
            spark: data.stats.spark,
            pulse: data.stats.pulse,
            maxESP: 1000,
            trackingSprite: this.sprite,
        });
    }

    createJumpTween(duration_ms: number) {
        const delay_ms = Math.random() * 500;
        return this.scene.tweens.add({
            targets: this,
            jumpY: -8,
            duration: duration_ms,
            yoyo: true,
            repeat: -1,
            ease: "Quad.easeOut",
            delay: delay_ms,
        });
    }

    protected frameUpdate(): void {
        super.frameUpdate();

        this.sprite.setPosition(this.currentPosition.x, this.currentPosition.y - this.jumpY);
        this.shadowSprite.setPosition(this.currentPosition.x + 32, this.currentPosition.y + 64);
        // this.emoticonEmitter.setPosition(this.currentPosition.x + 32, this.currentPosition.y + 16);
        // this.hpBar.setPosition(this.currentPosition.x, this.currentPosition.y);
        this.espBar.setPosition(this.currentPosition.x, this.currentPosition.y);

        this.actionSprite.setPosition(
            this.currentPosition.x,
            this.currentPosition.y
        );
        
        if (this.entityState === "dead") {
            this.currentActionType = "dead";
        }

        this.updateDirection();
    }


    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
        this.updateSVG();
        const currentAction = snapshot.data.actionPlan.currentAction;
        if (currentAction) {
            this.currentActionType = currentAction.type;
            const {texture, frame} = Emoticons.getTextureAndFrame(this.currentActionType);
            this.actionSprite.setTexture(texture);
            this.actionSprite.setFrame(frame);
        }

        // this.hpBar.updateHP(snapshot.data.stats.pulse);
        const {ecto, spark, pulse} = snapshot.data.stats;
        this.espBar.updateESP(ecto, spark, pulse);

        this.entityState = snapshot.data.state;

        // check for dead state
        if (snapshot.data.state === "dead" && !this.isDeathTriggered) {
            this.isDeathTriggered = true;

            this.renderDeathEffects();
        }
    }

    updateSVG() {
        if (!this.textureSet) {
            const svgMapItem = GotchiEntity.svgMap.get(this.gotchiId);
            if (!svgMapItem) return;
            if (svgMapItem.svgState === "ImageLoaded") {
                this.textureSet = {
                    front: `gotchi-${this.gotchiId}-front`,
                    left: `gotchi-${this.gotchiId}-left`,
                    right: `gotchi-${this.gotchiId}-right`,
                    back: `gotchi-${this.gotchiId}-back`,
                };
                this.sprite.stop();
                this.sprite.setTexture(this.textureSet.front);
                this.outlineEffect.rebuild();
            }
        }
    }

    protected updateDirection() {
        if (!this.textureSet) return;
        this.sprite.stop();

        var newTexture = "";

        switch (this.direction) {
            case "left":
                newTexture = this.textureSet.left;
                break;
            case "right":
                newTexture = this.textureSet.right;
                break;
            case "up":
                newTexture = this.textureSet.back;
                break;
            case "down":
            case "none":
            default:
                newTexture = this.textureSet.front;
                break;
        }

        this.sprite.setTexture(newTexture);
        this.outlineEffect.changeTexure(newTexture);
    }

    renderDeathEffects() {
        this.jumpTween.stop();        
        // this.hpBar.setVisible(false);
        this.espBar.setVisible(false);
        this.sprite.setAlpha(0.5);
        // this.emoticonEmitter.setAlpha(0.5);
    }

    destroy(): void {
        // this.hpBar.destroy();
        this.espBar.destroy();
        this.shadowSprite.destroy();
        this.actionSprite.destroy();
        // this.actionIcon.destroy();
        // this.emoticonEmitter.destroy();
        GotchiEntity.activeGotchis.delete(this.gotchiId); // Remove from activeGotchis
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
                const gotchiId = svgSet.id;
                const state = GotchiEntity.svgMap.get(gotchiId);
                if (state) {
                    state.svgState = "LoadingImage";
                    const { front, left, right, back } = svgSet;
                    this.loadGotchiSVG(scene, gotchiId, { front, left, right, back });
                }
            });
        } catch (error) {
            console.error("Failed to fetch bulk SVGs:", error);
        }
    }

    static async loadGotchiSVG(scene: Phaser.Scene, gotchiId: string, svgSet: { front: string; left: string; right: string; back: string }) {
        try {
            const views: ("front" | "left" | "right" | "back")[] = ["front", "left", "right", "back"];
            const loadedViews = new Set<string>();

            views.forEach((view) => {
                const svgDataUrl = `data:image/svg+xml;base64,${btoa(svgSet[view] || "")}`;
                const key = `gotchi-${gotchiId}-${view}`;
                scene.load.image(key, svgDataUrl);

                scene.load.once(`filecomplete-image-${key}`, () => {
                    loadedViews.add(view);
                    if (loadedViews.size === views.length) {
                        this.onAllImagesLoaded(gotchiId);
                    }
                });
            });

            scene.load.start();
        } catch (err) {
            console.error("Failed to load Gotchi SVG for UUID", gotchiId, ":", err);
        }
    }

    static onAllImagesLoaded(gotchiId: string) {
        const state = Array.from(GotchiEntity.svgMap.values()).find(
            (svgState) => svgState.gotchiId === gotchiId
        );
        if (!state) return;
        state.svgState = "ImageLoaded";
    }
}