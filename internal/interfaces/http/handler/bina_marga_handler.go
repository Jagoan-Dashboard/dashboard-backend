package handler

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"
	"building-report-backend/pkg/utils"
	"fmt"
	"strconv"
	"time"

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

    // Reporter
    req.ReporterName = c.FormValue("reporter_name")
    req.InstitutionUnit = c.FormValue("institution_unit")
    req.PhoneNumber = c.FormValue("phone_number")

    // Datetime
    reportDateTimeStr := c.FormValue("report_datetime")
    if reportDateTime, err := time.Parse(time.RFC3339, reportDateTimeStr); err == nil {
        req.ReportDateTime = reportDateTime
    } else if reportDateTime, err := time.Parse("2006-01-02 15:04:05", reportDateTimeStr); err == nil {
        req.ReportDateTime = reportDateTime
    } else {
        return response.BadRequest(c, "Invalid datetime format", fmt.Errorf("use RFC3339 or 2006-01-02 15:04:05"))
    }

    // Road info
    req.RoadName = c.FormValue("road_name")
    req.RoadType = c.FormValue("road_type")
    req.RoadClass = c.FormValue("road_class")

    if segmentLength, err := strconv.ParseFloat(c.FormValue("segment_length"), 64); err == nil {
        req.SegmentLength = segmentLength
    }

    // Coordinates
    if lat, err := strconv.ParseFloat(c.FormValue("latitude"), 64); err == nil {
        req.Latitude = lat
    }
    if lng, err := strconv.ParseFloat(c.FormValue("longitude"), 64); err == nil {
        req.Longitude = lng
    }

    // Pavement & damage
    req.PavementType = c.FormValue("pavement_type")
    req.DamageType = c.FormValue("damage_type")
    req.DamageLevel = c.FormValue("damage_level")

    if length, err := strconv.ParseFloat(c.FormValue("damaged_length"), 64); err == nil {
        req.DamagedLength = length
    }
    if width, err := strconv.ParseFloat(c.FormValue("damaged_width"), 64); err == nil {
        req.DamagedWidth = width
    }
    if totalArea, err := strconv.ParseFloat(c.FormValue("total_damaged_area"), 64); err == nil {
        req.TotalDamagedArea = totalArea
    }

    // Bridge (optional)
    req.BridgeName = c.FormValue("bridge_name")
    req.BridgeStructureType = c.FormValue("bridge_structure_type")
    req.BridgeDamageType = c.FormValue("bridge_damage_type")
    req.BridgeDamageLevel = c.FormValue("bridge_damage_level")

    // Traffic & urgency
    req.TrafficCondition = c.FormValue("traffic_condition")
    req.TrafficImpact = c.FormValue("traffic_impact")
    if dailyVolume, err := strconv.Atoi(c.FormValue("daily_traffic_volume")); err == nil {
        req.DailyTrafficVolume = dailyVolume
    }
    req.UrgencyLevel = c.FormValue("urgency_level")

    // Optional
    req.CauseOfDamage = c.FormValue("cause_of_damage")
    req.Notes = c.FormValue("notes")

    // === NORMALIZE DTO ===
    req.Normalize()

    // Validate
    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    // Photos
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }
    photos := form.File["photos"]
    if len(photos) < 2 {
        return response.BadRequest(c, "Minimum 2 photos required (before and damage detail)", nil)
    }
    for _, photo := range photos {
        if photo.Size > 10*1024*1024 {
            return response.BadRequest(c, "Photo file size too large (max 10MB)", nil)
        }
        ct := photo.Header.Get("Content-Type")
        if ct != "image/jpeg" && ct != "image/png" && ct != "image/jpg" {
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
    if page < 1 { page = 1 }
    if limit < 1 || limit > 100 { limit = 10 }

    filters := map[string]interface{}{
        "institution_unit":     c.Query("institution_unit"),
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
        "status":               c.Query("status"),
        "start_date":           c.Query("start_date"),
        "end_date":             c.Query("end_date"),
    }

    // === NORMALIZE FILTERS ===
    normalizeBinaMargaFilters(filters)

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

    // === NORMALIZE DTO ===
    req.Normalize()

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

    // === NORMALIZE DTO ===
    if n, ok := interface{}(req).(interface{ Normalize() }); ok {
        n.Normalize()
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

func (h *BinaMargaHandler) GetStatistics(c *fiber.Ctx) error {
    stats, err := h.binaMargaUseCase.GetStatistics(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve statistics", err)
    }

    return response.Success(c, "Statistics retrieved successfully", stats)
}

func (h *BinaMargaHandler) GetEmergencyReports(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    if limit < 1 || limit > 50 {
        limit = 10
    }
    
    reports, err := h.binaMargaUseCase.GetEmergencyReports(c.Context(), limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve emergency reports", err)
    }

    return response.Success(c, "Emergency reports retrieved successfully", reports)
}

func (h *BinaMargaHandler) GetBlockedRoads(c *fiber.Ctx) error {
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    if limit < 1 || limit > 50 {
        limit = 10
    }
    
    reports, err := h.binaMargaUseCase.GetBlockedRoads(c.Context(), limit)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve blocked roads", err)
    }

    return response.Success(c, "Blocked roads retrieved successfully", reports)
}

func (h *BinaMargaHandler) GetDamageByRoadType(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
    endDateStr := c.Query("end_date", time.Now().Format("2006-01-02"))
    
    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }
    
    if startDate.After(endDate) {
        return response.BadRequest(c, "Start date cannot be after end date", nil)
    }
    
    results, err := h.binaMargaUseCase.GetDamageByRoadType(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve damage statistics by road type", err)
    }

    return response.Success(c, "Damage statistics by road type retrieved successfully", map[string]interface{}{
        "data":       results,
        "start_date": startDate.Format("2006-01-02"),
        "end_date":   endDate.Format("2006-01-02"),
    })
}

func (h *BinaMargaHandler) GetDamageByLocation(c *fiber.Ctx) error {
    // Parse bounding box coordinates
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
    
    // Validate bounds
    if bounds["north"] <= bounds["south"] || bounds["east"] <= bounds["west"] {
        return response.BadRequest(c, "Invalid bounding box coordinates", nil)
    }
    
    results, err := h.binaMargaUseCase.GetDamageByLocation(c.Context(), bounds)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve damage statistics by location", err)
    }

    return response.Success(c, "Damage statistics by location retrieved successfully", map[string]interface{}{
        "data":   results,
        "bounds": bounds,
    })
}

