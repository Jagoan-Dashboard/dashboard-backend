package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type waterResourcesRepositoryImpl struct {
    db *gorm.DB
}

func NewWaterResourcesRepository(db *gorm.DB) repository.WaterResourcesRepository {
    return &waterResourcesRepositoryImpl{db: db}
}

// Keep GORM for complex operations with relations
func (r *waterResourcesRepositoryImpl) Create(ctx context.Context, report *entity.WaterResourcesReport) error {
    return r.db.WithContext(ctx).Create(report).Error
}

func (r *waterResourcesRepositoryImpl) Update(ctx context.Context, report *entity.WaterResourcesReport) error {
    report.UpdatedAt = time.Now()
    return r.db.WithContext(ctx).Save(report).Error
}

func (r *waterResourcesRepositoryImpl) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM water_resources_reports WHERE id = $1`
    return r.db.WithContext(ctx).Exec(query, id).Error
}

func (r *waterResourcesRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.WaterResourcesReport, error) {
    var report entity.WaterResourcesReport
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("id = ?", id).
        First(&report).Error
    
    if err != nil {
        return nil, err
    }
    return &report, nil
}

// Raw SQL implementations for all statistics methods

func (r *waterResourcesRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.WaterResourcesReport, int64, error) {
    var reports []*entity.WaterResourcesReport
    var total int64

    // Build dynamic query with filters
    baseQuery := "FROM water_resources_reports"
    whereClause := ""
    args := []interface{}{}
    argIndex := 1

    // Apply filters
    conditions := []string{}
    
    if institutionUnit, ok := filters["institution_unit"].(string); ok && institutionUnit != "" {
        conditions = append(conditions, fmt.Sprintf("institution_unit = $%d", argIndex))
        args = append(args, institutionUnit)
        argIndex++
    }
    if irrigationType, ok := filters["irrigation_type"].(string); ok && irrigationType != "" {
        conditions = append(conditions, fmt.Sprintf("irrigation_type = $%d", argIndex))
        args = append(args, irrigationType)
        argIndex++
    }
    if damageType, ok := filters["damage_type"].(string); ok && damageType != "" {
        conditions = append(conditions, fmt.Sprintf("damage_type = $%d", argIndex))
        args = append(args, damageType)
        argIndex++
    }
    if damageLevel, ok := filters["damage_level"].(string); ok && damageLevel != "" {
        conditions = append(conditions, fmt.Sprintf("damage_level = $%d", argIndex))
        args = append(args, damageLevel)
        argIndex++
    }
    if urgencyCategory, ok := filters["urgency_category"].(string); ok && urgencyCategory != "" {
        conditions = append(conditions, fmt.Sprintf("urgency_category = $%d", argIndex))
        args = append(args, urgencyCategory)
        argIndex++
    }
    if status, ok := filters["status"].(string); ok && status != "" {
        conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
        args = append(args, status)
        argIndex++
    }
    if irrigationArea, ok := filters["irrigation_area"].(string); ok && irrigationArea != "" {
        conditions = append(conditions, fmt.Sprintf("irrigation_area_name ILIKE $%d", argIndex))
        args = append(args, "%"+irrigationArea+"%")
        argIndex++
    }
    if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
        conditions = append(conditions, fmt.Sprintf("report_datetime >= $%d", argIndex))
        args = append(args, startDate)
        argIndex++
    }
    if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
        conditions = append(conditions, fmt.Sprintf("report_datetime <= $%d", argIndex))
        args = append(args, endDate)
        argIndex++
    }

    if len(conditions) > 0 {
        whereClause = " WHERE " + fmt.Sprintf("%s", conditions[0])
        for i := 1; i < len(conditions); i++ {
            whereClause += " AND " + conditions[i]
        }
    }

    // Count total
    countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
    err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count records: %w", err)
    }

    // Get data with pagination - using GORM for preload complexity
    query := r.db.WithContext(ctx).Model(&entity.WaterResourcesReport{})
    for i, condition := range conditions {
        if i == 0 {
            query = query.Where(condition, args[i])
        } else {
            query = query.Where(condition, args[i])
        }
    }

    err = query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("CASE WHEN urgency_category = 'MENDESAK' THEN 0 ELSE 1 END, " +
              "CASE WHEN damage_level = 'BERAT' THEN 0 WHEN damage_level = 'SEDANG' THEN 1 ELSE 2 END, " +
              "created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *waterResourcesRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.WaterResourcesReport, int64, error) {
    var reports []*entity.WaterResourcesReport
    var total int64

    // Count total
    countQuery := `SELECT COUNT(*) FROM water_resources_reports WHERE created_by = $1`
    err := r.db.WithContext(ctx).Raw(countQuery, userID).Scan(&total).Error
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count user reports: %w", err)
    }

    // Get data - using GORM for preload
    err = r.db.WithContext(ctx).
        Preload("Photos").
        Where("created_by = ?", userID).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *waterResourcesRepositoryImpl) FindByPriority(ctx context.Context, limit, offset int) ([]*entity.WaterResourcesReport, int64, error) {
    var reports []*entity.WaterResourcesReport
    var total int64

    // Count total with priority filter
    countQuery := `
        SELECT COUNT(*) 
        FROM water_resources_reports 
        WHERE status NOT IN ('COMPLETED', 'REJECTED')
    `
    err := r.db.WithContext(ctx).Raw(countQuery).Scan(&total).Error
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count priority reports: %w", err)
    }

    // Priority calculation query
    priorityQuery := `
        SELECT *, 
            (CASE WHEN urgency_category = 'MENDESAK' THEN 100 ELSE 0 END +
             CASE WHEN damage_level = 'BERAT' THEN 50 
                  WHEN damage_level = 'SEDANG' THEN 25 
                  ELSE 10 END +
             CASE WHEN affected_rice_field_area > 10 THEN 30 ELSE 0 END +
             CASE WHEN affected_farmers_count > 50 THEN 20 ELSE 0 END) as priority_score
        FROM water_resources_reports
        WHERE status NOT IN ('COMPLETED', 'REJECTED')
        ORDER BY priority_score DESC, created_at DESC
        LIMIT $1 OFFSET $2
    `

    err = r.db.WithContext(ctx).Raw(priorityQuery, limit, offset).Scan(&reports).Error
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get priority reports: %w", err)
    }

    return reports, total, err
}

func (r *waterResourcesRepositoryImpl) UpdateStatus(ctx context.Context, id string, status entity.WaterResourceStatus, notes string) error {
    query := `
        UPDATE water_resources_reports 
        SET status = $1, updated_at = $2, notes = CASE WHEN $3 != '' THEN $3 ELSE notes END
        WHERE id = $4
    `
    return r.db.WithContext(ctx).Exec(query, status, time.Now(), notes, id).Error
}

func (r *waterResourcesRepositoryImpl) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // 1. Total reports
    var totalReports int64
    err := r.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM water_resources_reports`).Scan(&totalReports).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count total reports: %w", err)
    }
    stats["total_reports"] = totalReports
    
    // 2. Urgent pending count
    var urgentCount int64
    query := `
        SELECT COUNT(*) 
        FROM water_resources_reports 
        WHERE urgency_category = $1 AND status NOT IN ('COMPLETED', 'REJECTED')
    `
    err = r.db.WithContext(ctx).Raw(query, "MENDESAK").Scan(&urgentCount).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count urgent reports: %w", err)
    }
    stats["urgent_pending"] = urgentCount
    
    // 3. Total affected area
    var totalArea float64
    err = r.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(affected_rice_field_area), 0) FROM water_resources_reports`).Scan(&totalArea).Error
    if err != nil {
        return nil, fmt.Errorf("failed to sum affected area: %w", err)
    }
    stats["total_affected_area_ha"] = totalArea
    
    // 4. Total affected farmers
    var totalFarmers int64
    err = r.db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(affected_farmers_count), 0) FROM water_resources_reports`).Scan(&totalFarmers).Error
    if err != nil {
        return nil, fmt.Errorf("failed to sum affected farmers: %w", err)
    }
    stats["total_affected_farmers"] = totalFarmers
    
    // 5. Damage types distribution
    var damageTypes []map[string]interface{}
    query = `
        SELECT damage_type, COUNT(*) as count 
        FROM water_resources_reports 
        GROUP BY damage_type 
        ORDER BY count DESC
    `
    err = r.db.WithContext(ctx).Raw(query).Scan(&damageTypes).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get damage types: %w", err)
    }
    stats["damage_types"] = damageTypes
    
    // 6. Irrigation types distribution
    var irrigationTypes []map[string]interface{}
    query = `
        SELECT irrigation_type, COUNT(*) as count 
        FROM water_resources_reports 
        GROUP BY irrigation_type 
        ORDER BY count DESC
    `
    err = r.db.WithContext(ctx).Raw(query).Scan(&irrigationTypes).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get irrigation types: %w", err)
    }
    stats["irrigation_types"] = irrigationTypes
    
    // 7. Status distribution
    var statusDist []map[string]interface{}
    query = `
        SELECT status, COUNT(*) as count 
        FROM water_resources_reports 
        GROUP BY status 
        ORDER BY count DESC
    `
    err = r.db.WithContext(ctx).Raw(query).Scan(&statusDist).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get status distribution: %w", err)
    }
    stats["status_distribution"] = statusDist
    
    // 8. Total estimated budget (pending only)
    var totalBudget float64
    query = `
        SELECT COALESCE(SUM(estimated_budget), 0) 
        FROM water_resources_reports 
        WHERE status NOT IN ('COMPLETED', 'REJECTED')
    `
    err = r.db.WithContext(ctx).Raw(query).Scan(&totalBudget).Error
    if err != nil {
        return nil, fmt.Errorf("failed to sum budget: %w", err)
    }
    stats["estimated_total_budget"] = totalBudget
    
    return stats, nil
}

