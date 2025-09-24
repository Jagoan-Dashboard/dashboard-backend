package handler

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"

	"github.com/gofiber/fiber/v2"
	
)

type AgricultureHandler struct {
    agricultureUseCase *usecase.AgricultureUseCase
}

func NewAgricultureHandler(agricultureUseCase *usecase.AgricultureUseCase) *AgricultureHandler {
    return &AgricultureHandler{
        agricultureUseCase: agricultureUseCase,
    }
}

func (h *AgricultureHandler) CreateReport(c *fiber.Ctx) error {
    var req dto.CreateAgricultureRequest
    
    // Parse basic information
    req.ExtensionOfficer = c.FormValue("extension_officer")
    req.FarmerName = c.FormValue("farmer_name")
    req.FarmerGroup = c.FormValue("farmer_group")
    req.FarmerGroupType = c.FormValue("farmer_group_type")
    req.Village = c.FormValue("village")
    req.District = c.FormValue("district")
    
    // Parse visit date
    visitDateStr := c.FormValue("visit_date")
    if visitDate, err := time.Parse("2006-01-02", visitDateStr); err == nil {
        req.VisitDate = visitDate
    } else {
        return response.BadRequest(c, "Invalid visit date format, use YYYY-MM-DD", err)
    }
    
    // Parse coordinates
    if lat, err := strconv.ParseFloat(c.FormValue("latitude"), 64); err == nil {
        req.Latitude = lat
    } else if c.FormValue("latitude") != "" {
        return response.BadRequest(c, "Invalid latitude format", err)
    }
    if lng, err := strconv.ParseFloat(c.FormValue("longitude"), 64); err == nil {
        req.Longitude = lng
    } else if c.FormValue("longitude") != "" {
        return response.BadRequest(c, "Invalid longitude format", err)
    }
    
    // Parse Food Crops (Pangan) section
    req.FoodCommodity = c.FormValue("food_commodity")
    if req.FoodCommodity != "" {
        req.FoodLandStatus = c.FormValue("food_land_status")
        if area, err := strconv.ParseFloat(c.FormValue("food_land_area"), 64); err == nil && area > 0 {
            req.FoodLandArea = area
        }
        req.FoodGrowthPhase = c.FormValue("food_growth_phase")
        if age, err := strconv.Atoi(c.FormValue("food_plant_age")); err == nil && age >= 0 {
            req.FoodPlantAge = age
        }
        req.FoodPlantingDate = c.FormValue("food_planting_date")
        req.FoodHarvestDate = c.FormValue("food_harvest_date")
        req.FoodDelayReason = c.FormValue("food_delay_reason")
        req.FoodTechnology = c.FormValue("food_technology")
    }
    
    // Parse Horticulture (Hortikultura) section
    req.HortiCommodity = c.FormValue("horti_commodity")
    if req.HortiCommodity != "" {
        req.HortiSubCommodity = c.FormValue("horti_sub_commodity")
        req.HortiLandStatus = c.FormValue("horti_land_status")
        if area, err := strconv.ParseFloat(c.FormValue("horti_land_area"), 64); err == nil && area > 0 {
            req.HortiLandArea = area
        }
        req.HortiGrowthPhase = c.FormValue("horti_growth_phase")
        if age, err := strconv.Atoi(c.FormValue("horti_plant_age")); err == nil && age >= 0 {
            req.HortiPlantAge = age
        }
        req.HortiPlantingDate = c.FormValue("horti_planting_date")
        req.HortiHarvestDate = c.FormValue("horti_harvest_date")
        req.HortiDelayReason = c.FormValue("horti_delay_reason")
        req.HortiTechnology = c.FormValue("horti_technology")
        req.PostHarvestProblems = c.FormValue("post_harvest_problems")
    }
    
    // Parse Plantation (Perkebunan) section
    req.PlantationCommodity = c.FormValue("plantation_commodity")
    if req.PlantationCommodity != "" {
        req.PlantationLandStatus = c.FormValue("plantation_land_status")
        if area, err := strconv.ParseFloat(c.FormValue("plantation_land_area"), 64); err == nil && area > 0 {
            req.PlantationLandArea = area
        }
        req.PlantationGrowthPhase = c.FormValue("plantation_growth_phase")
        if age, err := strconv.Atoi(c.FormValue("plantation_plant_age")); err == nil && age >= 0 {
            req.PlantationPlantAge = age
        }
        req.PlantationPlantingDate = c.FormValue("plantation_planting_date")
        req.PlantationHarvestDate = c.FormValue("plantation_harvest_date")
        req.PlantationDelayReason = c.FormValue("plantation_delay_reason")
        req.PlantationTechnology = c.FormValue("plantation_technology")
        req.ProductionProblems = c.FormValue("production_problems")
    }
    
    // Parse Pest and Disease section
    if hasPest, err := strconv.ParseBool(c.FormValue("has_pest_disease")); err == nil {
        req.HasPestDisease = hasPest
        if hasPest {
            req.PestDiseaseType = c.FormValue("pest_disease_type")
            req.PestDiseaseCommodity = c.FormValue("pest_disease_commodity")
            req.AffectedArea = c.FormValue("affected_area")
            req.ControlAction = c.FormValue("control_action")
        }
    }
    
    // Parse Weather and Environment section
    req.WeatherCondition = c.FormValue("weather_condition")
    req.WeatherImpact = c.FormValue("weather_impact")
    req.MainConstraint = c.FormValue("main_constraint")
    
    // Parse Farmer Needs and Aspirations section
    req.FarmerHope = c.FormValue("farmer_hope")
    req.TrainingNeeded = c.FormValue("training_needed")
    req.UrgentNeeds = c.FormValue("urgent_needs")
    req.WaterAccess = c.FormValue("water_access")
    req.Suggestions = c.FormValue("suggestions")

    // Validation: At least one commodity must be specified
    if req.FoodCommodity == "" && req.HortiCommodity == "" && req.PlantationCommodity == "" {
        return response.BadRequest(c, "At least one commodity (food/horticulture/plantation) must be specified", nil)
    }

    // Validate the request
    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    // Get user ID from context
    userID, ok := c.Locals("userID").(string)
    if !ok {
        return response.BadRequest(c, "Invalid user ID type", nil)
    }

    // Parse multipart form for photos
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }

    photos := form.File["photos"]
    if len(photos) < 1 {
        return response.BadRequest(c, "At least 1 photo required", nil)
    }

    // Validate photo files
    for _, photo := range photos {
        if photo.Size > 10*1024*1024 { // 10MB limit
            return response.BadRequest(c, "Photo file size too large (max 10MB)", nil)
        }
        
        contentType := photo.Header.Get("Content-Type")
        if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
            return response.BadRequest(c, "Only JPEG and PNG images are allowed", nil)
        }
    }

    report, err := h.agricultureUseCase.CreateReport(c.Context(), &req, photos, userID)
    if err != nil {
        return response.InternalError(c, "Failed to create agriculture report", err)
    }

    return response.Created(c, "Agriculture report created successfully", report)
}

