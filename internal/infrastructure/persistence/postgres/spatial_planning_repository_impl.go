package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"

	"gorm.io/gorm"
)

type spatialPlanningRepositoryImpl struct {
	db *gorm.DB
}

func NewSpatialPlanningRepository(db *gorm.DB) repository.SpatialPlanningRepository {
	return &spatialPlanningRepositoryImpl{db: db}
}

func (r *spatialPlanningRepositoryImpl) Create(ctx context.Context, report *entity.SpatialPlanningReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *spatialPlanningRepositoryImpl) Update(ctx context.Context, report *entity.SpatialPlanningReport) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *spatialPlanningRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.SpatialPlanningReport{}).Error
}

func (r *spatialPlanningRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.SpatialPlanningReport, error) {
	var report entity.SpatialPlanningReport
	err := r.db.WithContext(ctx).
		Preload("Photos").
		Where("id = ?", id).
		First(&report).Error

	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *spatialPlanningRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.SpatialPlanningReport, int64, error) {
	var reports []*entity.SpatialPlanningReport
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SpatialPlanningReport{})

	if institution, ok := filters["institution"].(string); ok && institution != "" {
		query = query.Where("institution = ?", institution)
	}
	if areaCategory, ok := filters["area_category"].(string); ok && areaCategory != "" {
		query = query.Where("area_category = ?", areaCategory)
	}
	if violationType, ok := filters["violation_type"].(string); ok && violationType != "" {
		query = query.Where("violation_type = ?", violationType)
	}
	if violationLevel, ok := filters["violation_level"].(string); ok && violationLevel != "" {
		query = query.Where("violation_level = ?", violationLevel)
	}
	if urgencyLevel, ok := filters["urgency_level"].(string); ok && urgencyLevel != "" {
		query = query.Where("urgency_level = ?", urgencyLevel)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("report_datetime >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("report_datetime <= ?", endDate)
	}

	query.Count(&total)

	err := query.
		Preload("Photos").
		Limit(limit).
		Offset(offset).
		Order("urgency_level DESC, created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

func (r *spatialPlanningRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.SpatialPlanningReport, int64, error) {
	var reports []*entity.SpatialPlanningReport
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.SpatialPlanningReport{})

	query.Count(&total)

	err := query.
		Preload("Photos").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

func (r *spatialPlanningRepositoryImpl) UpdateStatus(ctx context.Context, id string, status entity.SpatialReportStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.SpatialPlanningReport{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *spatialPlanningRepositoryImpl) CountByStatus(ctx context.Context, status entity.SpatialReportStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.SpatialPlanningReport{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *spatialPlanningRepositoryImpl) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total int64
	r.db.Model(&entity.SpatialPlanningReport{}).Count(&total)
	stats["total_reports"] = total

	var urgentCount int64
	r.db.Model(&entity.SpatialPlanningReport{}).
		Where("urgency_level = ?", entity.UrgencyMendesak).
		Count(&urgentCount)
	stats["urgent_reports"] = urgentCount

	type ViolationCount struct {
		Level string
		Count int64
	}
	var violationCounts []ViolationCount
	r.db.Model(&entity.SpatialPlanningReport{}).
		Select("violation_level as level, count(*) as count").
		Group("violation_level").
		Scan(&violationCounts)
	stats["violation_levels"] = violationCounts

	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	r.db.Model(&entity.SpatialPlanningReport{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts)
	stats["status_counts"] = statusCounts

	return stats, nil
}

func (r *spatialPlanningRepositoryImpl) GetTataRuangStatistics(ctx context.Context, areaCategory string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.WithContext(ctx).Model(&entity.SpatialPlanningReport{})
	if areaCategory != "" && areaCategory != "all" {
		query = query.Where("area_category = ?", areaCategory)
	}

	// Total reports
	var totalReports int64
	query.Count(&totalReports)
	stats["total_reports"] = totalReports

	// Since we don't have exact length/area measurements in current schema,
	// we'll calculate estimated values based on violation level and type
	// You can modify this based on actual field names if they exist

	// Estimated total violation length (hypothetical calculation)
	var estimatedTotalLength float64
	query.Select(`
        COALESCE(SUM(
            CASE 
                WHEN violation_level = 'BERAT' THEN 100.0
                WHEN violation_level = 'SEDANG' THEN 50.0
                WHEN violation_level = 'RINGAN' THEN 20.0
                ELSE 30.0
            END
        ), 0) as estimated_length
    `).Scan(&estimatedTotalLength)
	stats["estimated_total_length_m"] = estimatedTotalLength

	// Estimated total violation area (hypothetical calculation)
	var estimatedTotalArea float64
	query.Select(`
        COALESCE(SUM(
            CASE 
                WHEN violation_level = 'BERAT' AND violation_type LIKE '%SEMPADAN%' THEN 500.0
                WHEN violation_level = 'BERAT' THEN 1000.0
                WHEN violation_level = 'SEDANG' AND violation_type LIKE '%SEMPADAN%' THEN 200.0
                WHEN violation_level = 'SEDANG' THEN 400.0
                WHEN violation_level = 'RINGAN' THEN 100.0
                ELSE 150.0
            END
        ), 0) as estimated_area
    `).Scan(&estimatedTotalArea)
	stats["estimated_total_area_m2"] = estimatedTotalArea

	// Count urgent reports
	var urgentCount int64
	query.Where("urgency_level = ?", "MENDESAK").Count(&urgentCount)
	stats["urgent_reports_count"] = urgentCount

	return stats, nil
}

func (r *spatialPlanningRepositoryImpl) GetLocationDistribution(ctx context.Context, areaCategory string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
        SELECT 
            COALESCE(SPLIT_PART(address, ',', -1), 'Unknown') as district,
            COALESCE(SPLIT_PART(address, ',', -2), 'Unknown') as village,
            COUNT(*) as violation_count,
            AVG(latitude) as avg_latitude,
            AVG(longitude) as avg_longitude,
            COUNT(CASE WHEN urgency_level = 'MENDESAK' THEN 1 END) as urgent_count,
            COUNT(CASE WHEN violation_level = 'BERAT' THEN 1 END) as severe_count
        FROM spatial_planning_reports
        WHERE latitude IS NOT NULL AND longitude IS NOT NULL
    `

	args := []interface{}{}
	if areaCategory != "" && areaCategory != "all" {
		query += " AND area_category = ?"
		args = append(args, areaCategory)
	}

	query += `
        GROUP BY district, village
        HAVING AVG(latitude) IS NOT NULL AND AVG(longitude) IS NOT NULL
        ORDER BY violation_count DESC
    `

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
	return results, err
}

func (r *spatialPlanningRepositoryImpl) GetUrgencyLevelStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
        SELECT 
            urgency_level,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM spatial_planning_reports
    `

	args := []interface{}{}
	if areaCategory != "" && areaCategory != "all" {
		query += " WHERE area_category = ?"
		args = append(args, areaCategory)
	}

	query += `
        GROUP BY urgency_level
        ORDER BY count DESC
    `

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
	return results, err
}

func (r *spatialPlanningRepositoryImpl) GetViolationTypeStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
        SELECT 
            violation_type,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage,
            COUNT(CASE WHEN violation_level = 'BERAT' THEN 1 END) as severe_count,
            COUNT(CASE WHEN urgency_level = 'MENDESAK' THEN 1 END) as urgent_count
        FROM spatial_planning_reports
    `

	args := []interface{}{}
	if areaCategory != "" && areaCategory != "all" {
		query += " WHERE area_category = ?"
		args = append(args, areaCategory)
	}

	query += `
        GROUP BY violation_type
        ORDER BY count DESC
    `

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
	return results, err
}

func (r *spatialPlanningRepositoryImpl) GetViolationLevelStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
        SELECT 
            violation_level,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage,
            COUNT(CASE WHEN urgency_level = 'MENDESAK' THEN 1 END) as urgent_count
        FROM spatial_planning_reports
    `

	args := []interface{}{}
	if areaCategory != "" && areaCategory != "all" {
		query += " WHERE area_category = ?"
		args = append(args, areaCategory)
	}

	query += `
        GROUP BY violation_level
        ORDER BY 
            CASE violation_level 
                WHEN 'BERAT' THEN 1 
                WHEN 'SEDANG' THEN 2 
                WHEN 'RINGAN' THEN 3 
                ELSE 4 
            END
    `

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
	return results, err
}

func (r *spatialPlanningRepositoryImpl) GetAreaCategoryDistribution(ctx context.Context) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
        SELECT 
            area_category,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage,
            COUNT(CASE WHEN urgency_level = 'MENDESAK' THEN 1 END) as urgent_count,
            COUNT(CASE WHEN violation_level = 'BERAT' THEN 1 END) as severe_count
        FROM spatial_planning_reports
        GROUP BY area_category
        ORDER BY count DESC
    `

	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	return results, err
}

func (r *spatialPlanningRepositoryImpl) GetEnvironmentalImpactStatistics(ctx context.Context, areaCategory string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
        SELECT 
            environmental_impact,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage,
            COUNT(CASE WHEN violation_level = 'BERAT' THEN 1 END) as severe_count
        FROM spatial_planning_reports
    `

	args := []interface{}{}
	if areaCategory != "" && areaCategory != "all" {
		query += " WHERE area_category = ?"
		args = append(args, areaCategory)
	}

	query += `
        GROUP BY environmental_impact
        ORDER BY count DESC
    `

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
	return results, err
}