func (r *waterResourcesRepositoryImpl) GetDamageStatisticsByArea(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            irrigation_area_name,
            COUNT(*) as report_count,
            COALESCE(SUM(affected_rice_field_area), 0) as total_affected_area,
            COALESCE(SUM(affected_farmers_count), 0) as total_affected_farmers,
            COALESCE(SUM(estimated_budget), 0) as total_estimated_budget,
            COALESCE(AVG(estimated_length * estimated_width), 0) as avg_damage_area
        FROM water_resources_reports
        WHERE report_datetime BETWEEN $1 AND $2
        GROUP BY irrigation_area_name
        HAVING COUNT(*) > 0
        ORDER BY total_affected_area DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get damage statistics by area: %w", err)
    }
    
    return results, nil
}

func (r *waterResourcesRepositoryImpl) GetUrgentReports(ctx context.Context, limit int) ([]*entity.WaterResourcesReport, error) {
    var reports []*entity.WaterResourcesReport
    
    // Using GORM for preload complexity
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("urgency_category = ?", "MENDESAK").
        Where("status NOT IN ('COMPLETED', 'REJECTED')").
        Order("created_at DESC").
        Limit(limit).
        Find(&reports).Error
    
    if err != nil {
        return nil, fmt.Errorf("failed to get urgent reports: %w", err)
    }
    
    return reports, nil
}

