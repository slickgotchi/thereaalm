import { NavigationGrid } from "../navigation/NavigationGrid";
import { AltarEntity } from "./AltarEntity";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { BerrybushEntity } from "./BerrybushEntity";
import { GotchiEntity } from "./GotchiEntity";
import { LickquidatorEntity } from "./LickquidatorEntity";
import { ShopEntity } from "./ShopEntity";


export class EntityFactory {
    static create(scene: Phaser.Scene, snapshot: EntitySnapshot, navigationGrid: NavigationGrid): BaseEntity {
        const { id, zoneId, tileX, tileY, type, data } = snapshot;

        switch (type) {
            case "gotchi":
                return new GotchiEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "lickquidator":
                return new LickquidatorEntity(scene, id, zoneId, tileX, tileY, data);
            case "berrybush":
                return new BerrybushEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "shop":
                return new ShopEntity(scene, id, zoneId, tileX, tileY, data);
            case "altar":
                return new AltarEntity(scene, id, zoneId, tileX, tileY, data);
            default:
                return new BaseEntity(scene, id, zoneId, tileX, tileY, type, type, data);
        }
    }
}
