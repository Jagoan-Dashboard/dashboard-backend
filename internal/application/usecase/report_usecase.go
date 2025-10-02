package usecase

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/storage"
	"building-report-backend/pkg/utils"
	"context"
	"fmt"
	"mime/multipart"
	"time"
)

type ReportUseCase struct {
    reportRepo repository.ReportRepository
    storage    storage.StorageService
    cache      repository.CacheRepository
}

func NewReportUseCase(
    reportRepo repository.ReportRepository,
    storage storage.StorageService,
    cache repository.CacheRepository,
) *ReportUseCase {
    return &ReportUseCase{
        reportRepo: reportRepo,
        storage:    storage,
        cache:      cache,
    }
}

func (uc *ReportUseCase) CreateReport(ctx context.Context, req *dto.CreateReportRequest, photos []*multipart.FileHeader) (*entity.Report, error) {
    report := &entity.Report{
        ID:                   utils.GenerateULID(),
        ReporterName:         req.ReporterName,
        ReporterRole:         entity.ReporterRole(req.ReporterRole),
        Village:              req.Village,
        District:             req.District,
        BuildingName:         req.BuildingName,
        BuildingType:         entity.BuildingType(req.BuildingType),
        ReportStatus:         entity.ReportStatusType(req.ReportStatus),
        FundingSource:        entity.FundingSource(req.FundingSource),
        LastYearConstruction: req.LastYearConstruction,
        FullAddress:          req.FullAddress,
        Latitude:             req.Latitude,
        Longitude:            req.Longitude,
        FloorArea:            req.FloorArea,
        FloorCount:           req.FloorCount,
    }

    
    if req.WorkType != "" {
        workType := entity.WorkType(req.WorkType)
        report.WorkType = &workType
    }
    if req.ConditionAfterRehab != "" {
        condition := entity.ConditionAfterRehab(req.ConditionAfterRehab)
        report.ConditionAfterRehab = &condition
    }

    
    for i, photo := range photos {
        photoType := "overall"
        if i == 0 {
            photoType = "closeup"
        }

        photoURL, err := uc.storage.UploadFile(ctx, photo, "reports")
        if err != nil {
            return nil, err
        }

        report.Photos = append(report.Photos, entity.ReportPhoto{
            ID:       utils.GenerateULID(),
            PhotoURL:  photoURL,
            PhotoType: photoType,
        })
    }

    if err := uc.reportRepo.Create(ctx, report); err != nil {
        return nil, err
    }


    uc.cache.Delete(ctx, "reports:list")

    return report, nil
}

func (uc *ReportUseCase) GetReport(ctx context.Context, id string) (*entity.Report, error) {

    cacheKey := "report:" + id

    report, err := uc.reportRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, report, 3600) 

    return report, nil
}

func (uc *ReportUseCase) ListReports(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.PaginatedReportsResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.reportRepo.FindAll(ctx, limit, offset, filters)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedReportsResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *ReportUseCase) UpdateReport(ctx context.Context, id string, req *dto.UpdateReportRequest, userID string) (*entity.Report, error) {
    report, err := uc.reportRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    if report.CreatedBy != userID {
        return nil, ErrUnauthorized
    }

    
    if req.BuildingName != "" {
        report.BuildingName = req.BuildingName
    }
    if req.ReportStatus != "" {
        report.ReportStatus = entity.ReportStatusType(req.ReportStatus)
    }
    

    if err := uc.reportRepo.Update(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "report:"+id)
    uc.cache.Delete(ctx, "reports:list")

    return report, nil
}

func (uc *ReportUseCase) DeleteReport(ctx context.Context, id string, userID string) error {
    report, err := uc.reportRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    
    if report.CreatedBy != userID {
        return ErrUnauthorized
    }

    
    for _, photo := range report.Photos {
        uc.storage.DeleteFile(ctx, photo.PhotoURL)
    }

    if err := uc.reportRepo.Delete(ctx, id); err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "report:"+id)
    uc.cache.Delete(ctx, "reports:list")

    return nil
}

