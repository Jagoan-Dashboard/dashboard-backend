package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/storage"
	"building-report-backend/pkg/utils"
)

type SpatialPlanningUseCase struct {
	spatialRepo repository.SpatialPlanningRepository
	storage     storage.StorageService
	cache       repository.CacheRepository
}

func NewSpatialPlanningUseCase(
	spatialRepo repository.SpatialPlanningRepository,
	storage storage.StorageService,
	cache repository.CacheRepository,
) *SpatialPlanningUseCase {
	return &SpatialPlanningUseCase{
		spatialRepo: spatialRepo,
		storage:     storage,
		cache:       cache,
	}
}

func (uc *SpatialPlanningUseCase) CreateReport(ctx context.Context, req *dto.CreateSpatialPlanningRequest, photos []*multipart.FileHeader) (*entity.SpatialPlanningReport, error) {
	report := &entity.SpatialPlanningReport{
		ID:                  utils.GenerateULID(),
		ReporterName:        req.ReporterName,
		Institution:         entity.InstitutionType(req.Institution),
		PhoneNumber:         req.PhoneNumber,
		ReportDateTime:      req.ReportDateTime,
		AreaDescription:     req.AreaDescription,
		AreaCategory:        entity.AreaCategory(req.AreaCategory),
		ViolationType:       entity.SpatialViolationType(req.ViolationType),
		ViolationLevel:      entity.ViolationLevel(req.ViolationLevel),
		EnvironmentalImpact: entity.EnvironmentalImpact(req.EnvironmentalImpact),
		UrgencyLevel:        entity.UrgencyLevel(req.UrgencyLevel),
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		Address:             req.Address,
		Notes:               req.Notes,
		Status:              entity.SpatialStatusPending,
	}

	for i, photo := range photos {
		photoURL, err := uc.storage.UploadFile(ctx, photo, "spatial-planning")
		if err != nil {
			return nil, fmt.Errorf("failed to upload photo: %w", err)
		}

		caption := fmt.Sprintf("Photo %d", i+1)
		report.Photos = append(report.Photos, entity.SpatialPlanningPhoto{
			ID:       utils.GenerateULID(),
			PhotoURL: photoURL,
			Caption:  caption,
		})
	}

	if err := uc.spatialRepo.Create(ctx, report); err != nil {
		return nil, err
	}

	uc.cache.Delete(ctx, "spatial:list")
	uc.cache.Delete(ctx, "spatial:stats")

	return report, nil
}

func (uc *SpatialPlanningUseCase) GetReport(ctx context.Context, id string) (*entity.SpatialPlanningReport, error) {
    cacheKey := "spatial:" + id

    report, err := uc.spatialRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    uc.cache.Set(ctx, cacheKey, report, 3600)

    return report, nil
}

