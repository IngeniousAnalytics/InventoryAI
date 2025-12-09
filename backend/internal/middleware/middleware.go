package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func TenantMiddleware(c *fiber.Ctx) error {
	// Look for X-Tenant-ID header
	tenantID := c.Get("X-Tenant-ID")
	if tenantID == "" {
		// In production, we might resolve this from a subdomain or JWT context
		// For now, if missing in public routes, that's fine, but protected routes need it
		// c.Locals("tenant_id", "default")
		return c.Next()
	}

	c.Locals("tenant_id", tenantID)
	return c.Next()
}

// RateLimitMiddleware stub - in real implementation this connects to Redis
func RateLimitMiddleware(c *fiber.Ctx) error {
	// TODO: Implement actual Redis rate limiting logic here
	log.Println("[RateLimit] Checking limit for IP:", c.IP())
	return c.Next()
}
