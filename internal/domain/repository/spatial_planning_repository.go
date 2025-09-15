package repository

import (
    "context"
    "building-report-backend/internal/domain/entity"
    "github.com/google/uuid"
)

type SpatialPlanningRepository interface {
    Create(ctx context.Context, report *entity.SpatialPlanningReport) error
    Update(ctx context.Context, report *entity.SpatialPlanningReport) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.SpatialPlanningReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.SpatialPlanningReport, int64, error)
    FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.SpatialPlanningReport, int64, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status entity.SpatialReportStatus) error
    CountByStatus(ctx context.Context, status entity.SpatialReportStatus) (int64, error)
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
}