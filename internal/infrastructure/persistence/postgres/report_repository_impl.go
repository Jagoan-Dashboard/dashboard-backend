package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"	
	"fmt"
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
    query := `DELETE FROM reports WHERE id = $1`
    return r.db.WithContext(ctx).Exec(query, id).Error
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

// Raw SQL Methods for Statistics

func (r *reportRepositoryImpl) GetStatistics(ctx context.Context, buildingType string) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Build base query
    baseWhere := ""
    args := []interface{}{}
    argIndex := 1
    
    if buildingType != "" && buildingType != "all" {
        baseWhere = fmt.Sprintf("WHERE building_type = $%d", argIndex)
        args = append(args, buildingType)
        argIndex++
    }
    
    // 1. Total reports
    query := fmt.Sprintf(`SELECT COUNT(*) FROM reports %s`, baseWhere)
    var totalReports int64
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&totalReports).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count total reports: %w", err)
    }
    stats["total_reports"] = totalReports
    
    if totalReports == 0 {
        stats["average_floor_area"] = float64(0)
        stats["average_floor_count"] = float64(0)
        stats["damaged_buildings_count"] = int64(0)
        return stats, nil
    }
    
    // 2. Average floor area
    query = fmt.Sprintf(`SELECT COALESCE(AVG(floor_area), 0) FROM reports %s`, baseWhere)
    var avgFloorArea float64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&avgFloorArea).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate average floor area: %w", err)
    }
    stats["average_floor_area"] = avgFloorArea
    
    // 3. Average floor count
    query = fmt.Sprintf(`SELECT COALESCE(AVG(floor_count), 0) FROM reports %s`, baseWhere)
    var avgFloorCount float64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&avgFloorCount).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate average floor count: %w", err)
    }
    stats["average_floor_count"] = avgFloorCount
    
    // 4. Damaged buildings count (report_status = 'REHABILITASI')
    var damagedQuery string
    var damagedArgs []interface{}
    if buildingType != "" && buildingType != "all" {
        damagedQuery = `SELECT COUNT(*) FROM reports WHERE report_status = $1 AND building_type = $2`
        damagedArgs = []interface{}{"REHABILITASI", buildingType}
    } else {
        damagedQuery = `SELECT COUNT(*) FROM reports WHERE report_status = $1`
        damagedArgs = []interface{}{"REHABILITASI"}
    }
    
    var damagedCount int64
    err = r.db.WithContext(ctx).Raw(damagedQuery, damagedArgs...).Scan(&damagedCount).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count damaged buildings: %w", err)
    }
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
            COALESCE(AVG(latitude), 0) as avg_latitude,
            COALESCE(AVG(longitude), 0) as avg_longitude,
            COUNT(CASE WHEN report_status = 'REHABILITASI' THEN 1 END) as damaged_count
        FROM reports
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " WHERE building_type = $1"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY district, village
        ORDER BY building_count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get location statistics: %w", err)
    }
    
    return results, nil
}

func (r *reportRepositoryImpl) GetWorkTypeStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN work_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(work_type) = '' THEN 'EMPTY'
                ELSE work_type 
            END as work_type,
            COUNT(*) as count
        FROM reports
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " WHERE building_type = $1"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN work_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(work_type) = '' THEN 'EMPTY'
                ELSE work_type 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get work type statistics: %w", err)
    }
    
    return results, nil
}

func (r *reportRepositoryImpl) GetConditionAfterRehabStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN condition_after_rehab IS NULL THEN 'NOT_SET'
                WHEN TRIM(condition_after_rehab) = '' THEN 'EMPTY'
                ELSE condition_after_rehab 
            END as condition_after_rehab,
            COUNT(*) as count
        FROM reports
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " WHERE building_type = $1"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN condition_after_rehab IS NULL THEN 'NOT_SET'
                WHEN TRIM(condition_after_rehab) = '' THEN 'EMPTY'
                ELSE condition_after_rehab 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get condition statistics: %w", err)
    }
    
    return results, nil
}

func (r *reportRepositoryImpl) GetStatusStatistics(ctx context.Context, buildingType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN report_status IS NULL THEN 'NOT_SET'
                WHEN TRIM(report_status) = '' THEN 'EMPTY'
                ELSE report_status 
            END as report_status,
            COUNT(*) as count
        FROM reports
    `
    
    args := []interface{}{}
    if buildingType != "" && buildingType != "all" {
        query += " WHERE building_type = $1"
        args = append(args, buildingType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN report_status IS NULL THEN 'NOT_SET'
                WHEN TRIM(report_status) = '' THEN 'EMPTY'
                ELSE report_status 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get status statistics: %w", err)
    }
    
    return results, nil
}

func (r *reportRepositoryImpl) CountByBuildingType(ctx context.Context) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN building_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(building_type) = '' THEN 'EMPTY'
                ELSE building_type 
            END as building_type,
            COUNT(*) as count
        FROM reports
        GROUP BY 
            CASE 
                WHEN building_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(building_type) = '' THEN 'EMPTY'
                ELSE building_type 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count by building type: %w", err)
    }
    
    return results, nil
}