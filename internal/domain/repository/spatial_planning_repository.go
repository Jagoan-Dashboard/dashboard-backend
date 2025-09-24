package repository

import (
    "context"
    "building-report-backend/internal/domain/entity"
)

type SpatialPlanningRepository interface {
    Create(ctx context.Context, report *entity.SpatialPlanningReport) error
    Update(ctx context.Context, report *entity.SpatialPlanningReport) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.SpatialPlanningReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.SpatialPlanningReport, int64, error)
    FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.SpatialPlanningReport, int64, error)
    UpdateStatus(ctx context.Context, id string, status entity.SpatialReportStatus) error
    CountByStatus(ctx context.Context, status entity.SpatialReportStatus) (int64, error)
    GetStatistics(ctx context.Context) (map[string]interface{}, error)

    GetTataRuangStatistics(ctx context.Context, areaCategory string) (map[string]interface{}, error)
    GetLocationDistribution(ctx context.Context, areaCategory string) ([]map[string]interface{}, error)
    GetUrgencyLevelStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error)
    GetViolationTypeStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error)
    GetViolationLevelStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error)
    GetAreaCategoryDistribution(ctx context.Context) ([]map[string]interface{}, error)
    GetEnvironmentalImpactStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error)
}