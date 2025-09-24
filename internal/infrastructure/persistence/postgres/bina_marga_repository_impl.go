package postgres

import (
	"context"
	"fmt"
	"time"

	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type binaMargaRepositoryImpl struct {
	db *gorm.DB
}

func NewBinaMargaRepository(db *gorm.DB) repository.BinaMargaRepository {
	return &binaMargaRepositoryImpl{db: db}
}

func (r *binaMargaRepositoryImpl) Create(ctx context.Context, report *entity.BinaMargaReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *binaMargaRepositoryImpl) Update(ctx context.Context, report *entity.BinaMargaReport) error {
	report.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *binaMargaRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.BinaMargaReport{}).Error
}

func (r *binaMargaRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.BinaMargaReport, error) {
	var report entity.BinaMargaReport
	err := r.db.WithContext(ctx).
		Preload("Photos").
		Where("id = ?", id).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *binaMargaRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.BinaMargaReport, int64, error) {
	var (
		reports []*entity.BinaMargaReport
		total   int64
	)

	query := r.db.WithContext(ctx).Model(&entity.BinaMargaReport{})

	
	if v, ok := filters["institution_unit"].(string); ok && v != "" {
		query = query.Where("institution_unit = ?", v)
	}
	if v, ok := filters["road_type"].(string); ok && v != "" {
		query = query.Where("road_type = ?", v)
	}
	if v, ok := filters["road_class"].(string); ok && v != "" {
		query = query.Where("road_class = ?", v)
	}
	if v, ok := filters["damage_type"].(string); ok && v != "" {
		query = query.Where("damage_type = ?", v)
	}
	if v, ok := filters["damage_level"].(string); ok && v != "" {
		query = query.Where("damage_level = ?", v)
	}
	if v, ok := filters["urgency_level"].(string); ok && v != "" {
		query = query.Where("urgency_level = ?", v)
	}
	if v, ok := filters["traffic_impact"].(string); ok && v != "" {
		query = query.Where("traffic_impact = ?", v)
	}
	if v, ok := filters["status"].(string); ok && v != "" {
		query = query.Where("status = ?", v)
	}
	if v, ok := filters["road_name"].(string); ok && v != "" {
		query = query.Where("road_name ILIKE ?", "%"+v+"%")
	}
	
	if v, ok := filters["start_date"].(string); ok && v != "" {
		query = query.Where("report_datetime >= ?", v)
	}
	if v, ok := filters["end_date"].(string); ok && v != "" {
		query = query.Where("report_datetime <= ?", v)
	}

	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	
	orderExpr := "CASE " +
		"WHEN urgency_level = 'DARURAT' THEN 0 " +
		"WHEN urgency_level = 'TINGGI' THEN 1 " +
		"WHEN urgency_level = 'SEDANG' THEN 2 " +
		"ELSE 3 END, " +
		"CASE " +
		"WHEN traffic_impact = 'TERPUTUS' THEN 0 " +
		"WHEN traffic_impact = 'SANGAT_TERGANGGU' THEN 1 " +
		"WHEN traffic_impact = 'TERGANGGU' THEN 2 " +
		"ELSE 3 END, " +
		"created_at DESC"

	
	err := query.
		Preload("Photos").
		Limit(limit).
		Offset(offset).
		Order(orderExpr).
		Find(&reports).Error

	return reports, total, err
}

func (r *binaMargaRepositoryImpl) FindBlockedRoads(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error) {
	var reports []*entity.BinaMargaReport
	err := r.db.WithContext(ctx).
		Preload("Photos").
		Where("traffic_impact = ?", entity.TrafficImpactBlocked).
		Where("status NOT IN ('COMPLETED', 'REJECTED')").
		Order("created_at DESC").
		Limit(limit).
		Find(&reports).Error
	return reports, err
}