func (r *waterResourcesRepositoryImpl) CalculateTotalDamageArea(ctx context.Context) (float64, error) {
    var total float64
    query := `SELECT COALESCE(SUM(estimated_length * estimated_width), 0) FROM water_resources_reports`
    err := r.db.WithContext(ctx).Raw(query).Scan(&total).Error
    if err != nil {
        return 0, fmt.Errorf("failed to calculate total damage area: %w", err)
    }
    return total, nil
}

func (r *waterResourcesRepositoryImpl) CountAffectedFarmers(ctx context.Context) (int64, error) {
    var count int64
    query := `SELECT COALESCE(SUM(affected_farmers_count), 0) FROM water_resources_reports`
    err := r.db.WithContext(ctx).Raw(query).Scan(&count).Error
    if err != nil {
        return 0, fmt.Errorf("failed to count affected farmers: %w", err)
    }
    return count, nil
}

// Dashboard specific methods with Raw SQL
func (r *waterResourcesRepositoryImpl) GetSummaryKPIs(ctx context.Context, irrigationType string, startDate, endDate time.Time) (float64, float64, int64, error) {
    baseWhere := "WHERE report_datetime BETWEEN $1 AND $2"
    args := []interface{}{startDate, endDate}
    argIndex := 3
    
    if irrigationType != "" && irrigationType != "ALL" {
        baseWhere += fmt.Sprintf(" AND irrigation_type = $%d", argIndex)
        args = append(args, irrigationType)
    }
    
    // Total damage area (m2)
    var totalArea float64
    query := fmt.Sprintf(`SELECT COALESCE(SUM(estimated_length * estimated_width), 0) FROM water_resources_reports %s`, baseWhere)
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&totalArea).Error
    if err != nil {
        return 0, 0, 0, fmt.Errorf("failed to get total area: %w", err)
    }

    // Total rice field area (ha)
    var totalRice float64
    query = fmt.Sprintf(`SELECT COALESCE(SUM(affected_rice_field_area), 0) FROM water_resources_reports %s`, baseWhere)
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&totalRice).Error
    if err != nil {
        return 0, 0, 0, fmt.Errorf("failed to get total rice area: %w", err)
    }

    // Total reports
    var totalReports int64
    query = fmt.Sprintf(`SELECT COUNT(*) FROM water_resources_reports %s`, baseWhere)
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&totalReports).Error
    if err != nil {
        return 0, 0, 0, fmt.Errorf("failed to count total reports: %w", err)
    }

    return totalArea, totalRice, totalReports, nil
}

