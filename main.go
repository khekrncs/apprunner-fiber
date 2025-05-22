package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/khekrn/apprunner-fiber/handlers"
	"github.com/khekrn/apprunner-fiber/services"
)

const (
	port = "8080"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize services
	s3Service := services.NewS3Service()
	userService := services.NewUserService(s3Service)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	fileHandler := handlers.NewFileHandler(s3Service)

	// Your existing routes (unchanged)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "I'm running!"})
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	api := app.Group("/api/v1")

	// User management routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
	users.Get("/", userHandler.ListUsers)

	// File management routes
	files := api.Group("/files")
	files.Post("/upload/:userId", fileHandler.UploadFile)
	files.Get("/:userId/:filename", fileHandler.GetFile)
	files.Delete("/:userId/:filename", fileHandler.DeleteFile)
	files.Get("/:userId", fileHandler.ListUserFiles)

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
