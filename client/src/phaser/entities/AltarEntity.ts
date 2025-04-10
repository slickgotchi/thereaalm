
import { TILE_PIXELS, ZONE_TILES } from "../GameScene";
import { HPBar } from "../HPBar";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class AltarEntity extends BaseEntity {

    hpBar!: HPBar;

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "altar", texture: "static_entities", data});
            
            this.sprite.setFrame(4);
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
            this.sprite.setFrame(4);
            this.hpBar.setVisible(true);
        } else {
            this.sprite.setFrame(5);
            this.hpBar.setVisible(false);
        }
    }
}
