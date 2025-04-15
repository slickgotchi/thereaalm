package world

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// TiledMap represents the structure of a Tiled tilemap JSON file.
type TiledMap struct {
    Height      int     `json:"height"`
    Width       int     `json:"width"`
    Layers      []Layer `json:"layers"`
    Tilesets    []interface{} `json:"tilesets"`
    Type        string  `json:"type"`
    Version     string  `json:"version"`
    TileWidth   int     `json:"tilewidth"`
    TileHeight  int     `json:"tileheight"`
}

// Layer represents a layer in the tilemap.
type Layer struct {
    Data       []int         `json:"data"`
    Height     int           `json:"height"`
    Width      int           `json:"width"`
    Name       string        `json:"name"`
    Type       string        `json:"type"`
    Visible    bool          `json:"visible"`
    X          int           `json:"x"`
    Y          int           `json:"y"`
    Properties []Property    `json:"properties,omitempty"` // Optional field
}

// Property represents a custom property in a layer.
type Property struct {
    Name  string      `json:"name"`
    Type  string      `json:"type"`
    Value interface{} `json:"value"`
}

// LoadTilemap loads a Tiled tilemap from a JSON file and creates Impassable entities.
func LoadTilemap(filePath string, worldManager *WorldManager, zoneID int) error {
    // Read the JSON file
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return err
    }

    // Parse the JSON into a TiledMap struct
    var tiledMap TiledMap
    if err := json.Unmarshal(data, &tiledMap); err != nil {
        return err
    }

    // Validate the map dimensions
    if tiledMap.Width != 512 || tiledMap.Height != 512 {
        log.Printf("Warning: Tilemap dimensions are %dx%d, expected 512x512", tiledMap.Width, tiledMap.Height)
    }

    // Find the zone in the world manager
    zone := worldManager.Zones[zoneID]
    if zone == nil {
        log.Fatalf("Zone %d not found in world manager", zoneID)
    }

    // Process each layer
    for _, layer := range tiledMap.Layers {
        // Skip non-tile layers or layers without data
        if layer.Type != "tilelayer" || len(layer.Data) == 0 {
            continue
        }

        // Check if the layer is tagged as "Impassable"
        isObstacle := false
        for _, prop := range layer.Properties {
            if prop.Name == "isObstacle" && prop.Type == "bool" {
                if val, ok := prop.Value.(bool); ok && val {
                    isObstacle = true
                    break
                }
            }
        }

        if !isObstacle {
            // log.Printf("Layer %s is not marked as isObstacle (or missing properties)... continuing", layer.Name)
            continue
        }

        // The layer is impassable; process each tile
        log.Printf("Processing 'isObstacle' layer: %s", layer.Name)
        count := 0
        for i, tileID := range layer.Data {
            // Skip tiles with ID 0 (no tile present)
            if tileID == 0 {
                continue
            }

            // Convert the 1D index to 2D coordinates
            zoneWorldX, zoneWorldY := zone.GetPosition()
            x := zoneWorldX + i % layer.Width
            y := zoneWorldY + i / layer.Width

            // Create an obstacle
            zone.AddObstacle(x, y)
            // log.Println("Added obstacle at ", x, y)

            // increment count
            count++
        }
        log.Printf("Added %d 'isObstacle' cells in zone %d", count, zoneID)
    }

    log.Printf("Finished loading tilemap for zone %d", zoneID)
    return nil
}