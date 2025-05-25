package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Redis            *redis.Client
	Ctx              = context.Background()
	RateLimitEnabled bool
	RateLimitMax     int
	RateLimitWindow  time.Duration
)

func InitRateLimiting() {

	RateLimitEnabled = os.Getenv("RATE_LIMIT_ENABLED") == "true"
	if !RateLimitEnabled {
		fmt.Println("⚠️  Rate Limiting is DISABLED")
		return
	}

	maxStr := os.Getenv("RATE_LIMIT")
	max, err := strconv.Atoi(maxStr)
	if err != nil || max <= 0 {
		max = 100
	}
	RateLimitMax = max

	winStr := os.Getenv("RATE_LIMIT_WINDOW")
	sec, err := strconv.Atoi(winStr)
	if err != nil || sec <= 0 {
		sec = 60
	}
	RateLimitWindow = time.Duration(sec) * time.Second

	Redis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	if _, err := Redis.Ping(Ctx).Result(); err != nil {
		panic("❌ Failed to connect to Redis: " + err.Error())
	}

	fmt.Printf("✅ Rate Limiting ENABLED - %d req / %ds\n", RateLimitMax, int(RateLimitWindow.Seconds()))
}
