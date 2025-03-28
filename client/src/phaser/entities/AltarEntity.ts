
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class AltarEntity extends BaseEntity {

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, data: any) {
        super(scene, id, zoneId, tileX, tileY, "altar", "altar", data);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
    }
}
