
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class AltarEntity extends BaseEntity {

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, data: any) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "altar", texture: "static_entities", data});
            
            this.sprite.setFrame(4);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
    }
}
