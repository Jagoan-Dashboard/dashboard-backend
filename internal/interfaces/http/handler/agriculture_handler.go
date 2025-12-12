package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/interfaces/response"
	"building-report-backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type AgricultureHandler struct {
    agricultureUseCase *usecase.AgricultureUseCase
}

func (h *AgricultureHandler) ExportAlatPertanian(c *fiber.Ctx) error {
	// Parse query parameters
	district := c.Query("district", "")
	if district != "" {
		district = utils.NormalizeLocation(district)
	}

	equipmentType := strings.ToUpper(c.Query("equipment_type", ""))
	if equipmentType != "" && equipmentType != "PENGOLAH_GABAH" && equipmentType != "PERONTOK" && equipmentType != "MESIN" && equipmentType != "POMPA_AIR" {
		return response.BadRequest(c, "Invalid equipment_type. Must be one of: PENGOLAH_GABAH, PERONTOK, MESIN, POMPA_AIR", nil)
	}

	// Parse date range
	var startDate, endDate time.Time
	var err error

	startDateStr := c.Query("start_date", "")
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
		}
	}

	endDateStr := c.Query("end_date", "")
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
		}
	}

	// Generate Excel file
	excelData, err := h.exportAlatPertanianToExcel(c.Context(), district, equipmentType, startDate, endDate)
	if err != nil {
		return response.InternalError(c, "Failed to export data", err)
	}

	// Generate filename
	filename := h.generateAlatPertanianFilename(district, equipmentType, startDate, endDate)

	// Set response headers
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Set("Content-Length", fmt.Sprintf("%d", len(excelData)))

	return c.Send(excelData)
}

