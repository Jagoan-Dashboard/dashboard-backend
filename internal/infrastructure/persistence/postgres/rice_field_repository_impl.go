package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type riceFieldRepositoryImpl struct {
	db *gorm.DB
}

func NewRiceFieldRepository(db *gorm.DB) repository.RiceFieldRepository {
	return &riceFieldRepositoryImpl{db: db}
}

func (r *riceFieldRepositoryImpl) Create(ctx context.Context, riceField *entity.RiceField) error {
	return r.db.WithContext(ctx).Create(riceField).Error
}

func (r *riceFieldRepositoryImpl) Update(ctx context.Context, riceField *entity.RiceField) error {
	riceField.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(riceField).Error
}

func (r *riceFieldRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.RiceField{}).Error
}

func (r *riceFieldRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.RiceField, error) {
	var riceField entity.RiceField
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&riceField).Error

	if err != nil {
		return nil, err
	}
	return &riceField, nil
}

func (r *riceFieldRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.RiceField, int64, error) {
	var riceFields []*entity.RiceField
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.RiceField{})

	if district, ok := filters["district"].(string); ok && district != "" {
		query = query.Where("district = ?", district)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("date >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("date <= ?", endDate)
	}

	query.Count(&total)

	err := query.
		Limit(limit).
		Offset(offset).
		Order("date DESC, created_at DESC").
		Find(&riceFields).Error

	return riceFields, total, err
}

func (r *riceFieldRepositoryImpl) FindByDistrict(ctx context.Context, district string, limit, offset int) ([]*entity.RiceField, int64, error) {
	var riceFields []*entity.RiceField
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.RiceField{}).
		Where("district = ?", district)

	query.Count(&total)

	err := query.
		Limit(limit).
		Offset(offset).
		Order("date DESC, created_at DESC").
		Find(&riceFields).Error

	return riceFields, total, err
}

func (r *riceFieldRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.RiceField, int64, error) {
	var riceFields []*entity.RiceField
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.RiceField{}).
		Where("date BETWEEN ? AND ?", startDate, endDate)

	query.Count(&total)

	err := query.
		Limit(limit).
		Offset(offset).
		Order("date DESC, created_at DESC").
		Find(&riceFields).Error

	return riceFields, total, err
}

func (r *riceFieldRepositoryImpl) GetRiceFieldStatistics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalRainfedArea, totalIrrigatedArea float64
	var totalRecords int64

	// Calculate total areas
	err := r.db.WithContext(ctx).Model(&entity.RiceField{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(rainfed_rice_fields), 0) as total_rainfed_area, COALESCE(SUM(irrigated_rice_fields), 0) as total_irrigated_area").
		Scan(&struct {
			TotalRainfedArea   float64 `gorm:"column:total_rainfed_area"`
			TotalIrrigatedArea float64 `gorm:"column:total_irrigated_area"`
		}{totalRainfedArea, totalIrrigatedArea}).Error

	if err != nil {
		return nil, fmt.Errorf("failed to calculate total areas: %w", err)
	}

	// Count total records
	r.db.WithContext(ctx).Model(&entity.RiceField{}).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Count(&totalRecords)

	stats["total_rainfed_area"] = totalRainfedArea
	stats["total_irrigated_area"] = totalIrrigatedArea
	stats["total_rice_field_area"] = totalRainfedArea + totalIrrigatedArea
	stats["total_records"] = totalRecords

	// Calculate growth compared to previous period
	prevStartDate := startDate.AddDate(0, -1, 0)
	prevEndDate := endDate.AddDate(0, -1, 0)

	var prevTotalRainfedArea, prevTotalIrrigatedArea float64
	err = r.db.WithContext(ctx).Model(&entity.RiceField{}).
		Where("date BETWEEN ? AND ?", prevStartDate, prevEndDate).
		Select("COALESCE(SUM(rainfed_rice_fields), 0) as total_rainfed_area, COALESCE(SUM(irrigated_rice_fields), 0) as total_irrigated_area").
		Scan(&struct {
			TotalRainfedArea   float64 `gorm:"column:total_rainfed_area"`
			TotalIrrigatedArea float64 `gorm:"column:total_irrigated_area"`
		}{prevTotalRainfedArea, prevTotalIrrigatedArea}).Error

	if err != nil {
		return nil, fmt.Errorf("failed to calculate previous period areas: %w", err)
	}

	prevTotalArea := prevTotalRainfedArea + prevTotalIrrigatedArea
	currentTotalArea := totalRainfedArea + totalIrrigatedArea

	var rainfedGrowth, irrigatedGrowth, totalGrowth float64
	if prevTotalRainfedArea > 0 {
		rainfedGrowth = ((totalRainfedArea - prevTotalRainfedArea) / prevTotalRainfedArea) * 100
	} else if totalRainfedArea > 0 {
		rainfedGrowth = 100
	}

	if prevTotalIrrigatedArea > 0 {
		irrigatedGrowth = ((totalIrrigatedArea - prevTotalIrrigatedArea) / prevTotalIrrigatedArea) * 100
	} else if totalIrrigatedArea > 0 {
		irrigatedGrowth = 100
	}

	if prevTotalArea > 0 {
		totalGrowth = ((currentTotalArea - prevTotalArea) / prevTotalArea) * 100
	} else if currentTotalArea > 0 {
		totalGrowth = 100
	}

	stats["rainfed_growth"] = rainfedGrowth
	stats["irrigated_growth"] = irrigatedGrowth
	stats["total_growth"] = totalGrowth

	return stats, nil
}

