package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	Client *redis.Client
	Limit  int
	Window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Client: client,
		Limit:  limit,
		Window: window,
	}
}

func (rl *RateLimiter) LimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Identify by Token (User) or IP
		key := c.Locals("user_id")
		identifier := ""
		if key != nil {
			identifier = fmt.Sprintf("rate_limit:user:%v", key)
		} else {
			identifier = fmt.Sprintf("rate_limit:ip:%s", c.IP())
		}

		ctx := context.Background()

		// Simple Fixed Window Counter
		// Increment request count
		count, err := rl.Client.Incr(ctx, identifier).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Rate limit error"})
		}

		// Set expiration on first request
		if count == 1 {
			rl.Client.Expire(ctx, identifier, rl.Window)
		}

		if count > int64(rl.Limit) {
			ttl, _ := rl.Client.TTL(ctx, identifier).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":               "Rate limit exceeded",
				"retry_after_seconds": ttl.Seconds(),
			})
		}

		return c.Next()
	}
}
