package handler

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"
	"building-report-backend/pkg/utils"
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

	
	req.Normalize()

	if err := req.Validate(); err != nil {
		return response.ValidationError(c, err)
	}

	

	form, err := c.MultipartForm()
	if err != nil {
		return response.BadRequest(c, "Failed to parse multipart form", err)
	}

	photos := form.File["photos"]
	if len(photos) < 2 {
		return response.BadRequest(c, "Minimum 2 photos required", nil)
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

	normalizeReportFilters(filters)

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

	
	req.Normalize()

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
	raw := c.Query("building_type", "all")

	bt, ok := normalizeBuildingTypeParam(raw)
	if !ok {
		return response.BadRequest(
			c,
			"Invalid building type",
			fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"),
		)
	}

	overview, err := h.reportUseCase.GetTataBangunanOverview(c.Context(), bt)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve tata bangunan overview", err)
	}

	return response.Success(c, "Tata bangunan overview retrieved successfully", overview)
}


func (h *ReportHandler) GetBasicStatistics(c *fiber.Ctx) error {
	raw := c.Query("building_type", "all")

	bt, ok := normalizeBuildingTypeParam(raw)
	if !ok {
		return response.BadRequest(
			c,
			"Invalid building type",
			fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"),
		)
	}

	stats, err := h.reportUseCase.GetBasicStatistics(c.Context(), bt)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve basic statistics", err)
	}

	return response.Success(c, "Basic statistics retrieved successfully", stats)
}


func (h *ReportHandler) GetLocationDistribution(c *fiber.Ctx) error {
	raw := c.Query("building_type", "all")

	bt, ok := normalizeBuildingTypeParam(raw)
	if !ok {
		return response.BadRequest(
			c,
			"Invalid building type",
			fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"),
		)
	}

	locations, err := h.reportUseCase.GetLocationDistribution(c.Context(), bt)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve location distribution", err)
	}

	return response.Success(c, "Location distribution retrieved successfully", locations)
}


func (h *ReportHandler) GetWorkTypeStatistics(c *fiber.Ctx) error {
	raw := c.Query("building_type", "all")

	bt, ok := normalizeBuildingTypeParam(raw)
	if !ok {
		return response.BadRequest(
			c,
			"Invalid building type",
			fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"),
		)
	}

	workTypeStats, err := h.reportUseCase.GetWorkTypeStatistics(c.Context(), bt)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve work type statistics", err)
	}

	return response.Success(c, "Work type statistics retrieved successfully", workTypeStats)
}

func (h *ReportHandler) GetConditionAfterRehabStatistics(c *fiber.Ctx) error {
	raw := c.Query("building_type", "all")

	bt, ok := normalizeBuildingTypeParam(raw)
	if !ok {
		return response.BadRequest(
			c,
			"Invalid building type",
			fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"),
		)
	}

	conditionStats, err := h.reportUseCase.GetConditionAfterRehabStatistics(c.Context(), bt)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve condition after rehab statistics", err)
	}

	return response.Success(c, "Condition after rehab statistics retrieved successfully", conditionStats)
}

func (h *ReportHandler) GetStatusStatistics(c *fiber.Ctx) error {
	raw := c.Query("building_type", "all")

	bt, ok := normalizeBuildingTypeParam(raw)
	if !ok {
		return response.BadRequest(
			c,
			"Invalid building type",
			fmt.Errorf("building_type must be one of: all, SEKOLAH, PUSKESMAS_POSYANDU, PASAR, SARANA_OLAHRAGA, KANTOR_PEMERINTAH, FASILITAS_UMUM, LAINNYA"),
		)
	}

	statusStats, err := h.reportUseCase.GetStatusStatistics(c.Context(), bt)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve status statistics", err)
	}

	return response.Success(c, "Status statistics retrieved successfully", statusStats)
}

func (h *ReportHandler) GetBuildingTypeDistribution(c *fiber.Ctx) error {
	buildingTypeStats, err := h.reportUseCase.GetBuildingTypeDistribution(c.Context())
	if err != nil {
		return response.InternalError(c, "Failed to retrieve building type distribution", err)
	}

	return response.Success(c, "Building type distribution retrieved successfully", buildingTypeStats)
}



var reportLocationKeys = map[string]bool{
	"village":  true,
	"district": true,
}


var reportEnumKeys = map[string]bool{
	"building_type": true,
	"report_status": true,
}

func normalizeReportFilters(filters map[string]interface{}) {
	for k, v := range filters {
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		if reportLocationKeys[k] {
			filters[k] = utils.NormalizeLocation(s)
			continue
		}
		if reportEnumKeys[k] {
			filters[k] = utils.NormalizeEnum(s)
			continue
		}
	}
}


func normalizeBuildingTypeParam(bt string) (string, bool) {
	n := utils.NormalizeEnum(bt)
	if n == "" || n == "all" {
		return "", true
	}
	valid := map[string]bool{
		"sekolah":              true,
		"puskesmas_posyandu":   true,
		"pasar":                true,
		"sarana_olahraga":      true,
		"kantor_pemerintah":    true,
		"fasilitas_umum":       true,
		"lainnya":              true,
	}
	return n, valid[n]
}