func (h *AgricultureHandler) exportAlatPertanianToExcel(ctx context.Context, district, equipmentType string, startDate, endDate time.Time) ([]byte, error) {
	// Prepare filters
	filters := make(map[string]interface{})

	if district != "" {
		filters["district"] = district
	}

	if !startDate.IsZero() {
		filters["start_date"] = startDate.Format("2006-01-02")
	}

	if !endDate.IsZero() {
		filters["end_date"] = endDate.Format("2006-01-02")
	}

	// Fetch all reports
	reports,  err := h.agricultureUseCase.ListReports(ctx, 1, 10000, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Aggregate data per district per year
	type EquipmentKey struct {
		District string
		Year     int
	}

	type EquipmentCount struct {
		Latitude           float64
		Longitude          float64
		PengolahGabah      int
		PerontokMultiguna  int
		MesinPertanian     int
		PompaAir           int
	}

	equipmentMap := make(map[EquipmentKey]*EquipmentCount)

	for _, report := range reports.Reports {
		year := report.VisitDate.Year()
		key := EquipmentKey{
			District: report.District,
			Year:     year,
		}

		if _, exists := equipmentMap[key]; !exists {
			equipmentMap[key] = &EquipmentCount{
				Latitude:  report.Latitude,
				Longitude: report.Longitude,
			}
		}

		// Map technology to equipment type
		technologies := []string{
			string(report.FoodTechnology),
			string(report.HortiTechnology),
			string(report.PlantationTechnology),
		}

		for _, tech := range technologies {
			if tech == "" || tech == "TIDAK_ADA" {
				continue
			}

			techUpper := strings.ToUpper(tech)

			// Mapping real dari technology field
			switch {
			case strings.Contains(techUpper, "PENGOLAH_GABAH"), strings.Contains(techUpper, "PENGOLAH GABAH"):
				equipmentMap[key].PengolahGabah++
			case strings.Contains(techUpper, "PERONTOK"):
				equipmentMap[key].PerontokMultiguna++
			case strings.Contains(techUpper, "POMPA_AIR"), strings.Contains(techUpper, "POMPA AIR"), strings.Contains(techUpper, "POMPA"):
				equipmentMap[key].PompaAir++
			case strings.Contains(techUpper, "MESIN"), strings.Contains(techUpper, "ALAT"), 
				 strings.Contains(techUpper, "TRAKTOR"), strings.Contains(techUpper, "CULTIVATOR"):
				equipmentMap[key].MesinPertanian++
			default:
				// Technology lain yang tidak masuk kategori spesifik → masuk ke Mesin/Peralatan
				if tech != "" {
					equipmentMap[key].MesinPertanian++
				}
			}
		}
	}

	// Transform to export data
	type ExportData struct {
		Kecamatan          string
		Koordinat          string
		Tahun              int
		PengolahGabah      int
		PerontokMultiguna  int
		MesinPertanian     int
		PompaAir           int
	}

	var exportData []ExportData

	for key, count := range equipmentMap {
		// Filter by equipment type if specified
		if equipmentType != "" {
			switch equipmentType {
			case "PENGOLAH_GABAH":
				if count.PengolahGabah == 0 {
					continue
				}
			case "PERONTOK":
				if count.PerontokMultiguna == 0 {
					continue
				}
			case "MESIN":
				if count.MesinPertanian == 0 {
					continue
				}
			case "POMPA_AIR":
				if count.PompaAir == 0 {
					continue
				}
			}
		}

		exportData = append(exportData, ExportData{
			Kecamatan:         key.District,
			Koordinat:         fmt.Sprintf("%.6f, %.6f", count.Latitude, count.Longitude),
			Tahun:             key.Year,
			PengolahGabah:     count.PengolahGabah,
			PerontokMultiguna: count.PerontokMultiguna,
			MesinPertanian:    count.MesinPertanian,
			PompaAir:          count.PompaAir,
		})
	}

	// Sort by district and year
	// Simple bubble sort for demonstration
	for i := 0; i < len(exportData)-1; i++ {
		for j := 0; j < len(exportData)-i-1; j++ {
			if exportData[j].Kecamatan > exportData[j+1].Kecamatan ||
				(exportData[j].Kecamatan == exportData[j+1].Kecamatan && exportData[j].Tahun > exportData[j+1].Tahun) {
				exportData[j], exportData[j+1] = exportData[j+1], exportData[j]
			}
		}
	}

	// Create Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	sheetName := "alat_pertanian"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Write headers
	headers := []string{
		"Kecamatan",
		"Koordinat Lokasi",
		"Tahun",
		"Jumlah Alat Pengolah Gabah",
		"Jumlah Alat Perontok Multiguna",
		"Jumlah Mesin/Peralatan Pertanian",
		"Jumlah Pompa Air",
	}

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Column widths
	f.SetColWidth(sheetName, "A", "A", 20) // Kecamatan
	f.SetColWidth(sheetName, "B", "B", 25) // Koordinat
	f.SetColWidth(sheetName, "C", "C", 10) // Tahun
	f.SetColWidth(sheetName, "D", "D", 25) // Pengolah Gabah
	f.SetColWidth(sheetName, "E", "E", 28) // Perontok Multiguna
	f.SetColWidth(sheetName, "F", "F", 30) // Mesin Pertanian
	f.SetColWidth(sheetName, "G", "G", 18) // Pompa Air

	f.SetRowHeight(sheetName, 1, 30)

	// Data style
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Calibri",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
		},
	})

	numberStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Calibri",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
		},
		NumFmt: 1, // 0 format for integers
	})

	// Write data
	for i, row := range exportData {
		rowNum := i + 2

		// Column A: Kecamatan
		cellA := fmt.Sprintf("A%d", rowNum)
		f.SetCellValue(sheetName, cellA, row.Kecamatan)
		f.SetCellStyle(sheetName, cellA, cellA, dataStyle)

		// Column B: Koordinat
		cellB := fmt.Sprintf("B%d", rowNum)
		f.SetCellValue(sheetName, cellB, row.Koordinat)
		f.SetCellStyle(sheetName, cellB, cellB, dataStyle)

		// Column C: Tahun
		cellC := fmt.Sprintf("C%d", rowNum)
		f.SetCellValue(sheetName, cellC, row.Tahun)
		f.SetCellStyle(sheetName, cellC, cellC, dataStyle)

		// Column D: Pengolah Gabah
		cellD := fmt.Sprintf("D%d", rowNum)
		f.SetCellValue(sheetName, cellD, row.PengolahGabah)
		f.SetCellStyle(sheetName, cellD, cellD, numberStyle)

		// Column E: Perontok Multiguna
		cellE := fmt.Sprintf("E%d", rowNum)
		f.SetCellValue(sheetName, cellE, row.PerontokMultiguna)
		f.SetCellStyle(sheetName, cellE, cellE, numberStyle)

		// Column F: Mesin Pertanian
		cellF := fmt.Sprintf("F%d", rowNum)
		f.SetCellValue(sheetName, cellF, row.MesinPertanian)
		f.SetCellStyle(sheetName, cellF, cellF, numberStyle)

		// Column G: Pompa Air
		cellG := fmt.Sprintf("G%d", rowNum)
		f.SetCellValue(sheetName, cellG, row.PompaAir)
		f.SetCellStyle(sheetName, cellG, cellG, numberStyle)
	}

	// Freeze header row
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Add autofilter
	if len(exportData) > 0 {
		lastRow := len(exportData) + 1
		f.AutoFilter(sheetName, fmt.Sprintf("A1:G%d", lastRow), nil)
	}

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

