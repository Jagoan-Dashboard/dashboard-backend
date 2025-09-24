
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

    reportRoutes.Get("/tata-bangunan/overview", cont.ReportHandler.GetTataBangunanOverview)
    reportRoutes.Get("/tata-bangunan/basic-statistics", cont.ReportHandler.GetBasicStatistics)
    reportRoutes.Get("/tata-bangunan/location-distribution", cont.ReportHandler.GetLocationDistribution)
    reportRoutes.Get("/tata-bangunan/work-type-statistics", cont.ReportHandler.GetWorkTypeStatistics)
    reportRoutes.Get("/tata-bangunan/condition-statistics", cont.ReportHandler.GetConditionAfterRehabStatistics)
    reportRoutes.Get("/tata-bangunan/status-statistics", cont.ReportHandler.GetStatusStatistics)
    reportRoutes.Get("/tata-bangunan/building-type-distribution", cont.ReportHandler.GetBuildingTypeDistribution)

    
    spatialRoutes := protected.Group("/spatial-planning")
    spatialRoutes.Get("/", cont.SpatialPlanningHandler.ListReports)
    spatialRoutes.Get("/statistics", cont.SpatialPlanningHandler.GetStatistics)
    spatialRoutes.Get("/:id", cont.SpatialPlanningHandler.GetReport)
    spatialRoutes.Post("/", cont.SpatialPlanningHandler.CreateReport)
    spatialRoutes.Put("/:id", cont.SpatialPlanningHandler.UpdateReport)
    spatialRoutes.Delete("/:id", cont.SpatialPlanningHandler.DeleteReport)

    spatialRoutes.Get("/tata-ruang/overview", cont.SpatialPlanningHandler.GetTataRuangOverview)
    spatialRoutes.Get("/tata-ruang/basic-statistics", cont.SpatialPlanningHandler.GetTataRuangBasicStatistics)
    spatialRoutes.Get("/tata-ruang/location-distribution", cont.SpatialPlanningHandler.GetTataRuangLocationDistribution)
    spatialRoutes.Get("/tata-ruang/urgency-statistics", cont.SpatialPlanningHandler.GetUrgencyLevelStatistics)
    spatialRoutes.Get("/tata-ruang/violation-type-statistics", cont.SpatialPlanningHandler.GetViolationTypeStatistics)
    spatialRoutes.Get("/tata-ruang/violation-level-statistics", cont.SpatialPlanningHandler.GetViolationLevelStatistics)
    spatialRoutes.Get("/tata-ruang/area-category-distribution", cont.SpatialPlanningHandler.GetAreaCategoryDistribution)
    spatialRoutes.Get("/tata-ruang/environmental-impact-statistics", cont.SpatialPlanningHandler.GetEnvironmentalImpactStatistics)

    waterRoutes := protected.Group("/water-resources")
    waterRoutes.Get("/", cont.WaterResourcesHandler.ListReports)
    waterRoutes.Get("/overview", cont.WaterResourcesHandler.GetWaterResourcesOverview)
    waterRoutes.Get("/priority", cont.WaterResourcesHandler.ListByPriority)
    waterRoutes.Get("/statistics", cont.WaterResourcesHandler.GetStatistics)
    waterRoutes.Get("/urgent", cont.WaterResourcesHandler.GetUrgentReports)
    waterRoutes.Get("/damage-by-area", cont.WaterResourcesHandler.GetDamageByArea)
    waterRoutes.Get("/:id", cont.WaterResourcesHandler.GetReport)
    waterRoutes.Post("/", cont.WaterResourcesHandler.CreateReport)
    waterRoutes.Put("/:id", cont.WaterResourcesHandler.UpdateReport)
    waterRoutes.Delete("/:id", cont.WaterResourcesHandler.DeleteReport)
    waterRoutes.Get("/dashboard", cont.WaterResourcesHandler.GetDashboard)

    binaMargaRoutes := protected.Group("/bina-marga")
    binaMargaRoutes.Get("/", cont.BinaMargaHandler.ListReports)
    binaMargaRoutes.Get("/overview", cont.BinaMargaHandler.GetBinaMargaOverview)
    binaMargaRoutes.Get("/priority", cont.BinaMargaHandler.ListByPriority)
    binaMargaRoutes.Get("/statistics", cont.BinaMargaHandler.GetStatistics)
    binaMargaRoutes.Get("/emergency", cont.BinaMargaHandler.GetEmergencyReports)
    binaMargaRoutes.Get("/blocked", cont.BinaMargaHandler.GetBlockedRoads)
    binaMargaRoutes.Get("/damage-by-road-type", cont.BinaMargaHandler.GetDamageByRoadType)
    binaMargaRoutes.Get("/damage-by-location", cont.BinaMargaHandler.GetDamageByLocation)
    binaMargaRoutes.Get("/:id", cont.BinaMargaHandler.GetReport)
    binaMargaRoutes.Post("/", cont.BinaMargaHandler.CreateReport)
    binaMargaRoutes.Put("/:id", cont.BinaMargaHandler.UpdateReport)
    binaMargaRoutes.Delete("/:id", cont.BinaMargaHandler.DeleteReport)
    binaMargaRoutes.Get("/dashboard", cont.BinaMargaHandler.GetDashboard)

    agricultureRoutes := protected.Group("/agriculture")
    
    agricultureRoutes.Get("/executive/dashboard", cont.AgricultureHandler.GetExecutiveDashboard)
    agricultureRoutes.Get("/commodity/analysis", cont.AgricultureHandler.GetCommodityAnalysis)
    agricultureRoutes.Get("/food-crop/stats", cont.AgricultureHandler.GetFoodCropStats)
    agricultureRoutes.Get("/horticulture/stats", cont.AgricultureHandler.GetHorticultureStats)
    agricultureRoutes.Get("/plantation/stats", cont.AgricultureHandler.GetPlantationStats)
    agricultureRoutes.Get("/equipment/stats", cont.AgricultureHandler.GetAgriculturalEquipmentStats)
    agricultureRoutes.Get("/land-irrigation/stats", cont.AgricultureHandler.GetLandAndIrrigationStats)

    agricultureRoutes.Get("/", cont.AgricultureHandler.ListReports)
    agricultureRoutes.Get("/:id", cont.AgricultureHandler.GetReport)
    agricultureRoutes.Post("/", cont.AgricultureHandler.CreateReport)
    agricultureRoutes.Put("/:id", cont.AgricultureHandler.UpdateReport)
    agricultureRoutes.Delete("/:id", cont.AgricultureHandler.DeleteReport)
    
    agricultureRoutes.Get("/statistics/overview", cont.AgricultureHandler.GetStatistics)
    agricultureRoutes.Get("/statistics/commodity-production", cont.AgricultureHandler.GetCommodityProduction)
    agricultureRoutes.Get("/statistics/technology-adoption", cont.AgricultureHandler.GetTechnologyAdoptionStats)
    agricultureRoutes.Get("/statistics/farmer-needs", cont.AgricultureHandler.GetFarmerNeedsAnalysis)

    agricultureRoutes.Get("/extension-officer/:officer/reports", cont.AgricultureHandler.GetByExtensionOfficer)
    agricultureRoutes.Get("/extension-officer/performance", cont.AgricultureHandler.GetExtensionOfficerPerformance)

    agricultureRoutes.Get("/village/:village/reports", cont.AgricultureHandler.GetReportsByVillage)
    agricultureRoutes.Get("/reports/by-date-range", cont.AgricultureHandler.GetReportsByDateRange)

    agricultureRoutes.Get("/pest-disease/reports", cont.AgricultureHandler.GetPestDiseaseReports)

    agricultureRoutes.Get("/dashboard/summary", cont.AgricultureHandler.GetDashboardSummary)

    adminWaterRoutes := protected.Group("/admin/water-resources", middleware.RequireRole("ADMIN"))
    adminWaterRoutes.Put("/:id/status", cont.WaterResourcesHandler.UpdateStatus)

    adminRoutes := protected.Group("/admin", middleware.RequireRole("ADMIN"))
           adminRoutes.Get("/users", func(c *fiber.Ctx) error {
            return c.JSON(fiber.Map{"message": "Admin users list"})
        })

        adminSpatialRoutes := protected.Group("/admin/spatial-planning", middleware.RequireRole("ADMIN"))
        adminSpatialRoutes.Put("/:id/status", cont.SpatialPlanningHandler.UpdateStatus)

        
        adminBinaMargaRoutes := adminRoutes.Group("/bina-marga")
        adminBinaMargaRoutes.Put("/:id/status", cont.BinaMargaHandler.UpdateStatus)

        
        dashboardRoutes := protected.Group("/dashboard")
        dashboardRoutes.Get("/overview", func(c *fiber.Ctx) error {
            return c.JSON(fiber.Map{
                "message": "Dashboard overview endpoint - implement consolidated statistics",
            })
        })
        dashboardRoutes.Get("/urgent-reports", func(c *fiber.Ctx) error {
            return c.JSON(fiber.Map{
                "message": "All urgent reports across all modules",
            })
        })

        
    adminAgricultureRoutes := adminRoutes.Group("/agriculture")
    adminAgricultureRoutes.Get("/reports/all", func(c *fiber.Ctx) error {
        
        return c.JSON(fiber.Map{"message": "Admin view all agriculture reports"})
    })
    adminAgricultureRoutes.Get("/analytics/advanced", func(c *fiber.Ctx) error {
        
        return c.JSON(fiber.Map{"message": "Advanced agriculture analytics"})
    })
}