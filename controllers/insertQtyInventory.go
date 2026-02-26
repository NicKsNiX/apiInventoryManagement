package controllers

import (
	"inventory-management/database"
	"inventory-management/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type insertQtyRequest struct {
	Total    float64 `json:"total"`
	Employee string  `json:"employee"`
	IidID    int     `json:"IidId"`
}

func InsertQtyInventory(c *fiber.Ctx) error {
	var req insertQtyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.IidID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "IidId is required"})
	}
	if req.Employee == "" {
		req.Employee = "system"
	}

	// channel = 1 for manual adjustment (change if you use other channel codes)
	if err := models.InsertInventoryAdjustment(database.DB, req.IidID, req.Total, req.Employee, 2); err != nil {
		log.Printf("InsertInventoryAdjustment error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "inserted"})
} 

type insertAdjustQtyRequest struct {
	Total    float64 `json:"total"`
	Employee string  `json:"employee"`
	IidID    int     `json:"IidId"`
}

func InsertAdjustQtyInventory(c *fiber.Ctx) error {
	var req insertAdjustQtyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if req.IidID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "IidId is required"})
	}
	if req.Employee == "" {
		req.Employee = "system"
	}

	// channel = 1 for manual adjustment (change if you use other channel codes)
	if err := models.InsertAdjustInventory(database.DB, req.IidID, req.Total, req.Employee, 2); err != nil {
		log.Printf("InsertAdjustInventory error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "inserted"})
}
