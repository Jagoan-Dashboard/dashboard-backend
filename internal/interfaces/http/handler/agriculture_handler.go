package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"
	"building-report-backend/pkg/utils"

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

    req.ExtensionOfficer = c.FormValue("extension_officer")
    req.FarmerName = c.FormValue("farmer_name")
    req.FarmerGroup = c.FormValue("farmer_group")
    req.FarmerGroupType = c.FormValue("farmer_group_type")
    req.Village = c.FormValue("village")
    req.District = c.FormValue("district")

    visitDateStr := c.FormValue("visit_date")
    if visitDate, err := time.Parse("2006-01-02", visitDateStr); err == nil {
        req.VisitDate = visitDate
    } else {
        return response.BadRequest(c, "Invalid visit date format, use YYYY-MM-DD", err)
    }

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

    
    if hasPest, err := strconv.ParseBool(c.FormValue("has_pest_disease")); err == nil {
        req.HasPestDisease = hasPest
        if hasPest {
            req.PestDiseaseType = c.FormValue("pest_disease_type")
            req.PestDiseaseCommodity = c.FormValue("pest_disease_commodity")
            req.AffectedArea = c.FormValue("affected_area")
            req.ControlAction = c.FormValue("control_action")
        }
    }

    
    req.WeatherCondition = c.FormValue("weather_condition")
    req.WeatherImpact = c.FormValue("weather_impact")
    req.MainConstraint = c.FormValue("main_constraint")

    req.FarmerHope = c.FormValue("farmer_hope")
    req.TrainingNeeded = c.FormValue("training_needed")
    req.UrgentNeeds = c.FormValue("urgent_needs")
    req.WaterAccess = c.FormValue("water_access")
    req.Suggestions = c.FormValue("suggestions")

    if req.FoodCommodity == "" && req.HortiCommodity == "" && req.PlantationCommodity == "" {
        return response.BadRequest(c, "At least one commodity (food/horticulture/plantation) must be specified", nil)
    }

    
    req.Normalize()

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }
    photos := form.File["photos"]
    if len(photos) < 1 {
        return response.BadRequest(c, "At least 1 photo required", nil)
    }
    for _, photo := range photos {
        if photo.Size > 10*1024*1024 {
            return response.BadRequest(c, "Photo file size too large (max 10MB)", nil)
        }
        contentType := photo.Header.Get("Content-Type")
        if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
            return response.BadRequest(c, "Only JPEG and PNG images are allowed", nil)
        }
    }

    report, err := h.agricultureUseCase.CreateReport(c.Context(), &req, photos)
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

    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    filters := map[string]interface{}{
        "extension_officer":    c.Query("extension_officer"),
        "village":              c.Query("village"),
        "district":             c.Query("district"),
        "farmer_name":          c.Query("farmer_name"),
        "farmer_group":         c.Query("farmer_group"),
        "farmer_group_type":    c.Query("farmer_group_type"),
        "food_commodity":       c.Query("food_commodity"),
        "horti_commodity":      c.Query("horti_commodity"),
        "plantation_commodity": c.Query("plantation_commodity"),
        "main_constraint":      c.Query("main_constraint"),
        "weather_condition":    c.Query("weather_condition"),
        "water_access":         c.Query("water_access"),
        "start_date":           c.Query("start_date"),
        "end_date":             c.Query("end_date"),
    }

    if hasPestStr := c.Query("has_pest_disease"); hasPestStr != "" {
        if hasPest, err := strconv.ParseBool(hasPestStr); err == nil {
            filters["has_pest_disease"] = hasPest
        }
    }

    
    normalizeFilters(filters)

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

    
    req.Normalize()

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

func (h *AgricultureHandler) GetExecutiveDashboard(c *fiber.Ctx) error {
    commodityType := c.Query("commodity_type", "") 
    
    summary, err := h.agricultureUseCase.GetExecutiveSummary(c.Context(), commodityType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve executive summary", err)
    }

    return response.Success(c, "Executive summary retrieved successfully", summary)
}

func (h *AgricultureHandler) GetCommodityAnalysis(c *fiber.Ctx) error {
    rawName := c.Query("commodity_name")
    if rawName == "" {
        return response.BadRequest(c, "commodity_name parameter is required", nil)
    }
    commodityName := utils.NormalizeEnum(rawName)

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

    if analysis.TotalHarvestedArea == 0 {
        return response.Success(c, 
            fmt.Sprintf("No data found for commodity '%s' in the specified date range", commodityName), 
            analysis)
    }

    return response.Success(c, "Commodity analysis retrieved successfully", analysis)
}

func (h *AgricultureHandler) GetFoodCropStats(c *fiber.Ctx) error {
    commodityName := utils.NormalizeEnum(c.Query("commodity_name", ""))

    stats, err := h.agricultureUseCase.GetFoodCropStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve food crop statistics", err)
    }

    return response.Success(c, "Food crop statistics retrieved successfully", stats)
}

func (h *AgricultureHandler) GetHorticultureStats(c *fiber.Ctx) error {
    commodityName := utils.NormalizeEnum(c.Query("commodity_name", ""))

    stats, err := h.agricultureUseCase.GetHorticultureStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve horticulture statistics", err)
    }

    return response.Success(c, "Horticulture statistics retrieved successfully", stats)
}


func (h *AgricultureHandler) GetPlantationStats(c *fiber.Ctx) error {
    commodityName := utils.NormalizeEnum(c.Query("commodity_name", ""))

    stats, err := h.agricultureUseCase.GetPlantationStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve plantation statistics", err)
    }

    return response.Success(c, "Plantation statistics retrieved successfully", stats)
}


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

func (h *AgricultureHandler) HandlePanic(c *fiber.Ctx) {
    if r := recover(); r != nil {
        log.Printf("[PANIC] Agriculture handler panic: %v", r)
        response.InternalError(c, "Internal server error occurred", fmt.Errorf("%v", r))
    }
}

var locationKeys = map[string]bool{
    "extension_officer": true,
    "village":           true,
    "district":          true,
    "farmer_name":       true,
    "farmer_group":      true,
}


var enumKeys = map[string]bool{
    "farmer_group_type":     true,
    "food_commodity":        true,
    "horti_commodity":       true,
    "plantation_commodity":  true,
    "main_constraint":       true,
    "weather_condition":     true,
    "water_access":          true,
    
    
}


func normalizeFilters(filters map[string]interface{}) {
    for k, v := range filters {
        s, ok := v.(string)
        if !ok || s == "" {
            continue
        }
        if locationKeys[k] {
            filters[k] = utils.NormalizeLocation(s)
            continue
        }
        if enumKeys[k] {
            filters[k] = utils.NormalizeEnum(s)
            continue
        }
        
    }
}