func (uc *ReportUseCase) GetTataBangunanOverview(ctx context.Context, buildingType string) (*dto.TataBangunanOverviewResponse, error) {
    // Cache key based on building type
    cacheKey := fmt.Sprintf("tata_bangunan:overview:%s", buildingType)
    var response dto.TataBangunanOverviewResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    // Get basic statistics
    basicStatsRaw, err := uc.reportRepo.GetStatistics(ctx, buildingType)
    if err != nil {
        return nil, fmt.Errorf("failed to get basic statistics: %w", err)
    }
    
    response.BasicStats = dto.ReportStatisticsResponse{
        TotalReports:       basicStatsRaw["total_reports"].(int64),
        AverageFloorArea:   basicStatsRaw["average_floor_area"].(float64),
        AverageFloorCount:  basicStatsRaw["average_floor_count"].(float64),
        DamagedBuildings:   basicStatsRaw["damaged_buildings_count"].(int64),
    }

    // Get location distribution
    locationStats, err := uc.reportRepo.GetLocationStatistics(ctx, buildingType)
    if err != nil {
        return nil, fmt.Errorf("failed to get location statistics: %w", err)
    }
    
    for _, loc := range locationStats {
        response.LocationDistribution = append(response.LocationDistribution, dto.LocationStatisticsResponse{
            District:      loc["district"].(string),
            Village:       loc["village"].(string),
            BuildingCount: int(loc["building_count"].(int64)),
            AvgLatitude:   loc["avg_latitude"].(float64),
            AvgLongitude:  loc["avg_longitude"].(float64),
            DamagedCount:  int(loc["damaged_count"].(int64)),
        })
    }

    // Get status distribution
    statusStats, err := uc.reportRepo.GetStatusStatistics(ctx, buildingType)
    if err != nil {
        return nil, fmt.Errorf("failed to get status statistics: %w", err)
    }
    
    for _, status := range statusStats {
        response.StatusDistribution = append(response.StatusDistribution, dto.StatusStatisticsResponse{
            Status: status["report_status"].(string),
            Count:  status["count"].(int64),
        })
    }

    // Get work type distribution
    workTypeStats, err := uc.reportRepo.GetWorkTypeStatistics(ctx, buildingType)
    if err != nil {
        return nil, fmt.Errorf("failed to get work type statistics: %w", err)
    }
    
    for _, workType := range workTypeStats {
        response.WorkTypeDistribution = append(response.WorkTypeDistribution, dto.WorkTypeStatisticsResponse{
            WorkType: workType["work_type"].(string),
            Count:    workType["count"].(int64),
        })
    }

    // Get condition after rehab distribution
    conditionStats, err := uc.reportRepo.GetConditionAfterRehabStatistics(ctx, buildingType)
    if err != nil {
        return nil, fmt.Errorf("failed to get condition statistics: %w", err)
    }
    
    for _, condition := range conditionStats {
        response.ConditionDistribution = append(response.ConditionDistribution, dto.ConditionStatisticsResponse{
            Condition: condition["condition_after_rehab"].(string),
            Count:     condition["count"].(int64),
        })
    }

    // Get building type distribution (only if getting all types)
    if buildingType == "" || buildingType == "all" {
        buildingTypeStats, err := uc.reportRepo.CountByBuildingType(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed to get building type statistics: %w", err)
        }
        
        for _, bt := range buildingTypeStats {
            response.BuildingTypeDistribution = append(response.BuildingTypeDistribution, dto.BuildingTypeStatisticsResponse{
                BuildingType: bt["building_type"].(string),
                Count:        bt["count"].(int64),
            })
        }
    }

    // Cache the response for 5 minutes
    uc.cache.Set(ctx, cacheKey, &response, 300*time.Second)

    return &response, nil
}

func (uc *ReportUseCase) GetBasicStatistics(ctx context.Context, buildingType string) (*dto.ReportStatisticsResponse, error) {
    cacheKey := fmt.Sprintf("reports:basic_stats:%s", buildingType)
    var response dto.ReportStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    statsRaw, err := uc.reportRepo.GetStatistics(ctx, buildingType)
    if err != nil {
        return nil, err
    }

    response = dto.ReportStatisticsResponse{
        TotalReports:       statsRaw["total_reports"].(int64),
        AverageFloorArea:   statsRaw["average_floor_area"].(float64),
        AverageFloorCount:  statsRaw["average_floor_count"].(float64),
        DamagedBuildings:   statsRaw["damaged_buildings_count"].(int64),
    }

    uc.cache.Set(ctx, cacheKey, &response, 300*time.Second)
    return &response, nil
}