func (r *binaMargaRepositoryImpl) UpdateStatus(ctx context.Context, id string, status entity.BinaMargaStatus, notes string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if notes != "" {
		updates["notes"] = notes
	}
	return r.db.WithContext(ctx).
		Model(&entity.BinaMargaReport{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *binaMargaRepositoryImpl) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	
	var total int64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).Count(&total).Error
	stats["total_reports"] = total

	
	var emergencyCount int64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Where("urgency_level = ?", entity.RoadUrgencyEmergency).
		Where("status NOT IN ('COMPLETED', 'REJECTED')").
		Count(&emergencyCount).Error
	stats["emergency_reports"] = emergencyCount

	
	var blockedCount int64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Where("traffic_impact = ?", entity.TrafficImpactBlocked).
		Where("status NOT IN ('COMPLETED', 'REJECTED')").
		Count(&blockedCount).Error
	stats["blocked_roads"] = blockedCount

	
	var totalArea float64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("COALESCE(SUM(damaged_area), 0)").
		Scan(&totalArea).Error
	stats["total_damaged_area_sqm"] = totalArea

	var totalLength float64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("COALESCE(SUM(damaged_length), 0)").
		Scan(&totalLength).Error
	stats["total_damaged_length_m"] = totalLength

	
	type kvCount struct {
		Key   string `json:"key" gorm:"column:key"`
		Count int64  `json:"count" gorm:"column:count"`
	}

	var roadTypeCounts []struct {
		RoadType string `json:"road_type"`
		Count    int64  `json:"count"`
	}
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("road_type, COUNT(*) as count").
		Group("road_type").
		Scan(&roadTypeCounts).Error
	stats["road_type_distribution"] = roadTypeCounts

	var damageTypeCounts []struct {
		DamageType string `json:"damage_type"`
		Count      int64  `json:"count"`
	}
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("damage_type, COUNT(*) as count").
		Group("damage_type").
		Scan(&damageTypeCounts).Error
	stats["damage_type_distribution"] = damageTypeCounts

	var damageLevelCounts []struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("damage_level as level, COUNT(*) as count").
		Group("damage_level").
		Scan(&damageLevelCounts).Error
	stats["damage_level_counts"] = damageLevelCounts

	var urgencyLevelCounts []struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("urgency_level as level, COUNT(*) as count").
		Group("urgency_level").
		Scan(&urgencyLevelCounts).Error
	stats["urgency_level_counts"] = urgencyLevelCounts

	var statusCounts []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error
	stats["status_distribution"] = statusCounts

	var trafficImpactCounts []struct {
		Impact string `json:"impact"`
		Count  int64  `json:"count"`
	}
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select("traffic_impact as impact, COUNT(*) as count").
		Group("traffic_impact").
		Scan(&trafficImpactCounts).Error
	stats["traffic_impact_counts"] = trafficImpactCounts

	var totalBudget float64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Where("status NOT IN ('COMPLETED', 'REJECTED')").
		Select("COALESCE(SUM(estimated_budget), 0)").
		Scan(&totalBudget).Error
	stats["estimated_total_budget"] = totalBudget

	var avgRepairTime float64
	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Where("estimated_repair_time > 0").
		Select("COALESCE(AVG(estimated_repair_time), 0)").
		Scan(&avgRepairTime).Error
	stats["average_repair_time_days"] = avgRepairTime

	return stats, nil
}

func (r *binaMargaRepositoryImpl) GetDamageStatisticsByRoadType(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	const q = `
		SELECT 
			road_type,
			road_class,
			COUNT(*) AS report_count,
			SUM(damaged_area) AS total_damaged_area,
			SUM(damaged_length) AS total_damaged_length,
			SUM(estimated_budget) AS total_estimated_budget,
			AVG(estimated_repair_time) AS avg_repair_time,
			COUNT(CASE WHEN urgency_level = 'DARURAT' THEN 1 END) AS emergency_count
		FROM bina_marga_reports
		WHERE report_datetime BETWEEN ? AND ?
		GROUP BY road_type, road_class
		ORDER BY total_damaged_area DESC`
	err := r.db.WithContext(ctx).Raw(q, startDate, endDate).Scan(&results).Error
	return results, err
}