func (h *AgricultureHandler) generateAlatPertanianFilename(district, equipmentType string, startDate, endDate time.Time) string {
	timestamp := time.Now().Format("20060102_150405")

	var parts []string
	parts = append(parts, "export_alat_pertanian")

	if equipmentType != "" {
		parts = append(parts, strings.ToLower(strings.ReplaceAll(equipmentType, "_", "")))
	}

	if district != "" {
		parts = append(parts, strings.ToLower(strings.ReplaceAll(district, " ", "_")))
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		parts = append(parts, fmt.Sprintf("%s_to_%s", startDate.Format("20060102"), endDate.Format("20060102")))
	}

	parts = append(parts, timestamp)

	return strings.Join(parts, "_") + ".xlsx"
}


func NewAgricultureHandler(agricultureUseCase *usecase.AgricultureUseCase) *AgricultureHandler {
    return &AgricultureHandler{
        agricultureUseCase: agricultureUseCase,
    }
}

func (h *AgricultureHandler) CreateReport(c *fiber.Ctx) error {
    var req dto.CreateAgricultureRequest

    req.ExtensionOfficer = c.FormValue("extension_officer")
    req.FarmerName = c.FormValue("farmer_name")
    req.FarmerGroup = c.FormValue("farmer_group")
    req.FarmerGroupType = c.FormValue("farmer_group_type")
    req.Village = c.FormValue("village")
    req.District = c.FormValue("district")

    visitDateStr := c.FormValue("visit_date")
    if visitDate, err := time.Parse("2006-01-02", visitDateStr); err == nil {
        req.VisitDate = visitDate
    } else {
        return response.BadRequest(c, "Invalid visit date format, use YYYY-MM-DD", err)
    }

    if lat, err := strconv.ParseFloat(c.FormValue("latitude"), 64); err == nil {
        req.Latitude = lat
    } else if c.FormValue("latitude") != "" {
        return response.BadRequest(c, "Invalid latitude format", err)
    }
    if lng, err := strconv.ParseFloat(c.FormValue("longitude"), 64); err == nil {
        req.Longitude = lng
    } else if c.FormValue("longitude") != "" {
        return response.BadRequest(c, "Invalid longitude format", err)
    }

    
    req.FoodCommodity = c.FormValue("food_commodity")
    if req.FoodCommodity != "" {
        req.FoodLandStatus = c.FormValue("food_land_status")
        if area, err := strconv.ParseFloat(c.FormValue("food_land_area"), 64); err == nil && area > 0 {
            req.FoodLandArea = area
        }
        req.FoodGrowthPhase = c.FormValue("food_growth_phase")
        if age, err := strconv.Atoi(c.FormValue("food_plant_age")); err == nil && age >= 0 {
            req.FoodPlantAge = age
        }
        req.FoodPlantingDate = c.FormValue("food_planting_date")
        req.FoodHarvestDate = c.FormValue("food_harvest_date")
        req.FoodDelayReason = c.FormValue("food_delay_reason")
        req.FoodTechnology = c.FormValue("food_technology")
    }

    
    req.HortiCommodity = c.FormValue("horti_commodity")
    if req.HortiCommodity != "" {
        req.HortiSubCommodity = c.FormValue("horti_sub_commodity")
        req.HortiLandStatus = c.FormValue("horti_land_status")
        if area, err := strconv.ParseFloat(c.FormValue("horti_land_area"), 64); err == nil && area > 0 {
            req.HortiLandArea = area
        }
        req.HortiGrowthPhase = c.FormValue("horti_growth_phase")
        if age, err := strconv.Atoi(c.FormValue("horti_plant_age")); err == nil && age >= 0 {
            req.HortiPlantAge = age
        }
        req.HortiPlantingDate = c.FormValue("horti_planting_date")
        req.HortiHarvestDate = c.FormValue("horti_harvest_date")
        req.HortiDelayReason = c.FormValue("horti_delay_reason")
        req.HortiTechnology = c.FormValue("horti_technology")
        req.PostHarvestProblems = c.FormValue("post_harvest_problems")
    }

    
    req.PlantationCommodity = c.FormValue("plantation_commodity")
    if req.PlantationCommodity != "" {
        req.PlantationLandStatus = c.FormValue("plantation_land_status")
        if area, err := strconv.ParseFloat(c.FormValue("plantation_land_area"), 64); err == nil && area > 0 {
            req.PlantationLandArea = area
        }
        req.PlantationGrowthPhase = c.FormValue("plantation_growth_phase")
        if age, err := strconv.Atoi(c.FormValue("plantation_plant_age")); err == nil && age >= 0 {
            req.PlantationPlantAge = age
        }
        req.PlantationPlantingDate = c.FormValue("plantation_planting_date")
        req.PlantationHarvestDate = c.FormValue("plantation_harvest_date")
        req.PlantationDelayReason = c.FormValue("plantation_delay_reason")
        req.PlantationTechnology = c.FormValue("plantation_technology")
        req.ProductionProblems = c.FormValue("production_problems")
    }

    
    if hasPest, err := strconv.ParseBool(c.FormValue("has_pest_disease")); err == nil {
        req.HasPestDisease = hasPest
        if hasPest {
            req.PestDiseaseType = c.FormValue("pest_disease_type")
            req.PestDiseaseCommodity = c.FormValue("pest_disease_commodity")
            req.AffectedArea = c.FormValue("affected_area")
            req.ControlAction = c.FormValue("control_action")
        }
    }

    
    req.WeatherCondition = c.FormValue("weather_condition")
    req.WeatherImpact = c.FormValue("weather_impact")
    req.MainConstraint = c.FormValue("main_constraint")

    req.FarmerHope = c.FormValue("farmer_hope")
    req.TrainingNeeded = c.FormValue("training_needed")
    req.UrgentNeeds = c.FormValue("urgent_needs")
    req.WaterAccess = c.FormValue("water_access")
    req.Suggestions = c.FormValue("suggestions")

    if req.FoodCommodity == "" && req.HortiCommodity == "" && req.PlantationCommodity == "" {
        return response.BadRequest(c, "At least one commodity (food/horticulture/plantation) must be specified", nil)
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
    if len(photos) < 1 {
        return response.BadRequest(c, "At least 1 photo required", nil)
    }
    for _, photo := range photos {
        if photo.Size > 10*1024*1024 {
            return response.BadRequest(c, "Photo file size too large (max 10MB)", nil)
        }
        contentType := photo.Header.Get("Content-Type")
        if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/jpg" {
            return response.BadRequest(c, "Only JPEG and PNG images are allowed", nil)
        }
    }

    report, err := h.agricultureUseCase.CreateReport(c.Context(), &req, photos)
    if err != nil {
        return response.InternalError(c, "Failed to create agriculture report", err)
    }

    return response.Created(c, "Agriculture report created successfully", report)
}