func (h *BinaMargaHandler) GetDashboardSummary(c *fiber.Ctx) error {
    // Get overall statistics
    stats, err := h.binaMargaUseCase.GetStatistics(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve statistics", err)
    }

    // Get emergency reports
    emergencyReports, err := h.binaMargaUseCase.GetEmergencyReports(c.Context(), 5)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve emergency reports", err)
    }

    // Get blocked roads
    blockedRoads, err := h.binaMargaUseCase.GetBlockedRoads(c.Context(), 5)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve blocked roads", err)
    }

    summary := map[string]interface{}{
        "statistics":          stats,
        "emergency_reports":   emergencyReports,
        "blocked_roads":       blockedRoads,
        "generated_at":        time.Now(),
    }

    return response.Success(c, "Dashboard summary retrieved successfully", summary)
}

func (h *BinaMargaHandler) GetDashboard(c *fiber.Ctx) error {
    rawType := c.Query("road_type", "ALL")
    startStr := c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
    endStr := c.Query("end_date", time.Now().Format("2006-01-02"))

    // Normalisasi road_type (ALL/all → "")
    roadTypeNorm, ok := normalizeRoadTypeParam(rawType)
    if !ok {
        return response.BadRequest(c, "invalid road_type", fmt.Errorf("use: all, JALAN_NASIONAL, JALAN_PROVINSI, JALAN_KABUPATEN, JALAN_DESA"))
    }

    startDate, err := time.Parse("2006-01-02", startStr)
    if err != nil {
        return response.BadRequest(c, "invalid start_date format, use YYYY-MM-DD", err)
    }
    endDate, err := time.Parse("2006-01-02", endStr)
    if err != nil {
        return response.BadRequest(c, "invalid end_date format, use YYYY-MM-DD", err)
    }

    data, err := h.binaMargaUseCase.GetDashboard(c.Context(), roadTypeNorm, startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to build bina marga dashboard", err)
    }
    return response.Success(c, "Bina Marga dashboard", data)
}

