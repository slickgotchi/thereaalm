
import { TILE_PIXELS } from "../GameScene";
import { HPBar } from "../HPBar";
import { NavigationGrid } from "../navigation/NavigationGrid";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";


export class ShopEntity extends BaseEntity {

    hpBar!: HPBar;

    constructor(scene: Phaser.Scene, id: string, zoneId:number, tileX: number, tileY: number, 
        data: any, navigationGrid: NavigationGrid) {
        super({scene, id, zoneId, tileX, tileY, 
            type: "shop", texture: "static_entities", data});

        this.sprite.setFrame(3);
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

        if (snapshot.data.state === "active") {
            this.sprite.setFrame(3);
            this.hpBar.setVisible(true);
        } else {
            this.sprite.setFrame(11);
            this.hpBar.setVisible(false);
        }
    }
}
