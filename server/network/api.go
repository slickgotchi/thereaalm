package network

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"thereaalm/config"
	"thereaalm/interfaces"
	"thereaalm/stattypes"
	"thereaalm/world"

	"github.com/google/uuid"
)

// ZoneSnapshot represents the state of a zone, including Gotchi positions.
type ZoneSnapshot struct {
	EntitySnapshots []EntitySnapshot `json:"entitySnapshots"`
}

// GotchiSnapshot captures the position and ID of a Gotchi in a zone.
type EntitySnapshot struct {
	ID     uuid.UUID `json:"id"`
	ZoneID int       `json:"zoneId"`
	Type   string    `json:"type"`
	X      int       `json:"tileX"`
	Y      int       `json:"tileY"`
	Data   interface{} `json:"data"`
}

// ZoneMapResponse represents the structure of the zone map to be sent to the client.
type ZoneMapResponse struct {
	ZoneMap [][]string `json:"zoneMap"`
}

// StakeRequest represents the request body for staking/unstaking GHST.
type StakeRequest struct {
	UUID uuid.UUID `json:"uuid"`
	GHSTAmount int `json:"ghstAmount"`
	ZoneID int `json:"zoneId"`
}

// EatTreatRequest represents the request body for eating a treat.
type EatTreatRequest struct {
	TreatName string `json:"treatName"`
	ZoneID int `json:"zoneId"`
	UUID uuid.UUID `json:"uuid"`
}

// GotchiResponse represents the response for Gotchi-related actions.
type GotchiResponse struct {
	TreatTotal float64 `json:"treatTotal"`
	StakedGhst float64 `json:"stakedGhst"`
	Ecto int `json:"ecto"`
	Spark int `json:"spark"`
	Pulse int `json:"pulse"`
}

// StartAPIServer initializes the API server with the given world manager and port.
func StartAPIServer(worldManager *world.WorldManager, port string) {
	// Create a new ServeMux to handle routes explicitly
	mux := http.NewServeMux()

	// Register handlers with CORS middleware
	mux.HandleFunc("/zones/", withCORS(handleZoneSnapshot(worldManager)))
	mux.HandleFunc("/zonemap", withCORS(handleZoneMap()))
	mux.HandleFunc("/gotchi/stake", withCORS(handleStakeGotchi(worldManager)))
	mux.HandleFunc("/gotchi/unstake", withCORS(handleUnstakeGotchi(worldManager)))
	mux.HandleFunc("/gotchi/eat", withCORS(handleEatTreat(worldManager)))

	// Start the server
	log.Printf("Starting API server on port %s...", port)
	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()
}

// withCORS is a middleware that adds CORS headers to all responses and handles OPTIONS requests.
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the actual handler
		handler(w, r)
	}
}

// writeJSON encodes the given data as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError sends an error response with the given status code and message.
func writeError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s (status: %d)", message, statusCode)
	http.Error(w, message, statusCode)
}