func (h *BinaMargaHandler) GetBinaMargaOverview(c *fiber.Ctx) error {
    raw := c.Query("road_type", "all")

    // Normalisasi
    roadTypeNorm, ok := normalizeRoadTypeParam(raw)
    if !ok {
        return response.BadRequest(c, "Invalid road type",
            fmt.Errorf("road_type must be one of: all, JALAN_NASIONAL, JALAN_PROVINSI, JALAN_KABUPATEN, JALAN_DESA"))
    }

    overview, err := h.binaMargaUseCase.GetBinaMargaOverview(c.Context(), roadTypeNorm)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve bina marga overview", err)
    }

    return response.Success(c, "Bina marga overview retrieved successfully", overview)
}

// Helper functions untuk safe type conversion (reuse dari water resources)
func safeInt64(value interface{}) int64 {
    switch v := value.(type) {
    case int64:
        return v
    case float64:
        return int64(v)
    case int:
        return int64(v)
    case int32:
        return int64(v)
    default:
        return 0
    }
}

func safeFloat64(value interface{}) float64 {
    switch v := value.(type) {
    case float64:
        return v
    case float32:
        return float64(v)
    case int64:
        return float64(v)
    case int:
        return float64(v)
    default:
        return 0.0
    }
}

// ---- letakkan di file yang sama (di bawah type handler) ----

// Field "lokasi" → pakai NormalizeLocation
var bmLocationKeys = map[string]bool{
    "reporter_name": true,
    "road_name":     true,
    "bridge_name":   true,
}

// Field "enum/kategori" → pakai NormalizeEnum
var bmEnumKeys = map[string]bool{
    "institution_unit":    true,
    "road_type":           true,
    "road_class":          true,
    "pavement_type":       true,
    "damage_type":         true,
    "damage_level":        true,
    "bridge_structure_type": true,
    "bridge_damage_type":  true,
    "bridge_damage_level": true,
    "traffic_condition":   true,
    "traffic_impact":      true,
    "urgency_level":       true,
    "status":              true,
}

// Normalisasi semua filter string (biarkan tanggal/angka/bool apa adanya)
func normalizeBinaMargaFilters(filters map[string]interface{}) {
    for k, v := range filters {
        s, ok := v.(string)
        if !ok || s == "" {
            continue
        }
        if bmLocationKeys[k] {
            filters[k] = utils.NormalizeLocation(s)
            continue
        }
        if bmEnumKeys[k] {
            filters[k] = utils.NormalizeEnum(s)
            continue
        }
    }
}

// Normalisasi road_type dengan dukungan "all"
func normalizeRoadTypeParam(rt string) (queryValue string, ok bool) {
    n := utils.NormalizeEnum(rt)
    if n == "" || n == "all" {
        return "", true // kosongkan untuk "semua"
    }
    // valid values (setir sesuai domain kamu)
    valid := map[string]bool{
        "jalan_nasional":  true,
        "jalan_provinsi":  true,
        "jalan_kabupaten": true,
        "jalan_desa":      true,
    }
    if !valid[n] {
        return "", false
    }
    return n, true
}
