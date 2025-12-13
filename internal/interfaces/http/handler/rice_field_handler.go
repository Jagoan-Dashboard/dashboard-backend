package handler

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

const (
	MaxImportFileSize  = 10 << 20 // 10 MB
	MaxImportRows      = 10000
	RiceFieldSheetName = "lahan_pengairan"
	ImportTimeout      = 30 * time.Second
)

type RiceFieldHandler struct {
	useCase *usecase.RiceFieldUseCase
}

func NewRiceFieldHandler(useCase *usecase.RiceFieldUseCase) *RiceFieldHandler {
	return &RiceFieldHandler{useCase: useCase}
}

// Create handles rice field creation
func (h *RiceFieldHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateRiceFieldRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	ctx := c.Context()
	response, err := h.useCase.Create(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create rice field",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Rice field created successfully",
		"data":    response,
	})
}

// GetByID handles retrieving a single rice field
func (h *RiceFieldHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID is required",
		})
	}

	ctx := c.Context()
	response, err := h.useCase.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Rice field not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetAll handles retrieving paginated rice fields
func (h *RiceFieldHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	// Build filters
	filters := make(map[string]interface{})
	if district := c.Query("district"); district != "" {
		filters["district"] = district
	}
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = t
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = t
		}
	}

	ctx := c.Context()
	response, err := h.useCase.GetAll(ctx, page, perPage, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get rice fields",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// Update handles rice field update
func (h *RiceFieldHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID is required",
		})
	}

	var req dto.UpdateRiceFieldRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	ctx := c.Context()
	response, err := h.useCase.Update(ctx, id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update rice field",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rice field updated successfully",
		"data":    response,
	})
}

// Delete handles rice field deletion
func (h *RiceFieldHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID is required",
		})
	}

	ctx := c.Context()
	if err := h.useCase.Delete(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete rice field",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Rice field deleted successfully",
	})
}

// GetStatistics handles statistics retrieval
func (h *RiceFieldHandler) GetStatistics(c *fiber.Ctx) error {
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02")))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date", time.Now().Format("2006-01-02")))

	ctx := c.Context()
	stats, err := h.useCase.GetStatistics(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get statistics",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// GetAnalysis handles analysis retrieval
func (h *RiceFieldHandler) GetAnalysis(c *fiber.Ctx) error {
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02")))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date", time.Now().Format("2006-01-02")))

	ctx := c.Context()

	// Get statistics
	stats, err := h.useCase.GetStatistics(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get statistics",
			"error":   err.Error(),
		})
	}

	// Get distribution by district
	distribution, err := h.useCase.GetDistributionByDistrict(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get distribution",
			"error":   err.Error(),
		})
	}

	// Get map data
	mapData, err := h.useCase.GetMapData(ctx, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get map data",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"statistics":         stats,
			"distribution":       distribution,
			"map_data":          mapData,
			"date_range": fiber.Map{
				"start": startDate.Format("2006-01-02"),
				"end":   endDate.Format("2006-01-02"),
			},
		},
	})
}

// ============================================================================
// EXPORT - Generate Excel Template with Data
// ============================================================================

