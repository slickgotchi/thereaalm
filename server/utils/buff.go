// thereaalm/utils/buff.go
package utils

import (
	"thereaalm/interfaces"
)

func UpdateBuffMultiplier(consumer interfaces.IBuffConsumer, maxRange int) {
    zone := consumer.GetZone()
    if zone == nil {
        consumer.SetBuffMultiplier(1.0)
        return
    }
    x, y := consumer.GetPosition()
    nearbyEntities := zone.FindNearbyEntities(x, y, maxRange) // Max Altar range
    maxBuff := 1.0
    for _, entity := range nearbyEntities {
        if buffProvider, ok := entity.(interfaces.IBuffProvider); ok {
            if buffProvider.IsBuffActive() {
				ex, ey := entity.GetPosition()
                dist := zone.GetDistance(x, y, ex, ey)
                if dist <= buffProvider.GetBuffRange() {
                    buff := buffProvider.GetSpeedBuffMultiplier()
                    if buff > maxBuff {
                        maxBuff *= buff // multiply buffs
                    }
                }
            }
        }
    }
    consumer.SetBuffMultiplier(maxBuff)
}