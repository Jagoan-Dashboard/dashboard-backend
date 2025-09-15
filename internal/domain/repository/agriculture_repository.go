
package repository

import (
    "context"
    "time"
    "building-report-backend/internal/domain/entity"
    "github.com/google/uuid"
)

type AgricultureRepository interface {
    Create(ctx context.Context, report *entity.AgricultureReport) error
    Update(ctx context.Context, report *entity.AgricultureReport) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.AgricultureReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.AgricultureReport, int64, error)
    FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    FindByExtensionOfficer(ctx context.Context, extensionOfficer string, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    FindByVillage(ctx context.Context, village string, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    FindByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    
    
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
    GetCommodityProduction(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetExtensionOfficerPerformance(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetVillageProductionStats(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetPestDiseaseReports(ctx context.Context, limit int) ([]*entity.AgricultureReport, error)
    GetTechnologyAdoptionStats(ctx context.Context) (map[string]interface{}, error)
    GetFarmerNeedsAnalysis(ctx context.Context) (map[string]interface{}, error)
    
    
    CountTotalFarmers(ctx context.Context) (int64, error)
    CalculateTotalLandArea(ctx context.Context) (float64, error)
    CountReportsByCommodityType(ctx context.Context) (map[string]int64, error)
    CountReportsWithPestDisease(ctx context.Context) (int64, error)
    GetTopConstraints(ctx context.Context, limit int) ([]map[string]interface{}, error)
    GetTopFarmerHopes(ctx context.Context, limit int) ([]map[string]interface{}, error)
}