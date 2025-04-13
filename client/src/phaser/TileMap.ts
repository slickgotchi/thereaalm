import Phaser from "phaser";
import { TILE_PIXELS, ZONE_TILES } from "./GameScene";
import { Depth } from "./Depth";

export class TileMap {
    private map!: Phaser.Tilemaps.Tilemap;

    constructor(
        private scene: Phaser.Scene,
        private mapKey: string,
        private tilesetKey: string,
        private tilesetImage: string
    ) {}

    // Preload assets
    preload() {
        this.scene.load.tilemapTiledJSON(
            "yield_fields_2",
            `assets/tilemaps/maps/yield_fields_2.json`
        );
        this.scene.load.image(
            this.tilesetKey,
            `assets/tilemaps/tilesets/${this.tilesetImage}`
        );
    }

    // Create map and layers
    create() {
        this.map = this.scene.make.tilemap({ key: this.mapKey });

        // Load the tileset
        const tileset = this.map.addTilesetImage(this.tilesetKey);
        if (!tileset) {
            console.warn(
                `Tileset '${this.tilesetKey}' not found in map '${this.mapKey}'`
            );
            return;
        }

        // Iterate through all layers in the tilemap
        this.map.layers.forEach((layerData) => {
            const layer = this.map.createLayer(layerData.name, tileset);
            layer.setPosition(
                4 * ZONE_TILES * TILE_PIXELS,
                5 * ZONE_TILES * TILE_PIXELS
            );
            layer.setScale(2);
            layer.setDepth(Depth.TILES);
            if (layer) {
                console.log(`Created layer: ${layerData.name}`);
            }
        });
    }

    // Get the tilemap object
    getMap() {
        return this.map;
    }
}
