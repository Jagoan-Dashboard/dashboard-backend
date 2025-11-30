package repository

import (
	"building-report-backend/internal/domain/entity"
	"context"
	"time"
)

type RiceFieldRepository interface {
	Create(ctx context.Context, riceField *entity.RiceField) error
	Update(ctx context.Context, riceField *entity.RiceField) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.RiceField, error)
	FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.RiceField, int64, error)
	FindByDistrict(ctx context.Context, district string, limit, offset int) ([]*entity.RiceField, int64, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.RiceField, int64, error)
	
	// Rice field specific statistics
	GetRiceFieldStatistics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
	GetRiceFieldDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetIndividualRiceFieldDistribution(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetRiceFieldTrends(ctx context.Context, district string, years []int) ([]map[string]interface{}, error)
}