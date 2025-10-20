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

type WaterResourcesHandler struct {
    waterUseCase *usecase.WaterResourcesUseCase
}

func NewWaterResourcesHandler(waterUseCase *usecase.WaterResourcesUseCase) *WaterResourcesHandler {
    return &WaterResourcesHandler{
        waterUseCase: waterUseCase,
    }
}

func (h *WaterResourcesHandler) CreateReport(c *fiber.Ctx) error {
    var req dto.CreateWaterResourcesRequest
    
    
    req.ReporterName = c.FormValue("reporter_name")
    req.InstitutionUnit = c.FormValue("institution_unit")
    req.PhoneNumber = c.FormValue("phone_number")
    req.IrrigationAreaName = c.FormValue("irrigation_area_name")
    req.IrrigationType = c.FormValue("irrigation_type")
    req.DamageType = c.FormValue("damage_type")
    req.DamageLevel = c.FormValue("damage_level")
    req.UrgencyCategory = c.FormValue("urgency_category")
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
    if length, err := strconv.ParseFloat(c.FormValue("estimated_length"), 64); err == nil {
        req.EstimatedLength = length
    }
    if width, err := strconv.ParseFloat(c.FormValue("estimated_width"), 64); err == nil {
        req.EstimatedWidth = width
    }
    if volume, err := strconv.ParseFloat(c.FormValue("estimated_volume"), 64); err == nil {
        req.EstimatedVolume = volume
    }
    if area, err := strconv.ParseFloat(c.FormValue("affected_rice_field_area"), 64); err == nil {
        req.AffectedRiceFieldArea = area
    }
    if farmers, err := strconv.Atoi(c.FormValue("affected_farmers_count")); err == nil {
        req.AffectedFarmersCount = farmers
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
    if len(photos) < 2 {
        return response.BadRequest(c, "Minimum 2 photos required", nil)
    }

    report, err := h.waterUseCase.CreateReport(c.Context(), &req, photos)
    if err != nil {
        return response.InternalError(c, "Failed to create water resources report", err)
    }

    return response.Created(c, "Water resources report created successfully", report)
}

func (h *WaterResourcesHandler) GetReport(c *fiber.Ctx) error {
    id := c.Params("id")
    report, err := h.waterUseCase.GetReport(c.Context(), id)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *WaterResourcesHandler) ListReports(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    filters := map[string]interface{}{
        "institution_unit":  c.Query("institution_unit"),
        "irrigation_type":   c.Query("irrigation_type"),
        "irrigation_area":   c.Query("irrigation_area"),
        "damage_type":       c.Query("damage_type"),
        "damage_level":      c.Query("damage_level"),
        "urgency_category":  c.Query("urgency_category"),
        "status":           c.Query("status"),
        "start_date":       c.Query("start_date"),
        "end_date":         c.Query("end_date"),
    }

    result, err := h.waterUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports", err)
    }

    return response.Success(c, "Reports retrieved successfully", result)
}

func (h *WaterResourcesHandler) ListByPriority(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    result, err := h.waterUseCase.ListByPriority(c.Context(), page, limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve priority reports", err)
    }

    return response.Success(c, "Priority reports retrieved successfully", result)
}

func (h *WaterResourcesHandler) UpdateReport(c *fiber.Ctx) error {
    id := c.Params("id")
    var req dto.UpdateWaterResourcesRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    userID := c.Locals("userID").(string)

    report, err := h.waterUseCase.UpdateReport(c.Context(), id, &req, userID)
    if err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to update this report", err)
        }
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *WaterResourcesHandler) UpdateStatus(c *fiber.Ctx) error {
    id := c.Params("id")

    var req dto.UpdateWaterStatusRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    if err := h.waterUseCase.UpdateStatus(c.Context(), id, &req); err != nil {
        return response.InternalError(c, "Failed to update report status", err)
    }

    return response.Success(c, "Report status updated successfully", nil)
}

func (h *WaterResourcesHandler) DeleteReport(c *fiber.Ctx) error {
    id := c.Params("id")

    userID := c.Locals("userID").(string)

    if err := h.waterUseCase.DeleteReport(c.Context(), id, userID); err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to delete this report", err)
        }
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *WaterResourcesHandler) GetWaterResourcesOverview(c *fiber.Ctx) error {
    irrigationType := c.Query("irrigation_type", "all") // all, IRIGASI_PRIMER, IRIGASI_SEKUNDER, etc.
    
    // Validate irrigation type if needed
    validIrrigationTypes := map[string]bool{
        "all":                true,
        "ALL":                true,
        "IRIGASI_PRIMER":     true,
        "IRIGASI_SEKUNDER":   true,
        "IRIGASI_TERSIER":    true,
        "BENDUNG":            true,
        "EMBUNG_DAM":         true,
        "PINTU_AIR":          true,
        "SALURAN_DRAINASE":   true,
        "LAINNYA":            true,
    }
    
    if !validIrrigationTypes[irrigationType] {
        return response.BadRequest(c, "Invalid irrigation type", 
            fmt.Errorf("irrigation_type must be one of the valid types"))
    }
    
    // Convert "all" to empty string for repository layer
    queryIrrigationType := irrigationType
    if irrigationType == "all" || irrigationType == "ALL" {
        queryIrrigationType = ""
    }

    overview, err := h.waterUseCase.GetWaterResourcesOverview(c.Context(), queryIrrigationType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve water resources overview", err)
    }

    return response.Success(c, "Water resources overview retrieved successfully", overview)
}