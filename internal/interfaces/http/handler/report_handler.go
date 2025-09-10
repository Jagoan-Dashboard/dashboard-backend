// internal/interfaces/http/handler/report_handler.go
package handler

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

    // Validate request
    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    // Get user ID from context (set by auth middleware)
    userID := c.Locals("userID").(uuid.UUID)

    // Handle file uploads
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequest(c, "Failed to parse multipart form", err)
    }

    photos := form.File["photos"]
    if len(photos) < 2 {
        return response.BadRequest(c, "Minimum 2 photos required", nil)
    }

    report, err := h.reportUseCase.CreateReport(c.Context(), &req, photos, userID)
    if err != nil {
        return response.InternalError(c, "Failed to create report", err)
    }

    return response.Success(c, "Report created successfully", report)
}

func (h *ReportHandler) GetReport(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

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
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

    var req dto.UpdateReportRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    userID := c.Locals("userID").(uuid.UUID)

    report, err := h.reportUseCase.UpdateReport(c.Context(), id, &req, userID)
    if err != nil {
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}

func (h *ReportHandler) DeleteReport(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return response.BadRequest(c, "Invalid report ID", err)
    }

    userID := c.Locals("userID").(uuid.UUID)

    if err := h.reportUseCase.DeleteReport(c.Context(), id, userID); err != nil {
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}