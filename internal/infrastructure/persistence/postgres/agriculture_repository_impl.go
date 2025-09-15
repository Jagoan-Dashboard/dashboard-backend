
package postgres

import (
    "context"
    "time"
    "building-report-backend/internal/domain/entity"
    "building-report-backend/internal/domain/repository"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type agricultureRepositoryImpl struct {
    db *gorm.DB
}

func NewAgricultureRepository(db *gorm.DB) repository.AgricultureRepository {
    return &agricultureRepositoryImpl{db: db}
}

func (r *agricultureRepositoryImpl) Create(ctx context.Context, report *entity.AgricultureReport) error {
    return r.db.WithContext(ctx).Create(report).Error
}

func (r *agricultureRepositoryImpl) Update(ctx context.Context, report *entity.AgricultureReport) error {
    report.UpdatedAt = time.Now()
    return r.db.WithContext(ctx).Save(report).Error
}

func (r *agricultureRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.AgricultureReport{}).Error
}

func (r *agricultureRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.AgricultureReport, error) {
    var report entity.AgricultureReport
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("id = ?", id).
        First(&report).Error
    
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *agricultureRepositoryImpl) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.AgricultureReport, int64, error) {
    var reports []*entity.AgricultureReport
    var total int64

    query := r.db.WithContext(ctx).Model(&entity.AgricultureReport{})

    
    if extensionOfficer, ok := filters["extension_officer"].(string); ok && extensionOfficer != "" {
        query = query.Where("extension_officer = ?", extensionOfficer)
    }
    if village, ok := filters["village"].(string); ok && village != "" {
        query = query.Where("village = ?", village)
    }
    if district, ok := filters["district"].(string); ok && district != "" {
        query = query.Where("district = ?", district)
    }
    if farmerName, ok := filters["farmer_name"].(string); ok && farmerName != "" {
        query = query.Where("farmer_name ILIKE ?", "%"+farmerName+"%")
    }
    if foodCommodity, ok := filters["food_commodity"].(string); ok && foodCommodity != "" {
        query = query.Where("food_commodity = ?", foodCommodity)
    }
    if hortiCommodity, ok := filters["horti_commodity"].(string); ok && hortiCommodity != "" {
        query = query.Where("horti_commodity = ?", hortiCommodity)
    }
    if plantationCommodity, ok := filters["plantation_commodity"].(string); ok && plantationCommodity != "" {
        query = query.Where("plantation_commodity = ?", plantationCommodity)
    }
    if hasPestDisease, ok := filters["has_pest_disease"].(bool); ok {
        query = query.Where("has_pest_disease = ?", hasPestDisease)
    }
    if mainConstraint, ok := filters["main_constraint"].(string); ok && mainConstraint != "" {
        query = query.Where("main_constraint = ?", mainConstraint)
    }

    
    if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
        query = query.Where("visit_date >= ?", startDate)
    }
    if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
        query = query.Where("visit_date <= ?", endDate)
    }

    
    query.Count(&total)

    
    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("visit_date DESC, created_at DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *agricultureRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.AgricultureReport, int64, error) {
    var reports []*entity.AgricultureReport
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
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

func (r *agricultureRepositoryImpl) FindByExtensionOfficer(ctx context.Context, extensionOfficer string, limit, offset int) ([]*entity.AgricultureReport, int64, error) {
    var reports []*entity.AgricultureReport
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Where("extension_officer = ?", extensionOfficer)

    query.Count(&total)

    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("visit_date DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *agricultureRepositoryImpl) FindByVillage(ctx context.Context, village string, limit, offset int) ([]*entity.AgricultureReport, int64, error) {
    var reports []*entity.AgricultureReport
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Where("village = ?", village)

    query.Count(&total)

    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("visit_date DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *agricultureRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.AgricultureReport, int64, error) {
    var reports []*entity.AgricultureReport
    var total int64

    query := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", startDate, endDate)

    query.Count(&total)

    err := query.
        Preload("Photos").
        Limit(limit).
        Offset(offset).
        Order("visit_date DESC").
        Find(&reports).Error

    return reports, total, err
}

func (r *agricultureRepositoryImpl) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    
    var total int64
    r.db.Model(&entity.AgricultureReport{}).Count(&total)
    stats["total_reports"] = total
    
    
    var totalFarmers int64
    r.db.Model(&entity.AgricultureReport{}).
        Distinct("farmer_name").
        Count(&totalFarmers)
    stats["total_farmers"] = totalFarmers
    
    
    var totalLandArea float64
    r.db.Model(&entity.AgricultureReport{}).
        Select("COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0)").
        Scan(&totalLandArea)
    stats["total_land_area_ha"] = totalLandArea
    
    
    var foodCropReports int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("food_commodity != '' AND food_commodity IS NOT NULL").
        Count(&foodCropReports)
    stats["food_crop_reports"] = foodCropReports
    
    var horticultureReports int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("horti_commodity != '' AND horti_commodity IS NOT NULL").
        Count(&horticultureReports)
    stats["horticulture_reports"] = horticultureReports
    
    var plantationReports int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("plantation_commodity != '' AND plantation_commodity IS NOT NULL").
        Count(&plantationReports)
    stats["plantation_reports"] = plantationReports
    
    
    var pestDiseaseReports int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("has_pest_disease = true").
        Count(&pestDiseaseReports)
    stats["reports_with_pest_disease"] = pestDiseaseReports
    stats["pest_disease_percentage"] = float64(pestDiseaseReports) / float64(total) * 100
    
    
    var postHarvestProblemReports int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("post_harvest_problems != '' AND post_harvest_problems IS NOT NULL AND post_harvest_problems != 'TIDAK_ADA'").
        Count(&postHarvestProblemReports)
    stats["post_harvest_problem_reports"] = postHarvestProblemReports
    
    var productionProblemReports int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("production_problems != '' AND production_problems IS NOT NULL").
        Count(&productionProblemReports)
    stats["production_problem_reports"] = productionProblemReports
    
    
    type VillageCount struct {
        Village string `json:"village"`
        Count   int64  `json:"count"`
    }
    var villageCounts []VillageCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("village, count(*) as count").
        Group("village").
        Order("count DESC").
        Limit(10).
        Scan(&villageCounts)
    stats["village_distribution"] = villageCounts
    
    
    type ExtensionOfficerStats struct {
        ExtensionOfficer string `json:"extension_officer"`
        VisitCount       int64  `json:"visit_count"`
        FarmerCount      int64  `json:"farmer_count"`
    }
    var extensionStats []ExtensionOfficerStats
    r.db.Raw(`
        SELECT 
            extension_officer,
            COUNT(*) as visit_count,
            COUNT(DISTINCT farmer_name) as farmer_count
        FROM agriculture_reports 
        GROUP BY extension_officer 
        ORDER BY visit_count DESC
        LIMIT 10
    `).Scan(&extensionStats)
    stats["extension_officer_stats"] = extensionStats
    
    return stats, nil
}

func (r *agricultureRepositoryImpl) GetCommodityProduction(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            CASE 
                WHEN food_commodity IS NOT NULL AND food_commodity != '' THEN food_commodity
                WHEN horti_commodity IS NOT NULL AND horti_commodity != '' THEN horti_commodity
                WHEN plantation_commodity IS NOT NULL AND plantation_commodity != '' THEN plantation_commodity
                ELSE 'UNKNOWN'
            END as commodity,
            COUNT(*) as report_count,
            SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) as total_area,
            COUNT(DISTINCT farmer_name) as farmer_count,
            COUNT(DISTINCT village) as village_count
        FROM agriculture_reports
        WHERE visit_date BETWEEN ? AND ?
        GROUP BY commodity
        ORDER BY total_area DESC
    `
    
    err := r.db.WithContext(ctx).
        Raw(query, startDate, endDate).
        Scan(&results).Error
    
    return results, err
}

func (r *agricultureRepositoryImpl) GetExtensionOfficerPerformance(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            extension_officer,
            COUNT(*) as total_visits,
            COUNT(DISTINCT farmer_name) as farmers_visited,
            COUNT(DISTINCT village) as villages_covered,
            MAX(visit_date) as last_visit,
            COUNT(*) / GREATEST(EXTRACT(EPOCH FROM (? - ?)) / (30 * 24 * 3600), 1) as average_visits_per_month
        FROM agriculture_reports
        WHERE visit_date BETWEEN ? AND ?
        GROUP BY extension_officer
        ORDER BY total_visits DESC
    `
    
    err := r.db.WithContext(ctx).
        Raw(query, endDate, startDate, startDate, endDate).
        Scan(&results).Error
    
    return results, err
}

func (r *agricultureRepositoryImpl) GetVillageProductionStats(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            village,
            district,
            COUNT(*) as total_reports,
            COUNT(DISTINCT farmer_name) as farmer_count,
            SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) as total_land_area,
            COUNT(CASE WHEN has_pest_disease = true THEN 1 END) as pest_disease_reports,
            COUNT(DISTINCT extension_officer) as extension_officers
        FROM agriculture_reports
        WHERE visit_date BETWEEN ? AND ?
        GROUP BY village, district
        ORDER BY total_land_area DESC
    `
    
    err := r.db.WithContext(ctx).
        Raw(query, startDate, endDate).
        Scan(&results).Error
    
    return results, err
}

func (r *agricultureRepositoryImpl) GetPestDiseaseReports(ctx context.Context, limit int) ([]*entity.AgricultureReport, error) {
    var reports []*entity.AgricultureReport
    
    err := r.db.WithContext(ctx).
        Preload("Photos").
        Where("has_pest_disease = true").
        Order("visit_date DESC").
        Limit(limit).
        Find(&reports).Error
    
    return reports, err
}

func (r *agricultureRepositoryImpl) GetTechnologyAdoptionStats(ctx context.Context) (map[string]interface{}, error) {
    stats := make(map[string]interface{})
    
    
    type TechnologyCount struct {
        Technology string `json:"technology"`
        Count      int64  `json:"count"`
    }
    
    var foodTechCounts []TechnologyCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("food_technology as technology, count(*) as count").
        Where("food_technology IS NOT NULL AND food_technology != '' AND food_technology != 'TIDAK_ADA'").
        Group("food_technology").
        Order("count DESC").
        Scan(&foodTechCounts)
    stats["food_technology"] = foodTechCounts
    
    var hortiTechCounts []TechnologyCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("horti_technology as technology, count(*) as count").
        Where("horti_technology IS NOT NULL AND horti_technology != '' AND horti_technology != 'TIDAK_ADA'").
        Group("horti_technology").
        Order("count DESC").
        Scan(&hortiTechCounts)
    stats["horticulture_technology"] = hortiTechCounts
    
    var plantationTechCounts []TechnologyCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("plantation_technology as technology, count(*) as count").
        Where("plantation_technology IS NOT NULL AND plantation_technology != '' AND plantation_technology != 'TIDAK_ADA'").
        Group("plantation_technology").
        Order("count DESC").
        Scan(&plantationTechCounts)
    stats["plantation_technology"] = plantationTechCounts
    
    return stats, nil
}

func (r *agricultureRepositoryImpl) GetFarmerNeedsAnalysis(ctx context.Context) (map[string]interface{}, error) {
    analysis := make(map[string]interface{})
    
    type NeedCount struct {
        Need  string `json:"need"`
        Count int64  `json:"count"`
    }
    
    
    var constraintCounts []NeedCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("main_constraint as need, count(*) as count").
        Where("main_constraint IS NOT NULL AND main_constraint != ''").
        Group("main_constraint").
        Order("count DESC").
        Scan(&constraintCounts)
    analysis["main_constraints"] = constraintCounts
    
    
    var hopeCounts []NeedCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("farmer_hope as need, count(*) as count").
        Where("farmer_hope IS NOT NULL AND farmer_hope != ''").
        Group("farmer_hope").
        Order("count DESC").
        Scan(&hopeCounts)
    analysis["farmer_hopes"] = hopeCounts
    
    
    var trainingCounts []NeedCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("training_needed as need, count(*) as count").
        Where("training_needed IS NOT NULL AND training_needed != ''").
        Group("training_needed").
        Order("count DESC").
        Scan(&trainingCounts)
    analysis["training_needs"] = trainingCounts
    
    
    var urgentCounts []NeedCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("urgent_needs as need, count(*) as count").
        Where("urgent_needs IS NOT NULL AND urgent_needs != ''").
        Group("urgent_needs").
        Order("count DESC").
        Scan(&urgentCounts)
    analysis["urgent_needs"] = urgentCounts
    
    
    var waterAccessCounts []NeedCount
    r.db.Model(&entity.AgricultureReport{}).
        Select("water_access as need, count(*) as count").
        Where("water_access IS NOT NULL AND water_access != ''").
        Group("water_access").
        Order("count DESC").
        Scan(&waterAccessCounts)
    analysis["water_access"] = waterAccessCounts
    
    return analysis, nil
}

func (r *agricultureRepositoryImpl) CountTotalFarmers(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Distinct("farmer_name").
        Count(&count).Error
    return count, err
}

func (r *agricultureRepositoryImpl) CalculateTotalLandArea(ctx context.Context) (float64, error) {
    var total float64
    err := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Select("COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0)").
        Scan(&total).Error
    return total, err
}

func (r *agricultureRepositoryImpl) CountReportsByCommodityType(ctx context.Context) (map[string]int64, error) {
    counts := make(map[string]int64)
    
    var foodCount int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("food_commodity IS NOT NULL AND food_commodity != ''").
        Count(&foodCount)
    counts["food_crops"] = foodCount
    
    var hortiCount int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("horti_commodity IS NOT NULL AND horti_commodity != ''").
        Count(&hortiCount)
    counts["horticulture"] = hortiCount
    
    var plantationCount int64
    r.db.Model(&entity.AgricultureReport{}).
        Where("plantation_commodity IS NOT NULL AND plantation_commodity != ''").
        Count(&plantationCount)
    counts["plantation"] = plantationCount
    
    return counts, nil
}

func (r *agricultureRepositoryImpl) CountReportsWithPestDisease(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Where("has_pest_disease = true").
        Count(&count).Error
    return count, err
}

func (r *agricultureRepositoryImpl) GetTopConstraints(ctx context.Context, limit int) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    err := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Select("main_constraint as constraint, count(*) as count").
        Where("main_constraint IS NOT NULL AND main_constraint != ''").
        Group("main_constraint").
        Order("count DESC").
        Limit(limit).
        Scan(&results).Error
    
    return results, err
}

func (r *agricultureRepositoryImpl) GetTopFarmerHopes(ctx context.Context, limit int) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    err := r.db.WithContext(ctx).
        Model(&entity.AgricultureReport{}).
        Select("farmer_hope as hope, count(*) as count").
        Where("farmer_hope IS NOT NULL AND farmer_hope != ''").
        Group("farmer_hope").
        Order("count DESC").
        Limit(limit).
        Scan(&results).Error
    
    return results, err
}