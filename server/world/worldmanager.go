package world

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"thereaalm/action"
	"thereaalm/config"
	"thereaalm/entity"
	"thereaalm/entity/resourceentity"
	"thereaalm/interfaces"
	"thereaalm/types"
	"thereaalm/utils"
	"thereaalm/web3"
	"time"
)

const (
	ZoneTiles = 512
)

type WorldManager struct {
	Zones          []interfaces.IZone
	WorkerCount    int
	SpeedMultiplier float64       // 1.0 = normal, 2.0 = double speed
	GameTime        time.Duration // Simulated game time
	LastUpdate      time.Time     // Real time of last update
	SpawnAreas     []*SpawnArea  // Spawn areas loaded from tilemap
}

func NewWorldManager(workerCount int) *WorldManager {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
	}

	manager := &WorldManager{
		WorkerCount: workerCount,
		SpeedMultiplier: 1.0,
		GameTime:        0,
		LastUpdate:      time.Now(),
		SpawnAreas:      make([]*SpawnArea, 0),
	}

	// Initialize zones
	zoneID := 0
	for y, row := range config.ZoneMap {
		for x, zoneType := range row {
			if zoneType == "" {
				continue
			}
			zone := NewZone(manager, zoneID, ZoneTiles, ZoneTiles, x*ZoneTiles, y*ZoneTiles, 64)
			manager.Zones = append(manager.Zones, zone)
			zoneID++
		}
	}

	if len(manager.Zones) == 0 {
		log.Fatal("Error: No active zones initialized!")
		return manager
	}

	manager.SetSimulationSpeed(1)

	// Load the tilemap for zone 0
	spawnAreas, err := LoadTilemap("../shared/tilemaps/maps/yield_fields_2.json", manager, 42)
	if err != nil {
		log.Fatalf("Failed to load tilemap: %v", err)
	}
	manager.SpawnAreas = spawnAreas

	// load test entities
	manager.loadTestEntities()

	log.Printf("World initialized with %d active zones and %d spawn areas.", len(manager.Zones), len(manager.SpawnAreas))
	return manager
}