func (h *AgricultureHandler) GetReport(c *fiber.Ctx) error {
    idStr := c.Params("id")

    report, err := h.agricultureUseCase.GetReport(c.Context(), idStr)
    if err != nil {
        return response.NotFound(c, "Report not found", err)
    }

    return response.Success(c, "Report retrieved successfully", report)
}

func (h *AgricultureHandler) ListReports(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))

    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    filters := map[string]interface{}{
        "extension_officer":    c.Query("extension_officer"),
        "village":              c.Query("village"),
        "district":             c.Query("district"),
        "farmer_name":          c.Query("farmer_name"),
        "farmer_group":         c.Query("farmer_group"),
        "farmer_group_type":    c.Query("farmer_group_type"),
        "food_commodity":       c.Query("food_commodity"),
        "horti_commodity":      c.Query("horti_commodity"),
        "plantation_commodity": c.Query("plantation_commodity"),
        "main_constraint":      c.Query("main_constraint"),
        "weather_condition":    c.Query("weather_condition"),
        "water_access":         c.Query("water_access"),
        "start_date":           c.Query("start_date"),
        "end_date":             c.Query("end_date"),
    }

    if hasPestStr := c.Query("has_pest_disease"); hasPestStr != "" {
        if hasPest, err := strconv.ParseBool(hasPestStr); err == nil {
            filters["has_pest_disease"] = hasPest
        }
    }

    
    normalizeFilters(filters)

    result, err := h.agricultureUseCase.ListReports(c.Context(), page, limit, filters)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve reports", err)
    }

    return response.Success(c, "Reports retrieved successfully", result)
}

