// internal/application/usecase/report_usecase.go
package usecase

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/storage"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
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

func (uc *ReportUseCase) CreateReport(ctx context.Context, req *dto.CreateReportRequest, photos []*multipart.FileHeader, userID uuid.UUID) (*entity.Report, error) {
    report := &entity.Report{
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
        CreatedBy:            userID,
    }

    // Handle optional fields
    if req.WorkType != "" {
        workType := entity.WorkType(req.WorkType)
        report.WorkType = &workType
    }
    if req.ConditionAfterRehab != "" {
        condition := entity.ConditionAfterRehab(req.ConditionAfterRehab)
        report.ConditionAfterRehab = &condition
    }

    // Upload photos to MinIO
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
            PhotoURL:  photoURL,
            PhotoType: photoType,
        })
    }

    if err := uc.reportRepo.Create(ctx, report); err != nil {
        return nil, err
    }

    // Invalidate cache
    uc.cache.Delete(ctx, "reports:list")

    return report, nil
}

func (uc *ReportUseCase) GetReport(ctx context.Context, id uuid.UUID) (*entity.Report, error) {
    // Check cache first
    cacheKey := "report:" + id.String()
    
    report, err := uc.reportRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Cache the result
    uc.cache.Set(ctx, cacheKey, report, 3600) // 1 hour

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

func (uc *ReportUseCase) UpdateReport(ctx context.Context, id uuid.UUID, req *dto.UpdateReportRequest, userID uuid.UUID) (*entity.Report, error) {
    report, err := uc.reportRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Check permission
    if report.CreatedBy != userID {
        return nil, ErrUnauthorized
    }

    // Update fields
    if req.BuildingName != "" {
        report.BuildingName = req.BuildingName
    }
    if req.ReportStatus != "" {
        report.ReportStatus = entity.ReportStatusType(req.ReportStatus)
    }
    // ... update other fields

    if err := uc.reportRepo.Update(ctx, report); err != nil {
        return nil, err
    }

    // Invalidate cache
    uc.cache.Delete(ctx, "report:"+id.String())
    uc.cache.Delete(ctx, "reports:list")

    return report, nil
}

func (uc *ReportUseCase) DeleteReport(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
    report, err := uc.reportRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    // Check permission
    if report.CreatedBy != userID {
        return ErrUnauthorized
    }

    // Delete photos from MinIO
    for _, photo := range report.Photos {
        uc.storage.DeleteFile(ctx, photo.PhotoURL)
    }

    if err := uc.reportRepo.Delete(ctx, id); err != nil {
        return err
    }

    // Invalidate cache
    uc.cache.Delete(ctx, "report:"+id.String())
    uc.cache.Delete(ctx, "reports:list")

    return nil
}