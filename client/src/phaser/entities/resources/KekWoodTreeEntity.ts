
import { NavigationGrid } from "../../navigation/NavigationGrid";
import { ResourceIcon } from "../../ResourceIcon";
import { BaseEntity, EntitySnapshot } from "../BaseEntity";


export class KekWoodTreeEntity extends BaseEntity {
    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "kekwoodtree", texture: "static_entities", data});

        this.sprite.setFrame(2);
        navigationGrid.setPassable(tileX, tileY, false);
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);
    }
}