func (h *AgricultureHandler) UpdateReport(c *fiber.Ctx) error {
    idStr := c.Params("id")

    var req dto.UpdateAgricultureRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", err)
    }

    
    req.Normalize()

    if err := req.Validate(); err != nil {
        return response.ValidationError(c, err)
    }

    userID := c.Locals("userID").(string)

    report, err := h.agricultureUseCase.UpdateReport(c.Context(), idStr, &req, userID)
    if err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to update this report", err)
        }
        return response.InternalError(c, "Failed to update report", err)
    }

    return response.Success(c, "Report updated successfully", report)
}


func (h *AgricultureHandler) DeleteReport(c *fiber.Ctx) error {
    idStr := c.Params("id")

    userID := c.Locals("userID").(string)

    if err := h.agricultureUseCase.DeleteReport(c.Context(), idStr, userID); err != nil {
        if err == usecase.ErrUnauthorized {
            return response.Forbidden(c, "You don't have permission to delete this report", err)
        }
        return response.InternalError(c, "Failed to delete report", err)
    }

    return response.Success(c, "Report deleted successfully", nil)
}

func (h *AgricultureHandler) GetExecutiveDashboard(c *fiber.Ctx) error {
    commodityType := c.Query("commodity_type", "") 
    
    summary, err := h.agricultureUseCase.GetExecutiveSummary(c.Context(), commodityType)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve executive summary", err)
    }

    return response.Success(c, "Executive summary retrieved successfully", summary)
}

