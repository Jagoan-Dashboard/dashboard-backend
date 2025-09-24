package repository

import (
	"building-report-backend/internal/domain/entity"
	"context"
	"time"

	
)

type WaterResourcesRepository interface {
    Create(ctx context.Context, report *entity.WaterResourcesReport) error
    Update(ctx context.Context, report *entity.WaterResourcesReport) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.WaterResourcesReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.WaterResourcesReport, int64, error)
    FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.WaterResourcesReport, int64, error)
    FindByPriority(ctx context.Context, limit, offset int) ([]*entity.WaterResourcesReport, int64, error)
    UpdateStatus(ctx context.Context, id string, status entity.WaterResourceStatus, notes string) error
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
    GetDamageStatisticsByArea(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetUrgentReports(ctx context.Context, limit int) ([]*entity.WaterResourcesReport, error)
    CalculateTotalDamageArea(ctx context.Context) (float64, error)
    CountAffectedFarmers(ctx context.Context) (int64, error)
    GetWaterResourcesOverviewStats(ctx context.Context, irrigationType string) (map[string]interface{}, error)
    GetWaterLocationStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error)
    GetWaterUrgencyStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error) 
    GetWaterDamageTypeStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error)
    GetWaterDamageLevelStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error) 

    GetSummaryKPIs(ctx context.Context, irrigationType string, startDate, endDate time.Time) (totalAreaM2 float64, totalRiceHa float64, totalReports int64, err error)
    GroupCountBy(ctx context.Context, field, irrigationType string, startDate, endDate time.Time) ([]struct {
        Key   string
        Count int64
    }, error)
    GetMapPoints(ctx context.Context, irrigationType string, startDate, endDate time.Time) ([]struct {
        Latitude        float64
        Longitude       float64
        IrrigationArea  string
        DamageType      string
        DamageLevel     string
        UrgencyCategory string
    }, error)
}