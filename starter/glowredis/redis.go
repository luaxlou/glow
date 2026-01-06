package glowredis

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/luaxlou/glow/starter/glowapp"
	"github.com/luaxlou/glow/starter/glowapp/config"
	"github.com/redis/go-redis/v9"
)

var (
	client      *redis.Client
	initialized bool
	mu          sync.RWMutex
)

// Client returns the automatically initialized Redis client.
func Client() (*redis.Client, error) {
	mu.RLock()
	if initialized && client != nil {
		defer mu.RUnlock()
		return client, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if initialized && client != nil {
		return client, nil
	}

	appName := config.AppIdentity
	if appName == "" {
		return nil, fmt.Errorf("app identity not set. call sdk.Init() first")
	}

	// Convention: Use appName + "_cache" or similar for Redis?
	// For now, let's assume we request a generic "redis" resource or "cache"
	// But let's stick to the resource type "redis" and name = appName + "_redis"
	resourceName := appName + "_redis"

	log.Printf("Lazy initializing Redis Starter for %s...", appName)

	// Provision/Get Config
	cfg, err := config.ProvisionResource("redis", resourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to provision redis: %w", err)
	}

	// Extract Redis Config
	// Expecting structure: {"redis": {"addr": "...", "username": "...", "password": "...", "db": 0}}
	var rCfg struct {
		Addr     string `json:"addr"`
		Username string `json:"username"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}

	// Manual map decoding (viper/json generic map hell)
	if redisMap, ok := cfg["redis"].(map[string]any); ok {
		if v, ok := redisMap["addr"].(string); ok {
			rCfg.Addr = v
		}
		if v, ok := redisMap["username"].(string); ok {
			rCfg.Username = v
		}
		if v, ok := redisMap["password"].(string); ok {
			rCfg.Password = v
		}
		if v, ok := redisMap["db"].(float64); ok {
			rCfg.DB = int(v)
		} else if v, ok := redisMap["db"].(int); ok {
			rCfg.DB = v
		}
	}

	if rCfg.Addr == "" {
		// Fallback: maybe the config is flattened?
		// For now, fail if not found
		return nil, fmt.Errorf("redis address not found in config")
	}

	c := redis.NewClient(&redis.Options{
		Addr:     rCfg.Addr,
		Username: rCfg.Username,
		Password: rCfg.Password,
		DB:       rCfg.DB,
	})

	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	client = c
	initialized = true
	log.Println("Redis Starter initialized successfully.")

	glowapp.RegisterCleanup("Redis Starter", func() {
		if client != nil {
			log.Println("Closing Redis client...")
			client.Close()
		}
	})

	return client, nil
}

func Reload() {
	mu.Lock()
	defer mu.Unlock()
	if client != nil {
		client.Close()
		client = nil
	}
	initialized = false
}