func (r *binaMargaRepositoryImpl) GetDamageStatisticsByLocation(ctx context.Context, bounds map[string]float64) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	const q = `
		SELECT 
			road_name,
			road_type,
			road_class,
			latitude,
			longitude,
			damage_type,
			damage_level,
			urgency_level,
			traffic_impact,
			damaged_area,
			damaged_length,
			status,
			created_at
		FROM bina_marga_reports
		WHERE latitude BETWEEN ? AND ?
		  AND longitude BETWEEN ? AND ?
		ORDER BY urgency_level DESC, created_at DESC`
	err := r.db.WithContext(ctx).Raw(q, bounds["south"], bounds["north"], bounds["west"], bounds["east"]).Scan(&results).Error
	return results, err
}

func (r *binaMargaRepositoryImpl) CalculateTotalDamageArea(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entity.BinaMargaReport{}).
		Select("COALESCE(SUM(damaged_area), 0)").
		Scan(&total).Error
	return total, err
}

func (r *binaMargaRepositoryImpl) CalculateTotalDamageLength(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entity.BinaMargaReport{}).
		Select("COALESCE(SUM(damaged_length), 0)").
		Scan(&total).Error
	return total, err
}

func (r *binaMargaRepositoryImpl) CountReportsByUrgency(ctx context.Context, urgency entity.RoadUrgencyLevel) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.BinaMargaReport{}).
		Where("urgency_level = ?", urgency).
		Count(&count).Error
	return count, err
}

func (r *binaMargaRepositoryImpl) GetRepairTimeAnalysis(ctx context.Context) (map[string]interface{}, error) {
	analysis := make(map[string]interface{})

	type RepairTimeByLevel struct {
		DamageLevel   string  `json:"damage_level"`
		AvgRepairTime float64 `json:"avg_repair_time"`
		MinRepairTime int     `json:"min_repair_time"`
		MaxRepairTime int     `json:"max_repair_time"`
		Count         int64   `json:"count"`
	}
	var byLevel []RepairTimeByLevel
	r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select(`damage_level,
		        AVG(estimated_repair_time)  AS avg_repair_time,
		        MIN(estimated_repair_time)  AS min_repair_time,
		        MAX(estimated_repair_time)  AS max_repair_time,
		        COUNT(*)                    AS count`).
		Where("estimated_repair_time > 0").
		Group("damage_level").
		Scan(&byLevel)
	analysis["repair_time_by_level"] = byLevel

	type RepairTimeByClass struct {
		RoadClass     string  `json:"road_class"`
		AvgRepairTime float64 `json:"avg_repair_time"`
		Count         int64   `json:"count"`
	}
	var byClass []RepairTimeByClass
	r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Select(`road_class,
		        AVG(estimated_repair_time) AS avg_repair_time,
		        COUNT(*)                   AS count`).
		Where("estimated_repair_time > 0").
		Group("road_class").
		Scan(&byClass)
	analysis["repair_time_by_class"] = byClass

	return analysis, nil
}

func (r *binaMargaRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.BinaMargaReport, int64, error) {
	var (
		reports []*entity.BinaMargaReport
		total   int64
	)
	query := r.db.WithContext(ctx).
		Model(&entity.BinaMargaReport{}).
		Where("created_by = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Photos").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&reports).Error

	return reports, total, err
}

