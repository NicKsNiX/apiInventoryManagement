package controllers

import (
	"inventory-management/database"
	"inventory-management/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetInventoryInformation(c *fiber.Ctx) error {
	// รับค่า query string ชื่อ "search"
	search := c.Query("search", "") // ถ้าไม่มีจะได้ค่าเป็น ""

	// ส่ง search เข้าไปที่ model
	items, err := models.GetInventorySummary(database.DB, search)
	if err != nil {
		log.Printf("GetInventorySummary error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load inventory summary",
			"detail": err.Error(),
		})
	}
	return c.JSON(items)
}

func GetInventoryInformationReconfirm(c *fiber.Ctx) error {
	employee := c.Query("employee", "") // optional
	items, err := models.GetInventoryRowsAll(database.DB, employee)
	if err != nil {
		log.Printf("GetInventorySummaryAll error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load inventory summary",
			"detail": err.Error(),
		})
	}
	return c.JSON(items)
}

func GetInventoryInformationHisReconfirm(c *fiber.Ctx) error {
	employee := c.Query("employee", "") // optional
	items, err := models.GetInventoryHistory(database.DB, employee)
	if err != nil {
		log.Printf("GetInventoryHistory error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load inventory summary",
			"detail": err.Error(),
		})
	}
	return c.JSON(items)
}

func GetInventoryReconfirmCount(c *fiber.Ctx) error {
	employee := c.Query("employee", "") // optional
	count, err := models.GetInventoryReconfirmCount(database.DB, employee)
	if err != nil {
		log.Printf("GetInventoryReconfirmCount error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "failed to load inventory summary",
			"detail": err.Error(),
		})
	}
	return c.JSON(count)
}
