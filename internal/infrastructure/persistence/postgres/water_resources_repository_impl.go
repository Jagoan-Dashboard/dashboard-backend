package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"
	"time"

	"gorm.io/gorm"
)

type waterResourcesRepositoryImpl struct {
    db *gorm.DB
}

func NewWaterResourcesRepository(db *gorm.DB) repository.WaterResourcesRepository {
    return &waterResourcesRepositoryImpl{db: db}
}

func (r *waterResourcesRepositoryImpl) Create(ctx context.Context, report *entity.WaterResourcesReport) error {
    return r.db.WithContext(ctx).Create(report).Error
}

func (r *waterResourcesRepositoryImpl) Update(ctx context.Context, report *entity.WaterResourcesReport) error {
    report.UpdatedAt = time.Now()
    return r.db.WithContext(ctx).Save(report).Error
}

func (r *waterResourcesRepositoryImpl) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.WaterResourcesReport{}).Error
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

func (r *waterResourcesRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.WaterResourcesReport, int64, error) {
    var reports []*entity.WaterResourcesReport
    var total int64

    query := r.db.WithContext(ctx).Model(&entity.WaterResourcesReport{})

    
    if institutionUnit, ok := filters["institution_unit"].(string); ok && institutionUnit != "" {
        query = query.Where("institution_unit = ?", institutionUnit)
    }
    if irrigationType, ok := filters["irrigation_type"].(string); ok && irrigationType != "" {
        query = query.Where("irrigation_type = ?", irrigationType)
    }
    if damageType, ok := filters["damage_type"].(string); ok && damageType != "" {
        query = query.Where("damage_type = ?", damageType)
    }
    if damageLevel, ok := filters["damage_level"].(string); ok && damageLevel != "" {
        query = query.Where("damage_level = ?", damageLevel)
    }
    if urgencyCategory, ok := filters["urgency_category"].(string); ok && urgencyCategory != "" {
        query = query.Where("urgency_category = ?", urgencyCategory)
    }
    if status, ok := filters["status"].(string); ok && status != "" {
        query = query.Where("status = ?", status)
    }
    
    
    if irrigationArea, ok := filters["irrigation_area"].(string); ok && irrigationArea != "" {
        query = query.Where("irrigation_area_name ILIKE ?", "%"+irrigationArea+"%")
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
        Order("CASE WHEN urgency_category = 'MENDESAK' THEN 0 ELSE 1 END, " +
              "CASE WHEN damage_level = 'BERAT' THEN 0 WHEN damage_level = 'SEDANG' THEN 1 ELSE 2 END, " +
              "created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *waterResourcesRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.WaterResourcesReport, int64, error) {
    var reports []*entity.WaterResourcesReport
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.WaterResourcesReport{}).
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

func (r *waterResourcesRepositoryImpl) FindByPriority(ctx context.Context, limit, offset int) ([]*entity.WaterResourcesReport, int64, error) {
    var reports []*entity.WaterResourcesReport
    var total int64

    
    prioritySQL := `
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
        LIMIT ? OFFSET ?
    `

    err := r.db.WithContext(ctx).
        Raw(prioritySQL, limit, offset).
        Scan(&reports).Error

    
    r.db.Model(&entity.WaterResourcesReport{}).
        Where("status NOT IN ('COMPLETED', 'REJECTED')").
        Count(&total)

    return reports, total, err
}

func (r *waterResourcesRepositoryImpl) UpdateStatus(ctx context.Context, id string, status entity.WaterResourceStatus, notes string) error {
    updates := map[string]interface{}{
        "status":     status,
        "updated_at": time.Now(),
    }
    if notes != "" {
        updates["notes"] = notes
    }
    
    return r.db.WithContext(ctx).
        Model(&entity.WaterResourcesReport{}).
        Where("id = ?", id).
        Updates(updates).Error
}

func (r *waterResourcesRepositoryImpl) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    
    var total int64
    r.db.Model(&entity.WaterResourcesReport{}).Count(&total)
    stats["total_reports"] = total
    
    
    var urgentCount int64
    r.db.Model(&entity.WaterResourcesReport{}).
        Where("urgency_category = ?", entity.UrgencyCategoryMendesak).
        Where("status NOT IN ('COMPLETED', 'REJECTED')").
        Count(&urgentCount)
    stats["urgent_pending"] = urgentCount
    
    
    var totalArea float64
    r.db.Model(&entity.WaterResourcesReport{}).
        Select("COALESCE(SUM(affected_rice_field_area), 0)").
        Scan(&totalArea)
    stats["total_affected_area_ha"] = totalArea
    
    
    var totalFarmers int64
    r.db.Model(&entity.WaterResourcesReport{}).
        Select("COALESCE(SUM(affected_farmers_count), 0)").
        Scan(&totalFarmers)
    stats["total_affected_farmers"] = totalFarmers
    
    
    type DamageTypeCount struct {
        DamageType string `json:"damage_type"`
        Count      int64  `json:"count"`
    }
    var damageTypeCounts []DamageTypeCount
    r.db.Model(&entity.WaterResourcesReport{}).
        Select("damage_type, count(*) as count").
        Group("damage_type").
        Scan(&damageTypeCounts)
    stats["damage_types"] = damageTypeCounts
    
    
    type IrrigationTypeCount struct {
        IrrigationType string `json:"irrigation_type"`
        Count          int64  `json:"count"`
    }
    var irrigationTypeCounts []IrrigationTypeCount
    r.db.Model(&entity.WaterResourcesReport{}).
        Select("irrigation_type, count(*) as count").
        Group("irrigation_type").
        Scan(&irrigationTypeCounts)
    stats["irrigation_types"] = irrigationTypeCounts
    
    
    type StatusCount struct {
        Status string `json:"status"`
        Count  int64  `json:"count"`
    }
    var statusCounts []StatusCount
    r.db.Model(&entity.WaterResourcesReport{}).
        Select("status, count(*) as count").
        Group("status").
        Scan(&statusCounts)
    stats["status_distribution"] = statusCounts
    
    
    var totalBudget float64
    r.db.Model(&entity.WaterResourcesReport{}).
        Where("status NOT IN ('COMPLETED', 'REJECTED')").
        Select("COALESCE(SUM(estimated_budget), 0)").
        Scan(&totalBudget)
    stats["estimated_total_budget"] = totalBudget
    
    return stats, nil
}

func (r *waterResourcesRepositoryImpl) GetDamageStatisticsByArea(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            irrigation_area_name,
            COUNT(*) as report_count,
            SUM(affected_rice_field_area) as total_affected_area,
            SUM(affected_farmers_count) as total_affected_farmers,
            SUM(estimated_budget) as total_estimated_budget,
            AVG(estimated_length * estimated_width) as avg_damage_area
        FROM water_resources_reports
        WHERE report_datetime BETWEEN ? AND ?
        GROUP BY irrigation_area_name
        ORDER BY total_affected_area DESC
    `
    
    err := r.db.WithContext(ctx).
        Raw(query, startDate, endDate).
        Scan(&results).Error
    
    return results, err
}

func (r *waterResourcesRepositoryImpl) GetUrgentReports(ctx context.Context, limit int) ([]*entity.WaterResourcesReport, error) {
    var reports []*entity.WaterResourcesReport
    
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("urgency_category = ?", entity.UrgencyCategoryMendesak).
        Where("status NOT IN ('COMPLETED', 'REJECTED')").
        Order("created_at DESC").
        Limit(limit).
        Find(&reports).Error
    
    return reports, err
}

func (r *waterResourcesRepositoryImpl) CalculateTotalDamageArea(ctx context.Context) (float64, error) {
    var total float64
    err := r.db.WithContext(ctx).
        Model(&entity.WaterResourcesReport{}).
        Select("COALESCE(SUM(estimated_length * estimated_width), 0)").
        Scan(&total).Error
    return total, err
}

func (r *waterResourcesRepositoryImpl) CountAffectedFarmers(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&entity.WaterResourcesReport{}).
        Select("COALESCE(SUM(affected_farmers_count), 0)").
        Scan(&count).Error
    return count, err
}
type keyCountRow struct {
    Key   string `gorm:"column:key"`
    Count int64  `gorm:"column:count"`
}

func (r *waterResourcesRepositoryImpl) baseScoped(ctx context.Context, irrigationType string, startDate, endDate time.Time) *gorm.DB {
    q := r.db.WithContext(ctx).Model(&entity.WaterResourcesReport{}).
        Where("report_datetime BETWEEN ? AND ?", startDate, endDate)
    if irrigationType != "" && irrigationType != "ALL" {
        q = q.Where("irrigation_type = ?", irrigationType)
    }
    return q
}

func (r *waterResourcesRepositoryImpl) GetSummaryKPIs(ctx context.Context, irrigationType string, startDate, endDate time.Time) (float64, float64, int64, error) {
    // Total area (m2) = sum(length*width)
    var totalArea float64
    if err := r.baseScoped(ctx, irrigationType, startDate, endDate).
        Select("COALESCE(SUM(estimated_length * estimated_width), 0)").
        Scan(&totalArea).Error; err != nil {
        return 0, 0, 0, err
    }

    // Total rice field area (ha)
    var totalRice float64
    if err := r.baseScoped(ctx, irrigationType, startDate, endDate).
        Select("COALESCE(SUM(affected_rice_field_area), 0)").
        Scan(&totalRice).Error; err != nil {
        return 0, 0, 0, err
    }

    // Total reports
    var totalReports int64
    if err := r.baseScoped(ctx, irrigationType, startDate, endDate).
        Count(&totalReports).Error; err != nil {
        return 0, 0, 0, err
    }

    return totalArea, totalRice, totalReports, nil
}

func (r *waterResourcesRepositoryImpl) GroupCountBy(ctx context.Context, field, irrigationType string, startDate, endDate time.Time) ([]struct {
    Key   string
    Count int64
}, error) {
    var rows []keyCountRow
    q := r.baseScoped(ctx, irrigationType, startDate, endDate).
        Select(field+" as key, COUNT(*) as count").
        Group(field).
        Order("count DESC")
    if err := q.Scan(&rows).Error; err != nil {
        return nil, err
    }
    // cast ke slice anonim yang sesuai interface
    out := make([]struct {
        Key   string
        Count int64
    }, len(rows))
    for i, r0 := range rows {
        out[i] = struct {
            Key   string
            Count int64
        }{Key: r0.Key, Count: r0.Count}
    }
    return out, nil
}

func (r *waterResourcesRepositoryImpl) GetMapPoints(ctx context.Context, irrigationType string, startDate, endDate time.Time) ([]struct {
    Latitude        float64
    Longitude       float64
    IrrigationArea  string
    DamageType      string
    DamageLevel     string
    UrgencyCategory string
}, error) {
    type row struct {
        Latitude        float64
        Longitude       float64
        IrrigationArea  string `gorm:"column:irrigation_area_name"`
        DamageType      string
        DamageLevel     string
        UrgencyCategory string `gorm:"column:urgency_category"`
    }
    var rows []row
    if err := r.baseScoped(ctx, irrigationType, startDate, endDate).
        Select("latitude, longitude, irrigation_area_name, damage_type, damage_level, urgency_category").
        Where("latitude IS NOT NULL AND longitude IS NOT NULL").
        Scan(&rows).Error; err != nil {
        return nil, err
    }
    out := make([]struct {
        Latitude        float64
        Longitude       float64
        IrrigationArea  string
        DamageType      string
        DamageLevel     string
        UrgencyCategory string
    }, len(rows))
    for i, v := range rows {
        out[i] = struct {
            Latitude        float64
            Longitude       float64
            IrrigationArea  string
            DamageType      string
            DamageLevel     string
            UrgencyCategory string
        }{
            Latitude:        v.Latitude,
            Longitude:       v.Longitude,
            IrrigationArea:  v.IrrigationArea,
            DamageType:      v.DamageType,
            DamageLevel:     v.DamageLevel,
            UrgencyCategory: v.UrgencyCategory,
        }
    }
    return out, nil
}