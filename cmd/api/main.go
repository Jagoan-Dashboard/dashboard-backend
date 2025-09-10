// cmd/api/main.go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    
    "your-module/pkg/config"
    "your-module/pkg/database"
    "your-module/pkg/cache"
    "your-module/pkg/storage"
    "your-module/pkg/container"
    "your-module/internal/interfaces/http/router"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Initialize database
    db, err := database.NewPostgresDB(cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize Redis
    redisClient := cache.NewRedisClient(cfg.Redis)

    // Initialize MinIO
    minioClient, err := storage.NewMinioClient(cfg.Minio)
    if err != nil {
        log.Fatal("Failed to connect to MinIO:", err)
    }

    // Initialize dependency injection container
    cont := container.NewContainer(cfg, db, redisClient, minioClient)

    // Initialize Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: customErrorHandler,
    })

    // Middleware
    app.Use(logger.New())
    app.Use(recover.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: cfg.App.AllowedOrigins,
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
        AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
    }))

    // Setup routes
    router.SetupRoutes(app, cont)

    // Start server
    log.Printf("Server starting on port %s", cfg.App.Port)
    if err := app.Listen(":" + cfg.App.Port); err != nil {
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

    return c.Status(code).JSON(fiber.Map{
        "success": false,
        "message": message,
        "error":   err.Error(),
    })
}