func (h *RiceFieldHandler) ExportRiceFields(c *fiber.Ctx) error {
	// Parse query parameters for filtering
	district := c.Query("district")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	yearStr := c.Query("year")

	// Build filters
	filters := make(map[string]interface{})

	if district != "" {
		filters["district"] = district
	}

	var startDate, endDate time.Time
	var err error

	if yearStr != "" {
		year, _ := strconv.Atoi(yearStr)
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
		filters["start_date"] = startDate
		filters["end_date"] = endDate
	} else {
		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err == nil {
				filters["start_date"] = startDate
			}
		}
		if endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err == nil {
				filters["end_date"] = endDate
			}
		}
	}

	// Fetch data using usecase
	ctx := c.Context()
	result, err := h.useCase.GetAll(ctx, 1, 10000, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch rice field data",
			"error":   err.Error(),
		})
	}

	// Create Excel file
	f := excelize.NewFile()
	defer f.Close()

	sheetName := RiceFieldSheetName
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create Excel sheet",
			"error":   err.Error(),
		})
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Set headers with exact names from template
	headers := []string{
		"Kecamatan",                                       // A
		"Koordinat Lokasi",                                // B
		"Tahun",                                           // C
		"Luas Sawah Irigasi",                              // D
		"Luas Sawah Tadah Hujan",                          // E
		"Luas Lahan Sawah ",                               // F (note the space!)
		"Luas Lahan Tegal/Kebun",                          // G
		"Luas Lahan Ladang/Huma",                          // H
		"Luas Lahan yang Sementara Tidak Diusahakan",     // I
		"Luas Lahan Bukan Sawah",                          // J
		"Total Luas Lahan",                                // K
	}

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style headers
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	f.SetCellStyle(sheetName, "A1", "K1", headerStyle)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 20) // Kecamatan
	f.SetColWidth(sheetName, "B", "B", 25) // Koordinat
	f.SetColWidth(sheetName, "C", "C", 10) // Tahun
	f.SetColWidth(sheetName, "D", "K", 18) // All area columns

	// Write data
	for i, rf := range result.RiceFields {
		row := i + 2

		// Format coordinates
		koordinat := ""
		if rf.Latitude != 0 && rf.Longitude != 0 {
			koordinat = fmt.Sprintf("%.6f, %.6f", rf.Latitude, rf.Longitude)
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), rf.District)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), koordinat)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), rf.Year)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), rf.IrrigatedRiceFields)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), rf.RainfedRiceFields)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), rf.TotalRiceFieldArea)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), rf.DryfieldArea)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), rf.ShiftingCultivationArea)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), rf.UnusedLandArea)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), rf.TotalNonRiceFieldArea)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), rf.TotalLandArea)
	}

	// Generate filename
	filename := "lahan_pengairan_template.xlsx"
	if district != "" {
		filename = fmt.Sprintf("lahan_pengairan_%s.xlsx", strings.ToLower(district))
	} else if yearStr != "" {
		filename = fmt.Sprintf("lahan_pengairan_%s.xlsx", yearStr)
	}

	// Set response headers
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Write to response
	return f.Write(c.Response().BodyWriter())
}

// ============================================================================
// IMPORT - Process Excel Upload
// ============================================================================

type ImportRiceFieldResult struct {
	TotalRows    int                `json:"total_rows"`
	SuccessCount int                `json:"success_count"`
	FailedCount  int                `json:"failed_count"`
	SkippedCount int                `json:"skipped_count"`
	FailedRows   []ImportFailedRow  `json:"failed_rows,omitempty"`
	SkippedRows  []ImportSkippedRow `json:"skipped_rows,omitempty"`
}

type ImportFailedRow struct {
	RowNumber int      `json:"row_number"`
	Data      string   `json:"data"`
	Errors    []string `json:"errors"`
}

type ImportSkippedRow struct {
	RowNumber int    `json:"row_number"`
	Reason    string `json:"reason"`
}

type RiceFieldImportData struct {
	Kecamatan                string
	Latitude                 float64
	Longitude                float64
	Tahun                    int
	LuasSawahIrigasi         float64
	LuasSawahTadahHujan      float64
	LuasLahanSawah           float64
	LuasLahanTegal           float64
	LuasLahanLadang          float64
	LuasLahanTidakDiusahakan float64
	LuasLahanBukanSawah      float64
	TotalLuasLahan           float64
}

