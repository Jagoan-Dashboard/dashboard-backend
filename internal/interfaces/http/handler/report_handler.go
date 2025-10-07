package handler

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	
)

type ReportHandler struct {
    reportUseCase *usecase.ReportUseCase
}

func NewReportHandler(reportUseCase *usecase.ReportUseCase) *ReportHandler {
    return &ReportHandler{
        reportUseCase: reportUseCase,
    }
}

func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
    var req dto.CreateReportRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
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

    // Photos are required only for rehabilitation reports, not for new construction
    isPembangunanBaru := req.ReportStatus == "PEMBANGUNAN_BARU"
    if !isPembangunanBaru && len(photos) < 2 {
        return response.BadRequest(c, "Minimum 2 photos required for rehabilitation reports", nil)
    }

    report, err := h.reportUseCase.CreateReport(c.Context(), &req, photos)
    if err != nil {
        return response.InternalError(c, "Failed to create report", err)
    }

    return response.Success(c, "Report created successfully", report)
}

func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
    id := c.Params("id")
  

    report, err := h.reportUseCase.GetReport(c.Context(), id)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *ReportHandler) ListReports(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    filters := map[string]interface{}{
        "village":       c.Query("village"),
        "district":      c.Query("district"),
        "building_type": c.Query("building_type"),
        "report_status": c.Query("report_status"),
    }

    result, err := h.reportUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports", err)
    }

    return response.Success(c, "Reports retrieved successfully", result)
}

func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error {
    id := c.Params("id")
   

    var req dto.UpdateReportRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    userID := c.Locals("userID").(string)

    report, err := h.reportUseCase.UpdateReport(c.Context(), id, &req, userID)
    if err != nil {
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *ReportHandler) DeleteReport(c *fiber.Ctx) error {
    id := c.Params("id")
    

    userID := c.Locals("userID").(string)

    if err := h.reportUseCase.DeleteReport(c.Context(), id, userID); err != nil {
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *ReportHandler) GetTataBangunanOverview(c *fiber.Ctx) error {
    buildingType := c.Query("building_type", "all") // all, SEKOLAH, PUSKESMAS, PASAR, etc.
    
    // Validate building type
    validBuildingTypes := map[string]bool{
        "all":                    true,
        "SEKOLAH":               true,
        "PUSKESMAS_POSYANDU":    true,
        "PASAR":                 true,
        "SARANA_OLAHRAGA":       true,
        "KANTOR_PEMERINTAH":     true,
        "FASILITAS_UMUM":        true,
        "LAINNYA":               true,
    }
    
    if !validBuildingTypes[buildingType] {
        return response.BadRequest(c, "Invalid building type", fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"))
    }
    
    // Convert "all" to empty string for repository layer
    queryBuildingType := buildingType
    if buildingType == "all" {
        queryBuildingType = ""
    }

    overview, err := h.reportUseCase.GetTataBangunanOverview(c.Context(), queryBuildingType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve tata bangunan overview", err)
    }

    return response.Success(c, "Tata bangunan overview retrieved successfully", overview)
}

// GetBasicStatistics handles basic statistics endpoint
func (h *ReportHandler) GetBasicStatistics(c *fiber.Ctx) error {
    buildingType := c.Query("building_type", "all")
    
    if buildingType == "all" {
        buildingType = ""
    }

    stats, err := h.reportUseCase.GetBasicStatistics(c.Context(), buildingType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve basic statistics", err)
    }

    return response.Success(c, "Basic statistics retrieved successfully", stats)
}

// GetLocationDistribution handles location distribution for mapping
func (h *ReportHandler) GetLocationDistribution(c *fiber.Ctx) error {
    buildingType := c.Query("building_type", "all")
    
    if buildingType == "all" {
        buildingType = ""
    }

    locations, err := h.reportUseCase.GetLocationDistribution(c.Context(), buildingType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve location distribution", err)
    }

    return response.Success(c, "Location distribution retrieved successfully", locations)
}

// GetWorkTypeStatistics handles work type statistics
func (h *ReportHandler) GetWorkTypeStatistics(c *fiber.Ctx) error {
    buildingType := c.Query("building_type", "all")
    
    if buildingType == "all" {
        buildingType = ""
    }

    workTypeStats, err := h.reportUseCase.GetWorkTypeStatistics(c.Context(), buildingType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve work type statistics", err)
    }

    return response.Success(c, "Work type statistics retrieved successfully", workTypeStats)
}

// GetConditionAfterRehabStatistics handles condition after rehab statistics
func (h *ReportHandler) GetConditionAfterRehabStatistics(c *fiber.Ctx) error {
    buildingType := c.Query("building_type", "all")
    
    if buildingType == "all" {
        buildingType = ""
    }

    conditionStats, err := h.reportUseCase.GetConditionAfterRehabStatistics(c.Context(), buildingType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve condition after rehab statistics", err)
    }

    return response.Success(c, "Condition after rehab statistics retrieved successfully", conditionStats)
}

// GetStatusStatistics handles status statistics
func (h *ReportHandler) GetStatusStatistics(c *fiber.Ctx) error {
    buildingType := c.Query("building_type", "all")
    
    if buildingType == "all" {
        buildingType = ""
    }

    statusStats, err := h.reportUseCase.GetStatusStatistics(c.Context(), buildingType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve status statistics", err)
    }

    return response.Success(c, "Status statistics retrieved successfully", statusStats)
}

// GetBuildingTypeDistribution handles building type distribution
func (h *ReportHandler) GetBuildingTypeDistribution(c *fiber.Ctx) error {
    buildingTypeStats, err := h.reportUseCase.GetBuildingTypeDistribution(c.Context())
    if err != nil {
        return response.InternalError(c, "Failed to retrieve building type distribution", err)
    }

    return response.Success(c, "Building type distribution retrieved successfully", buildingTypeStats)
}