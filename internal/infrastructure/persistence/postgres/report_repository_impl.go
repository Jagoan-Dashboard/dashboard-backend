
package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reportRepositoryImpl struct {
    db *gorm.DB
}

func NewReportRepository(db *gorm.DB) repository.ReportRepository {
    return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) Create(ctx context.Context, report *entity.Report) error {
    return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepositoryImpl) Update(ctx context.Context, report *entity.Report) error {
    return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.Report{}).Error
}

func (r *reportRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Report, error) {
    var report entity.Report
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("id = ?", id).
        First(&report).Error
    
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *reportRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.Report, int64, error) {
    var reports []*entity.Report
    var total int64

    query := r.db.WithContext(ctx).Model(&entity.Report{})

    
    if village, ok := filters["village"].(string); ok && village != "" {
        query = query.Where("village ILIKE ?", "%"+village+"%")
    }
    if district, ok := filters["district"].(string); ok && district != "" {
        query = query.Where("district ILIKE ?", "%"+district+"%")
    }
    if buildingType, ok := filters["building_type"].(string); ok && buildingType != "" {
        query = query.Where("building_type = ?", buildingType)
    }
    if reportStatus, ok := filters["report_status"].(string); ok && reportStatus != "" {
        query = query.Where("report_status = ?", reportStatus)
    }

    
    query.Count(&total)

    
    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *reportRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Report, int64, error) {
    var reports []*entity.Report
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.Report{}).
        Where("created_by = ?", userID)

    query.Count(&total)

    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&reports).Error

    return reports, total, err
}