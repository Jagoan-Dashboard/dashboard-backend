package handler

import (
	"fmt"
	"strconv"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"

	"github.com/gofiber/fiber/v2"
)

type SpatialPlanningHandler struct {
    spatialUseCase *usecase.SpatialPlanningUseCase
}

func NewSpatialPlanningHandler(spatialUseCase *usecase.SpatialPlanningUseCase) *SpatialPlanningHandler {
    return &SpatialPlanningHandler{
        spatialUseCase: spatialUseCase,
    }
}

func (h *SpatialPlanningHandler) CreateReport(c *fiber.Ctx) error {
    var req dto.CreateSpatialPlanningRequest
    
    
    req.ReporterName = c.FormValue("reporter_name")
    req.Institution = c.FormValue("institution")
    req.PhoneNumber = c.FormValue("phone_number")
    req.AreaDescription = c.FormValue("area_description")
    req.AreaCategory = c.FormValue("area_category")
    req.ViolationType = c.FormValue("violation_type")
    req.ViolationLevel = c.FormValue("violation_level")
    req.EnvironmentalImpact = c.FormValue("environmental_impact")
    req.UrgencyLevel = c.FormValue("urgency_level")
    req.Address = c.FormValue("address")
    req.Notes = c.FormValue("notes")
    
    
    reportDateTimeStr := c.FormValue("report_datetime")
    if reportDateTime, err := time.Parse(time.RFC3339, reportDateTimeStr); err == nil {
        req.ReportDateTime = reportDateTime
    } else {
        
        if reportDateTime, err := time.Parse("2006-01-02 15:04:05", reportDateTimeStr); err == nil {
            req.ReportDateTime = reportDateTime
        } else {
            return response.BadRequest(c, "Invalid datetime format", err)
        }
    }
    
    
    if lat, err := strconv.ParseFloat(c.FormValue("latitude"), 64); err == nil {
        req.Latitude = lat
    }
    if lng, err := strconv.ParseFloat(c.FormValue("longitude"), 64); err == nil {
        req.Longitude = lng
    }

    
    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    
    // userID := c.Locals("userID").(string)

    
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }

    photos := form.File["photos"]
    if len(photos) < 1 {
        return response.BadRequest(c, "Minimum 1 photo required", nil)
    }

    report, err := h.spatialUseCase.CreateReport(c.Context(), &req, photos)
    if err != nil {
        return response.InternalError(c, "Failed to create spatial planning report", err)
    }

    return response.Created(c, "Spatial planning report created successfully", report)
}

func (h *SpatialPlanningHandler) GetReport(c *fiber.Ctx) error {
    id := c.Params("id")

    report, err := h.spatialUseCase.GetReport(c.Context(), id)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *SpatialPlanningHandler) ListReports(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    filters := map[string]interface{}{
        "institution":      c.Query("institution"),
        "area_category":    c.Query("area_category"),
        "violation_type":   c.Query("violation_type"),
        "violation_level":  c.Query("violation_level"),
        "urgency_level":    c.Query("urgency_level"),
        "status":          c.Query("status"),
        "start_date":      c.Query("start_date"),
        "end_date":        c.Query("end_date"),
    }

    result, err := h.spatialUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports", err)
    }

    return response.Success(c, "Reports retrieved successfully", result)
}

func (h *SpatialPlanningHandler) UpdateReport(c *fiber.Ctx) error {
    id:= c.Params("id")
   
    var req dto.UpdateSpatialPlanningRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    userID := c.Locals("userID").(string)

    report, err := h.spatialUseCase.UpdateReport(c.Context(), id, &req, userID)
    if err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to update this report", err)
        }
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *SpatialPlanningHandler) UpdateStatus(c *fiber.Ctx) error {
    id := c.Params("id")

    var req dto.UpdateSpatialStatusRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    if err := h.spatialUseCase.UpdateStatus(c.Context(), id, &req); err != nil {
        return response.InternalError(c, "Failed to update report status", err)
    }

    return response.Success(c, "Report status updated successfully", nil)
}

func (h *SpatialPlanningHandler) DeleteReport(c *fiber.Ctx) error {
    id := c.Params("id")
   
    userID := c.Locals("userID").(string)

    if err := h.spatialUseCase.DeleteReport(c.Context(), id, userID); err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to delete this report", err)
        }
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *SpatialPlanningHandler) GetStatistics(c *fiber.Ctx) error {
    stats, err := h.spatialUseCase.GetStatistics(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve statistics", err)
    }

    return response.Success(c, "Statistics retrieved successfully", stats)
}

// GetTataRuangOverview handles the overview endpoint for tata ruang page
func (h *SpatialPlanningHandler) GetTataRuangOverview(c *fiber.Ctx) error {
    areaCategory := c.Query("area_category", "all")
    
    // Validate area category based on your entity.AreaCategory constants
    validAreaCategories := map[string]bool{
        "all":                              true,
        "KAWASAN_CAGAR_BUDAYA":            true,
        "KAWASAN_HUTAN":                   true,
        "KAWASAN_PARIWISATA":              true,
        "KAWASAN_PERKEBUNAN":              true,
        "KAWASAN_PERMUKIMAN":              true,
        "KAWASAN_PERTAHANAN_KEAMANAN":     true,
        "KAWASAN_PERUNTUKAN_INDUSTRI":     true,
        "KAWASAN_PERUNTUKAN_PERTAMBANGAN": true,
        "KAWASAN_TANAMAN_PANGAN":          true,
        "KAWASAN_TRANSPORTASI":            true,
        "LAINNYA":                         true,
    }
    
    if !validAreaCategories[areaCategory] {
        return response.BadRequest(c, "Invalid area category", fmt.Errorf("area_category must be one of: all, KAWASAN_CAGAR_BUDAYA, KAWASAN_HUTAN, KAWASAN_PARIWISATA, KAWASAN_PERKEBUNAN, KAWASAN_PERMUKIMAN, KAWASAN_PERTAHANAN_KEAMANAN, KAWASAN_PERUNTUKAN_INDUSTRI, KAWASAN_PERUNTUKAN_PERTAMBANGAN, KAWASAN_TANAMAN_PANGAN, KAWASAN_TRANSPORTASI, LAINNYA"))
    }
    
    // Convert "all" to empty string for repository layer
    queryAreaCategory := areaCategory
    if areaCategory == "all" {
        queryAreaCategory = ""
    }

    overview, err := h.spatialUseCase.GetTataRuangOverview(c.Context(), queryAreaCategory)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve tata ruang overview", err)
    }

    return response.Success(c, "Tata ruang overview retrieved successfully", overview)
}