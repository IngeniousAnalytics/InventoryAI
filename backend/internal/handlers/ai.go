package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/inventory_ai/backend/internal/queue"
)

type AIHandler struct {
	Publisher *queue.Publisher
}

func NewAIHandler(pub *queue.Publisher) *AIHandler {
	return &AIHandler{Publisher: pub}
}

func (h *AIHandler) QueueImageAnalysis(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// In a real app, you'd upload the file to S3/Local first.
	// Here we mock the URL or accept base64 in body for simplicity,
	// OR (better) we say "Image uploaded to /tmp/..."

	// For this MVP, let's assume client sends a URL or we just mock an ID
	// Real world: Multipart upload -> Save -> Get Path -> Queue Path

	job := queue.AIJob{
		ImageID:   "img_" + time.Now().Format("20060102150405"),
		UserID:    userID,
		ImageURL:  "https://example.com/sample_inventory.jpg", // Placeholder
		Timestamp: time.Now().Unix(),
	}

	if err := h.Publisher.PublishJob(job); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to queue job"})
	}

	return c.JSON(fiber.Map{
		"message": "Image queued for processing",
		"job_id":  job.ImageID,
	})
}
