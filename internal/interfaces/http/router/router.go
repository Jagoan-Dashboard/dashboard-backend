
package router

import (
    "building-report-backend/internal/interfaces/http/middleware"
    "building-report-backend/pkg/container"
    
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cont *container.Container) {
    
    api := app.Group("/api/v1")

    
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "message": "Service is running",
        })
    })

    
    authRoutes := api.Group("/auth")
    authRoutes.Post("/register", cont.AuthHandler.Register)
    authRoutes.Post("/login", cont.AuthHandler.Login)

    
    protected := api.Use(middleware.AuthMiddleware(cont.AuthService))
    
    
    protected.Get("/profile", cont.AuthHandler.GetProfile)

    
    reportRoutes := protected.Group("/reports")
    reportRoutes.Get("/", cont.ReportHandler.ListReports)
    reportRoutes.Get("/:id", cont.ReportHandler.GetReport)
    reportRoutes.Post("/", cont.ReportHandler.CreateReport)
    reportRoutes.Put("/:id", cont.ReportHandler.UpdateReport)
    reportRoutes.Delete("/:id", cont.ReportHandler.DeleteReport)

    
    spatialRoutes := protected.Group("/spatial-planning")
    spatialRoutes.Get("/", cont.SpatialPlanningHandler.ListReports)
    spatialRoutes.Get("/statistics", cont.SpatialPlanningHandler.GetStatistics)
    spatialRoutes.Get("/:id", cont.SpatialPlanningHandler.GetReport)
    spatialRoutes.Post("/", cont.SpatialPlanningHandler.CreateReport)
    spatialRoutes.Put("/:id", cont.SpatialPlanningHandler.UpdateReport)
    spatialRoutes.Delete("/:id", cont.SpatialPlanningHandler.DeleteReport)

    waterRoutes := protected.Group("/water-resources")
    waterRoutes.Get("/", cont.WaterResourcesHandler.ListReports)
waterRoutes.Get("/priority", cont.WaterResourcesHandler.ListByPriority)
waterRoutes.Get("/statistics", cont.WaterResourcesHandler.GetStatistics)
waterRoutes.Get("/urgent", cont.WaterResourcesHandler.GetUrgentReports)
waterRoutes.Get("/damage-by-area", cont.WaterResourcesHandler.GetDamageByArea)
waterRoutes.Get("/:id", cont.WaterResourcesHandler.GetReport)
waterRoutes.Post("/", cont.WaterResourcesHandler.CreateReport)
waterRoutes.Put("/:id", cont.WaterResourcesHandler.UpdateReport)
waterRoutes.Delete("/:id", cont.WaterResourcesHandler.DeleteReport)


adminWaterRoutes := protected.Group("/admin/water-resources", middleware.RequireRole("ADMIN"))
adminWaterRoutes.Put("/:id/status", cont.WaterResourcesHandler.UpdateStatus)

    
    adminRoutes := protected.Group("/admin", middleware.RequireRole("ADMIN"))
    adminRoutes.Get("/users", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Admin users list"})
    })

    adminSpatialRoutes := protected.Group("/admin/spatial-planning", middleware.RequireRole("ADMIN"))
    adminSpatialRoutes.Put("/:id/status", cont.SpatialPlanningHandler.UpdateStatus)
}