
package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"
	"gorm.io/gorm"
)

type reportRepositoryImpl struct {
    db *gorm.DB
}

func NewReportRepository(db *gorm.DB) repository.ReportRepository {
    return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) Create(ctx context.Context, report *entity.Report) error {
    return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepositoryImpl) Update(ctx context.Context, report *entity.Report) error {
    return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepositoryImpl) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.Report{}).Error
}

func (r *reportRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Report, error) {
    var report entity.Report
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("id = ?", id).
        First(&report).Error
    
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *reportRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.Report, int64, error) {
    var reports []*entity.Report
    var total int64

    query := r.db.WithContext(ctx).Model(&entity.Report{})

    
    if village, ok := filters["village"].(string); ok && village != "" {
        query = query.Where("village ILIKE ?", "%"+village+"%")
    }
    if district, ok := filters["district"].(string); ok && district != "" {
        query = query.Where("district ILIKE ?", "%"+district+"%")
    }
    if buildingType, ok := filters["building_type"].(string); ok && buildingType != "" {
        query = query.Where("building_type = ?", buildingType)
    }
    if reportStatus, ok := filters["report_status"].(string); ok && reportStatus != "" {
        query = query.Where("report_status = ?", reportStatus)
    }

    
    query.Count(&total)

    
    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *reportRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Report, int64, error) {
    var reports []*entity.Report
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.Report{}).
        Where("created_by = ?", userID)

    query.Count(&total)

    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *reportRepositoryImpl) GetStatistics(ctx context.Context, buildingType string) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    query := r.db.WithContext(ctx).Model(&entity.Report{})
    if buildingType != "" && buildingType != "all" {
        query = query.Where("building_type = ?", buildingType)
    }
    
    // Total reports
    var totalReports int64
    query.Count(&totalReports)
    stats["total_reports"] = totalReports
    
    // Average floor area
    var avgFloorArea float64
    query.Select("COALESCE(AVG(floor_area), 0)").Scan(&avgFloorArea)
    stats["average_floor_area"] = avgFloorArea
    
    // Average floor count
    var avgFloorCount float64
    query.Select("COALESCE(AVG(floor_count), 0)").Scan(&avgFloorCount)
    stats["average_floor_count"] = avgFloorCount
    
    // Count damaged buildings (assuming buildings needing rehabilitation are "damaged")
    var damagedCount int64
    query.Where("report_status = ?", "REHABILITASI").Count(&damagedCount)
    stats["damaged_buildings_count"] = damagedCount
    
    return stats, nil
}

func (r *reportRepositoryImpl) GetLocationStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            district,
            village,
            COUNT(*) as building_count,
            AVG(latitude) as avg_latitude,
            AVG(longitude) as avg_longitude,
            COUNT(CASE WHEN report_status = 'REHABILITASI' THEN 1 END) as damaged_count
        FROM reports
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " WHERE building_type = ?"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY district, village
        ORDER BY building_count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *reportRepositoryImpl) GetWorkTypeStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            work_type,
            COUNT(*) as count
        FROM reports
        WHERE work_type IS NOT NULL AND work_type != ''
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " AND building_type = ?"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY work_type
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *reportRepositoryImpl) GetConditionAfterRehabStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            condition_after_rehab,
            COUNT(*) as count
        FROM reports
        WHERE condition_after_rehab IS NOT NULL AND condition_after_rehab != ''
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " AND building_type = ?"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY condition_after_rehab
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *reportRepositoryImpl) GetStatusStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            report_status,
            COUNT(*) as count
        FROM reports
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " WHERE building_type = ?"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY report_status
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *reportRepositoryImpl) CountByBuildingType(ctx context.Context) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            building_type,
            COUNT(*) as count
        FROM reports
        GROUP BY building_type
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
    return results, err
}