// handleZoneSnapshot returns a handler for the /zones/{id}/snapshot endpoint.
func handleZoneSnapshot(worldManager *world.WorldManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Only allow GET requests
		if r.Method != http.MethodGet {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract zone ID from URL path (e.g., /zones/0/snapshot)
		parts := strings.Split(r.URL.Path, "/")
		// log.Printf("URL parts: %v", parts)

		// Expected path: /zones/0/snapshot
		if len(parts) != 4 || parts[1] != "zones" || parts[3] != "snapshot" {
			writeError(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}

		zoneID, err := strconv.Atoi(parts[2])
		if err != nil {
			writeError(w, "Invalid Zone ID", http.StatusBadRequest)
			return
		}

		// Find the zone
		zone := worldManager.Zones[zoneID]
		if zone == nil {
			writeError(w, "Zone not found: "+strconv.Itoa(zoneID), http.StatusNotFound)
			return
		}

		// Collect entities
		var snapshot ZoneSnapshot
		for _, entity := range zone.GetEntities() {
			x, y := entity.GetPosition()
			snapshot.EntitySnapshots = append(snapshot.EntitySnapshots, EntitySnapshot{
				ID:     entity.GetUUID(),
				ZoneID: zoneID,
				Type: 	  entity.GetType(),
				X:        x,
				Y:        y,
				Data: entity.GetSnapshotData(),
			})
		}

		// Respond with JSON
		writeJSON(w, snapshot)
	}
}

// handleZoneMap returns a handler for the /zonemap endpoint.
func handleZoneMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Only allow GET requests
		if r.Method != http.MethodGet {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Ensure the endpoint is exactly /zonemap
		if r.URL.Path != "/zonemap" {
			writeError(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}

		// Create the response with the zone map
		response := ZoneMapResponse{
			ZoneMap: config.ZoneMap,
		}

		// Respond with JSON
		writeJSON(w, response)
	}
}

// handleStakeGotchi handles staking GHST for a Gotchi.
func handleStakeGotchi(worldManager *world.WorldManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req StakeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.GHSTAmount <= 0 {
			writeError(w, "Amount must be a positive number", http.StatusBadRequest)
			return
		}

		// For simplicity, assume the Gotchi is in zone 0 and we find the first Gotchi
		zone := worldManager.Zones[req.ZoneID]
		if zone == nil {
			writeError(w, "Zone not found", http.StatusNotFound)
			return
		}

		gotchiEntity := zone.GetEntityByUUID(req.UUID)

		if gotchiEntity == nil {
			writeError(w, "Gotchi not found", http.StatusNotFound)
			return
		}

		gotchiStats, _ := gotchiEntity.(interfaces.IStats)
		gotchiStats.DeltaStat(stattypes.StakedGHST, float64(req.GHSTAmount))

		// Respond with updated Gotchi data
		response := GotchiResponse{
			TreatTotal: gotchiStats.GetStat(stattypes.TreatTotal),
			StakedGhst: gotchiStats.GetStat(stattypes.StakedGHST),
			Ecto: int(gotchiStats.GetStat(stattypes.Ecto)),
			Spark: int(gotchiStats.GetStat(stattypes.Spark)),
			Pulse: int(gotchiStats.GetStat(stattypes.Pulse)),
		}
		writeJSON(w, response)
	}
}

// handleUnstakeGotchi handles unstaking GHST for a Gotchi.
func handleUnstakeGotchi(worldManager *world.WorldManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req StakeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.GHSTAmount <= 0 {
			writeError(w, "Amount must be a positive number", http.StatusBadRequest)
			return
		}

		// For simplicity, assume the Gotchi is in zone 0 and we find the first Gotchi
		zone := worldManager.Zones[req.ZoneID]
		if zone == nil {
			writeError(w, "Zone not found", http.StatusNotFound)
			return
		}

		gotchiEntity := zone.GetEntityByUUID(req.UUID)

		if gotchiEntity == nil {
			writeError(w, "Gotchi not found", http.StatusNotFound)
			return
		}

		gotchiStats, _ := gotchiEntity.(interfaces.IStats)
		gotchiStats.DeltaStat(stattypes.StakedGHST, float64(-req.GHSTAmount))

		if gotchiStats.GetStat(stattypes.StakedGHST) < 0 {
			gotchiStats.SetStat(stattypes.StakedGHST, 0)
		}

		// Respond with updated Gotchi data
		response := GotchiResponse{
			TreatTotal: gotchiStats.GetStat(stattypes.TreatTotal),
			StakedGhst: gotchiStats.GetStat(stattypes.StakedGHST),
			Ecto: int(gotchiStats.GetStat(stattypes.Ecto)),
			Spark: int(gotchiStats.GetStat(stattypes.Spark)),
			Pulse: int(gotchiStats.GetStat(stattypes.Pulse)),
		}
		writeJSON(w, response)
	}
}

// handleEatTreat handles a Gotchi eating a treat.
func handleEatTreat(worldManager *world.WorldManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req EatTreatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate treat name
		validTreats := map[string]int{
			"Sushi Roll": 500,
			"Coconut":    500,
			"Candy":      500,
		}
		_, ok := validTreats[req.TreatName]
		if !ok {
			writeError(w, "Invalid treat name", http.StatusBadRequest)
			return
		}

		// For simplicity, assume the Gotchi is in zone 0 and we find the first Gotchi
		zone := worldManager.Zones[req.ZoneID]
		if zone == nil {
			writeError(w, "Zone not found", http.StatusNotFound)
			return
		}

		gotchiEntity := zone.GetEntityByUUID(req.UUID)

		if gotchiEntity == nil {
			writeError(w, "Gotchi not found", http.StatusNotFound)
			return
		}

		gotchiStats, _ := gotchiEntity.(interfaces.IStats)

		if req.TreatName == "Sushi Roll" {
			gotchiStats.DeltaStat(stattypes.Ecto, 100)
		} else if req.TreatName == "Coconut" {
			gotchiStats.DeltaStat(stattypes.Spark, 100)
		} else if req.TreatName == "Candy" {
			gotchiStats.DeltaStat(stattypes.Pulse, 100)
		}
		
		// remove some treat
		gotchiStats.DeltaStat(stattypes.TreatTotal, -float64(validTreats[req.TreatName]))

		// Respond with updated Gotchi data
		response := GotchiResponse{
			TreatTotal: gotchiStats.GetStat(stattypes.TreatTotal),
			StakedGhst: gotchiStats.GetStat(stattypes.StakedGHST),
			Ecto: int(gotchiStats.GetStat(stattypes.Ecto)),
			Spark: int(gotchiStats.GetStat(stattypes.Spark)),
			Pulse: int(gotchiStats.GetStat(stattypes.Pulse)),
		}
		writeJSON(w, response)
	}
}