func (h *AgricultureHandler) GetCommodityAnalysis(c *fiber.Ctx) error {
    rawName := c.Query("commodity_name")
    if rawName == "" {
        return response.BadRequest(c, "commodity_name parameter is required", nil)
    }
    commodityName := utils.NormalizeEnum(rawName)

    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "start_date and end_date parameters are required", nil)
    }

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    analysis, err := h.agricultureUseCase.GetCommodityAnalysis(c.Context(), startDate, endDate, commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve commodity analysis", err)
    }

    if analysis.TotalHarvestedArea == 0 {
        return response.Success(c, 
            fmt.Sprintf("No data found for commodity '%s' in the specified date range", commodityName), 
            analysis)
    }

    return response.Success(c, "Commodity analysis retrieved successfully", analysis)
}

func (h *AgricultureHandler) GetFoodCropStats(c *fiber.Ctx) error {
    commodityName := utils.NormalizeEnum(c.Query("commodity_name", ""))

    stats, err := h.agricultureUseCase.GetFoodCropStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve food crop statistics", err)
    }

    return response.Success(c, "Food crop statistics retrieved successfully", stats)
}

func (h *AgricultureHandler) GetHorticultureStats(c *fiber.Ctx) error {
    commodityName := utils.NormalizeEnum(c.Query("commodity_name", ""))

    stats, err := h.agricultureUseCase.GetHorticultureStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve horticulture statistics", err)
    }

    return response.Success(c, "Horticulture statistics retrieved successfully", stats)
}


func (h *AgricultureHandler) GetPlantationStats(c *fiber.Ctx) error {
    commodityName := utils.NormalizeEnum(c.Query("commodity_name", ""))

    stats, err := h.agricultureUseCase.GetPlantationStats(c.Context(), commodityName)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve plantation statistics", err)
    }

    return response.Success(c, "Plantation statistics retrieved successfully", stats)
}


func (h *AgricultureHandler) GetAgriculturalEquipmentStats(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "start_date and end_date parameters are required", nil)
    }

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    stats, err := h.agricultureUseCase.GetAgriculturalEquipmentStats(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve agricultural equipment statistics", err)
    }

    return response.Success(c, "Agricultural equipment statistics retrieved successfully", stats)
}



func (h *AgricultureHandler) GetLandAndIrrigationStats(c *fiber.Ctx) error {
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")
    
    if startDateStr == "" || endDateStr == "" {
        return response.BadRequest(c, "start_date and end_date parameters are required", nil)
    }

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
    }
    
    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
    }

    stats, err := h.agricultureUseCase.GetLandAndIrrigationStats(c.Context(), startDate, endDate)
    if err != nil {
        return response.InternalError(c, "Failed to retrieve land and irrigation statistics", err)
    }

    return response.Success(c, "Land and irrigation statistics retrieved successfully", stats)
}

func (h *AgricultureHandler) HandlePanic(c *fiber.Ctx) {
    if r := recover(); r != nil {
        log.Printf("[PANIC] Agriculture handler panic: %v", r)
        response.InternalError(c, "Internal server error occurred", fmt.Errorf("%v", r))
    }
}

var locationKeys = map[string]bool{
    "extension_officer": true,
    "village":           true,
    "district":          true,
    "farmer_name":       true,
    "farmer_group":      true,
}


var enumKeys = map[string]bool{
    "farmer_group_type":     true,
    "food_commodity":        true,
    "horti_commodity":       true,
    "plantation_commodity":  true,
    "main_constraint":       true,
    "weather_condition":     true,
    "water_access":          true,
    
    
}


func normalizeFilters(filters map[string]interface{}) {
    for k, v := range filters {
        s, ok := v.(string)
        if !ok || s == "" {
            continue
        }
        if locationKeys[k] {
            filters[k] = utils.NormalizeLocation(s)
            continue
        }
        if enumKeys[k] {
            filters[k] = utils.NormalizeEnum(s)
            continue
        }
        
    }
}

