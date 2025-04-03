
import { NavigationGrid } from "../navigation/NavigationGrid";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class ShopEntity extends BaseEntity {

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "shop", texture: "static_entities", data});

        this.sprite.setFrame(3);
        navigationGrid.setPassable(tileX, tileY, false);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
    }
}
