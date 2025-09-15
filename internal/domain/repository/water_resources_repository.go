
package repository

import (
    "context"
    "time"
    "building-report-backend/internal/domain/entity"
    "github.com/google/uuid"
)

type WaterResourcesRepository interface {
    Create(ctx context.Context, report *entity.WaterResourcesReport) error
    Update(ctx context.Context, report *entity.WaterResourcesReport) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.WaterResourcesReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.WaterResourcesReport, int64, error)
    FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.WaterResourcesReport, int64, error)
    FindByPriority(ctx context.Context, limit, offset int) ([]*entity.WaterResourcesReport, int64, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status entity.WaterResourceStatus, notes string) error
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
    GetDamageStatisticsByArea(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetUrgentReports(ctx context.Context, limit int) ([]*entity.WaterResourcesReport, error)
    CalculateTotalDamageArea(ctx context.Context) (float64, error)
    CountAffectedFarmers(ctx context.Context) (int64, error)
}