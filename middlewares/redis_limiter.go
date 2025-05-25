package middlewares

import (
	"fmt"

	"tawtheeq-backend/config"
	"tawtheeq-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func RedisRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !config.RateLimitEnabled {
			return c.Next()
		}

		ip := c.IP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := config.Redis.Incr(config.Ctx, key).Result()
		if err != nil {
			utils.HandleError(err, "Rate limiter failed", utils.Error)
			return c.Status(500).JSON(fiber.Map{"error": "Rate limiter failed"})
		}

		if count == 1 {
			config.Redis.Expire(config.Ctx, key, config.RateLimitWindow)
		}

		if int(count) > config.RateLimitMax {
			ttl, _ := config.Redis.TTL(config.Ctx, key).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":        "Too many requests",
				"try_again_in": ttl.Seconds(),
			})
		}

		return c.Next()
	}
}
