
import { NavigationGrid } from "../navigation/NavigationGrid";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class ShopEntity extends BaseEntity {

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "shop", texture: "shop", data});

        navigationGrid.setPassable(tileX, tileY, false);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
    }
}
