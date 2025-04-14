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
	"thereaalm/stattypes"
	"thereaalm/types"
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

    // load tilemap
    // Load the tilemap for zone 0
    if err := LoadTilemap("../shared/tilemaps/maps/yield_fields_2.json", manager, 42); err != nil {
        log.Fatalf("Failed to load tilemap: %v", err)
    }

    // load test entities
    manager.loadTestEntities()

    log.Printf("World initialized with %d active zones.", len(manager.Zones))
    return manager
}

func (wm *WorldManager) loadTestEntities() {
    // grab a few gotchis to use
    gotchisMap, err := web3.FetchGotchisByIDs([]string{"4285", "19005", "21550",
    "8281", "5401"})
    if err != nil {
        fmt.Printf("Error fetching gotchis: %v\n", err)
        return
    }

    // entity generation
    zoneX, zoneY := wm.Zones[42].GetPosition()

    // gotchis
    for i := 0; i < 400; i++ {
        // use one of our gotchisMap
        idSlice := []string{"4285", "19005", "21550", "8281", "5401"}
        index := rand.Intn(5)

        gotchiId := idSlice[index]

        posX := rand.Intn(ZoneTiles) + zoneX
        posY := rand.Intn(ZoneTiles) + zoneY

        jobSlice := []string{"mercenary", "farmer",
            "minerjack", "builder", "explorer"}

        generateGenericGotchi(wm, 42, posX, posY, 
            gotchisMap[gotchiId], jobSlice[index])
    }

    log.Println("here")

    // lickvoids
    for i := 0; i < 100; i++ {
        posX := rand.Intn(ZoneTiles) + zoneX
        posY := rand.Intn(ZoneTiles) + zoneY

        if wm.Zones[42].IsTileOccupied(posX, posY) {
            continue
        }

        lickvoid := entity.NewLickVoid(posX, posY)
        wm.Zones[42].AddEntity(lickvoid)
        lickvoid.SpawnInterval_s = 5
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

    // // gotchis
    // generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["4285"])
    // generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["19005"])
    // generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["21550"])
    // generateGenericGotchi(wm, 42, 30+zoneX, 30+zoneY, gotchisMap["8281"])
    // generateGenericGotchi(wm, 42, 30+zoneX, 30+zoneY, gotchisMap["5401"])

    // // lickquidators
    // generateGenericLickquidator(wm, 42, 9+zoneX, 14+zoneY)
    // generateGenericLickquidator(wm, 42, 18+zoneX, 23+zoneY)
    // generateGenericLickquidator(wm, 42, 15+zoneX, 19+zoneY)

    // // altar
    // altar := entity.NewAltar(19+zoneX, 12+zoneY)
    // wm.Zones[42].AddEntity(altar)
    // altar.SetStat(stattypes.Pulse, 20)

    // // lickvoid
    // lickvoid := entity.NewLickVoid(25+zoneX, 14+zoneY)
    // wm.Zones[42].AddEntity(lickvoid)
    // lickvoid.SpawnInterval_s = 5
}

func generateGenericLickquidator(wm *WorldManager, zoneID int, x, y int) {
    lickquidator := entity.NewLickquidator(x, y)
    wm.Zones[zoneID].AddEntity(lickquidator)

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

func generateGenericGotchi(wm *WorldManager, zoneID int, x, y int, 
    subgraphData web3.SubgraphGotchiData, job string) {

    newGotchi := entity.NewGotchi(x, y, subgraphData)
    wm.Zones[zoneID].AddEntity(newGotchi)
    newGotchi.SetStat(stattypes.Pulse, 300)
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