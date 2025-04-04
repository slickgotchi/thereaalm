import { EmoticonEmitter } from "../emoticons/EmoticonEmitter";
import { ZONE_TILES } from "../GameScene";
import { HPBar } from "../HPBar";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { OutlineEffect } from "../OutlineEffect";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { TweenableEntity } from "./TweenableEntity";



export class LickquidatorEntity extends TweenableEntity {
    emoticonEmitter!: EmoticonEmitter;
    lastEmoticonEmitTime_ms: number = 0;
    emoticonEmitInterval_ms: number = 3000;
    jumpY = 0;
    shadowSprite!: Phaser.GameObjects.Sprite;

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "lickquidator", texture: "lickquidator", data, navigationGrid});
        
        // Add animation
        this.sprite.setFrame(1);

        this.shadowSprite = scene.add.sprite(this.sprite.x + 32, this.sprite.y + 64, "shadow");
        this.shadowSprite.setDepth(this.sprite.depth - 1);
        this.shadowSprite.setOrigin(0.5, 0.5);
        this.shadowSprite.setAlpha(0.3);
        this.shadowSprite.setScale(0.8);

        const delay_ms = Math.random() * 500;
        scene.tweens.add({
            targets: this,
            jumpY: -8,
            duration: 150,
            yoyo: true,
            repeat: -1,
            ease: "Quad.easeOut",
            delay: delay_ms,
        });

        this.emoticonEmitter = new EmoticonEmitter(scene, tileX * ZONE_TILES, tileY * ZONE_TILES);
    }

    protected frameUpdate(): void {
        super.frameUpdate();

        this.sprite.setPosition(this.currentPosition.x, this.currentPosition.y - this.jumpY);
        this.shadowSprite.setPosition(this.currentPosition.x + 32, this.currentPosition.y + 64);
        this.emoticonEmitter.setPosition(this.currentPosition.x + 32, this.currentPosition.y + 16);

        if (this.tweenWorker.getIsTweening()) {
            this.lastEmoticonEmitTime_ms = 0;
        }

        const currTime_ms = Date.now();
        if (
            currTime_ms - this.lastEmoticonEmitTime_ms > this.emoticonEmitInterval_ms &&
            !this.tweenWorker.getIsTweening()
        ) {
            this.lastEmoticonEmitTime_ms = currTime_ms;
            this.emoticonEmitter.emit("attack", 240);
        }
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
        
        // update direction frame
    }

    protected updateDirection(): void {
        switch (this.direction) {
            case "left":
                this.sprite.setFrame(3);
                break;
            case "right":
                this.sprite.setFrame(4);
                break;
            case "up":
                this.sprite.setFrame(0);
                break;
            case "down":
            case "none":
            default:
                this.sprite.setFrame(1);
                break;
        }
    }
}
