
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class ShopEntity extends BaseEntity {

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, data: any) {
        super(scene, id, zoneId, tileX, tileY, "shop", "shop", data);
    }

    update(snapshot: EntitySnapshot) {
        super.update(snapshot);
    }
}