func (r *binaMargaRepositoryImpl) FindByPriority(ctx context.Context, limit, offset int) ([]*entity.BinaMargaReport, int64, error) {
	var (
		reports []*entity.BinaMargaReport
		total   int64
	)

	const prioritySQL = `
		SELECT *,
			( CASE WHEN urgency_level = 'DARURAT' THEN 100
			       WHEN urgency_level = 'TINGGI'  THEN 75
			       WHEN urgency_level = 'SEDANG'  THEN 50
			       ELSE 25 END
			  + CASE WHEN damage_level = 'BERAT' THEN 50
			         WHEN damage_level = 'SEDANG' THEN 30
			         ELSE 15 END
			  + CASE WHEN road_class = 'ARTERI'   THEN 40
			         WHEN road_class = 'KOLEKTOR' THEN 30
			         WHEN road_class = 'LOKAL'    THEN 20
			         ELSE 10 END
			  + CASE WHEN traffic_impact = 'TERPUTUS'         THEN 60
			         WHEN traffic_impact = 'SANGAT_TERGANGGU'  THEN 40
			         WHEN traffic_impact = 'TERGANGGU'         THEN 20
			         ELSE 5 END
			  + CASE WHEN damaged_area > 100 THEN 25
			         WHEN damaged_area > 50  THEN 15
			         ELSE 0 END ) AS priority_score
		FROM bina_marga_reports
		WHERE status NOT IN ('COMPLETED', 'REJECTED')
		ORDER BY priority_score DESC, created_at DESC
		LIMIT ? OFFSET ?`

	err := r.db.WithContext(ctx).Raw(prioritySQL, limit, offset).Scan(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	_ = r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
		Where("status NOT IN ('COMPLETED', 'REJECTED')").
		Count(&total).Error

	return reports, total, nil
}

func (r *binaMargaRepositoryImpl) FindEmergencyReports(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error) {
	var reports []*entity.BinaMargaReport
	err := r.db.WithContext(ctx).
		Preload("Photos").
		Where("urgency_level = ?", entity.RoadUrgencyEmergency).
		Where("status NOT IN ('COMPLETED', 'REJECTED')").
		Order("created_at DESC").
		Limit(limit).
		Find(&reports).Error
	return reports, err
}


func (r *binaMargaRepositoryImpl) baseScoped(ctx context.Context, roadType string, startDate, endDate time.Time) *gorm.DB {
    q := r.db.WithContext(ctx).Model(&entity.BinaMargaReport{}).
        Where("report_datetime BETWEEN ? AND ?", startDate, endDate)
    if roadType != "" && roadType != "ALL" {
        q = q.Where("road_type = ?", roadType)
    }
    return q
}

func (r *binaMargaRepositoryImpl) GetKPIs(ctx context.Context, roadType string, startDate, endDate time.Time) (float64, float64, float64, int64, error) {
    var avgSeg float64
    if err := r.baseScoped(ctx, roadType, startDate, endDate).
        Select("COALESCE(AVG(segment_length), 0)").
        Scan(&avgSeg).Error; err != nil {
        return 0, 0, 0, 0, err
    }

    var avgArea float64
    if err := r.baseScoped(ctx, roadType, startDate, endDate).
        Select("COALESCE(AVG(COALESCE(total_damaged_area, damaged_area)), 0)").
        Scan(&avgArea).Error; err != nil {
        return 0, 0, 0, 0, err
    }

    var avgTraffic float64
    if err := r.baseScoped(ctx, roadType, startDate, endDate).
        Select("COALESCE(AVG(daily_traffic_volume), 0)").
        Scan(&avgTraffic).Error; err != nil {
        return 0, 0, 0, 0, err
    }

    var total int64
    if err := r.baseScoped(ctx, roadType, startDate, endDate).
        Count(&total).Error; err != nil {
        return 0, 0, 0, 0, err
    }

    return avgSeg, avgArea, avgTraffic, total, nil
}

type keyCountRow struct {
	Key   string
	Count int64
}

func (r *binaMargaRepositoryImpl) GroupCountBy(ctx context.Context, column, roadType string, startDate, endDate time.Time, onlyBridge, onlyRoad bool) ([]struct {
	Key string
	Count int64
}, error) {
	q := r.baseScoped(ctx, roadType, startDate, endDate)
	if onlyBridge {
		q = q.Where("bridge_name IS NOT NULL AND bridge_name <> ''")
	}
	if onlyRoad {
		q = q.Where("(bridge_name IS NULL OR bridge_name = '')")
	}

	var rows []keyCountRow
	if err := q.Select(column+" as key, COUNT(*) as count").
		Group(column).
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]struct{ Key string; Count int64 }, len(rows))
	for i, r0 := range rows {
		out[i] = struct{ Key string; Count int64 }{Key: r0.Key, Count: r0.Count}
	}
	return out, nil
}