func (uc *SpatialPlanningUseCase) ListReports(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.PaginatedSpatialReportsResponse, error) {
    offset := (page - 1) * limit

    reports, total, err := uc.spatialRepo.FindAll(ctx, limit, offset, filters)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedSpatialReportsResponse{
        Reports:    reports,
        Total:      total,
        Page:       page,
        PerPage:    limit,
        TotalPages: (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *SpatialPlanningUseCase) UpdateReport(ctx context.Context, id string, req *dto.UpdateSpatialPlanningRequest, userID string) (*entity.SpatialPlanningReport, error) {
	report, err := uc.spatialRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// if report.CreatedBy != userID {

	//     return nil, ErrUnauthorized
	// }

	if req.AreaDescription != "" {
		report.AreaDescription = req.AreaDescription
	}
	if req.AreaCategory != "" {
		report.AreaCategory = entity.AreaCategory(req.AreaCategory)
	}
	if req.ViolationType != "" {
		report.ViolationType = entity.SpatialViolationType(req.ViolationType)
	}
	if req.ViolationLevel != "" {
		report.ViolationLevel = entity.ViolationLevel(req.ViolationLevel)
	}
	if req.EnvironmentalImpact != "" {
		report.EnvironmentalImpact = entity.EnvironmentalImpact(req.EnvironmentalImpact)
	}
	if req.UrgencyLevel != "" {
		report.UrgencyLevel = entity.UrgencyLevel(req.UrgencyLevel)
	}
	if req.Latitude != 0 {
		report.Latitude = req.Latitude
	}
	if req.Longitude != 0 {
		report.Longitude = req.Longitude
	}
	if req.Address != "" {
		report.Address = req.Address
	}
	if req.Notes != "" {
		report.Notes = req.Notes
	}
	if req.Status != "" {
		report.Status = entity.SpatialReportStatus(req.Status)
	}

	if err := uc.spatialRepo.Update(ctx, report); err != nil {
		return nil, err
	}

	uc.cache.Delete(ctx, "spatial:"+id)
	uc.cache.Delete(ctx, "spatial:list")

	return report, nil
}

func (uc *SpatialPlanningUseCase) UpdateStatus(ctx context.Context, id string, req *dto.UpdateSpatialStatusRequest) error {
	report, err := uc.spatialRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	report.Status = entity.SpatialReportStatus(req.Status)
	if req.Notes != "" {
		report.Notes = req.Notes
	}

	err = uc.spatialRepo.UpdateStatus(ctx, id, entity.SpatialReportStatus(req.Status))
	if err != nil {
		return err
	}

	uc.cache.Delete(ctx, "spatial:"+id)
	uc.cache.Delete(ctx, "spatial:stats")

	return nil
}

func (uc *SpatialPlanningUseCase) DeleteReport(ctx context.Context, id string, userID string) error {
	report, err := uc.spatialRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// if report.CreatedBy != userID {

	// 	return ErrUnauthorized
	// }

	for _, photo := range report.Photos {
		uc.storage.DeleteFile(ctx, photo.PhotoURL)
	}

	if err := uc.spatialRepo.Delete(ctx, id); err != nil {
		return err
	}

	uc.cache.Delete(ctx, "spatial:"+id)
	uc.cache.Delete(ctx, "spatial:list")
	uc.cache.Delete(ctx, "spatial:stats")

	return nil
}

func (uc *SpatialPlanningUseCase) GetStatistics(ctx context.Context) (*dto.SpatialStatisticsResponse, error) {

	cacheKey := "spatial:stats"
	var stats dto.SpatialStatisticsResponse

	err := uc.cache.Get(ctx, cacheKey, &stats)
	if err == nil {
		return &stats, nil
	}

	rawStats, err := uc.spatialRepo.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	response := &dto.SpatialStatisticsResponse{
		TotalReports:  rawStats["total_reports"].(int64),
		UrgentReports: rawStats["urgent_reports"].(int64),
	}

	if violationLevels, ok := rawStats["violation_levels"].([]interface{}); ok {
		for _, v := range violationLevels {
			if m, ok := v.(map[string]interface{}); ok {
				response.ViolationLevels = append(response.ViolationLevels, m)
			}
		}
	}

	if statusCounts, ok := rawStats["status_counts"].([]interface{}); ok {
		for _, v := range statusCounts {
			if m, ok := v.(map[string]interface{}); ok {
				response.StatusCounts = append(response.StatusCounts, m)
			}
		}
	}

	uc.cache.Set(ctx, cacheKey, response, 300)

	return response, nil
}

func (uc *SpatialPlanningUseCase) GetTataRuangOverview(ctx context.Context, areaCategory string) (*dto.TataRuangOverviewResponse, error) {
	// Cache key based on area category
	cacheKey := fmt.Sprintf("tata_ruang:overview:%s", areaCategory)
	var response dto.TataRuangOverviewResponse

	err := uc.cache.Get(ctx, cacheKey, &response)
	if err == nil {
		return &response, nil
	}

	// Get basic statistics
	basicStatsRaw, err := uc.spatialRepo.GetTataRuangStatistics(ctx, areaCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic statistics: %w", err)
	}

	response.BasicStats = dto.TataRuangBasicStatistics{
		TotalReports:          basicStatsRaw["total_reports"].(int64),
		EstimatedTotalLengthM: basicStatsRaw["estimated_total_length_m"].(float64),
		EstimatedTotalAreaM2:  basicStatsRaw["estimated_total_area_m2"].(float64),
		UrgentReportsCount:    basicStatsRaw["urgent_reports_count"].(int64),
	}

	// Get location distribution
	locationStats, err := uc.spatialRepo.GetLocationDistribution(ctx, areaCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get location distribution: %w", err)
	}

	for _, loc := range locationStats {
		response.LocationDistribution = append(response.LocationDistribution, dto.TataRuangLocationDistribution{
			District:       loc["district"].(string),
			Village:        loc["village"].(string),
			ViolationCount: int(loc["violation_count"].(int64)),
			AvgLatitude:    loc["avg_latitude"].(float64),
			AvgLongitude:   loc["avg_longitude"].(float64),
			UrgentCount:    int(loc["urgent_count"].(int64)),
			SevereCount:    int(loc["severe_count"].(int64)),
		})
	}

	// Get urgency level statistics
	urgencyStats, err := uc.spatialRepo.GetUrgencyLevelStatistics(ctx, areaCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get urgency statistics: %w", err)
	}

	for _, urgency := range urgencyStats {
		response.UrgencyStatistics = append(response.UrgencyStatistics, dto.TataRuangUrgencyStatistics{
			UrgencyLevel: urgency["urgency_level"].(string),
			Count:        urgency["count"].(int64),
			Percentage:   urgency["percentage"].(float64),
		})
	}

	// Get violation type statistics
	violationTypeStats, err := uc.spatialRepo.GetViolationTypeStatistics(ctx, areaCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get violation type statistics: %w", err)
	}

	for _, vt := range violationTypeStats {
		response.ViolationTypeStatistics = append(response.ViolationTypeStatistics, dto.TataRuangViolationTypeStatistics{
			ViolationType: vt["violation_type"].(string),
			Count:         vt["count"].(int64),
			Percentage:    vt["percentage"].(float64),
			SevereCount:   int(vt["severe_count"].(int64)),
			UrgentCount:   int(vt["urgent_count"].(int64)),
		})
	}

	// Get violation level statistics
	violationLevelStats, err := uc.spatialRepo.GetViolationLevelStatistics(ctx, areaCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get violation level statistics: %w", err)
	}

	for _, vl := range violationLevelStats {
		response.ViolationLevelStatistics = append(response.ViolationLevelStatistics, dto.TataRuangViolationLevelStatistics{
			ViolationLevel: vl["violation_level"].(string),
			Count:          vl["count"].(int64),
			Percentage:     vl["percentage"].(float64),
			UrgentCount:    int(vl["urgent_count"].(int64)),
		})
	}

	// Get area category distribution (only if getting all categories)
	if areaCategory == "" || areaCategory == "all" {
		areaCategoryStats, err := uc.spatialRepo.GetAreaCategoryDistribution(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get area category distribution: %w", err)
		}

		for _, ac := range areaCategoryStats {
			response.AreaCategoryDistribution = append(response.AreaCategoryDistribution, dto.TataRuangAreaCategoryDistribution{
				AreaCategory: ac["area_category"].(string),
				Count:        ac["count"].(int64),
				Percentage:   ac["percentage"].(float64),
				UrgentCount:  int(ac["urgent_count"].(int64)),
				SevereCount:  int(ac["severe_count"].(int64)),
			})
		}
	}

	// Get environmental impact statistics
	environmentalStats, err := uc.spatialRepo.GetEnvironmentalImpactStatistics(ctx, areaCategory)
	if err != nil {
		return nil, fmt.Errorf("failed to get environmental impact statistics: %w", err)
	}

	for _, env := range environmentalStats {
		response.EnvironmentalImpactStatistics = append(response.EnvironmentalImpactStatistics, dto.TataRuangEnvironmentalImpactStatistics{
			EnvironmentalImpact: env["environmental_impact"].(string),
			Count:               env["count"].(int64),
			Percentage:          env["percentage"].(float64),
			SevereCount:         int(env["severe_count"].(int64)),
		})
	}

	// Cache the response for 5 minutes
	uc.cache.Set(ctx, cacheKey, &response, 300*time.Second)

	return &response, nil
}
