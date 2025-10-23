
package repository

import (
    "context"
    "time"
    "building-report-backend/internal/domain/entity"
)

type BinaMargaRepository interface {
    Create(ctx context.Context, report *entity.BinaMargaReport) error
    Update(ctx context.Context, report *entity.BinaMargaReport) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.BinaMargaReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.BinaMargaReport, int64, error)
    FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.BinaMargaReport, int64, error)
    FindByPriority(ctx context.Context, limit, offset int) ([]*entity.BinaMargaReport, int64, error)
    FindEmergencyReports(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error)
    FindBlockedRoads(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error)
    UpdateStatus(ctx context.Context, id string, status entity.BinaMargaStatus, notes string) error
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
    GetDamageStatisticsByRoadType(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetDamageStatisticsByLocation(ctx context.Context, bounds map[string]float64) ([]map[string]interface{}, error)
    CalculateTotalDamageArea(ctx context.Context) (float64, error)
    CalculateTotalDamageLength(ctx context.Context) (float64, error)
    CountReportsByUrgency(ctx context.Context, urgency entity.RoadUrgencyLevel) (int64, error)
    GetRepairTimeAnalysis(ctx context.Context) (map[string]interface{}, error)
    GetBinaMargaOverviewStats(ctx context.Context, roadType string) (map[string]interface{}, error)
    GetBinaMargaLocationStats(ctx context.Context, roadType string) ([]map[string]interface{}, error)
    GetBinaMargaPriorityStats(ctx context.Context, roadType string) ([]map[string]interface{}, error)
    GetBinaMargaRoadDamageLevelStats(ctx context.Context, roadType string) ([]map[string]interface{}, error)
    GetBinaMargaBridgeDamageLevelStats(ctx context.Context, roadType string) ([]map[string]interface{}, error)
    GetBinaMargaTopRoadDamageTypes(ctx context.Context, roadType string) ([]map[string]interface{}, error)
     GetBinaMargaTopBridgeDamageTypes(ctx context.Context, roadType string) ([]map[string]interface{}, error)
    
    GetKPIs(ctx context.Context, roadType string, startDate, endDate time.Time) (avgSegLen, avgDamageArea, avgDailyTraffic float64, totalReports int64, err error)

    // generic group-by untuk kolom tertentu, dengan opsi filter jembatan
    GroupCountBy(ctx context.Context, column, roadType string, startDate, endDate time.Time, onlyBridge, onlyRoad bool) ([]struct {
        Key   string
        Count int64
    }, error)

    GetMapPoints(ctx context.Context, roadType string, startDate, endDate time.Time) ([]struct {
        Latitude           float64
        Longitude          float64
        District           string
        RoadName           string
        RoadType           string
        DamageType         string
        DamageLevel        string
        BridgeName         *string
        BridgeSection      *string
        BridgeDamageType   *string
        BridgeDamageLevel  *string
        UrgencyLevel       string
    }, error)
}