func (h *AgricultureHandler) ExportKomoditas(c *fiber.Ctx) error {
	// Parse query parameters
	commodityName := c.Query("commodity_name", "")
	if commodityName != "" {
		commodityName = utils.NormalizeEnum(commodityName)
	}

	commodityType := strings.ToUpper(c.Query("commodity_type", ""))
	if commodityType != "" && commodityType != "PANGAN" && commodityType != "HORTIKULTURA" && commodityType != "PERKEBUNAN" {
		return response.BadRequest(c, "Invalid commodity_type. Must be one of: PANGAN, HORTIKULTURA, PERKEBUNAN", nil)
	}

	district := c.Query("district", "")
	if district != "" {
		district = utils.NormalizeLocation(district)
	}

	// Parse date range
	var startDate, endDate time.Time
	var err error

	startDateStr := c.Query("start_date", "")
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD", err)
		}
	}

	endDateStr := c.Query("end_date", "")
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD", err)
		}
	}

	// Generate Excel file
	excelData, err := h.exportKomoditasToExcel(c.Context(), commodityName, commodityType, district, startDate, endDate)
	if err != nil {
		return response.InternalError(c, "Failed to export data", err)
	}

	// Generate filename
	filename := h.generateExportFilename(commodityName, commodityType, district, startDate, endDate)

	// Set response headers
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Set("Content-Length", fmt.Sprintf("%d", len(excelData)))

	return c.Send(excelData)
}

