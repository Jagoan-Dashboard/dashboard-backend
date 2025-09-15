
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    
    "building-report-backend/pkg/config"
    "building-report-backend/pkg/database"
    "building-report-backend/pkg/cache"
    "building-report-backend/pkg/storage"
    "building-report-backend/pkg/container"
    "building-report-backend/internal/interfaces/http/router"
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
    })

    
    app.Use(logger.New())
    app.Use(recover.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: cfg.App.AllowedOrigins,
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
        AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
    }))

    
    router.SetupRoutes(app, cont)

    
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