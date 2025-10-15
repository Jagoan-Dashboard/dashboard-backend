package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"building-report-backend/internal/interfaces/http/router"
	"building-report-backend/pkg/cache"
	"building-report-backend/pkg/config"
	"building-report-backend/pkg/container"
	"building-report-backend/pkg/database"
	"building-report-backend/pkg/storage"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	redisClient := cache.NewRedisClient(cfg.Redis)
	minioClient, err := storage.NewMinioClient(cfg.Minio)
	if err != nil {
		log.Fatal("Failed to connect to MinIO:", err)
	}

	cont := container.NewContainer(cfg, db, redisClient, minioClient)

	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		BodyLimit:    10 * 1024 * 1024, // 10 MB
	})

	// 1. Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} | ${latency}\n",
	}))

	// 2. Recover middleware
	app.Use(recover.New())

	// 3. Debug middleware (optional - bisa dihapus di production)
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin != "" {
			log.Printf("📍 Request from: %s | Method: %s | Path: %s", origin, c.Method(), c.Path())
		}
		return c.Next()
	})

	// 4. CORS middleware - HARUS SEBELUM ROUTES
	origins := os.Getenv("APP_ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000"
	}

	// Trim spaces dari setiap origin (safety measure)
	originList := strings.Split(origins, ",")
	for i, origin := range originList {
		originList[i] = strings.TrimSpace(origin)
	}
	cleanedOrigins := strings.Join(originList, ",")

	log.Printf("🌐 CORS Allowed Origins: %s", cleanedOrigins)

	app.Use(cors.New(cors.Config{
    AllowOrigins:     "*",
    AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
    AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Cache-Control,Pragma,Expires,X-Requested-With",
    ExposeHeaders:    "Content-Length,Content-Type,Authorization", 
    AllowCredentials: false,
    MaxAge:           86400,
}))

	// 5. Handle preflight requests
	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	// 6. Setup routes - SETELAH CORS
	router.SetupRoutes(app, cont)

	// 7. 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route not found",
		})
	})

	// Start server
	port := cfg.App.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📝 Environment: %s", os.Getenv("APP_ENV"))
	
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Log error untuk debugging
	log.Printf("❌ Error: %v | Path: %s | Method: %s", err, c.Path(), c.Method())

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}