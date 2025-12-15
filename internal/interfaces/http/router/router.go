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
    users := api.Group("/users")
    users.Use(middleware.AuthMiddleware(cont.AuthService))
    
    users.Get("/", cont.AuthHandler.GetAllUsers)
    users.Get("/:id", cont.AuthHandler.GetUserByID)
    users.Post("/", cont.AuthHandler.CreateUser)
    users.Put("/:id", cont.AuthHandler.UpdateUser)
    users.Delete("/:id", cont.AuthHandler.DeleteUser)

    api.Post("/reports", cont.ReportHandler.CreateReport)
    api.Post("/spatial-planning", cont.SpatialPlanningHandler.CreateReport)
    api.Post("/water-resources", cont.WaterResourcesHandler.CreateReport)
    api.Post("/bina-marga", cont.BinaMargaHandler.CreateReport)
    api.Post("/agriculture", cont.AgricultureHandler.CreateReport)

    
    //protected := api.Use(middleware.AuthMiddleware(cont.AuthService))
    
    
    api.Get("/profile", cont.AuthHandler.GetProfile)

    reportRoutes := api.Group("/reports")
    reportRoutes.Get("/", cont.ReportHandler.ListReports)
    reportRoutes.Get("/:id", cont.ReportHandler.GetReport)
    reportRoutes.Put("/:id", cont.ReportHandler.UpdateReport)
    reportRoutes.Delete("/:id", cont.ReportHandler.DeleteReport)

    reportRoutes.Get("/tata-bangunan/overview", cont.ReportHandler.GetTataBangunanOverview)
    
    spatialRoutes := api.Group("/spatial-planning")
    spatialRoutes.Get("/", cont.SpatialPlanningHandler.ListReports)
    spatialRoutes.Get("/statistics", cont.SpatialPlanningHandler.GetStatistics)
    spatialRoutes.Get("/:id", cont.SpatialPlanningHandler.GetReport)
    spatialRoutes.Put("/:id", cont.SpatialPlanningHandler.UpdateReport)
    spatialRoutes.Delete("/:id", cont.SpatialPlanningHandler.DeleteReport)

    spatialRoutes.Get("/tata-ruang/overview", cont.SpatialPlanningHandler.GetTataRuangOverview)

    waterRoutes := api.Group("/water-resources")
    waterRoutes.Get("/", cont.WaterResourcesHandler.ListReports)
    waterRoutes.Get("/overview", cont.WaterResourcesHandler.GetWaterResourcesOverview)

    binaMargaRoutes := api.Group("/bina-marga")
    binaMargaRoutes.Get("/", cont.BinaMargaHandler.ListReports)
    binaMargaRoutes.Get("/overview", cont.BinaMargaHandler.GetBinaMargaOverview)

    agricultureRoutes := api.Group("/agriculture")
    agricultureRoutes.Post("/komoditas/import", cont.AgricultureHandler.ImportKomoditas)
    agricultureRoutes.Post("/alat-pertanian/import", cont.AgricultureHandler.ImportAlatPertanian)

    agricultureRoutes.Get("/komoditas/export", cont.AgricultureHandler.ExportKomoditas)
    agricultureRoutes.Get("/alat-pertanian/export", cont.AgricultureHandler.ExportAlatPertanian)
    agricultureRoutes.Get("/executive/dashboard", cont.AgricultureHandler.GetExecutiveDashboard)
    agricultureRoutes.Get("/commodity/analysis", cont.AgricultureHandler.GetCommodityAnalysis)
    agricultureRoutes.Get("/food-crop/stats", cont.AgricultureHandler.GetFoodCropStats)
    agricultureRoutes.Get("/horticulture/stats", cont.AgricultureHandler.GetHorticultureStats)
    agricultureRoutes.Get("/plantation/stats", cont.AgricultureHandler.GetPlantationStats)
    agricultureRoutes.Get("/equipment/stats", cont.AgricultureHandler.GetAgriculturalEquipmentStats)
    agricultureRoutes.Get("/land-irrigation/stats", cont.AgricultureHandler.GetLandAndIrrigationStats)

    agricultureRoutes.Get("/", cont.AgricultureHandler.ListReports)
    
    agricultureRoutes.Get("/:id", cont.AgricultureHandler.GetReport)
    agricultureRoutes.Put("/:id", cont.AgricultureHandler.UpdateReport)
    agricultureRoutes.Delete("/:id", cont.AgricultureHandler.DeleteReport)

    riceFieldRoutes := api.Group("/rice-fields")
	riceFieldRoutes.Post("/lahan-pengairan/import", cont.RiceFieldHandler.ImportRiceFields)
    riceFieldRoutes.Get("/lahan-engairan/export", cont.RiceFieldHandler.ExportRiceFields)
    riceFieldRoutes.Post("/", cont.RiceFieldHandler.Create)
	riceFieldRoutes.Get("/", cont.RiceFieldHandler.GetAll)
	riceFieldRoutes.Get("/:id", cont.RiceFieldHandler.GetByID)
	riceFieldRoutes.Put("/:id", cont.RiceFieldHandler.Update)
	riceFieldRoutes.Delete("/:id", cont.RiceFieldHandler.Delete)
	
	riceFieldRoutes.Get("/stats", cont.RiceFieldHandler.GetStatistics)
	riceFieldRoutes.Get("/analysis", cont.RiceFieldHandler.GetAnalysis)

    executiveRoutes := api.Group("/executive")
    economyRoutes := executiveRoutes.Group("/economy")
    economyRoutes.Get("/overview", cont.ExecutiveHandler.GetEkonomiOverview)

    populationRoutes := executiveRoutes.Group("/population")
    populationRoutes.Get("/overview", cont.ExecutiveHandler.GetPopulationOverview)

    povertyRoutes := executiveRoutes.Group("/poverty")
    povertyRoutes.Get("/overview", cont.ExecutiveHandler.GetPovertyOverview)

    employmentRoutes := executiveRoutes.Group("/employment")
    employmentRoutes.Get("/overview", cont.ExecutiveHandler.GetEmploymentOverview)

    educationRoutes := executiveRoutes.Group("/education")
    educationRoutes.Get("/overview", cont.ExecutiveHandler.GetEducationOverview)
}