func (h *AgricultureHandler) GetReport(c *fiber.Ctx) error {
    idStr := c.Params("id")

    report, err := h.agricultureUseCase.GetReport(c.Context(), idStr)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *AgricultureHandler) ListReports(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    // Validate pagination parameters
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    filters := map[string]interface{}{
        "extension_officer":       c.Query("extension_officer"),
        "village":                 c.Query("village"),
        "district":                c.Query("district"),
        "farmer_name":             c.Query("farmer_name"),
        "farmer_group":            c.Query("farmer_group"),
        "farmer_group_type":       c.Query("farmer_group_type"),
        "food_commodity":          c.Query("food_commodity"),
        "horti_commodity":         c.Query("horti_commodity"),
        "plantation_commodity":    c.Query("plantation_commodity"),
        "main_constraint":         c.Query("main_constraint"),
        "weather_condition":       c.Query("weather_condition"),
        "water_access":            c.Query("water_access"),
        "start_date":              c.Query("start_date"),
        "end_date":                c.Query("end_date"),
    }

    // Parse boolean filters
    if hasPestStr := c.Query("has_pest_disease"); hasPestStr != "" {
        if hasPest, err := strconv.ParseBool(hasPestStr); err == nil {
            filters["has_pest_disease"] = hasPest
        }
    }

    result, err := h.agricultureUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports", err)
    }

    return response.Success(c, "Reports retrieved successfully", result)
}

