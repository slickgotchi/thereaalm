// server/network/api.go
package network

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"thereaalm/entity/entities"
	"thereaalm/world"
)

type ZoneSnapshot struct {
	Gotchis []GotchiSnapshot `json:"gotchis"`
}

type GotchiSnapshot struct {
	UUID string `json:"uuid"`
	GotchiID string `json:"gotchiId"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func StartAPIServer(worldManager *world.WorldManager, port string) {
	http.HandleFunc("/zones/", func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request for debugging
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Handle GET requests
		if r.Method != http.MethodGet {
			log.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract zone ID from URL path (e.g., /zones/0/snapshot)
		parts := strings.Split(r.URL.Path, "/")
		log.Printf("URL parts: %v", parts)
		// Expected path: /zones/0/snapshot
		// parts should be: ["", "zones", "0", "snapshot"]
		if len(parts) != 4 || parts[1] != "zones" || parts[3] != "snapshot" {
			log.Printf("Invalid endpoint: %s", r.URL.Path)
			http.Error(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}

		zoneID, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("Invalid zone ID: %s", parts[2])
			http.Error(w, "Invalid Zone ID", http.StatusBadRequest)
			return
		}

		// Find the zone
		zone := worldManager.Zones[zoneID]
		if zone == nil {
			log.Printf("Zone not found: %d", zoneID)
			http.Error(w, "Zone not found: "+strconv.Itoa(zoneID), http.StatusNotFound)
			return
		}

		// Collect Gotchi positions
		var snapshot ZoneSnapshot
		for _, entity := range zone.Entities {
			gotchi, ok := entity.(*entities.GotchiEntity)
			if !ok {
				continue
			}

			snapshot.Gotchis = append(snapshot.Gotchis, GotchiSnapshot{
				UUID: gotchi.UUID,
				GotchiID: gotchi.Gotchi.SubgraphData.ID,
				X:  gotchi.Position.X,
				Y:  gotchi.Position.Y,
			})
		}

		// Respond with JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snapshot); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	log.Printf("Starting API server on port %s...", port)
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()
}