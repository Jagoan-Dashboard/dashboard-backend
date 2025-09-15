
package postgres

import (
    "context"
    "building-report-backend/internal/domain/entity"
    "building-report-backend/internal/domain/repository"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type spatialPlanningRepositoryImpl struct {
    db *gorm.DB
}

func NewSpatialPlanningRepository(db *gorm.DB) repository.SpatialPlanningRepository {
    return &spatialPlanningRepositoryImpl{db: db}
}

func (r *spatialPlanningRepositoryImpl) Create(ctx context.Context, report *entity.SpatialPlanningReport) error {
    return r.db.WithContext(ctx).Create(report).Error
}

func (r *spatialPlanningRepositoryImpl) Update(ctx context.Context, report *entity.SpatialPlanningReport) error {
    return r.db.WithContext(ctx).Save(report).Error
}

func (r *spatialPlanningRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.SpatialPlanningReport{}).Error
}

func (r *spatialPlanningRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.SpatialPlanningReport, error) {
    var report entity.SpatialPlanningReport
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("id = ?", id).
        First(&report).Error
    
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *spatialPlanningRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.SpatialPlanningReport, int64, error) {
    var reports []*entity.SpatialPlanningReport
    var total int64

    query := r.db.WithContext(ctx).Model(&entity.SpatialPlanningReport{})

    
    if institution, ok := filters["institution"].(string); ok && institution != "" {
        query = query.Where("institution = ?", institution)
    }
    if areaCategory, ok := filters["area_category"].(string); ok && areaCategory != "" {
        query = query.Where("area_category = ?", areaCategory)
    }
    if violationType, ok := filters["violation_type"].(string); ok && violationType != "" {
        query = query.Where("violation_type = ?", violationType)
    }
    if violationLevel, ok := filters["violation_level"].(string); ok && violationLevel != "" {
        query = query.Where("violation_level = ?", violationLevel)
    }
    if urgencyLevel, ok := filters["urgency_level"].(string); ok && urgencyLevel != "" {
        query = query.Where("urgency_level = ?", urgencyLevel)
    }
    if status, ok := filters["status"].(string); ok && status != "" {
        query = query.Where("status = ?", status)
    }

    
    if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
        query = query.Where("report_datetime >= ?", startDate)
    }
    if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
        query = query.Where("report_datetime <= ?", endDate)
    }

    
    query.Count(&total)

    
    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("urgency_level DESC, created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *spatialPlanningRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.SpatialPlanningReport, int64, error) {
    var reports []*entity.SpatialPlanningReport
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.SpatialPlanningReport{}).
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

func (r *spatialPlanningRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SpatialReportStatus) error {
    return r.db.WithContext(ctx).
        Model(&entity.SpatialPlanningReport{}).
        Where("id = ?", id).
        Update("status", status).Error
}

func (r *spatialPlanningRepositoryImpl) CountByStatus(ctx context.Context, status entity.SpatialReportStatus) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&entity.SpatialPlanningReport{}).
        Where("status = ?", status).
        Count(&count).Error
    return count, err
}

func (r *spatialPlanningRepositoryImpl) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    
    var total int64
    r.db.Model(&entity.SpatialPlanningReport{}).Count(&total)
    stats["total_reports"] = total
    
    
    var urgentCount int64
    r.db.Model(&entity.SpatialPlanningReport{}).
        Where("urgency_level = ?", entity.UrgencyMendesak).
        Count(&urgentCount)
    stats["urgent_reports"] = urgentCount
    
    
    type ViolationCount struct {
        Level string
        Count int64
    }
    var violationCounts []ViolationCount
    r.db.Model(&entity.SpatialPlanningReport{}).
        Select("violation_level as level, count(*) as count").
        Group("violation_level").
        Scan(&violationCounts)
    stats["violation_levels"] = violationCounts
    
    
    type StatusCount struct {
        Status string
        Count  int64
    }
    var statusCounts []StatusCount
    r.db.Model(&entity.SpatialPlanningReport{}).
        Select("status, count(*) as count").
        Group("status").
        Scan(&statusCounts)
    stats["status_counts"] = statusCounts
    
    return stats, nil
}