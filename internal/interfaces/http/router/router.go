// internal/interfaces/http/router/router.go
package router

import (
    "building-report-backend/internal/interfaces/http/middleware"
    "building-report-backend/pkg/container"
    
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cont *container.Container) {
    // API v1 group
    api := app.Group("/api/v1")

    // Health check
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "message": "Service is running",
        })
    })

    // Auth routes (public)
    authRoutes := api.Group("/auth")
    authRoutes.Post("/register", cont.AuthHandler.Register)
    authRoutes.Post("/login", cont.AuthHandler.Login)

    // Protected routes
    protected := api.Use(middleware.AuthMiddleware(cont.AuthService))
    
    // User routes
    protected.Get("/profile", cont.AuthHandler.GetProfile)

    // Report routes
    reportRoutes := protected.Group("/reports")
    reportRoutes.Get("/", cont.ReportHandler.ListReports)
    reportRoutes.Get("/:id", cont.ReportHandler.GetReport)
    reportRoutes.Post("/", cont.ReportHandler.CreateReport)
    reportRoutes.Put("/:id", cont.ReportHandler.UpdateReport)
    reportRoutes.Delete("/:id", cont.ReportHandler.DeleteReport)

    // Admin only routes
    adminRoutes := protected.Group("/admin", middleware.RequireRole("ADMIN"))
    adminRoutes.Get("/users", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Admin users list"})
    })
}