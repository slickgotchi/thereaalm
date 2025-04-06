import { Depth } from "../Depth";
import { TILE_PIXELS } from "../GameScene";
import { OutlineEffect } from "../OutlineEffect";

export interface EntitySnapshot {
    id: string;
    zoneId: number;
    type: string;
    tileX: number;
    tileY: number;
    data: any;
}

interface Props {
    scene: Phaser.Scene;
    id: string;
    zoneId: number;
    tileX: number;
    tileY: number;
    type: string;
    texture: string;
    data: any;
}

export class BaseEntity {
    id: string;
    zoneId: number;
    scene: Phaser.Scene;
    tileX: number;
    tileY: number;
    x: number;
    y: number;
    type: string;
    data: any;
    private isSelected: boolean = false;

    // must be destroyed on death
    sprite: Phaser.GameObjects.Sprite;
    outlineEffect: OutlineEffect;

    constructor(props: Props) {
        const {scene, id, zoneId, tileX, tileY, type, texture, data} = props;
        
        this.scene = scene;
        this.id = id;
        this.zoneId = Number(zoneId);
        this.tileX = tileX;
        this.tileY = tileY;
        this.x = tileX * TILE_PIXELS;
        this.y = tileY * TILE_PIXELS;
        this.type = type;

        this.sprite = scene.add.sprite(this.x, this.y, texture)
            .setOrigin(0, 0)
            .setDepth(Depth.ENTITIES)
            .setInteractive({ useHandCursor: true })
            .setData("entity", this); // Link sprite to this entity

        this.outlineEffect = new OutlineEffect({
            scene: this.scene,
            target: this.sprite,
        });

        this.data = data;
        this.setupHoverEvents();
        // console.log(`[${this.id}] Entity constructed`);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        const { tileX, tileY, data } = snapshot;
        this.tileX = tileX;
        this.tileY = tileY;
        this.x = tileX * TILE_PIXELS;
        this.y = tileY * TILE_PIXELS;
        this.data = data;
        this.sprite.setPosition(this.x, this.y);
    }

    destroy() {
        this.sprite?.destroy();
        this.outlineEffect?.destroy();
    }

    setSelected(selected: boolean) {
        this.isSelected = selected;
        if (this.isSelected) {
            this.outlineEffect.showOutline();
        } else {
            this.outlineEffect.hideOutline();
        }
    }

    private setupHoverEvents() {
        this.sprite.on("pointerover", () => {
            if (!this.isSelected) this.outlineEffect.showOutline();
        });
        this.sprite.on("pointerout", () => {
            if (!this.isSelected) this.outlineEffect.hideOutline();
        });
    }
}