func (h *AgricultureHandler) exportKomoditasToExcel(ctx context.Context, commodityName, commodityType, district string, startDate, endDate time.Time) ([]byte, error) {
	// Prepare filters
	filters := make(map[string]interface{})

	if district != "" {
		filters["district"] = district
	}

	if !startDate.IsZero() {
		filters["start_date"] = startDate.Format("2006-01-02")
	}

	if !endDate.IsZero() {
		filters["end_date"] = endDate.Format("2006-01-02")
	}

	// Fetch all reports
	reports,  err := h.agricultureUseCase.ListReports(ctx, 1, 10000, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reports: %w", err)
	}

	// Transform to export data
	type ExportData struct {
		Komoditas      string
		Kecamatan      string
		Koordinat      string
		Tahun          int
		Produksi       float64
		JumlahProduksi float64
		LuasPanen      float64
		Produktivitas  float64
	}

	var exportData []ExportData

	for _, report := range reports.Reports {
		// Process Food Crops
		if report.FoodCommodity != "" {
			if commodityType == "" || commodityType == "PANGAN" {
				if commodityName == "" || string(report.FoodCommodity) == commodityName {
					produksi := report.FoodLandArea * 5.0
					exportData = append(exportData, ExportData{
						Komoditas:      string(report.FoodCommodity),
						Kecamatan:      report.District,
						Koordinat:      fmt.Sprintf("%.6f, %.6f", report.Latitude, report.Longitude),
						Tahun:          report.VisitDate.Year(),
						Produksi:       produksi,
						JumlahProduksi: produksi,
						LuasPanen:      report.FoodLandArea,
						Produktivitas:  5.0,
					})
				}
			}
		}

		// Process Horticulture
		if report.HortiCommodity != "" {
			if commodityType == "" || commodityType == "HORTIKULTURA" {
				commodityNameVal := string(report.HortiCommodity)
				if report.HortiSubCommodity != "" {
					commodityNameVal = report.HortiSubCommodity
				}

				if commodityName == "" || commodityNameVal == commodityName {
					produksi := report.HortiLandArea * 10.0
					exportData = append(exportData, ExportData{
						Komoditas:      commodityNameVal,
						Kecamatan:      report.District,
						Koordinat:      fmt.Sprintf("%.6f, %.6f", report.Latitude, report.Longitude),
						Tahun:          report.VisitDate.Year(),
						Produksi:       produksi,
						JumlahProduksi: produksi,
						LuasPanen:      report.HortiLandArea,
						Produktivitas:  10.0,
					})
				}
			}
		}

		// Process Plantation
		if report.PlantationCommodity != "" {
			if commodityType == "" || commodityType == "PERKEBUNAN" {
				if commodityName == "" || string(report.PlantationCommodity) == commodityName {
					produksi := report.PlantationLandArea * 2.0
					exportData = append(exportData, ExportData{
						Komoditas:      string(report.PlantationCommodity),
						Kecamatan:      report.District,
						Koordinat:      fmt.Sprintf("%.6f, %.6f", report.Latitude, report.Longitude),
						Tahun:          report.VisitDate.Year(),
						Produksi:       produksi,
						JumlahProduksi: produksi,
						LuasPanen:      report.PlantationLandArea,
						Produktivitas:  2.0,
					})
				}
			}
		}
	}

	// Create Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	sheetName := "komoditas"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Write headers
	headers := []string{
		"Komoditas",
		"Kecamatan",
		"Koordinat Lokasi",
		"Tahun",
		"Produksi (ton)",
		"Jumlah Produksi (Ton/Ha)",
		"Luas Panen (ha)",
		"Produktivitas (Ton/Ha)",
	}

	// Header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Column widths
	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 25)
	f.SetColWidth(sheetName, "D", "D", 10)
	f.SetColWidth(sheetName, "E", "E", 15)
	f.SetColWidth(sheetName, "F", "F", 20)
	f.SetColWidth(sheetName, "G", "G", 15)
	f.SetColWidth(sheetName, "H", "H", 20)

	f.SetRowHeight(sheetName, 1, 25)

	// Data style
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Calibri",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
		},
	})

	numberStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Calibri",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
		},
		CustomNumFmt: stringPtr("#,##0.00"),
	})

	// Write data
	for i, row := range exportData {
		rowNum := i + 2

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row.Komoditas)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNum), fmt.Sprintf("A%d", rowNum), dataStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row.Kecamatan)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowNum), fmt.Sprintf("B%d", rowNum), dataStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), row.Koordinat)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("C%d", rowNum), dataStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), row.Tahun)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", rowNum), fmt.Sprintf("D%d", rowNum), dataStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), row.Produksi)
		f.SetCellStyle(sheetName, fmt.Sprintf("E%d", rowNum), fmt.Sprintf("E%d", rowNum), numberStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), row.JumlahProduksi)
		f.SetCellStyle(sheetName, fmt.Sprintf("F%d", rowNum), fmt.Sprintf("F%d", rowNum), numberStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), row.LuasPanen)
		f.SetCellStyle(sheetName, fmt.Sprintf("G%d", rowNum), fmt.Sprintf("G%d", rowNum), numberStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("H%d", rowNum), row.Produktivitas)
		f.SetCellStyle(sheetName, fmt.Sprintf("H%d", rowNum), fmt.Sprintf("H%d", rowNum), numberStyle)
	}

	// Freeze header row
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Add autofilter
	if len(exportData) > 0 {
		lastRow := len(exportData) + 1
		f.AutoFilter(sheetName, fmt.Sprintf("A1:H%d", lastRow), nil)
	}

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

func (h *AgricultureHandler) generateExportFilename(commodityName, commodityType, district string, startDate, endDate time.Time) string {
	timestamp := time.Now().Format("20060102_150405")
	
	var parts []string
	parts = append(parts, "export_komoditas")
	
	if commodityType != "" {
		parts = append(parts, strings.ToLower(commodityType))
	}
	
	if commodityName != "" {
		parts = append(parts, strings.ToLower(strings.ReplaceAll(commodityName, " ", "_")))
	}
	
	if district != "" {
		parts = append(parts, strings.ToLower(strings.ReplaceAll(district, " ", "_")))
	}
	
	if !startDate.IsZero() && !endDate.IsZero() {
		parts = append(parts, fmt.Sprintf("%s_to_%s", startDate.Format("20060102"), endDate.Format("20060102")))
	}
	
	parts = append(parts, timestamp)
	
	return strings.Join(parts, "_") + ".xlsx"
}

func stringPtr(s string) *string {
	return &s
}