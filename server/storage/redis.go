// redis client for live memory
package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitRedis initializes Redis connection.
func InitRedis() {
    RDB = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
}

// SaveGotchiMind saves Gotchi state in Redis.
func SaveGotchiMind(id string, mind string) error {
    return RDB.Set(context.Background(), "gotchi:"+id, mind, 0).Err()
}