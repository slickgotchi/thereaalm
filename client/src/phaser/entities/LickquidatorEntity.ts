import { HPBar } from "../HPBar";
import { OutlineEffect } from "../OutlineEffect";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";



export class LickquidatorEntity extends BaseEntity {
    // hpBar: HPBar;

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, data: any) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "lickquidator", texture: "lickquidator", data});
        
        // Add animation
        this.sprite.setFrame(1);

        // // Add HP Bar
        // this.hpBar = new HPBar({scene: scene, 
        //     x: this.x,
        //     y: this.y,
        //     currentHP: data.currentHp,
        //     maxHP: data.maxHp, 
        //     trackingSprite: this.sprite
        // });
        // this.hpBar.updateHP(data.currentHp);

        // this.outlineEffect = new OutlineEffect({
        //     scene: this.scene,
        //     target: this.sprite,
        // });
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
        // this.hpBar.updateHP(data.currentHp); // Example: Update HP bar dynamically
    }
}
