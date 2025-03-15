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
    ZoneWidth  = 128
    ZoneHeight = 128
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
            zone := NewZone(zoneID, ZoneWidth, ZoneHeight, x*ZoneWidth, y*ZoneHeight)
            manager.Zones = append(manager.Zones, zone)
            zoneID++
        }
    }

    if len(manager.Zones) == 0 {
        log.Fatal("Error: No active zones initialized!")
    }

    // Load initial Gotchis into Zone 0
    manager.loadInitialGotchis(5)

    log.Printf("World initialized with %d active zones.", len(manager.Zones))
    return manager
}

func (wm *WorldManager) loadInitialGotchis(batchSize int) {
    log.Println("Loading Gotchis from subgraph...")

    gotchiData := storage.GetLatestDatabaseGotchiEntities(batchSize)
    if len(gotchiData) == 0 {
        log.Println("No Gotchis loaded.")
        return
    }

    if len(wm.Zones) == 0 {
        log.Fatal("Error: No zones available to add Gotchis.")
    }

    for _, g := range gotchiData {
        x := rand.Intn(ZoneWidth) + wm.Zones[0].X
        y := rand.Intn(ZoneHeight) + wm.Zones[0].Y
        gotchi := entities.NewGotchi(0, x, y, g)
        wm.Zones[0].AddEntity(gotchi)
    }

    log.Printf("%d Gotchis placed in Zone 0.\n", len(gotchiData))
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