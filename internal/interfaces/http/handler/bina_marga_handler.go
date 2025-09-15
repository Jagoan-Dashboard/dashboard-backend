
package handler

import (
    "strconv"
    "time"
    
    "building-report-backend/internal/application/dto"
    "building-report-backend/internal/application/usecase"
    "building-report-backend/internal/interfaces/response"
    
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type BinaMargaHandler struct {
    binaMargaUseCase *usecase.BinaMargaUseCase
}

func NewBinaMargaHandler(binaMargaUseCase *usecase.BinaMargaUseCase) *BinaMargaHandler {
    return &BinaMargaHandler{
        binaMargaUseCase: binaMargaUseCase,
    }
}

func (h *BinaMargaHandler) CreateReport(c *fiber.Ctx) error {
    var req dto.CreateBinaMargaRequest
    
    
    req.ReporterName = c.FormValue("reporter_name")
    req.InstitutionUnit = c.FormValue("institution_unit")
    req.PhoneNumber = c.FormValue("phone_number")
    req.RoadName = c.FormValue("road_name")
    req.RoadType = c.FormValue("road_type")
    req.RoadClass = c.FormValue("road_class")
    req.DamageType = c.FormValue("damage_type")
    req.DamageLevel = c.FormValue("damage_level")
    req.TrafficImpact = c.FormValue("traffic_impact")
    req.UrgencyLevel = c.FormValue("urgency_level")
    req.CauseOfDamage = c.FormValue("cause_of_damage")
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
    if length, err := strconv.ParseFloat(c.FormValue("damaged_length"), 64); err == nil {
        req.DamagedLength = length
    }
    if width, err := strconv.ParseFloat(c.FormValue("damaged_width"), 64); err == nil {
        req.DamagedWidth = width
    }

    
    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    
    userID := c.Locals("userID").(uuid.UUID)

    
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }

    photos := form.File["photos"]
    if len(photos) < 2 {
        return response.BadRequest(c, "Minimum 2 photos required (before and after damage)", nil)
    }

    report, err := h.binaMargaUseCase.CreateReport(c.Context(), &req, photos, userID)
    if err != nil {
        return response.InternalError(c, "Failed to create bina marga report", err)
    }

    return response.Created(c, "Bina Marga report created successfully", report)
}

func (h *BinaMargaHandler) GetReport(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

    report, err := h.binaMargaUseCase.GetReport(c.Context(), id)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *BinaMargaHandler) ListReports(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    filters := map[string]interface{}{
        "institution_unit": c.Query("institution_unit"),
        "road_type":        c.Query("road_type"),
        "road_class":       c.Query("road_class"),
        "road_name":        c.Query("road_name"),
        "damage_type":      c.Query("damage_type"),
        "damage_level":     c.Query("damage_level"),
        "urgency_level":    c.Query("urgency_level"),
        "traffic_impact":   c.Query("traffic_impact"),
        "status":          c.Query("status"),
        "start_date":      c.Query("start_date"),
        "end_date":        c.Query("end_date"),
    }

    result, err := h.binaMargaUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports", err)
    }

    return response.Success(c, "Reports retrieved successfully", result)
}

func (h *BinaMargaHandler) ListByPriority(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    result, err := h.binaMargaUseCase.ListByPriority(c.Context(), page, limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve priority reports", err)
    }

    return response.Success(c, "Priority reports retrieved successfully", result)
}

func (h *BinaMargaHandler) UpdateReport(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

    var req dto.UpdateBinaMargaRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    userID := c.Locals("userID").(uuid.UUID)

    report, err := h.binaMargaUseCase.UpdateReport(c.Context(), id, &req, userID)
    if err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to update this report", err)
        }
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *BinaMargaHandler) UpdateStatus(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

    var req dto.UpdateBinaMargaStatusRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    if err := h.binaMargaUseCase.UpdateStatus(c.Context(), id, &req); err != nil {
        return response.InternalError(c, "Failed to update report status", err)
    }

    return response.Success(c, "Report status updated successfully", nil)
}

func (h *BinaMargaHandler) DeleteReport(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

    userID := c.Locals("userID").(uuid.UUID)

    if err := h.binaMargaUseCase.DeleteReport(c.Context(), id, userID); err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to delete this report", err)
        }
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *BinaMargaHandler) GetStatistics(c *fiber.Ctx) error {
    stats, err := h.binaMargaUseCase.GetStatistics(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve statistics", err)
    }

    return response.Success(c, "Statistics retrieved successfully", stats)
}

func (h *BinaMargaHandler) GetEmergencyReports(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    reports, err := h.binaMargaUseCase.GetEmergencyReports(c.Context(), limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve emergency reports", err)
    }

    return response.Success(c, "Emergency reports retrieved successfully", reports)
}

func (h *BinaMargaHandler) GetBlockedRoads(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    reports, err := h.binaMargaUseCase.GetBlockedRoads(c.Context(), limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve blocked roads", err)
    }

    return response.Success(c, "Blocked roads retrieved successfully", reports)
}

func (h *BinaMargaHandler) GetDamageByRoadType(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
    endDateStr := c.Query("end_date", time.Now().Format("2006-01-02"))
    
    startDate, _ := time.Parse("2006-01-02", startDateStr)
    endDate, _ := time.Parse("2006-01-02", endDateStr)
    
    results, err := h.binaMargaUseCase.GetDamageByRoadType(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve damage statistics by road type", err)
    }

    return response.Success(c, "Damage statistics by road type retrieved successfully", results)
}

func (h *BinaMargaHandler) GetDamageByLocation(c *fiber.Ctx) error {
    
    bounds := map[string]float64{}
    
    if north, err := strconv.ParseFloat(c.Query("north"), 64); err == nil {
        bounds["north"] = north
    }
    if south, err := strconv.ParseFloat(c.Query("south"), 64); err == nil {
        bounds["south"] = south
    }
    if east, err := strconv.ParseFloat(c.Query("east"), 64); err == nil {
        bounds["east"] = east
    }
    if west, err := strconv.ParseFloat(c.Query("west"), 64); err == nil {
        bounds["west"] = west
    }
    
    if len(bounds) != 4 {
        return response.BadRequest(c, "Missing bounding box parameters (north, south, east, west)", nil)
    }
    
    results, err := h.binaMargaUseCase.GetDamageByLocation(c.Context(), bounds)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve damage statistics by location", err)
    }

    return response.Success(c, "Damage statistics by location retrieved successfully", results)
}