func (h *AgricultureHandler) UpdateReport(c *fiber.Ctx) error {
    idStr := c.Params("id")


    var req dto.UpdateAgricultureRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    userID := c.Locals("userID").(string)

    report, err := h.agricultureUseCase.UpdateReport(c.Context(), idStr, &req, userID)
    if err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to update this report", err)
        }
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *AgricultureHandler) DeleteReport(c *fiber.Ctx) error {
    idStr := c.Params("id")

    userID := c.Locals("userID").(string)

    if err := h.agricultureUseCase.DeleteReport(c.Context(), idStr, userID); err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to delete this report", err)
        }
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *AgricultureHandler) GetStatistics(c *fiber.Ctx) error {
    stats, err := h.agricultureUseCase.GetStatistics(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve statistics", err)
    }

    return response.Success(c, "Statistics retrieved successfully", stats)
}

func (h *AgricultureHandler) GetByExtensionOfficer(c *fiber.Ctx) error {
    extensionOfficer := c.Params("officer")
    if extensionOfficer == "" {
        return response.BadRequest(c, "Extension officer name is required", nil)
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 50 {
        limit = 10
    }

    result, err := h.agricultureUseCase.GetByExtensionOfficer(c.Context(), extensionOfficer, page, limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports by extension officer", err)
    }

    return response.Success(c, "Reports by extension officer retrieved successfully", result)
}

func (h *AgricultureHandler) GetPestDiseaseReports(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    
    if limit < 1 || limit > 100 {
        limit = 20
    }

    reports, err := h.agricultureUseCase.GetPestDiseaseReports(c.Context(), limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve pest disease reports", err)
    }

    return response.Success(c, "Pest disease reports retrieved successfully", reports)
}

func (h *AgricultureHandler) GetCommodityProduction(c *fiber.Ctx) error {
    // Default to last 3 months
    endDate := time.Now()
    startDate := endDate.AddDate(0, -3, 0)
    
    if startDateStr := c.Query("start_date"); startDateStr != "" {
        if parsedDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
            startDate = parsedDate
        } else {
            return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
        }
    }
    
    if endDateStr := c.Query("end_date"); endDateStr != "" {
        if parsedDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
            endDate = parsedDate
        } else {
            return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
        }
    }
    
    // Validate date range
    if startDate.After(endDate) {
        return response.BadRequest(c, "Start date cannot be after end date", nil)
    }
    
    if endDate.Sub(startDate) > 365*24*time.Hour {
        return response.BadRequest(c, "Date range cannot exceed 1 year", nil)
    }

    results, err := h.agricultureUseCase.GetCommodityProduction(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve commodity production statistics", err)
    }

    return response.Success(c, "Commodity production statistics retrieved successfully", map[string]interface{}{
        "data":       results,
        "start_date": startDate.Format("2006-01-02"),
        "end_date":   endDate.Format("2006-01-02"),
    })
}

func (h *AgricultureHandler) GetExtensionOfficerPerformance(c *fiber.Ctx) error {
    // Default to last 6 months
    endDate := time.Now()
    startDate := endDate.AddDate(0, -6, 0)
    
    if startDateStr := c.Query("start_date"); startDateStr != "" {
        if parsedDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
            startDate = parsedDate
        }
    }
    
    if endDateStr := c.Query("end_date"); endDateStr != "" {
        if parsedDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
            endDate = parsedDate
        }
    }

    results, err := h.agricultureUseCase.GetExtensionOfficerPerformance(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve extension officer performance", err)
    }

    return response.Success(c, "Extension officer performance retrieved successfully", map[string]interface{}{
        "data":       results,
        "start_date": startDate.Format("2006-01-02"),
        "end_date":   endDate.Format("2006-01-02"),
    })
}

func (h *AgricultureHandler) GetTechnologyAdoptionStats(c *fiber.Ctx) error {
    stats, err := h.agricultureUseCase.GetTechnologyAdoptionStats(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve technology adoption statistics", err)
    }

    return response.Success(c, "Technology adoption statistics retrieved successfully", stats)
}

func (h *AgricultureHandler) GetFarmerNeedsAnalysis(c *fiber.Ctx) error {
    analysis, err := h.agricultureUseCase.GetFarmerNeedsAnalysis(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve farmer needs analysis", err)
    }

    return response.Success(c, "Farmer needs analysis retrieved successfully", analysis)
}

