package handler

import (
    "strconv"
    "time"
    "fmt"
    "building-report-backend/internal/application/dto"
    "building-report-backend/internal/application/usecase"
    "building-report-backend/internal/interfaces/response"
    
    "github.com/gofiber/fiber/v2"
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
    
    // Parse basic reporter information
    req.ReporterName = c.FormValue("reporter_name")
    req.InstitutionUnit = c.FormValue("institution_unit")
    req.PhoneNumber = c.FormValue("phone_number")
    
    // Parse datetime
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
    
    req.District = c.FormValue("district")
    
    // Parse road information
    req.RoadName = c.FormValue("road_name")
    req.RoadType = c.FormValue("road_type")
    req.RoadClass = c.FormValue("road_class")
    
    // Parse segment length
    if segmentLength, err := strconv.ParseFloat(c.FormValue("segment_length"), 64); err == nil {
        req.SegmentLength = segmentLength
    }
    
    // Parse coordinates
    if lat, err := strconv.ParseFloat(c.FormValue("latitude"), 64); err == nil {
        req.Latitude = lat
    }
    if lng, err := strconv.ParseFloat(c.FormValue("longitude"), 64); err == nil {
        req.Longitude = lng
    }
    
    // Parse pavement and damage information
    req.PavementType = c.FormValue("pavement_type")
    req.DamageType = c.FormValue("damage_type")
    req.DamageLevel = c.FormValue("damage_level")
    
    // Parse damage dimensions
    if length, err := strconv.ParseFloat(c.FormValue("damaged_length"), 64); err == nil {
        req.DamagedLength = length
    }
    if width, err := strconv.ParseFloat(c.FormValue("damaged_width"), 64); err == nil {
        req.DamagedWidth = width
    }
    if totalArea, err := strconv.ParseFloat(c.FormValue("total_damaged_area"), 64); err == nil {
        req.TotalDamagedArea = totalArea
    }
    
    // Parse bridge information (optional)
    req.BridgeName = c.FormValue("bridge_name")
    req.BridgeStructureType = c.FormValue("bridge_structure_type")
    req.BridgeDamageType = c.FormValue("bridge_damage_type")
    req.BridgeDamageLevel = c.FormValue("bridge_damage_level")
    
    // Parse traffic and urgency information
    req.TrafficCondition = c.FormValue("traffic_condition")
    req.TrafficImpact = c.FormValue("traffic_impact")
    if dailyVolume, err := strconv.Atoi(c.FormValue("daily_traffic_volume")); err == nil {
        req.DailyTrafficVolume = dailyVolume
    }
    req.UrgencyLevel = c.FormValue("urgency_level")
    
    // Parse optional fields
    req.CauseOfDamage = c.FormValue("cause_of_damage")
    req.Notes = c.FormValue("notes")

    // Validate the request
    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    // Get user ID from context
    // userID := c.Locals("userID").(string)

    // Parse multipart form for photos
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }

    photos := form.File["photos"]
    if len(photos) < 2 {
        return response.BadRequest(c, "Minimum 2 photos required (before and damage detail)", nil)
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

    report, err := h.binaMargaUseCase.CreateReport(c.Context(), &req, photos)
    if err != nil {
        return response.InternalError(c, "Failed to create bina marga report", err)
    }

    return response.Created(c, "Bina Marga report created successfully", report)
}

func (h *BinaMargaHandler) GetReport(c *fiber.Ctx) error {
    id := c.Params("id")
   
    report, err := h.binaMargaUseCase.GetReport(c.Context(), id)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *BinaMargaHandler) ListReports(c *fiber.Ctx) error {
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
        "institution_unit":      c.Query("institution_unit"),
        "road_type":            c.Query("road_type"),
        "road_class":           c.Query("road_class"),
        "road_name":            c.Query("road_name"),
        "pavement_type":        c.Query("pavement_type"),
        "damage_type":          c.Query("damage_type"),
        "damage_level":         c.Query("damage_level"),
        "urgency_level":        c.Query("urgency_level"),
        "traffic_impact":       c.Query("traffic_impact"),
        "traffic_condition":    c.Query("traffic_condition"),
        "bridge_name":          c.Query("bridge_name"),
        "status":              c.Query("status"),
        "start_date":          c.Query("start_date"),
        "end_date":            c.Query("end_date"),
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

    var req dto.UpdateBinaMargaRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    userID := c.Locals("userID").(string)

    report, err := h.binaMargaUseCase.UpdateReport(c.Context(), idStr, &req, userID)
    if err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to update this report", err)
        }
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *BinaMargaHandler) UpdateStatus(c *fiber.Ctx) error {
    id := c.Params("id")
    

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
    id := c.Params("id")
    

    userID := c.Locals("userID").(string)

    if err := h.binaMargaUseCase.DeleteReport(c.Context(), id, userID); err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to delete this report", err)
        }
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *BinaMargaHandler) GetBinaMargaOverview(c *fiber.Ctx) error {
    roadType := c.Query("road_type", "all") // all, JALAN_NASIONAL, JALAN_PROVINSI, etc.
    
    // Validate road type
    validRoadTypes := map[string]bool{
        "all":             true,
        "ALL":             true,
        "JALAN_NASIONAL":  true,
        "JALAN_PROVINSI":  true,
        "JALAN_KABUPATEN": true,
        "JALAN_DESA":      true,
    }
    
    if !validRoadTypes[roadType] {
        return response.BadRequest(c, "Invalid road type", 
            fmt.Errorf("road_type must be one of: all, JALAN_NASIONAL, JALAN_PROVINSI, JALAN_KABUPATEN, JALAN_DESA"))
    }
    
    // Convert "all" to empty string for repository layer
    queryRoadType := roadType
    if roadType == "all" || roadType == "ALL" {
        queryRoadType = ""
    }

    overview, err := h.binaMargaUseCase.GetBinaMargaOverview(c.Context(), queryRoadType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve bina marga overview", err)
    }

    return response.Success(c, "Bina marga overview retrieved successfully", overview)
}