
package repository

import (
	"building-report-backend/internal/domain/entity"
	"context"
)

type ReportRepository interface {
    Create(ctx context.Context, report *entity.Report) error
    Update(ctx context.Context, report *entity.Report) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.Report, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.Report, int64, error)
    FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Report, int64, error)

    GetStatistics(ctx context.Context, buildingType string) (map[string]interface{}, error)
    GetLocationStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error)
    GetWorkTypeStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error)
    GetConditionAfterRehabStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error)
    GetStatusStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error)
    CountByBuildingType(ctx context.Context) ([]map[string]interface{}, error)
}