func (h *AgricultureHandler) GetReportsByVillage(c *fiber.Ctx) error {
    village := c.Params("village")
    if village == "" {
        return response.BadRequest(c, "Village name is required", nil)
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    // Create filter for village
    filters := map[string]interface{}{
        "village": village,
    }
    
    result, err := h.agricultureUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports by village", err)
    }

    return response.Success(c, "Reports by village retrieved successfully", result)
}

func (h *AgricultureHandler) GetReportsByDateRange(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "Both start_date and end_date are required", nil)
    }

    _, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    _, err = time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    filters := map[string]interface{}{
        "start_date": startDateStr,
        "end_date":   endDateStr,
    }
    
    result, err := h.agricultureUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports by date range", err)
    }

    return response.Success(c, "Reports by date range retrieved successfully", result)
}

func (h *AgricultureHandler) GetDashboardSummary(c *fiber.Ctx) error {
    // Get overall statistics
    stats, err := h.agricultureUseCase.GetStatistics(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve statistics", err)
    }

    // Get recent pest disease reports
    pestReports, err := h.agricultureUseCase.GetPestDiseaseReports(c.Context(), 5)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve pest disease reports", err)
    }

    // Get farmer needs analysis
    farmerNeeds, err := h.agricultureUseCase.GetFarmerNeedsAnalysis(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve farmer needs", err)
    }

    // Get technology adoption stats
    technologyStats, err := h.agricultureUseCase.GetTechnologyAdoptionStats(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve technology adoption stats", err)
    }

    summary := map[string]interface{}{
        "statistics":            stats,
        "recent_pest_reports":   pestReports,
        "farmer_needs_summary":  farmerNeeds,
        "technology_adoption":   technologyStats,
        "generated_at":          time.Now(),
    }

    return response.Success(c, "Dashboard summary retrieved successfully", summary)
}

func (h *AgricultureHandler) GetReportsByCommodity(c *fiber.Ctx) error {
    commodityType := c.Query("type") // food, horticulture, plantation
    commodity := c.Query("commodity")
    
    if commodityType == "" {
        return response.BadRequest(c, "Commodity type parameter is required (food, horticulture, plantation)", nil)
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    filters := map[string]interface{}{}
    
    switch commodityType {
    case "food":
        if commodity != "" {
            filters["food_commodity"] = commodity
        } else {
            // Get all food crop reports
            filters["has_food_commodity"] = true
        }
    case "horticulture":
        if commodity != "" {
            filters["horti_commodity"] = commodity
        } else {
            filters["has_horti_commodity"] = true
        }
    case "plantation":
        if commodity != "" {
            filters["plantation_commodity"] = commodity
        } else {
            filters["has_plantation_commodity"] = true
        }
    default:
        return response.BadRequest(c, "Invalid commodity type. Use: food, horticulture, or plantation", nil)
    }

    result, err := h.agricultureUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports by commodity", err)
    }

    return response.Success(c, "Reports by commodity retrieved successfully", result)
}

// Commodity Analysis Endpoints

// Executive Dashboard Endpoints

func (h *AgricultureHandler) GetExecutiveDashboard(c *fiber.Ctx) error {
    summary, err := h.agricultureUseCase.GetExecutiveSummary(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve executive summary", err)
    }

    return response.Success(c, "Executive summary retrieved successfully", summary)
}

// Commodity Analysis Endpoints

func (h *AgricultureHandler) GetCommodityAnalysis(c *fiber.Ctx) error {
    commodityName := c.Query("commodity_name")
    if commodityName == "" {
        return response.BadRequest(c, "commodity_name parameter is required", nil)
    }

    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "start_date and end_date parameters are required", nil)
    }

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    analysis, err := h.agricultureUseCase.GetCommodityAnalysis(c.Context(), startDate, endDate, commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve commodity analysis", err)
    }

    return response.Success(c, "Commodity analysis retrieved successfully", analysis)
}

// Food Crop Endpoints

func (h *AgricultureHandler) GetFoodCropStats(c *fiber.Ctx) error {
    commodityName := c.Query("commodity_name", "")
    
    stats, err := h.agricultureUseCase.GetFoodCropStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve food crop statistics", err)
    }

    return response.Success(c, "Food crop statistics retrieved successfully", stats)
}

// Horticulture Endpoints

