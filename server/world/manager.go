// thereaalm/world/manager.go
package world

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"thereaalm/config"
	"thereaalm/entity/entities"
	"thereaalm/storage"
	"time"
)

const (
    ZoneTiles  = 256
    // ZoneHeight = 128
)

type WorldManager struct {
    Zones       []*Zone
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
            zone := NewZone(zoneID, ZoneTiles, ZoneTiles, x*ZoneTiles, y*ZoneTiles)
            manager.Zones = append(manager.Zones, zone)
            zoneID++
        }
    }

    if len(manager.Zones) == 0 {
        log.Fatal("Error: No active zones initialized!")
    }

    // Load initial Gotchis into Zone 0
    manager.loadInitialGotchis()

    log.Printf("World initialized with %d active zones.", len(manager.Zones))
    return manager
}

func (wm *WorldManager) loadInitialGotchis() {
    log.Println("Loading Gotchis from subgraph...")

    gotchiData := storage.GetLatestDatabaseGotchiEntities()
    if len(gotchiData) == 0 {
        log.Println("No Gotchis loaded.")
        return
    }

    if len(wm.Zones) == 0 {
        log.Fatal("Error: No zones available to add Gotchis.")
    }

    // place gotchis across all available zones, in ZoneMap
    // start with a known seed
    r := rand.New(rand.NewSource(123))

    for _, g := range gotchiData {
        // pick random zone
        zoneIndex := r.Intn(len(wm.Zones))

        // pick random location in the zone
        x := r.Intn(ZoneTiles) + wm.Zones[zoneIndex].X
        y := r.Intn(ZoneTiles) + wm.Zones[zoneIndex].Y

        // create new gotchi
        gotchi := entities.NewGotchi(0, x, y, g)
        wm.Zones[zoneIndex].AddEntity(gotchi)
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
        wm.updateZonesParallel()
    }
}

func (wm *WorldManager) updateZonesParallel() {
    var wg sync.WaitGroup
    jobs := make(chan *Zone, len(wm.Zones))

    for i := 0; i < wm.WorkerCount; i++ {
        wg.Add(1)
        go wm.zoneWorker(jobs, &wg)
    }

    for _, zone := range wm.Zones {
        jobs <- zone
    }
    close(jobs)

    wg.Wait()
}

func (wm *WorldManager) zoneWorker(jobs <-chan *Zone, wg *sync.WaitGroup) {
    defer wg.Done()

    for zone := range jobs {
        zone.Update()
    }
}