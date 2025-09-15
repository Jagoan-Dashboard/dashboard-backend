
package repository

import (
    "context"
    "time"
    "building-report-backend/internal/domain/entity"
    "github.com/google/uuid"
)

type BinaMargaRepository interface {
    Create(ctx context.Context, report *entity.BinaMargaReport) error
    Update(ctx context.Context, report *entity.BinaMargaReport) error
    Delete(ctx context.Context, id uuid.UUID) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.BinaMargaReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.BinaMargaReport, int64, error)
    FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.BinaMargaReport, int64, error)
    FindByPriority(ctx context.Context, limit, offset int) ([]*entity.BinaMargaReport, int64, error)
    FindEmergencyReports(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error)
    FindBlockedRoads(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BinaMargaStatus, notes string) error
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
    GetDamageStatisticsByRoadType(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetDamageStatisticsByLocation(ctx context.Context, bounds map[string]float64) ([]map[string]interface{}, error)
    CalculateTotalDamageArea(ctx context.Context) (float64, error)
    CalculateTotalDamageLength(ctx context.Context) (float64, error)
    CountReportsByUrgency(ctx context.Context, urgency entity.RoadUrgencyLevel) (int64, error)
    GetRepairTimeAnalysis(ctx context.Context) (map[string]interface{}, error)
}