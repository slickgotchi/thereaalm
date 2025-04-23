import { EmoticonEmitter } from "../emoticons/EmoticonEmitter";
import { Emoticons } from "../emoticons/Emoticons";
import { TILE_PIXELS, ZONE_TILES } from "../GameScene";
import { HPBar } from "../HPBar";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { OutlineEffect } from "../OutlineEffect";
import { VFXManager } from "../VFXManager";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { TweenableEntity } from "./TweenableEntity";



export class LickquidatorEntity extends TweenableEntity {
    lastEmoticonEmitTime_ms: number = 0;
    emoticonEmitInterval_ms: number = 3000;
    jumpY = 0;

    // must be destroyed on death
    actionSprite!: Phaser.GameObjects.Sprite;
    shadowSprite!: Phaser.GameObjects.Sprite;
    hpBar!: HPBar;

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "lickquidator", texture: "lickquidator", data, navigationGrid});
        
        // set lick direction frame (down initially)
        this.sprite.setFrame(1);

        // shadow
        this.shadowSprite = scene.add.sprite(this.sprite.x + 32, this.sprite.y + 64, "shadow");
        this.shadowSprite.setDepth(this.sprite.depth - 1);
        this.shadowSprite.setOrigin(0.5, 0.5);
        this.shadowSprite.setAlpha(0.3);
        this.shadowSprite.setScale(0.8);

        const {texture,frame} = Emoticons.getTextureAndFrame("attack");
        this.actionSprite = scene.add.sprite(this.sprite.x, this.sprite.y, 
            texture, frame);
        this.actionSprite.setDepth(this.sprite.depth + 1);
        this.actionSprite.setOrigin(0, 0.5);
        this.actionSprite.setAlpha(1);
        this.actionSprite.setScale(10/48);

        // float anim
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

        // hp bar 
        this.hpBar = new HPBar({
            scene, 
            x: tileX*TILE_PIXELS, 
            y: tileY*TILE_PIXELS, 
            currentHP: data.stats.pulse,
            maxHP: data.stats.maxpulse,
            trackingSprite: this.sprite
        });
    }

    protected frameUpdate(): void {
        super.frameUpdate();

        // set positions
        this.sprite.setPosition(this.currentPosition.x, this.currentPosition.y - this.jumpY);
        this.shadowSprite.setPosition(this.currentPosition.x + 32, this.currentPosition.y + 64);
        this.hpBar.setPosition(this.currentPosition.x, this.currentPosition.y);
        this.actionSprite.setPosition(this.currentPosition.x, this.currentPosition.y);

        // reset last emit time during tween so the emoticon emits as soon as we
        // finish a movement tween
        if (this.tweenWorker.getIsTweening()) {
            this.lastEmoticonEmitTime_ms = 0;
        }

        // // emoticon emission
        // const currTime_ms = Date.now();
        // if (
        //     currTime_ms - this.lastEmoticonEmitTime_ms > this.emoticonEmitInterval_ms &&
        //     !this.tweenWorker.getIsTweening()
        // ) {
        //     this.lastEmoticonEmitTime_ms = currTime_ms;
        //     this.emoticonEmitter.emit("attack", 240);
        // }

        // update facing direction
        this.updateDirection();

    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
        
        // update hp
        this.hpBar.updateHP(snapshot.data.stats.pulse);


    }

    destroy(): void {
        // Play explosion animation
        VFXManager.getInstance(this.scene)
            .playLickExplosion(this.currentPosition.x, this.currentPosition.y);

        // this.emoticonEmitter.destroy();
        this.actionSprite.destroy();
        this.shadowSprite.destroy(); 
        this.hpBar.destroy();

        super.destroy();
    }

    // this is auto called by TweenableEntity callback function
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
