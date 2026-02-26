package main

import (
	"fmt"
	"inventory-management/controllers"
	"inventory-management/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = database.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Database connected successfully!")
	app.Post("/api/login", controllers.Login)
	app.Get("/api/getMenuinventory", controllers.GetMenuinventory)
	app.Get("/api/getInventoryInfo", controllers.GetInventoryInfo)
	app.Get("/api/inventoryInformation", controllers.GetInventoryInformation)
	app.Get("/api/inventoryInformationReconfirm", controllers.GetInventoryInformationReconfirm)
	app.Get("/api/inventoryInformationHisReconfirm", controllers.GetInventoryInformationHisReconfirm)
	app.Post("/api/updateQtyInventory", controllers.UpdateQtyInventory)
	app.Post("/api/updateReconfirmQtyInventory", controllers.UpdateReconfirmQtyInventory)
	app.Post("/api/updateNoconfirmQtyInventory", controllers.UpdateNoconfirmQtyInventory)
	app.Post("/api/insertQtyInventory", controllers.InsertQtyInventory)
	app.Post("/api/insertAdjustQtyInventory", controllers.InsertAdjustQtyInventory)
	app.Post("/api/confirmEditQtyAdjust", controllers.ConfirmEditQtyAdjust)
	app.Get("/api/getInventoryCheckList", controllers.GetInventoryCheckList)
	app.Get("/api/showHistoryAuditor", controllers.GetShowHistoryAuditor)
	app.Get("/api/getInventoryReconfirmCount", controllers.GetInventoryReconfirmCount)

	app.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "4001"
	}

	log.Println("Server running at http://0.0.0.0:" + port)
	log.Fatal(app.Listen(":4002"))
}
