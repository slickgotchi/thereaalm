// thereaalm/world/manager.go
package world

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"thereaalm/action"
	"thereaalm/config"
	"thereaalm/entity"
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
    // log.Println("Loading Gotchis from subgraph...")

    // OUR TEST SCENARIO
    // - create gotchi, shop and berrybush
    // - gotchi gathers berrys from berrybush for 10 seconds
    // - gotchi sells berrys to shop for 5 seconds
    // - gotchi gathers berrys from berrybush for 10 seconds
    // - gotchi sells berrys to shop for 5 seconds
    // - berry bush out of berries
    // - gotchi idles for 20 seconds
    // - berry bush replenished with berries
    // - repeat

    // grab a few gotchis to use
    gotchisMap, err := web3.FetchGotchisByIDs([]string{"4285", "19005", "21550"})
    if err != nil {
        fmt.Printf("Error fetching gotchis: %v\n", err)
        return
    }

    zoneX := wm.Zones[42].X
    zoneY := wm.Zones[42].Y

    // ENTITIES
    bush := entity.NewBerryBush(42, 12+zoneX, 10+zoneY)
    wm.Zones[42].AddEntity(bush)

    // obstacles
    // wm.Zones[42].AddEntity(entity.NewBerryBush(42, 9+zoneX, 8+zoneY))
    // wm.Zones[42].AddEntity(entity.NewBerryBush(42, 9+zoneX, 9+zoneY))
    wm.Zones[42].AddEntity(entity.NewBerryBush(42, 9+zoneX, 10+zoneY))
    // wm.Zones[42].AddEntity(entity.NewBerryBush(42, 9+zoneX, 11+zoneY))
    // wm.Zones[42].AddEntity(entity.NewBerryBush(42, 9+zoneX, 12+zoneY))

    // shop
    shop := entity.NewShop(42, 6+zoneX, 12+zoneY)
    wm.Zones[42].AddEntity(shop)

    gotchiA := entity.NewGotchi(42, 18+zoneX, 14+zoneY, gotchisMap["4285"])
    wm.Zones[42].AddEntity(gotchiA)

    // gotchiB := entity.NewGotchi(42, 20+zoneX, 10+zoneY, gotchisMap["19005"])
    // wm.Zones[42].AddEntity(gotchiB)

    // gotchiC := entity.NewGotchi(42, 14+zoneX, 19+zoneY, gotchisMap["21550"])
    // wm.Zones[42].AddEntity(gotchiC)

    // lickquidator := entity.NewLickquidator(42, 9+zoneX, 14+zoneY)
    // wm.Zones[42].AddEntity(lickquidator)

    altar := entity.NewAltar(42, 15+zoneX, 12+zoneY)
    wm.Zones[42].AddEntity(altar)

    // ACTIONS
    gotchiA.AddAction(action.NewHarvestAction(gotchiA, bush, 0.5))
    gotchiA.AddAction(action.NewTradeAction(gotchiA, shop, 0.5, "SellAllForGold"))   // FUTURE: we pass a TradeOffer rather than "SellAllForGold"
    // gotchiA.AddAction(action.NewAttackAction(gotchiA, lickquidator, 0.3))
    // gotchiA.AddAction(action.NewRoamAction(gotchiA, 0.1))

    // gotchiB.AddAction(action.NewRoamAction(gotchiB, 0.1))

    // gotchiC.AddAction(action.NewRoamAction(gotchiC, 0.1))

    // lickquidator.AddAction(action.NewAttackAction(lickquidator, gotchiA, 0.2))
    // lickquidator.AddAction(action.NewAttackAction(lickquidator, altar, 0.8))









    // gotchiData := storage.GetLatestDatabaseGotchiEntities(1)
    // if len(gotchiData) == 0 {
    //     log.Println("No Gotchis loaded.")
    //     return
    // }

    // // TEMPORARY: start only in zone 42 for now
    // zoneIndex := 42

    // // zoneX := wm.Zones[zoneIndex].X
    // // zoneY := wm.Zones[zoneIndex].Y

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
    //     gotchi := entities.NewGotchi(42, x, y)
    //     wm.Zones[zoneIndex].AddEntity(gotchi)

    //     // Set initial action sequence
    //     actionSequence := []action.IAction{
    //         action.NewGatherAction(bush.UUID), // Gather (berries: 5 -> 3)
    //         action.NewSellAction(shop.UUID),   // Sell (inventory: 2 -> 0)
    //         action.NewGatherAction(bush.UUID), // Gather (berries: 3 -> 1)
    //         action.NewSellAction(shop.UUID),   // Sell (inventory: 2 -> 0)
    //         action.NewGatherAction(bush.UUID), // Gather (berries: 1 -> 0)
    //         action.NewSellAction(shop.UUID),   // Sell (inventory: 1 -> 0)
    //         action.NewIdleAction(),            // Idle (bush depleted, nothing to do)
    //     }
    //     gotchi.SetActionSequence(actionSequence)
    // }
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