func (wm *WorldManager) loadTestEntities() {
	// use one of our gotchisMap
	idSlice := []string{"4285", "19005", "21550", "8281", "5401",
		"13401", "2325", "1699", "22313", "4895",
		"15434", "19553", "22473", "19450", "5410",
		"928", "22566"}

	// grab a few gotchis to use
	gotchisMap, err := web3.FetchGotchisByIDs(idSlice)
	if err != nil {
		fmt.Printf("Error fetching gotchis: %v\n", err)
		return
	}

	// lets start by placing entities in zone 42 only for now
	zone := wm.Zones[42]
	zoneWorldX, zoneWorldY := zone.GetPosition()

	log.Println("Zone minimum: ", zoneWorldX, zoneWorldY)
	log.Println("Zone maximum: ", zoneWorldX+zone.GetWidth(), zoneWorldY+zone.GetHeight())

	// gotchis
	for i := 0; i < 400; i++ {
		var posX, posY int
		var found bool

		// Try to find a spawn area for gotchis
		var spawnArea *SpawnArea
		for _, sa := range wm.SpawnAreas {
			if sa.EntitySpawnType == "gotchi" {
				spawnArea = sa
				break
			}
		}

		if spawnArea != nil {
			posX, posY, found = spawnArea.GetRandomPosition()
			if !found {
				log.Println("No valid position found in gotchi spawn area, falling back to random position")
			}
		} else {
			log.Println("No gotchi spawn area found, falling back to random position")
		}

		// Fallback to random position if no spawn area or no valid position
		if spawnArea == nil || !found {
			for attempts := 0; attempts < 10; attempts++ {
				posX = rand.Intn(ZoneTiles) + zoneWorldX
				posY = rand.Intn(ZoneTiles) + zoneWorldY
				if zone.IsPositionAvailable(posX, posY) {
					found = true
					break
				}
			}
			if !found {
				log.Println("No available random position found for gotchi after 10 attempts")
				continue
			}
		}

		if !zone.IsPositionAvailable(posX, posY) {
			emptyX, emptyY, found := wm.FindNearbyAvailablePosition(posX, posY, 10, 0)
			if found {
				posX = emptyX
				posY = emptyY
			} else {
				log.Println("No nearby available position found for gotchi")
				continue
			}
		}

		idIndex := rand.Intn(len(idSlice))
		gotchiId := idSlice[idIndex]

		jobSlice := []string{"mercenary", "farmer",
			"minerjack", "builder", "explorer"}
		jobIndex := rand.Intn(5)

		_, exists := gotchisMap[gotchiId]
		if !exists {
			generateGenericGotchi(wm, posX, posY,
				web3.DefaultSubgraphGotchiData, jobSlice[jobIndex])
		} else {
			generateGenericGotchi(wm, posX, posY,
				gotchisMap[gotchiId], jobSlice[jobIndex])
		}
	}

	// lickvoids
	for i := 0; i < 100; i++ {
		var posX, posY int
		var found bool

		// Try to find a spawn area for lickvoids
		var spawnArea *SpawnArea
		for _, sa := range wm.SpawnAreas {
			if sa.EntitySpawnType == "lickvoid" {
				spawnArea = sa
				break
			}
		}

		if spawnArea != nil {
			posX, posY, found = spawnArea.GetRandomPosition()
			if !found {
				log.Println("No valid position found in lickvoid spawn area, falling back to random position")
			}
		} else {
			log.Println("No lickvoid spawn area found, falling back to random position")
		}

		// Fallback to random position if no spawn area or no valid position
		if spawnArea == nil || !found {
			for attempts := 0; attempts < 10; attempts++ {
				posX = rand.Intn(ZoneTiles) + zoneWorldX
				posY = rand.Intn(ZoneTiles) + zoneWorldY
				if zone.IsPositionAvailable(posX, posY) {
					found = true
					break
				}
			}
			if !found {
				log.Println("No available random position found for lickvoid after 10 attempts")
				continue
			}
		}

		if !zone.IsPositionAvailable(posX, posY) {
			log.Println("Position not available for lickvoid")
			continue
		}

		lickvoid := entity.NewLickVoid(posX, posY)
		wm.AddEntity(lickvoid)
		lickvoid.SpawnInterval_s = 5
	}

	// lickquidators
	for i := 0; i < 100; i++ {
		var posX, posY int
		var found bool

		// Try to find a spawn area for lickquidators
		var spawnArea *SpawnArea
		for _, sa := range wm.SpawnAreas {
			if sa.EntitySpawnType == "lickquidator" {
				spawnArea = sa
				break
			}
		}

		if spawnArea != nil {
			posX, posY, found = spawnArea.GetRandomPosition()
			if !found {
				log.Println("No valid position found in lickquidator spawn area, falling back to random position")
			}
		} else {
			log.Println("No lickquidator spawn area found, falling back to random position")
		}

		// Fallback to random position if no spawn area or no valid position
		if spawnArea == nil || !found {
			for attempts := 0; attempts < 10; attempts++ {
				posX = rand.Intn(ZoneTiles) + zoneWorldX
				posY = rand.Intn(ZoneTiles) + zoneWorldY
				if zone.IsPositionAvailable(posX, posY) {
					found = true
					break
				}
			}
			if !found {
				log.Println("No available random position found for lickquidator after 10 attempts")
				continue
			}
		}

		if !zone.IsPositionAvailable(posX, posY) {
			log.Println("Position not available for lickquidator")
			continue
		}

		generateGenericLickquidator(wm, posX, posY)
	}

	// resources
	for i := 0; i < 800; i++ {
		// pick rando resource
		resourceSlice := []string{"fomoberrybush", "kekwoodtree", "alphaslateboulders"}
		resource := resourceSlice[rand.Intn(3)]
		var posX, posY int
		var found bool

		// Try to find a spawn area for resources
		var spawnArea *SpawnArea
		for _, sa := range wm.SpawnAreas {
			if sa.EntitySpawnType == resource {
				spawnArea = sa
				break
			}
		}

		if spawnArea != nil {
			posX, posY, found = spawnArea.GetRandomPosition()
			if !found {
				log.Printf("No valid position found in %s spawn area, falling back to random position", resource)
			}
		} else {
			log.Printf("No %s spawn area found, falling back to random position", resource)
		}

		// Fallback to random position if no spawn area or no valid position
		if spawnArea == nil || !found {
			for attempts := 0; attempts < 10; attempts++ {
				posX = rand.Intn(ZoneTiles) + zoneWorldX
				posY = rand.Intn(ZoneTiles) + zoneWorldY
				if zone.IsPositionAvailable(posX, posY) {
					found = true
					break
				}
			}
			if !found {
				log.Printf("No available random position found for %s after 10 attempts", resource)
				continue
			}
		}

		if !zone.IsPositionAvailable(posX, posY) {
			emptyX, emptyY, found := wm.FindNearbyAvailablePosition(posX, posY, 10, 0)
			if found {
				posX = emptyX
				posY = emptyY
			} else {
				log.Printf("No nearby available position found for %s", resource)
				continue
			}
		}

		if resource == "fomoberrybush" {
			wm.AddEntity(resourceentity.NewFomoBerryBush(posX, posY))
		} else if resource == "kekwoodtree" {
			wm.AddEntity(resourceentity.NewKekWoodTree(posX, posY))
		} else {
			wm.AddEntity(resourceentity.NewAlphaSlateBoulders(posX, posY))
		}
	}

	// altars
	for i := 0; i < 100; i++ {
		var posX, posY int
		var found bool

		// Try to find a spawn area for altars
		var spawnArea *SpawnArea
		for _, sa := range wm.SpawnAreas {
			if sa.EntitySpawnType == "altar" {
				spawnArea = sa
				break
			}
		}

		if spawnArea != nil {
			posX, posY, found = spawnArea.GetRandomPosition()
			if !found {
				log.Println("No valid position found in altar spawn area, falling back to random position")
			}
		} else {
			log.Println("No altar spawn area found, falling back to random position")
		}

		// Fallback to random position if no spawn area or no valid position
		if spawnArea == nil || !found {
			for attempts := 0; attempts < 10; attempts++ {
				posX = rand.Intn(ZoneTiles) + zoneWorldX
				posY = rand.Intn(ZoneTiles) + zoneWorldY
				if zone.IsPositionAvailable(posX, posY) {
					found = true
					break
				}
			}
			if !found {
				log.Println("No available random position found for altar after 10 attempts")
				continue
			}
		}

		if !wm.IsPositionAvailable(posX, posY) {
			log.Println("Position not available for altar")
			continue
		}

		wm.AddEntity(entity.NewAltar(posX, posY))
	}

	// shops
	for i := 0; i < 50; i++ {
		var posX, posY int
		var found bool

		// Try to find a spawn area for shops
		var spawnArea *SpawnArea
		for _, sa := range wm.SpawnAreas {
			if sa.EntitySpawnType == "shop" {
				spawnArea = sa
				break
			}
		}

		if spawnArea != nil {
			posX, posY, found = spawnArea.GetRandomPosition()
			if !found {
				log.Println("No valid position found in shop spawn area, falling back to random position")
			}
		} else {
			log.Println("No shop spawn area found, falling back to random position")
		}

		// Fallback to random position if no spawn area or no valid position
		if spawnArea == nil || !found {
			for attempts := 0; attempts < 10; attempts++ {
				posX = rand.Intn(ZoneTiles) + zoneWorldX
				posY = rand.Intn(ZoneTiles) + zoneWorldY
				if zone.IsPositionAvailable(posX, posY) {
					found = true
					break
				}
			}
			if !found {
				log.Println("No available random position found for shop after 10 attempts")
				continue
			}
		}

		if !wm.IsPositionAvailable(posX, posY) {
			log.Println("Position not available for shop")
			continue
		}

		wm.AddEntity(entity.NewShop(posX, posY))
	}
}

