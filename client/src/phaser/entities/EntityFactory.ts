import { NavigationGrid } from "../navigation/NavigationGrid";
import { AltarEntity } from "./AltarEntity";
import { BaseEntity, EntitySnapshot } from "./BaseEntity";
import { GotchiEntity } from "./GotchiEntity";
import { LickquidatorEntity } from "./LickquidatorEntity";
import { ShopEntity } from "./ShopEntity";
import { KekWoodTreeEntity } from "./resources/KekWoodTreeEntity";
import { AlphaSlateBouldersEntity } from "./resources/AlphaSlateBouldersEntity";
import { FomoBerryBushEntity } from "./resources/FomoBerryBushEntity";
import { LickVoidEntity } from "./LickVoidEntity";


export class EntityFactory {
    static create(scene: Phaser.Scene, snapshot: EntitySnapshot, navigationGrid: NavigationGrid): BaseEntity {
        const { id, zoneId, tileX, tileY, type, data } = snapshot;

        switch (type) {
            case "gotchi":
                return new GotchiEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "lickquidator":
                return new LickquidatorEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "fomoberrybush":
                return new FomoBerryBushEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "kekwoodtree":
                return new KekWoodTreeEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "alphaslateboulders":
                return new AlphaSlateBouldersEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "shop":
                return new ShopEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "altar":
                return new AltarEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            case "lickvoid":
                return new LickVoidEntity(scene, id, zoneId, tileX, tileY, data, navigationGrid);
            default:
                return new BaseEntity({scene, id, zoneId, tileX, tileY, type, texture: "", data});
        }
    }
}
