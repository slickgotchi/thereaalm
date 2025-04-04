
import { NavigationGrid } from "../navigation/NavigationGrid";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class AltarEntity extends BaseEntity {

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "altar", texture: "static_entities", data});
            
            this.sprite.setFrame(4);
            navigationGrid.setPassable(tileX, tileY, false);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
    }
}