func (h *AgricultureHandler) GetHorticultureStats(c *fiber.Ctx) error {
    commodityName := c.Query("commodity_name", "")
    
    stats, err := h.agricultureUseCase.GetHorticultureStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve horticulture statistics", err)
    }

    return response.Success(c, "Horticulture statistics retrieved successfully", stats)
}

// Plantation Endpoints

func (h *AgricultureHandler) GetPlantationStats(c *fiber.Ctx) error {
    commodityName := c.Query("commodity_name", "")
    
    stats, err := h.agricultureUseCase.GetPlantationStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve plantation statistics", err)
    }

    return response.Success(c, "Plantation statistics retrieved successfully", stats)
}

// Agricultural Equipment Endpoints

func (h *AgricultureHandler) GetAgriculturalEquipmentStats(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "start_date and end_date parameters are required", nil)
    }

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    stats, err := h.agricultureUseCase.GetAgriculturalEquipmentStats(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve agricultural equipment statistics", err)
    }

    return response.Success(c, "Agricultural equipment statistics retrieved successfully", stats)
}

// Land and Irrigation Endpoints

func (h *AgricultureHandler) GetLandAndIrrigationStats(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "start_date and end_date parameters are required", nil)
    }

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    stats, err := h.agricultureUseCase.GetLandAndIrrigationStats(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve land and irrigation statistics", err)
    }

    return response.Success(c, "Land and irrigation statistics retrieved successfully", stats)
}

// Helper Methods

// ValidateDateRange validates if end date is after start date
func (h *AgricultureHandler) ValidateDateRange(startDate, endDate time.Time) error {
    if endDate.Before(startDate) {
        return errors.New("end_date cannot be before start_date")
    }
    
    // Check if date range is not more than 2 years
    if endDate.Sub(startDate) > 2*365*24*time.Hour {
        return errors.New("date range cannot exceed 2 years")
    }
    
    return nil
}

// ParseAndValidateDate parses date string and validates format
func (h *AgricultureHandler) ParseAndValidateDate(dateStr, fieldName string) (time.Time, error) {
    if dateStr == "" {
        return time.Time{}, fmt.Errorf("%s parameter is required", fieldName)
    }
    
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid %s format, use YYYY-MM-DD", fieldName)
    }
    
    // Check if date is not in the future
    if date.After(time.Now()) {
        return time.Time{}, fmt.Errorf("%s cannot be in the future", fieldName)
    }
    
    // Check if date is not too old (more than 10 years)
    if date.Before(time.Now().AddDate(-10, 0, 0)) {
        return time.Time{}, fmt.Errorf("%s cannot be more than 10 years ago", fieldName)
    }
    
    return date, nil
}

// GetCommodityType returns commodity type based on commodity name
func (h *AgricultureHandler) GetCommodityType(commodityName string) string {
    foodCrops := []string{"padi", "jagung", "kedelai", "ubi jalar", "ubi kayu", "kacang tanah"}
    horticulture := []string{"cabai", "tomat", "wortel", "bawang merah", "bawang putih", "kentang"}
    plantation := []string{"kelapa", "kopi", "kakao", "karet", "kelapa sawit", "tebu"}
    
    commodityLower := strings.ToLower(commodityName)
    
    for _, crop := range foodCrops {
        if strings.Contains(commodityLower, crop) {
            return "FOOD"
        }
    }
    
    for _, crop := range horticulture {
        if strings.Contains(commodityLower, crop) {
            return "HORTICULTURE"
        }
    }
    
    for _, crop := range plantation {
        if strings.Contains(commodityLower, crop) {
            return "PLANTATION"
        }
    }
    
    return "UNKNOWN"
}

// LogRequestMetrics logs request metrics for monitoring
func (h *AgricultureHandler) LogRequestMetrics(c *fiber.Ctx, endpoint string, duration time.Duration) {
    log.Printf("[METRICS] Endpoint: %s, Method: %s, Duration: %v, IP: %s, UserAgent: %s",
        endpoint,
        c.Method(),
        duration,
        c.IP(),
        c.Get("User-Agent"),
    )
}

// HandlePanic recovers from panic and returns proper error response
func (h *AgricultureHandler) HandlePanic(c *fiber.Ctx) {
    if r := recover(); r != nil {
        log.Printf("[PANIC] Agriculture handler panic: %v", r)
        response.InternalError(c, "Internal server error occurred", fmt.Errorf("%v", r))
    }
}