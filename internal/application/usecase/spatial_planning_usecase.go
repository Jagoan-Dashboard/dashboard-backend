
package usecase

import (
    "context"
    "mime/multipart"
    "fmt"
    
    "building-report-backend/internal/application/dto"
    "building-report-backend/internal/domain/entity"
    "building-report-backend/internal/domain/repository"
    "building-report-backend/internal/infrastructure/storage"
    
    "github.com/google/uuid"
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

func (uc *SpatialPlanningUseCase) CreateReport(ctx context.Context, req *dto.CreateSpatialPlanningRequest, photos []*multipart.FileHeader, userID uuid.UUID) (*entity.SpatialPlanningReport, error) {
    report := &entity.SpatialPlanningReport{
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
        CreatedBy:           userID,
    }

    
    for i, photo := range photos {
        photoURL, err := uc.storage.UploadFile(ctx, photo, "spatial-planning")
        if err != nil {
            return nil, fmt.Errorf("failed to upload photo: %w", err)
        }

        caption := fmt.Sprintf("Photo %d", i+1)
        report.Photos = append(report.Photos, entity.SpatialPlanningPhoto{
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

func (uc *SpatialPlanningUseCase) GetReport(ctx context.Context, id uuid.UUID) (*entity.SpatialPlanningReport, error) {
    cacheKey := "spatial:" + id.String()
    
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
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *SpatialPlanningUseCase) UpdateReport(ctx context.Context, id uuid.UUID, req *dto.UpdateSpatialPlanningRequest, userID uuid.UUID) (*entity.SpatialPlanningReport, error) {
    report, err := uc.spatialRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    if report.CreatedBy != userID {
        
        
        return nil, ErrUnauthorized
    }

    
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

    
    uc.cache.Delete(ctx, "spatial:"+id.String())
    uc.cache.Delete(ctx, "spatial:list")

    return report, nil
}

func (uc *SpatialPlanningUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateSpatialStatusRequest) error {
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

    
    uc.cache.Delete(ctx, "spatial:"+id.String())
    uc.cache.Delete(ctx, "spatial:stats")

    return nil
}

func (uc *SpatialPlanningUseCase) DeleteReport(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
    report, err := uc.spatialRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    
    if report.CreatedBy != userID {
        return ErrUnauthorized
    }

    
    for _, photo := range report.Photos {
        uc.storage.DeleteFile(ctx, photo.PhotoURL)
    }

    if err := uc.spatialRepo.Delete(ctx, id); err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "spatial:"+id.String())
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