func (r *riceFieldRepositoryImpl) GetRiceFieldDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			district,
			COALESCE(SUM(rainfed_rice_fields), 0) as total_rainfed_area,
			COALESCE(SUM(irrigated_rice_fields), 0) as total_irrigated_area,
			COALESCE(SUM(rainfed_rice_fields), 0) + COALESCE(SUM(irrigated_rice_fields), 0) as total_area,
			COUNT(*) as count
		FROM rice_fields
		WHERE date BETWEEN $1 AND $2
		GROUP BY district
		ORDER BY total_area DESC
	`

	err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get rice field distribution by district: %w", err)
	}

	return results, nil
}

func (r *riceFieldRepositoryImpl) GetIndividualRiceFieldDistribution(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			id,
			latitude,
			longitude,
			district,
			TO_CHAR(date, 'YYYY-MM-DD') as date,
			rainfed_rice_fields,
			irrigated_rice_fields,
			(rainfed_rice_fields + irrigated_rice_fields) as total_area
		FROM rice_fields
		WHERE date BETWEEN $1 AND $2
		ORDER BY date DESC
	`

	err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get individual rice field distribution: %w", err)
	}

	return results, nil
}

func (r *riceFieldRepositoryImpl) GetRiceFieldTrends(ctx context.Context, district string, years []int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	for _, year := range years {
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

		query := `
			SELECT 
				$1 as year,
				COALESCE(SUM(rainfed_rice_fields), 0) as total_rainfed_area,
				COALESCE(SUM(irrigated_rice_fields), 0) as total_irrigated_area,
				COALESCE(SUM(rainfed_rice_fields), 0) + COALESCE(SUM(irrigated_rice_fields), 0) as total_area,
				COUNT(*) as count
			FROM rice_fields
			WHERE date BETWEEN $2 AND $3
		`

		if district != "" {
			query += " AND district = $4"
		}

		var row map[string]interface{}
		var err error

		if district != "" {
			err = r.db.WithContext(ctx).Raw(query, year, startDate, endDate, district).Scan(&row).Error
		} else {
			err = r.db.WithContext(ctx).Raw(query, year, startDate, endDate).Scan(&row).Error
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get rice field trend for year %d: %w", year, err)
		}

		results = append(results, row)
	}

	return results, nil
}
