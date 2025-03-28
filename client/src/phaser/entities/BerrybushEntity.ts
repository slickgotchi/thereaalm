
import { NavigationGrid } from "../navigation/NavigationGrid";
import { ResourceIcon } from "../ResourceIcon";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class BerrybushEntity extends BaseEntity {
    // resourceIcon: ResourceIcon;

    constructor(scene: Phaser.Scene, id: string, zoneId: number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super(scene, id, zoneId, tileX, tileY, "berrybush", "berrybush", data);

        navigationGrid.setPassable(tileX, tileY, false);

        // this.resourceIcon = new ResourceIcon({
        //     scene: this.scene,
        //     x: this.x,
        //     y: this.y,
        //     iconTexture: "berry_icon",
        //     resourceCount: 100,
        //     trackingSprite: this.sprite,
        // });
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);

        // this.data = snapshot.data;

        // this.resourceIcon.updateResourceCount(snapshot.data.berryCount);
    }
}