func (r *waterResourcesRepositoryImpl) GroupCountBy(ctx context.Context, field, irrigationType string, startDate, endDate time.Time) ([]struct {
    Key   string
    Count int64
}, error) {
    baseWhere := "WHERE report_datetime BETWEEN $1 AND $2"
    args := []interface{}{startDate, endDate}
    argIndex := 3
    
    if irrigationType != "" && irrigationType != "ALL" {
        baseWhere += fmt.Sprintf(" AND irrigation_type = $%d", argIndex)
        args = append(args, irrigationType)
    }
    
    query := fmt.Sprintf(`
        SELECT %s as key, COUNT(*) as count 
        FROM water_resources_reports %s 
        GROUP BY %s 
        ORDER BY count DESC
    `, field, baseWhere, field)
    
    var results []struct {
        Key   string
        Count int64
    }
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to group count by %s: %w", field, err)
    }
    
    return results, nil
}

func (r *waterResourcesRepositoryImpl) GetMapPoints(ctx context.Context, irrigationType string, startDate, endDate time.Time) ([]struct {
    Latitude        float64
    Longitude       float64
    IrrigationArea  string
    DamageType      string
    DamageLevel     string
    UrgencyCategory string
}, error) {
    baseWhere := "WHERE report_datetime BETWEEN $1 AND $2 AND latitude IS NOT NULL AND longitude IS NOT NULL"
    args := []interface{}{startDate, endDate}
    argIndex := 3
    
    if irrigationType != "" && irrigationType != "ALL" {
        baseWhere += fmt.Sprintf(" AND irrigation_type = $%d", argIndex)
        args = append(args, irrigationType)
    }
    
    query := fmt.Sprintf(`
        SELECT 
            latitude, 
            longitude, 
            irrigation_area_name as irrigation_area, 
            damage_type, 
            damage_level, 
            urgency_category
        FROM water_resources_reports %s
    `, baseWhere)
    
    var results []struct {
        Latitude        float64
        Longitude       float64
        IrrigationArea  string
        DamageType      string
        DamageLevel     string
        UrgencyCategory string
    }
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get map points: %w", err)
    }
    
    return results, nil
}

