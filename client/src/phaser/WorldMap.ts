// src/phaser/WorldMap.ts
import Phaser from "phaser";
import { TILE_PIXELS, ZONE_TILES } from "./GameScene";

// Define the structure of the zone map response from the API
interface ZoneMapResponse {
    zoneMap: string[][];
}

export class WorldMap {
    private scene: Phaser.Scene;
    private zoneColors: Map<string, number> = new Map(); // Maps zone names to colors
    private zonePixelSize: number; // Will be 128 * 32 = 4096 pixels

    constructor(scene: Phaser.Scene) {
        this.scene = scene;
        this.zonePixelSize = ZONE_TILES * TILE_PIXELS; // 4096 pixels

        // Assign colors to each unique zone name
        this.assignZoneColors();
    }

    // Assign a unique color to each zone name (excluding empty zones)
    private assignZoneColors() {
        // List of unique zone names (you can extract this dynamically from the zone map later)
        const zoneNames = [
            "rofl_reef",
            "north_beach",
            "aalpha_river_valley",
            "laughing_peaks",
            "broken_line",
            "defi_desert",
            "genesis_rocks",
            "shelnot_pass",
            "open_steppe",
            "caaverns",
            "daark_forest",
            "phaantastic_grounds",
            "alpha_river_valley",
            "maagma_springs",
            "tree_of_fud",
            "yield_fields",
            "poly_lakes",
            "mount_oomf",
            "the_arena",
            "aalpha_lake",
            "lickquidator_ruins",
            "the_infinity_cliffs",
            "south_beach",
            "the_citaadel",
        ];

        // Assign a color to each zone (using hex color codes)
        const colors = [
            0xff0000, // red
            0x00ff00, // green
            0x0000ff, // blue
            0xffff00, // yellow
            0xff00ff, // magenta
            0x00ffff, // cyan
            0xffa500, // orange
            0x800080, // purple
            0x008000, // dark green
            0x000080, // navy
            0x808000, // olive
            0x800000, // maroon
            0x008080, // teal
            0xff4500, // orange red
            0x9400d3, // dark violet
            0x00ced1, // dark turquoise
            0x228b22, // forest green
            0x20b2aa, // light sea green
            0x87ceeb, // sky blue
            0xff69b4, // hot pink
            0x4682b4, // steel blue
            0x9acd32, // yellow green
            0xb8860b, // dark goldenrod
            0x2e8b57, // sea green
            0x9932cc, // dark orchid
        ];

        // Map each zone name to a color
        zoneNames.forEach((zone, index) => {
            if (index < colors.length) {
                this.zoneColors.set(zone, colors[index]);
            } else {
                // Fallback color if we run out of colors
                this.zoneColors.set(zone, 0xcccccc); // gray
            }
        });
    }

    // Fetch the zone map from the server
    private async fetchZoneMap(): Promise<string[][]> {
        try {
            const response = await fetch("http://localhost:8080/zonemap", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                throw new Error(
                    `Failed to fetch zone map: ${response.status} ${response.statusText}`
                );
            }

            const data: ZoneMapResponse = await response.json();
            return data.zoneMap;
        } catch (error) {
            console.error("Error fetching zone map:", error);
            // Return an empty map as a fallback
            return [];
        }
    }

    // Draw the world map
    public async draw(): Promise<number> {
        // Fetch the zone map
        const zoneMap = await this.fetchZoneMap();
        if (zoneMap.length === 0) {
            console.error("No zone map data to draw");
            return 0;
        }

        // Iterate through the zone map (10 rows x 10 columns)
        for (let row = 0; row < zoneMap.length; row++) {
            for (let col = 0; col < zoneMap[row].length; col++) {
                const zoneName = zoneMap[row][col];
                if (!zoneName) continue; // Skip empty zones

                // Calculate the top-left position of the zone rectangle
                const x = col * this.zonePixelSize; // 4096 pixels per zone
                const y = row * this.zonePixelSize;

                // Get the color for this zone
                const color = this.zoneColors.get(zoneName) || 0xcccccc; // Default to gray if no color assigned

                // Draw the rectangle
                const graphics = this.scene.add.graphics();
                graphics.fillStyle(color, 0.5); // 0.5 transparency
                graphics.fillRect(x, y, this.zonePixelSize, this.zonePixelSize);
                graphics.setDepth(0); // Render at the bottom

                // Add the zone name as text in the top-left corner of the rectangle
                this.scene.add
                    .text(x + 10, y + 10, zoneName, {
                        fontSize: "128px",
                        color: "#ffffff",
                        fontStyle: "bold",
                    })
                    .setDepth(1); // Text above the rectangle
            }
        }

        return zoneMap.length;
    }
}
