// handler/menuHandler.go
package controllers

import (
	"inventory-management/database"
	"inventory-management/models" // นำเข้า model ที่สร้างไว้
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ฟังก์ชันสำหรับจัดการกับ API Request
// handler.go
func GetMenuinventory(c *fiber.Ctx) error {
	// รับค่าจาก query (เช่น ?employee=02035&spg_id=3)
	emp := c.Query("employee")
	spgIDStr := c.Query("spg_id")

	if emp == "" || spgIDStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "missing employee or spg_id"})
	}

	spgID, err := strconv.Atoi(spgIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid spg_id"})
	}

	menuDetails, err := models.GetMenuDetails(c.Context(), spgID, emp)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch menu details"})
	}
	return c.Status(200).JSON(menuDetails)
}

type updateQtyRequest struct {
	MitID    int     `json:"MitId"`
	Qty      float64 `json:"total"` // changed to float64 to accept decimals
	Employee string  `json:"employee"`
	IidId    int     `json:"IidId"`
}

func UpdateQtyInventory(c *fiber.Ctx) error {
	var req updateQtyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.MitID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "mit_id is required"})
	}

	// set a default employee if not provided (optional)
	if req.Employee == "" {
		req.Employee = "SYSTEM"
	}

	if err := models.UpdateItemQuantity(database.DB, req.MitID, req.Qty, req.Employee, req.IidId); err != nil {
		log.Printf("UpdateItemQuantity error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "quantity updated"})
}

type UpdateReconfirmQtyRequest struct {
	Qty       float64 `json:"qty"`
	Employee  string  `json:"employee"`
	ItemCode  string  `json:"itemCode"`
	Warehouse string  `json:"warehouse"`
	IidID     int64   `json:"iid_id"`
}

func UpdateReconfirmQtyInventory(c *fiber.Ctx) error {
	var req UpdateReconfirmQtyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Qty == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "qty is required"})
	}

	if req.Employee == "" {
		req.Employee = "SYSTEM"
	}

	if err := models.UpdateReconfirmQty(database.DB, req.Qty, req.Employee, req.ItemCode, req.Warehouse, req.IidID); err != nil {
		log.Printf("UpdateReconfirmQty error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "quantity updated"})
}

func UpdateNoconfirmQtyInventory(c *fiber.Ctx) error {
	var req UpdateReconfirmQtyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Qty == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "qty is required"})
	}

	// set a default employee if not provided (optional)
	if req.Employee == "" {
		req.Employee = "SYSTEM"
	}

	if err := models.UpdateNoconfirmQty(database.DB, req.Qty, req.Employee, req.ItemCode, req.Warehouse, req.IidID); err != nil {
		log.Printf("UpdateReconfirmQty error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "quantity updated"})
}

func ConfirmEditQtyAdjust(c *fiber.Ctx) error {
	type ConfirmEditQtyAdjustRequest struct {
		Qty       float64 `json:"qty"`
		Employee  string  `json:"employee"`
		ItemCode  string  `json:"itemCode"`
		Warehouse string  `json:"warehouse"`
		IidID     int64   `json:"iid_id"`
	}
	var req ConfirmEditQtyAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := models.ConfirmEditQtyAdjust(database.DB, req.Qty, req.Employee, req.ItemCode, req.Warehouse, req.IidID); err != nil {
		log.Printf("ConfirmEditQtyAdjust error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "quantity updated"})
}