func (r *binaMargaRepositoryImpl) GetMapPoints(ctx context.Context, roadType string, startDate, endDate time.Time) ([]struct {
    Latitude           float64
    Longitude          float64
    RoadName           string
    RoadType           string
    DamageType         string
    DamageLevel        string
    BridgeName         *string
    BridgeDamageType   *string
    BridgeDamageLevel  *string
    UrgencyLevel       string
}, error) {
    type row struct {
        Latitude          float64
        Longitude         float64
        RoadName          string
        RoadType          string
        DamageType        string
        DamageLevel       string
        BridgeName        *string
        BridgeDamageType  *string
        BridgeDamageLevel *string
        UrgencyLevel      string
    }

    var rows []row
    if err := r.baseScoped(ctx, roadType, startDate, endDate).
        Select("latitude, longitude, road_name, road_type, damage_type, damage_level, bridge_name, bridge_damage_type, bridge_damage_level, urgency_level").
        Where("latitude IS NOT NULL AND longitude IS NOT NULL").
        Scan(&rows).Error; err != nil {
        return nil, err
    }

    out := make([]struct {
        Latitude           float64
        Longitude          float64
        RoadName           string
        RoadType           string
        DamageType         string
        DamageLevel        string
        BridgeName         *string
        BridgeDamageType   *string
        BridgeDamageLevel  *string
        UrgencyLevel       string
    }, len(rows))
    for i, v := range rows {
        out[i] = struct {
            Latitude           float64
            Longitude          float64
            RoadName           string
            RoadType           string
            DamageType         string
            DamageLevel        string
            BridgeName         *string
            BridgeDamageType   *string
            BridgeDamageLevel  *string
            UrgencyLevel       string
        }{
            Latitude: v.Latitude, Longitude: v.Longitude,
            RoadName: v.RoadName, RoadType: v.RoadType,
            DamageType: v.DamageType, DamageLevel: v.DamageLevel,
            BridgeName: v.BridgeName, BridgeDamageType: v.BridgeDamageType, BridgeDamageLevel: v.BridgeDamageLevel,
            UrgencyLevel: v.UrgencyLevel,
        }
    }
    return out, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaOverviewStats(ctx context.Context, roadType string) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    // Build base query
    baseWhere := ""
    args := []interface{}{}
    
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        baseWhere = "WHERE road_type = $1"
        args = append(args, roadType)
    }
    
    // 1. Average segment length (m)
    query := fmt.Sprintf(`SELECT COALESCE(AVG(segment_length), 0) FROM bina_marga_reports %s`, baseWhere)
    var avgSegmentLength float64
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&avgSegmentLength).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate avg segment length: %w", err)
    }
    stats["avg_segment_length_m"] = avgSegmentLength
    
    // 2. Average damage area (m2) - use total_damaged_area or fallback to damaged_area
    query = fmt.Sprintf(`
        SELECT COALESCE(AVG(COALESCE(total_damaged_area, damaged_area)), 0) 
        FROM bina_marga_reports %s
    `, baseWhere)
    var avgDamageArea float64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&avgDamageArea).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate avg damage area: %w", err)
    }
    stats["avg_damage_area_m2"] = avgDamageArea
    
    // 3. Average daily traffic volume
    query = fmt.Sprintf(`SELECT COALESCE(AVG(daily_traffic_volume), 0) FROM bina_marga_reports %s`, baseWhere)
    var avgTrafficVolume float64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&avgTrafficVolume).Error
    if err != nil {
        return nil, fmt.Errorf("failed to calculate avg traffic volume: %w", err)
    }
    stats["avg_daily_traffic_volume"] = avgTrafficVolume
    
    // 4. Total infrastructure reports count
    query = fmt.Sprintf(`SELECT COUNT(*) FROM bina_marga_reports %s`, baseWhere)
    var totalReports int64
    err = r.db.WithContext(ctx).Raw(query, args...).Scan(&totalReports).Error
    if err != nil {
        return nil, fmt.Errorf("failed to count total reports: %w", err)
    }
    stats["total_infrastructure_reports"] = totalReports
    
    return stats, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaLocationStats(ctx context.Context, roadType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            road_name,
            latitude,
            longitude,
            damage_type,
            damage_level,
            urgency_level,
            traffic_impact,
            COALESCE(total_damaged_area, damaged_area) as damaged_area
        FROM bina_marga_reports
        WHERE latitude IS NOT NULL AND longitude IS NOT NULL
    `
    
    args := []interface{}{}
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        query += " AND road_type = $1"
        args = append(args, roadType)
    }
    
    query += " ORDER BY damaged_area DESC"
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get location statistics: %w", err)
    }
    
    return results, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaPriorityStats(ctx context.Context, roadType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN urgency_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(urgency_level) = '' THEN 'EMPTY'
                ELSE urgency_level 
            END as priority_level,
            COUNT(*) as count
        FROM bina_marga_reports
    `
    
    args := []interface{}{}
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        query += " WHERE road_type = $1"
        args = append(args, roadType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN urgency_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(urgency_level) = '' THEN 'EMPTY'
                ELSE urgency_level 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get priority statistics: %w", err)
    }
    
    return results, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaRoadDamageLevelStats(ctx context.Context, roadType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN damage_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_level) = '' THEN 'EMPTY'
                ELSE damage_level 
            END as damage_level,
            COUNT(*) as count
        FROM bina_marga_reports
        WHERE (bridge_name IS NULL OR bridge_name = '')
    `
    
    args := []interface{}{}
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        query += " AND road_type = $1"
        args = append(args, roadType)
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
        return nil, fmt.Errorf("failed to get road damage level statistics: %w", err)
    }
    
    return results, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaBridgeDamageLevelStats(ctx context.Context, roadType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN bridge_damage_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(bridge_damage_level) = '' THEN 'EMPTY'
                ELSE bridge_damage_level 
            END as damage_level,
            COUNT(*) as count
        FROM bina_marga_reports
        WHERE bridge_name IS NOT NULL AND bridge_name != ''
    `
    
    args := []interface{}{}
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        query += " AND road_type = $1"
        args = append(args, roadType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN bridge_damage_level IS NULL THEN 'NOT_SET'
                WHEN TRIM(bridge_damage_level) = '' THEN 'EMPTY'
                ELSE bridge_damage_level 
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get bridge damage level statistics: %w", err)
    }
    
    return results, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaTopRoadDamageTypes(ctx context.Context, roadType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN damage_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_type) = '' THEN 'EMPTY'
                ELSE damage_type 
            END as damage_type,
            COUNT(*) as count
        FROM bina_marga_reports
        WHERE (bridge_name IS NULL OR bridge_name = '')
    `
    
    args := []interface{}{}
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        query += " AND road_type = $1"
        args = append(args, roadType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN damage_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(damage_type) = '' THEN 'EMPTY'
                ELSE damage_type 
            END
        ORDER BY count DESC
        LIMIT 10
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get top road damage types: %w", err)
    }
    
    return results, nil
}

func (r *binaMargaRepositoryImpl) GetBinaMargaTopBridgeDamageTypes(ctx context.Context, roadType string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN bridge_damage_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(bridge_damage_type) = '' THEN 'EMPTY'
                ELSE bridge_damage_type 
            END as damage_type,
            COUNT(*) as count
        FROM bina_marga_reports
        WHERE bridge_name IS NOT NULL AND bridge_name != ''
    `
    
    args := []interface{}{}
    if roadType != "" && roadType != "all" && roadType != "ALL" {
        query += " AND road_type = $1"
        args = append(args, roadType)
    }
    
    query += `
        GROUP BY 
            CASE 
                WHEN bridge_damage_type IS NULL THEN 'NOT_SET'
                WHEN TRIM(bridge_damage_type) = '' THEN 'EMPTY'
                ELSE bridge_damage_type 
            END
        ORDER BY count DESC
        LIMIT 10
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get top bridge damage types: %w", err)
    }
    
    return results, nil
}
