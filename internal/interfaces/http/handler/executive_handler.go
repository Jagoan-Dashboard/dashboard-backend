
package handler

import (
    "building-report-backend/internal/application/usecase"
    "building-report-backend/internal/interfaces/response"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
)

type ExecutiveHandler struct {
    executiveUseCase *usecase.ExecutiveUseCase
}

func NewExecutiveHandler(execUseCase *usecase.ExecutiveUseCase) *ExecutiveHandler {
    return &ExecutiveHandler{
        executiveUseCase: execUseCase,
    }
}

func (h *ExecutiveHandler) GetEkonomiOverview(c *fiber.Ctx) error {
    tahunStr := c.Query("year")
    if tahunStr == "" {
        return response.BadRequest(c, "Parameter tahun wajib diisi", nil)
    }

    tahun, err := strconv.Atoi(tahunStr)
    if err != nil {
        return response.BadRequest(c, "Format tahun tidak valid", err)
    }

    if tahun < 1900 || tahun > 2100 {
        return response.BadRequest(c, "Tahun tidak valid", nil)
    }

    result, err := h.executiveUseCase.GetEkonomiOverview(c.Context(), tahun)
    if err != nil {
        return response.InternalError(c, "Gagal mengambil data ekonomi", err)
    }

    return response.Success(c, "Data ekonomi berhasil diambil", result)
}

func (h *ExecutiveHandler) GetPopulationOverview(c *fiber.Ctx) error {
    tahunStr := c.Query("year")
    if tahunStr == "" {
        return response.BadRequest(c, "Parameter year wajib diisi", nil)
    }

    tahun, err := strconv.Atoi(tahunStr)
    if err != nil {
        return response.BadRequest(c, "Format year tidak valid", err)
    }

    if tahun < 1900 || tahun > 2100 {
        return response.BadRequest(c, "Year tidak valid", nil)
    }

    result, err := h.executiveUseCase.GetPopulationOverview(c.Context(), tahun)
    if err != nil {
        return response.InternalError(c, "Gagal mengambil data demografi", err)
    }

    return response.Success(c, "Data demografi berhasil diambil", result)
}

func (h *ExecutiveHandler) GetPovertyOverview(c *fiber.Ctx) error {
    tahunStr := c.Query("year")
    if tahunStr == "" {
        return response.BadRequest(c, "Parameter year wajib diisi", nil)
    }

    tahun, err := strconv.Atoi(tahunStr)
    if err != nil {
        return response.BadRequest(c, "Format year tidak valid", err)
    }

    if tahun < 1900 || tahun > 2100 {
        return response.BadRequest(c, "Year tidak valid", nil)
    }

    result, err := h.executiveUseCase.GetPovertyOverview(c.Context(), tahun)
    if err != nil {
        return response.InternalError(c, "Gagal mengambil data kemiskinan", err)
    }

    return response.Success(c, "Data kemiskinan berhasil diambil", result)
}


func (h *ExecutiveHandler) GetEmploymentOverview(c *fiber.Ctx) error {
    tahunStr := c.Query("year")
    if tahunStr == "" {
        return response.BadRequest(c, "Parameter year wajib diisi", nil)
    }

    tahun, err := strconv.Atoi(tahunStr)
    if err != nil {
        return response.BadRequest(c, "Format year tidak valid", err)
    }

    if tahun < 1900 || tahun > 2100 {
        return response.BadRequest(c, "Year tidak valid", nil)
    }

    result, err := h.executiveUseCase.GetEmploymentOverview(c.Context(), tahun)
    if err != nil {
        return response.InternalError(c, "Gagal mengambil data ketenagakerjaan", err)
    }

    return response.Success(c, "Data ketenagakerjaan berhasil diambil", result)
}

func (h *ExecutiveHandler) GetEducationOverview(c *fiber.Ctx) error {
    tahunStr := c.Query("year")
    if tahunStr == "" {
        return response.BadRequest(c, "Parameter year wajib diisi", nil)
    }

    tahun, err := strconv.Atoi(tahunStr)
    if err != nil {
        return response.BadRequest(c, "Format year tidak valid", err)
    }

    if tahun < 1900 || tahun > 2100 {
        return response.BadRequest(c, "Year tidak valid", nil)
    }

    result, err := h.executiveUseCase.GetEducationOverview(c.Context(), tahun)
    if err != nil {
        return response.InternalError(c, "Gagal mengambil data pendidikan", err)
    }

    return response.Success(c, "Data pendidikan berhasil diambil", result)
}