func generateGenericLickquidator(wm *WorldManager, x, y int) {
	lickquidator := entity.NewLickquidator(x, y)
	wm.AddEntity(lickquidator)

	lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.3,
		&types.TargetSpec{
			TargetType: "gotchi",
			TargetCriterion: "nearest",
		}))
	lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.3,
		&types.TargetSpec{
			TargetType: "altar",
			TargetCriterion: "nearest",
		}))
	lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.3,
		&types.TargetSpec{
			TargetType: "shop",
			TargetCriterion: "nearest",
		}))
}

func generateGenericGotchi(wm *WorldManager, x, y int,
	subgraphData web3.SubgraphGotchiData, job string) {

	newGotchi := entity.NewGotchi(x, y, subgraphData)
	wm.AddEntity(newGotchi)
	newGotchi.Job = job

	// determine action profile
	var profile GotchiBehaviorProfile
	if job == "mercenary" {
		profile = MercenaryProfile
	} else if job == "farmer" {
		profile = FarmerProfile
	} else if job == "minerjack" {
		profile = MinerJackProfile
	} else if job == "builder" {
		profile = BuilderProfile
	} else {
		profile = ExplorerProfile
	}

	profile(newGotchi)
}

// generateBerryBushes generates 100 unique berry bushes in a 100x100 area within the specified zone
func generateBerryBushes(wm *WorldManager, zoneID int, zoneX, zoneY int) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Track occupied positions to avoid duplicates
	occupied := make(map[[2]int]bool)
	bushesToGenerate := 400
	zone := wm.Zones[zoneID]

	for len(occupied) < bushesToGenerate {
		// Generate random coordinates within 100x100 area
		x := zoneX + rand.Intn(60) // 0+zoneX to 100+zoneX
		y := zoneY + rand.Intn(60) // 0+zoneY to 100+zoneY

		// Check if position is already occupied
		pos := [2]int{x, y}
		if !occupied[pos] {
			occupied[pos] = true
			zone.AddEntity(resourceentity.NewFomoBerryBush(x, y))
		}
	}
}