// Overview specific methods with Raw SQL
func (r *waterResourcesRepositoryImpl) GetWaterResourcesOverviewStats(ctx context.Context, irrigationType string) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Build base query
    baseWhere := ""
    args := []interface{}{}
    
    if irrigationType != "" && irrigationType != "all" && irrigationType != "ALL" {
        baseWhere = "WHERE irrigation_type = $1"
        args = append(args, irrigationType)
    }
    
    // 1. Total damage volume (estimated_length * estimated_width) in m²
    query := fmt.Sprintf(`SELECT COALESCE(SUM(estimated_length * estimated_width), 0) FROM water_resources_reports %s`, baseWhere)
    var totalDamageVolume float64
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&totalDamageVolume).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate total damage volume: %w", err)
    }
    stats["total_damage_volume_m2"] = totalDamageVolume
    
    // 2. Total rice field area affected (ha)
    query = fmt.Sprintf(`SELECT COALESCE(SUM(affected_rice_field_area), 0) FROM water_resources_reports %s`, baseWhere)
    var totalRiceFieldArea float64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&totalRiceFieldArea).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate total rice field area: %w", err)
    }
    stats["total_rice_field_area_ha"] = totalRiceFieldArea
    
    // 3. Total damaged reports count
    query = fmt.Sprintf(`SELECT COUNT(*) FROM water_resources_reports %s`, baseWhere)
    var totalReports int64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&totalReports).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count total reports: %w", err)
    }
    stats["total_damaged_reports"] = totalReports
    
    return stats, nil
}

func (r *waterResourcesRepositoryImpl) GetWaterLocationStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            irrigation_area_name,
            COUNT(*) as report_count,
            COALESCE(AVG(latitude), 0) as avg_latitude,
            COALESCE(AVG(longitude), 0) as avg_longitude,
            COALESCE(SUM(affected_rice_field_area), 0) as total_affected_area,
            COALESCE(SUM(affected_farmers_count), 0) as total_affected_farmers
        FROM water_resources_reports
    `
    
    args := []interface{}{}
    if irrigationType != "" && irrigationType != "all" && irrigationType != "ALL" {
        query += " WHERE irrigation_type = $1"
        args = append(args, irrigationType)
    }
    
    query += `
        GROUP BY irrigation_area_name
        HAVING COUNT(*) > 0
        ORDER BY report_count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get location statistics: %w", err)
    }
    
    return results, nil
}

func (r *waterResourcesRepositoryImpl) GetWaterUrgencyStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN urgency_category IS NULL THEN 'NOT_SET'
                WHEN TRIM(urgency_category) = '' THEN 'EMPTY'
                ELSE urgency_category 
            END as urgency_category,
            COUNT(*) as count
        FROM water_resources_reports
    `
    
    args := []interface{}{}
    if irrigationType != "" && irrigationType != "all" && irrigationType != "ALL" {
        query += " WHERE irrigation_type = $1"
        args = append(args, irrigationType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN urgency_category IS NULL THEN 'NOT_SET'
                WHEN TRIM(urgency_category) = '' THEN 'EMPTY'
                ELSE urgency_category 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get urgency statistics: %w", err)
    }
    
    return results, nil
}

func (r *waterResourcesRepositoryImpl) GetWaterDamageTypeStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN damage_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_type) = '' THEN 'EMPTY'
                ELSE damage_type 
            END as damage_type,
            COUNT(*) as count
        FROM water_resources_reports
    `
    
    args := []interface{}{}
    if irrigationType != "" && irrigationType != "all" && irrigationType != "ALL" {
        query += " WHERE irrigation_type = $1"
        args = append(args, irrigationType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN damage_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_type) = '' THEN 'EMPTY'
                ELSE damage_type 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get damage type statistics: %w", err)
    }
    
    return results, nil
}

func (r *waterResourcesRepositoryImpl) GetWaterDamageLevelStats(ctx context.Context, irrigationType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN damage_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_level) = '' THEN 'EMPTY'
                ELSE damage_level 
            END as damage_level,
            COUNT(*) as count
        FROM water_resources_reports
    `
    
    args := []interface{}{}
    if irrigationType != "" && irrigationType != "all" && irrigationType != "ALL" {
        query += " WHERE irrigation_type = $1"
        args = append(args, irrigationType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN damage_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_level) = '' THEN 'EMPTY'
                ELSE damage_level 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get damage level statistics: %w", err)
    }
    
    return results, nil
}