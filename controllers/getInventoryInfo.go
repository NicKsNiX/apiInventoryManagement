package controllers

import (
	"database/sql"
	"errors"
	"inventory-management/database"
	"inventory-management/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// controllers/getInventoryInfo.go
func GetInventoryInfo(c *fiber.Ctx) error {
	item_cd := c.Query("item_cd")
	if item_cd == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item_cd is required"})
	}
	if len(item_cd) > 25 {
		item_cd = item_cd[25:] // ตัด 25 ตัวแรกออก
	}
	inv, err := models.GetItemDetailsByItem(database.DB, item_cd)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// ✅ 404 ชัดเจน
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "inventory not found",
				"message": "No data found",
				"code":    "NOT_FOUND",
			})
		}
		log.Printf("Error fetching inventory data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "internal server error",
			"message": err.Error(),
			"code":    "INTERNAL_ERROR",
		})
	}

	return c.Status(fiber.StatusOK).JSON(inv)
}

func GetInventoryCheckList(c *fiber.Ctx) error {
	employee := c.Query("employee", "") // optional
	items, err := models.GetInventoryCheckList(database.DB, employee)
	if err != nil {
		log.Printf("GetInventoryCheckList error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load inventory checklist",
			"detail": err.Error(),
		})
	}
	return c.JSON(items)
}

func GetShowHistoryAuditor(c *fiber.Ctx) error {
	employee := c.Query("employee", "") // optional
	items, err := models.GetShowHistoryAuditor(database.DB, employee)
	if err != nil {
		log.Printf("GetShowHistoryAuditor error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load inventory checklist",
			"detail": err.Error(),
		})
	}
	return c.JSON(items)
}