func (wm *WorldManager) Run() {
	log.Printf("World is running with %d workers...", wm.WorkerCount)
	go wm.updateLoop()
}

func (wm *WorldManager) updateLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		dt_s := now.Sub(wm.LastUpdate).Seconds()
		wm.LastUpdate = now
		wm.updateZonesParallel(dt_s)
	}
}

func (wm *WorldManager) updateZonesParallel(dt_s float64) {
	scaledDt := dt_s * wm.SpeedMultiplier
	wm.GameTime += time.Duration(scaledDt * float64(time.Second))

	var wg sync.WaitGroup
	jobs := make(chan interfaces.IZone, len(wm.Zones))

	for i := 0; i < wm.WorkerCount; i++ {
		wg.Add(1)
		go wm.zoneWorker(jobs, dt_s*scaledDt, &wg)
	}

	for _, zone := range wm.Zones {
		jobs <- zone
	}
	close(jobs)

	wg.Wait()
}

func (wm *WorldManager) zoneWorker(jobs <-chan interfaces.IZone, dt_s float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for zone := range jobs {
		zone.Update(dt_s)
	}
}

// Time access methods
func (wm *WorldManager) Now() time.Duration {
	return wm.GameTime
}

func (wm *WorldManager) Since(startTime time.Duration) time.Duration {
	if wm.GameTime < startTime {
		return 0 // Prevent negative durations
	}
	return wm.GameTime - startTime
}

func (wm *WorldManager) SetSimulationSpeed(multiplier float64) {
	if multiplier <= 0 {
		log.Println("WARNING: Speed multiplier must be positive, ignoring:", multiplier)
		return
	}
	wm.SpeedMultiplier = multiplier
	log.Printf("Simulation speed set to %.2fx", multiplier)
}

func (wm *WorldManager) IsPositionAvailable(x, y int) bool {
	zone := wm.getZoneForPosition(x, y)

	if zone == nil {
		return false
	}

	return zone.IsPositionAvailable(x, y)
}

func (wm *WorldManager) AddEntity(e interfaces.IEntity) {
	ex, ey := e.GetPosition()

	eZone := wm.getZoneForPosition(ex, ey)

	if eZone != nil {
		eZone.AddEntity(e)
	}
}

func (wm *WorldManager) RemoveEntity(e interfaces.IEntity) {
	eZone := e.GetZone()

	if eZone != nil {
		eZone.RemoveEntity(e)
	}
}

func (wm *WorldManager) getZoneForPosition(x, y int) interfaces.IZone {
	var zone interfaces.IZone
	for _, z := range wm.Zones {
		zoneWorldX, zoneWorldY := z.GetPosition()
		if x >= zoneWorldX && y >= zoneWorldY &&
			x < zoneWorldX+z.GetWidth() && y < zoneWorldY+z.GetHeight() {
			zone = z
			break
		}
	}

	return zone
}

func (wm *WorldManager) FindNearbyAvailablePosition(x, y, radius, minimumGap int) (int, int, bool) {
	// get zone for position
	zone := wm.getZoneForPosition(x, y)
	if zone == nil {
		log.Println("No zone available at the position: ", x, y)
		return 0, 0, false
	}

	// use the zones utility function
	emptyX, emptyY, found := zone.FindNearbyAvailablePosition(x, y, radius, minimumGap)
	if found {
		return emptyX, emptyY, true
	}

	return 0, 0, false
}