func (h *RiceFieldHandler) ImportRiceFields(c *fiber.Ctx) error {
	// Validate file upload
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
		})
	}

	// Validate file size
	if file.Size > MaxImportFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("File size exceeds maximum limit of %d MB", MaxImportFileSize/(1<<20)),
		})
	}

	// Validate file type
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid file type. Only .xlsx files are allowed",
		})
	}

	// Open uploaded file
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open uploaded file",
			"error":   err.Error(),
		})
	}
	defer fileContent.Close()

	// Read Excel file
	f, err := excelize.OpenReader(fileContent)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to read Excel file",
			"error":   err.Error(),
		})
	}
	defer f.Close()

	// Process import with timeout
	ctx, cancel := context.WithTimeout(c.Context(), ImportTimeout)
	defer cancel()

	result := h.processRiceFieldImport(ctx, f)

	// Build response message
	message := fmt.Sprintf("Import completed. Success: %d, Failed: %d, Skipped: %d",
		result.SuccessCount, result.FailedCount, result.SkippedCount)

	statusCode := fiber.StatusOK
	if result.FailedCount > 0 || result.SkippedCount > 0 {
		statusCode = fiber.StatusPartialContent
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": result.SuccessCount > 0,
		"message": message,
		"data":    result,
	})
}

func (h *RiceFieldHandler) processRiceFieldImport(ctx context.Context, f *excelize.File) *ImportRiceFieldResult {
	result := &ImportRiceFieldResult{
		FailedRows:  make([]ImportFailedRow, 0),
		SkippedRows: make([]ImportSkippedRow, 0),
	}

	// Get rows from sheet
	rows, err := f.GetRows(RiceFieldSheetName)
	if err != nil {
		result.FailedRows = append(result.FailedRows, ImportFailedRow{
			RowNumber: 0,
			Data:      "Sheet validation",
			Errors:    []string{fmt.Sprintf("Failed to get rows from sheet '%s': %v", RiceFieldSheetName, err)},
		})
		return result
	}

	if len(rows) == 0 {
		result.FailedRows = append(result.FailedRows, ImportFailedRow{
			RowNumber: 0,
			Data:      "Sheet validation",
			Errors:    []string{"Sheet is empty"},
		})
		return result
	}

	// Validate headers
	headers := rows[0]
	expectedHeaders := []string{
		"Kecamatan",
		"Koordinat Lokasi",
		"Tahun",
		"Luas Sawah Irigasi",
		"Luas Sawah Tadah Hujan",
		"Luas Lahan Sawah ", // Note: has trailing space in template!
	}

	for i, expected := range expectedHeaders {
		if i >= len(headers) || !strings.Contains(headers[i], strings.TrimSpace(expected)) {
			result.FailedRows = append(result.FailedRows, ImportFailedRow{
				RowNumber: 1,
				Data:      "Header validation",
				Errors:    []string{fmt.Sprintf("Invalid template: missing or incorrect column '%s'", expected)},
			})
			return result
		}
	}

	// Process data rows
	dataRows := rows[1:]
	result.TotalRows = len(dataRows)

	for i, row := range dataRows {
		rowNumber := i + 2 // Excel row number (header is row 1)

		// Skip empty rows
		if isEmptyRow(row) {
			result.SkippedCount++
			result.SkippedRows = append(result.SkippedRows, ImportSkippedRow{
				RowNumber: rowNumber,
				Reason:    "Empty row",
			})
			continue
		}

		// Parse row data
		importData, errors := parseRiceFieldRow(row, rowNumber)
		if len(errors) > 0 {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, ImportFailedRow{
				RowNumber: rowNumber,
				Data:      rowToString(row),
				Errors:    errors,
			})
			continue
		}

		// Convert to CreateRequest
		createReq := riceFieldImportToCreateRequest(importData)

		// Save using usecase
		_, err := h.useCase.Create(ctx, createReq)
		if err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, ImportFailedRow{
				RowNumber: rowNumber,
				Data:      rowToString(row),
				Errors:    []string{fmt.Sprintf("Database error: %v", err)},
			})
			continue
		}

		result.SuccessCount++
	}

	return result
}

