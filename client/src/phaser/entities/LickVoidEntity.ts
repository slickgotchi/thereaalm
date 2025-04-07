
import { TILE_PIXELS } from "../GameScene";
import { HPBar } from "../HPBar";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class LickVoidEntity extends BaseEntity {

    hpBar!: HPBar;

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "lickvoid", texture: "static_entities", data});
            
            this.sprite.setFrame(8);
            navigationGrid.setPassable(tileX, tileY, false);

            this.hpBar = new HPBar({
                scene,
                x: tileX*TILE_PIXELS,
                y: tileY*TILE_PIXELS,
                currentHP: data.stats.pulse,
                maxHP: data.stats.maxpulse,
                trackingSprite: this.sprite,
            });
    }

    snapshotUpdate(snapshot: EntitySnapshot) {
        super.snapshotUpdate(snapshot);

        this.hpBar.updateHP(snapshot.data.stats.pulse);

        if (snapshot.data.state === "active") {
            this.sprite.setFrame(8);
        } else {
            this.sprite.setFrame(9);
        }
    }
}