/*
package network

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"thereaalm/config"
	"thereaalm/world"

	"github.com/google/uuid"
)

// ZoneSnapshot represents the state of a zone, including Gotchi positions.
type ZoneSnapshot struct {
	EntitySnapshots []EntitySnapshot `json:"entitySnapshots"`
}

// GotchiSnapshot captures the position and ID of a Gotchi in a zone.
type EntitySnapshot struct {
	ID     uuid.UUID `json:"id"`
	ZoneID int `json:"zoneId"`
	Type 	 string `json:"type"`
	X        int    `json:"tileX"`
	Y        int    `json:"tileY"`
	Data interface{}	`json:"data"`
}

// ZoneMapResponse represents the structure of the zone map to be sent to the client.
type ZoneMapResponse struct {
	ZoneMap [][]string `json:"zoneMap"`
}

// StartAPIServer initializes the API server with the given world manager and port.
func StartAPIServer(worldManager *world.WorldManager, port string) {
	// Create a new ServeMux to handle routes explicitly
	mux := http.NewServeMux()

	// Register handlers with CORS middleware
	mux.HandleFunc("/zones/", withCORS(handleZoneSnapshot(worldManager)))
	mux.HandleFunc("/zonemap", withCORS(handleZoneMap()))

	// Start the server
	log.Printf("Starting API server on port %s...", port)
	go func() {
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("API server failed: %v", err)
		}
	}()
}

// withCORS is a middleware that adds CORS headers to all responses and handles OPTIONS requests.
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle CORS preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the actual handler
		handler(w, r)
	}
}

// writeJSON encodes the given data as JSON and writes it to the response.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError sends an error response with the given status code and message.
func writeError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s (status: %d)", message, statusCode)
	http.Error(w, message, statusCode)
}

// handleZoneSnapshot returns a handler for the /zones/{id}/snapshot endpoint.
func handleZoneSnapshot(worldManager *world.WorldManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Only allow GET requests
		if r.Method != http.MethodGet {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract zone ID from URL path (e.g., /zones/0/snapshot)
		parts := strings.Split(r.URL.Path, "/")
		// log.Printf("URL parts: %v", parts)

		// Expected path: /zones/0/snapshot
		if len(parts) != 4 || parts[1] != "zones" || parts[3] != "snapshot" {
			writeError(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}

		zoneID, err := strconv.Atoi(parts[2])
		if err != nil {
			writeError(w, "Invalid Zone ID", http.StatusBadRequest)
			return
		}

		// Find the zone
		zone := worldManager.Zones[zoneID]
		if zone == nil {
			writeError(w, "Zone not found: "+strconv.Itoa(zoneID), http.StatusNotFound)
			return
		}

		// Collect entities
		var snapshot ZoneSnapshot
		for _, entity := range zone.GetEntities() {
			x, y := entity.GetPosition()
			snapshot.EntitySnapshots = append(snapshot.EntitySnapshots, EntitySnapshot{
				ID:     entity.GetUUID(),
				ZoneID: zoneID,
				Type: 	  entity.GetType(),
				X:        x,
				Y:        y,
				Data: entity.GetSnapshotData(),
			})
		}

		// Respond with JSON
		writeJSON(w, snapshot)
	}
}

// handleZoneMap returns a handler for the /zonemap endpoint.
func handleZoneMap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Only allow GET requests
		if r.Method != http.MethodGet {
			writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Ensure the endpoint is exactly /zonemap
		if r.URL.Path != "/zonemap" {
			writeError(w, "Invalid endpoint", http.StatusBadRequest)
			return
		}

		// Create the response with the zone map
		response := ZoneMapResponse{
			ZoneMap: config.ZoneMap,
		}

		// Respond with JSON
		writeJSON(w, response)
	}
}
	*/