func parseRiceFieldRow(row []string, rowNumber int) (*RiceFieldImportData, []string) {
	errors := make([]string, 0)
	data := &RiceFieldImportData{}

	// Helper function to get cell value safely
	getCell := func(index int) string {
		if index < len(row) {
			return strings.TrimSpace(row[index])
		}
		return ""
	}

	// Column A: Kecamatan (REQUIRED)
	data.Kecamatan = getCell(0)
	if data.Kecamatan == "" {
		errors = append(errors, "Kecamatan is required")
	}

	// Column B: Koordinat Lokasi (OPTIONAL)
	koordinat := getCell(1)
	if koordinat != "" {
		coords := strings.Split(koordinat, ",")
		if len(coords) == 2 {
			lat, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
			lng, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
			if err1 == nil && err2 == nil {
				data.Latitude = lat
				data.Longitude = lng
			} else {
				errors = append(errors, "Invalid coordinate format. Use: latitude, longitude")
			}
		} else {
			errors = append(errors, "Invalid coordinate format. Use: latitude, longitude")
		}
	}

	// Column C: Tahun (OPTIONAL, default to current year)
	tahunStr := getCell(2)
	if tahunStr != "" {
		tahun, err := strconv.Atoi(tahunStr)
		if err != nil {
			errors = append(errors, "Invalid year format. Must be a number")
		} else {
			data.Tahun = tahun
		}
	} else {
		data.Tahun = time.Now().Year()
	}

	// Helper function to parse float with comma handling
	parseFloat := func(value string) (float64, error) {
		value = strings.ReplaceAll(value, ",", "")
		if value == "" {
			return 0, nil
		}
		return strconv.ParseFloat(value, 64)
	}

	// Parse all area columns
	if val, err := parseFloat(getCell(3)); err == nil {
		data.LuasSawahIrigasi = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Sawah Irigasi: %v", err))
	}

	if val, err := parseFloat(getCell(4)); err == nil {
		data.LuasSawahTadahHujan = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Sawah Tadah Hujan: %v", err))
	}

	if val, err := parseFloat(getCell(5)); err == nil {
		data.LuasLahanSawah = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Lahan Sawah: %v", err))
	}

	if val, err := parseFloat(getCell(6)); err == nil {
		data.LuasLahanTegal = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Lahan Tegal/Kebun: %v", err))
	}

	if val, err := parseFloat(getCell(7)); err == nil {
		data.LuasLahanLadang = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Lahan Ladang/Huma: %v", err))
	}

	if val, err := parseFloat(getCell(8)); err == nil {
		data.LuasLahanTidakDiusahakan = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Lahan Tidak Diusahakan: %v", err))
	}

	if val, err := parseFloat(getCell(9)); err == nil {
		data.LuasLahanBukanSawah = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Luas Lahan Bukan Sawah: %v", err))
	}

	if val, err := parseFloat(getCell(10)); err == nil {
		data.TotalLuasLahan = val
	} else {
		errors = append(errors, fmt.Sprintf("Invalid Total Luas Lahan: %v", err))
	}

	return data, errors
}

func riceFieldImportToCreateRequest(data *RiceFieldImportData) *dto.CreateRiceFieldRequest {
	// Create date from year
	date := time.Date(data.Tahun, 1, 1, 0, 0, 0, 0, time.UTC)

	return &dto.CreateRiceFieldRequest{
		District:                data.Kecamatan,
		Latitude:                data.Latitude,
		Longitude:               data.Longitude,
		Date:                    date,
		Year:                    data.Tahun,
		IrrigatedRiceFields:     data.LuasSawahIrigasi,
		RainfedRiceFields:       data.LuasSawahTadahHujan,
		TotalRiceFieldArea:      data.LuasLahanSawah,
		DryfieldArea:            data.LuasLahanTegal,
		ShiftingCultivationArea: data.LuasLahanLadang,
		UnusedLandArea:          data.LuasLahanTidakDiusahakan,
		TotalNonRiceFieldArea:   data.LuasLahanBukanSawah,
		TotalLandArea:           data.TotalLuasLahan,
		DataSource:              "import",
	}
}

// Helper functions
func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func rowToString(row []string) string {
	parts := make([]string, 0)
	for i, cell := range row {
		if i > 10 { // Only show first 11 columns
			break
		}
		if cell != "" {
			parts = append(parts, cell)
		}
	}
	return strings.Join(parts, " | ")
}