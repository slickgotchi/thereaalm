// thereaalm/world/manager.go
package world

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"thereaalm/action"
	"thereaalm/action/resourceaction"
	"thereaalm/config"
	"thereaalm/entity"
	"thereaalm/entity/resourceentity"
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
    Zones       []*types.Zone
    WorkerCount int
}

func NewWorldManager(workerCount int) *WorldManager {
    if workerCount <= 0 {
        workerCount = runtime.NumCPU()
    }

    manager := &WorldManager{
        WorkerCount: workerCount,
    }

    // Initialize zones
    zoneID := 0
    for y, row := range config.ZoneMap {
        for x, zoneType := range row {
            if zoneType == "" {
                continue
            }
            zone := types.NewZone(zoneID, ZoneTiles, ZoneTiles, x*ZoneTiles, y*ZoneTiles, 64)
            manager.Zones = append(manager.Zones, zone)
            zoneID++
        }
    }

    if len(manager.Zones) == 0 {
        log.Fatal("Error: No active zones initialized!")
        return manager
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

    zoneX := wm.Zones[42].X
    zoneY := wm.Zones[42].Y

    // ENTITIES
    // resources
    bush := resourceentity.NewFomoBerryBush(42, 12+zoneX, 10+zoneY)
    wm.Zones[42].AddEntity(bush)

    tree := resourceentity.NewKekWoodTree(42, 15+zoneX, 9+zoneY)
    wm.Zones[42].AddEntity(tree)

    boulders := resourceentity.NewAlphaSlateBoulders(42, 13+zoneX, 12+zoneY)
    wm.Zones[42].AddEntity(boulders)

    // generateBerryBushes(wm, 42, zoneX, zoneY)

    // shop
    shop := entity.NewShop(42, 6+zoneX, 12+zoneY)
    wm.Zones[42].AddEntity(shop)

    // gotchis
    generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["4285"])
    generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["19005"])
    generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["21550"])
    generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["8281"])
    generateGenericGotchi(wm, 42, 10+zoneX, 10+zoneY, gotchisMap["21550"])

    // lickquidators
    generateGenericLickquidator(wm, 42, 9+zoneX, 14+zoneY)
    generateGenericLickquidator(wm, 42, 18+zoneX, 23+zoneY)
    generateGenericLickquidator(wm, 42, 15+zoneX, 19+zoneY)

    // altar
    altar := entity.NewAltar(42, 19+zoneX, 12+zoneY)
    wm.Zones[42].AddEntity(altar)

    // // TEMPORARY: start only in zone 42 for now
    // zoneIndex := 42

    // // place gotchis across all available zones, in ZoneMap
    // // start with a known seed
    // // r := rand.New(rand.NewSource(123))

    // for _, gd := range gotchiData {
    //     // pick random zone
    //     // zoneIndex := r.Intn(len(wm.Zones))


    //     // pick random location in the zone
    //     // x := r.Intn(ZoneTiles) + wm.Zones[zoneIndex].X
    //     // y := r.Intn(ZoneTiles) + wm.Zones[zoneIndex].Y
    //     x := wm.Zones[zoneIndex].X + 5
    //     y := wm.Zones[zoneIndex].Y + 8

    //     // create new gotchi
    // }
}

func generateGenericLickquidator(wm *WorldManager, zoneID int, x, y int) {
    lickquidator := entity.NewLickquidator(zoneID, x, y)
    wm.Zones[42].AddEntity(lickquidator)

    lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.2,
        &action.TargetSpec{
            TargetType: "gotchi",
            TargetCriterion: "nearest",
        }))
    lickquidator.AddActionToPlan(action.NewAttackAction(lickquidator, nil, 0.8,
        &action.TargetSpec{
            TargetType: "altar",
            TargetCriterion: "nearest",
        }))
}

func generateGenericGotchi(wm *WorldManager, zoneID int, x, y int, 
    subgraphData web3.SubgraphGotchiData) {

    newGotchi := entity.NewGotchi(zoneID, x, y, subgraphData)
    wm.Zones[42].AddEntity(newGotchi)

    // ACTIONS
    newGotchi.AddActionToPlan(
        resourceaction.NewForageAction(newGotchi, nil, 0.3, 
            &action.TargetSpec{
                TargetType: "fomoberrybush",
                TargetCriterion: "nearest",
            }))
    newGotchi.AddActionToPlan(
        resourceaction.NewChopAction(newGotchi, nil, 0.3, 
            &action.TargetSpec{
                TargetType: "kekwoodtree",
                TargetCriterion: "nearest",
            }))
    newGotchi.AddActionToPlan(
        resourceaction.NewMineAction(newGotchi, nil, 0.3, 
            &action.TargetSpec{
                TargetType: "alphaslateboulders",
                TargetCriterion: "nearest",
            }))
    newGotchi.AddActionToPlan(
        action.NewSellAction(newGotchi, nil, 0.5, 
            &action.TargetSpec{
                TargetType: "shop",
                TargetCriterion: "nearest",
            }))
    newGotchi.AddActionToPlan(
        action.NewAttackAction(newGotchi, nil, 0.3, 
            &action.TargetSpec{
                TargetType: "lickquidator",
                TargetCriterion: "nearest",
            }))
    newGotchi.AddActionToPlan(
        action.NewRoamAction(newGotchi, nil, 0.1, nil))
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
            zone.AddEntity(resourceentity.NewFomoBerryBush(zoneID, x, y))
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

    lastUpdate := time.Now()
    for range ticker.C {
        now := time.Now()
        dt_s := now.Sub(lastUpdate).Seconds() // Delta time in seconds
        lastUpdate = now
        wm.updateZonesParallel(dt_s)
    }
}

func (wm *WorldManager) updateZonesParallel(dt_s float64) {
    var wg sync.WaitGroup
    jobs := make(chan *types.Zone, len(wm.Zones))

    for i := 0; i < wm.WorkerCount; i++ {
        wg.Add(1)
        go wm.zoneWorker(jobs, dt_s, &wg)
    }

    for _, zone := range wm.Zones {
        jobs <- zone
    }
    close(jobs)

    wg.Wait()
}

func (wm *WorldManager) zoneWorker(jobs <-chan *types.Zone, dt_s float64, wg *sync.WaitGroup) {
    defer wg.Done()

    for zone := range jobs {
        zone.Update(dt_s)
    }
}