func (uc *ReportUseCase) GetLocationDistribution(ctx context.Context, buildingType string) ([]dto.LocationStatisticsResponse, error) {
    cacheKey := fmt.Sprintf("reports:location_dist:%s", buildingType)
    var response []dto.LocationStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return response, nil
    }

    locationStats, err := uc.reportRepo.GetLocationStatistics(ctx, buildingType)
    if err != nil {
        return nil, err
    }

    for _, loc := range locationStats {
        response = append(response, dto.LocationStatisticsResponse{
            District:      loc["district"].(string),
            Village:       loc["village"].(string),
            BuildingCount: int(loc["building_count"].(int64)),
            AvgLatitude:   loc["avg_latitude"].(float64),
            AvgLongitude:  loc["avg_longitude"].(float64),
            DamagedCount:  int(loc["damaged_count"].(int64)),
        })
    }

    uc.cache.Set(ctx, cacheKey, response, 300*time.Second)
    return response, nil
}


func (uc *ReportUseCase) GetWorkTypeStatistics(ctx context.Context, buildingType string) ([]dto.WorkTypeStatisticsResponse, error) {
    cacheKey := fmt.Sprintf("reports:work_type_stats:%s", buildingType)
    var response []dto.WorkTypeStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return response, nil
    }

    workTypeStats, err := uc.reportRepo.GetWorkTypeStatistics(ctx, buildingType)
    if err != nil {
        return nil, err
    }

    for _, workType := range workTypeStats {
        response = append(response, dto.WorkTypeStatisticsResponse{
            WorkType: workType["work_type"].(string),
            Count:    workType["count"].(int64),
        })
    }

    uc.cache.Set(ctx, cacheKey, response, 300*time.Second)
    return response, nil
}

func (uc *ReportUseCase) GetConditionAfterRehabStatistics(ctx context.Context, buildingType string) ([]dto.ConditionStatisticsResponse, error) {
    cacheKey := fmt.Sprintf("reports:condition_stats:%s", buildingType)
    var response []dto.ConditionStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return response, nil
    }

    conditionStats, err := uc.reportRepo.GetConditionAfterRehabStatistics(ctx, buildingType)
    if err != nil {
        return nil, err
    }

    for _, condition := range conditionStats {
        response = append(response, dto.ConditionStatisticsResponse{
            Condition: condition["condition_after_rehab"].(string),
            Count:     condition["count"].(int64),
        })
    }

    uc.cache.Set(ctx, cacheKey, response, 300*time.Second)
    return response, nil
}

func (uc *ReportUseCase) GetStatusStatistics(ctx context.Context, buildingType string) ([]dto.StatusStatisticsResponse, error) {
    cacheKey := fmt.Sprintf("reports:status_stats:%s", buildingType)
    var response []dto.StatusStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return response, nil
    }

    statusStats, err := uc.reportRepo.GetStatusStatistics(ctx, buildingType)
    if err != nil {
        return nil, err
    }

    for _, status := range statusStats {
        response = append(response, dto.StatusStatisticsResponse{
            Status: status["report_status"].(string),
            Count:  status["count"].(int64),
        })
    }

    uc.cache.Set(ctx, cacheKey, response, 300*time.Second)
    return response, nil
}

func (uc *ReportUseCase) GetBuildingTypeDistribution(ctx context.Context) ([]dto.BuildingTypeStatisticsResponse, error) {
    cacheKey := "reports:building_type_dist"
    var response []dto.BuildingTypeStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return response, nil
    }

    buildingTypeStats, err := uc.reportRepo.CountByBuildingType(ctx)
    if err != nil {
        return nil, err
    }

    for _, bt := range buildingTypeStats {
        response = append(response, dto.BuildingTypeStatisticsResponse{
            BuildingType: bt["building_type"].(string),
            Count:        bt["count"].(int64),
        })
    }

    uc.cache.Set(ctx, cacheKey, response, 600*time.Second) // Cache for 10 minutes
    return response, nil
}