func (wm *WorldManager) GetDistance(x1, y1, x2, y2 int) int {
	return utils.Abs(x1-x2) + utils.Abs(y1-y2)
}

/*

// thereaalm/world/manager.go
package world

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"thereaalm/action"
	"thereaalm/config"
	"thereaalm/entity"
	"thereaalm/entity/resourceentity"
	"thereaalm/interfaces"
	"thereaalm/types"
	"thereaalm/utils"
	"thereaalm/web3"

	// "thereaalm/storage"
	"time"
)

const (
    ZoneTiles  = 512
    // ZoneHeight = 128
)

type WorldManager struct {
    Zones       []interfaces.IZone
    WorkerCount int
    SpeedMultiplier float64       // 1.0 = normal, 2.0 = double speed
    GameTime        time.Duration // Simulated game time
    LastUpdate      time.Time     // Real time of last update
    SpawnAreas []*SpawnArea
}

func NewWorldManager(workerCount int) *WorldManager {
    if workerCount <= 0 {
        workerCount = runtime.NumCPU()
    }

    manager := &WorldManager{
        WorkerCount: workerCount,
        SpeedMultiplier: 1.0,
        GameTime:        0,
        LastUpdate:      time.Now(),
        SpawnAreas: make([]*SpawnArea, 0),
    }

    // Initialize zones
    zoneID := 0
    for y, row := range config.ZoneMap {
        for x, zoneType := range row {
            if zoneType == "" {
                continue
            }
            zone := NewZone(manager, zoneID, ZoneTiles, ZoneTiles, x*ZoneTiles, y*ZoneTiles, 64)
            manager.Zones = append(manager.Zones, zone)
            zoneID++
        }
    }

    if len(manager.Zones) == 0 {
        log.Fatal("Error: No active zones initialized!")
        return manager
    }

    manager.SetSimulationSpeed(1)

	// Load the tilemap for zone 0
	spawnAreas, err := LoadTilemap("../shared/tilemaps/maps/yield_fields_2.json", manager, 42)
	if err != nil {
		log.Fatalf("Failed to load tilemap: %v", err)
	}
	manager.SpawnAreas = spawnAreas

    // load test entities
    manager.loadTestEntities()

    log.Printf("World initialized with %d active zones.", len(manager.Zones))
    return manager
}

func (wm *WorldManager) loadTestEntities() {
    // use one of our gotchisMap
    idSlice := []string{"4285", "19005", "21550", "8281", "5401",
                        "13401", "2325", "1699", "22313"}

    // grab a few gotchis to use
    gotchisMap, err := web3.FetchGotchisByIDs(idSlice)
    if err != nil {
        fmt.Printf("Error fetching gotchis: %v\n", err)
        return
    }

    // lets start by placing entities in zone 42 only for now
    zone := wm.Zones[42]
    zoneWorldX, zoneWorldY := zone.GetPosition()

    log.Println("Zone minimum: ", zoneWorldX, zoneWorldY)
    log.Println("Zone maximum: ", zoneWorldX + zone.GetWidth(), zoneWorldY + zone.GetHeight())

    // gotchis
    for i := 0; i < 400; i++ {

        idIndex := rand.Intn(len(idSlice))

        gotchiId := idSlice[idIndex]

        posX := rand.Intn(ZoneTiles) + zoneWorldX
        posY := rand.Intn(ZoneTiles) + zoneWorldY

        // log.Println("Try make gotchi at: ", posX, posY)

        if !zone.IsPositionAvailable(posX, posY) {
            // log.Println("Position not available, try find neary available one...")
            emptyX, emptyY, found := wm.FindNearbyAvailablePosition(posX, posY, 10, 0)
            if found {
                // log.Println("Found position at: ", posX, posY)
                posX = emptyX
                posY = emptyY
            } else {
                continue
            }
        }

        jobSlice := []string{"mercenary", "farmer",
            "minerjack", "builder", "explorer"}
        jobIndex := rand.Intn(5)

        _, exists := gotchisMap[gotchiId]
        if !exists {
            generateGenericGotchi(wm, posX, posY, 
                web3.DefaultSubgraphGotchiData, jobSlice[jobIndex])
        } else {
            generateGenericGotchi(wm, posX, posY, 
                gotchisMap[gotchiId], jobSlice[jobIndex])
        }

    }

    // lickvoids
    for i := 0; i < 100; i++ {
        posX := rand.Intn(ZoneTiles) + zoneWorldX
        posY := rand.Intn(ZoneTiles) + zoneWorldY

        if !wm.IsPositionAvailable(posX, posY) {
            continue
        }

        lickvoid := entity.NewLickVoid(posX, posY)
        wm.AddEntity(lickvoid)
        lickvoid.SpawnInterval_s = 5
    }

    // resources
    for i := 0; i < 800; i++ {
        // pick rando resource
        resourceSlice := []string{"fomoberrybush", "kekwoodtree", "alphaslateboulders" }
        resource := resourceSlice[rand.Intn(3)]

        posX := rand.Intn(ZoneTiles) + zoneWorldX
        posY := rand.Intn(ZoneTiles) + zoneWorldY

        if !zone.IsPositionAvailable(posX, posY) {
            // log.Println("Position not available, try find neary available one...")
            emptyX, emptyY, found := wm.FindNearbyAvailablePosition(posX, posY, 10, 0)
            if found {
                // log.Println("Found position at: ", posX, posY)
                posX = emptyX
                posY = emptyY
            } else {
                continue
            }
        }

        if resource == "fomoberrybush" {
            wm.AddEntity(resourceentity.NewFomoBerryBush(posX, posY))
        } else if resource == "kekwoodtree" {
            wm.AddEntity(resourceentity.NewKekWoodTree(posX, posY))
        } else {
            wm.AddEntity(resourceentity.NewAlphaSlateBoulders(posX, posY))
        }
    }

    // altars
    for i := 0; i < 100; i++ {
        posX := rand.Intn(ZoneTiles) + zoneWorldX
        posY := rand.Intn(ZoneTiles) + zoneWorldY

        if !wm.IsPositionAvailable(posX, posY) {
            continue
        }

        wm.AddEntity(entity.NewAltar(posX, posY))
    }

    // shops
    for i := 0; i < 50; i++ {
        posX := rand.Intn(ZoneTiles) + zoneWorldX
        posY := rand.Intn(ZoneTiles) + zoneWorldY

        if !wm.IsPositionAvailable(posX, posY) {
            continue
        }

        wm.AddEntity(entity.NewShop(posX, posY))
    }

    // ENTITIES
    // resources
    // bush := resourceentity.NewFomoBerryBush(12+zoneX, 10+zoneY)
    // wm.Zones[42].AddEntity(bush)

    // tree := resourceentity.NewKekWoodTree(15+zoneX, 9+zoneY)
    // wm.Zones[42].AddEntity(tree)

    // boulders := resourceentity.NewAlphaSlateBoulders(13+zoneX, 12+zoneY)
    // wm.Zones[42].AddEntity(boulders)

    // generateBerryBushes(wm, 42, zoneX, zoneY)

    // // shop
    // shop := entity.NewShop(6+zoneX, 12+zoneY)
    // wm.Zones[42].AddEntity(shop)

    // // altar
    // altar := entity.NewAltar(19+zoneX, 12+zoneY)
    // wm.Zones[42].AddEntity(altar)
    // altar.SetStat(stattypes.Pulse, 20)
}

func generateGenericLickquidator(wm *WorldManager, x, y int) {
    lickquidator := entity.NewLickquidator(x, y)
    wm.AddEntity(lickquidator)

    lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.3,
        &types.TargetSpec{
            TargetType: "gotchi",
            TargetCriterion: "nearest",
        }))
    lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.3,
        &types.TargetSpec{
            TargetType: "altar",
            TargetCriterion: "nearest",
        }))
    lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.3,
        &types.TargetSpec{
            TargetType: "shop",
            TargetCriterion: "nearest",
        }))
}

func generateGenericGotchi(wm *WorldManager, x, y int, 
    subgraphData web3.SubgraphGotchiData, job string) {

    newGotchi := entity.NewGotchi(x, y, subgraphData)
    wm.AddEntity(newGotchi)
    newGotchi.Job = job

    // determine action profile
    var profile GotchiBehaviorProfile
    if job == "mercenary" {
        profile = MercenaryProfile
    } else if job == "farmer" {
        profile = FarmerProfile
    } else if job == "minerjack" {
        profile = MinerJackProfile
    } else if job == "builder" {
        profile = BuilderProfile
    } else {
        profile = ExplorerProfile
    }

    profile(newGotchi)
}

// generateBerryBushes generates 100 unique berry bushes in a 100x100 area within the specified zone
func generateBerryBushes(wm *WorldManager, zoneID int, zoneX, zoneY int) {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())

    // Track occupied positions to avoid duplicates
    occupied := make(map[[2]int]bool)
    bushesToGenerate := 400
    zone := wm.Zones[zoneID]

    for len(occupied) < bushesToGenerate {
        // Generate random coordinates within 100x100 area
        x := zoneX + rand.Intn(60) // 0+zoneX to 100+zoneX
        y := zoneY + rand.Intn(60) // 0+zoneY to 100+zoneY

        // Check if position is already occupied
        pos := [2]int{x, y}
        if !occupied[pos] {
            occupied[pos] = true
            zone.AddEntity(resourceentity.NewFomoBerryBush(x, y))
        }
    }
}

func (wm *WorldManager) Run() {
    log.Printf("World is running with %d workers...", wm.WorkerCount)
    go wm.updateLoop()
}

func (wm *WorldManager) updateLoop() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        now := time.Now()
        dt_s := now.Sub(wm.LastUpdate).Seconds()
        wm.LastUpdate = now
        wm.updateZonesParallel(dt_s)
    }
}

func (wm *WorldManager) updateZonesParallel(dt_s float64) {
    scaledDt := dt_s * wm.SpeedMultiplier
    wm.GameTime += time.Duration(scaledDt * float64(time.Second))

    var wg sync.WaitGroup
    jobs := make(chan interfaces.IZone, len(wm.Zones))

    for i := 0; i < wm.WorkerCount; i++ {
        wg.Add(1)
        go wm.zoneWorker(jobs, dt_s * scaledDt, &wg)
    }

    for _, zone := range wm.Zones {
        jobs <- zone
    }
    close(jobs)

    wg.Wait()
}

func (wm *WorldManager) zoneWorker(jobs <-chan interfaces.IZone, dt_s float64, wg *sync.WaitGroup) {
    defer wg.Done()

    for zone := range jobs {
        zone.Update(dt_s)
    }
}

// Time access methods
func (wm *WorldManager) Now() time.Duration {
    return wm.GameTime
}

func (wm *WorldManager) Since(startTime time.Duration) time.Duration {
    if wm.GameTime < startTime {
        return 0 // Prevent negative durations
    }
    return wm.GameTime - startTime
}

func (wm *WorldManager) SetSimulationSpeed(multiplier float64) {
    if multiplier <= 0 {
        log.Println("WARNING: Speed multiplier must be positive, ignoring:", multiplier)
        return
    }
    wm.SpeedMultiplier = multiplier
    log.Printf("Simulation speed set to %.2fx", multiplier)
}

func (wm *WorldManager) IsPositionAvailable(x, y int) bool {
    zone := wm.getZoneForPosition(x, y)

    if zone == nil {
        return false
    }

    return zone.IsPositionAvailable(x, y)
}

func (wm *WorldManager) AddEntity(e interfaces.IEntity) {
    ex, ey := e.GetPosition()

    eZone := wm.getZoneForPosition(ex, ey)

    if eZone != nil {
        eZone.AddEntity(e)
    }
}

func (wm *WorldManager) RemoveEntity(e interfaces.IEntity) {
    eZone := e.GetZone()
    
    if eZone != nil {
        eZone.RemoveEntity(e)
    }
}

func (wm *WorldManager) getZoneForPosition(x, y int) interfaces.IZone {
    var zone interfaces.IZone
    for _, z := range wm.Zones {
        zoneWorldX, zoneWorldY := z.GetPosition()
        if x >= zoneWorldX && y >= zoneWorldY && 
        x < zoneWorldX + z.GetWidth() && y < zoneWorldY + z.GetHeight() {
            zone = z
            break    
        }
    }

    return zone
}

func (wm *WorldManager) FindNearbyAvailablePosition(x, y, radius, minimumGap int) (int, int, bool) {
    // get zone for position
    zone := wm.getZoneForPosition(x, y)
    if zone == nil {
        log.Println("No zone available at the position: ", x, y)
        return 0, 0, false
    }

    // use the zones utility function
    emptyX, emptyY, found := zone.FindNearbyAvailablePosition(x, y, radius, minimumGap)
    if found {
        return emptyX, emptyY, true
    }

    return 0, 0, false
}

func (wm *WorldManager) GetDistance(x1, y1, x2, y2 int) int {
    return utils.Abs(x1 - x2) + utils.Abs(y1 - y2)
}
   */