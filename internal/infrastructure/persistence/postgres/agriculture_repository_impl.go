package postgres

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"
	"fmt"
    "strings"
	"time"

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

func (r *agricultureRepositoryImpl) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.AgricultureReport{}).Error
}

func (r *agricultureRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.AgricultureReport, error) {
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

func (r *agricultureRepositoryImpl) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.AgricultureReport, int64, error) {
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


func (r *agricultureRepositoryImpl) GetExecutiveSummary(ctx context.Context) (map[string]interface{}, error) {
    summary := make(map[string]interface{})
    
    // Total land area
    var totalLandArea float64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Select("COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0)").
        Scan(&totalLandArea)
    summary["total_land_area"] = totalLandArea
    
    // Pest disease reports count
    var pestReports int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("has_pest_disease = true").
        Count(&pestReports)
    summary["pest_disease_reports"] = pestReports
    
    // Total extension reports
    var totalReports int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Count(&totalReports)
    summary["total_extension_reports"] = totalReports
    
    return summary, nil
}

func (r *agricultureRepositoryImpl) GetCommodityDistributionByDistrict(ctx context.Context) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            latitude, longitude, village, district,
            CASE 
                WHEN food_commodity IS NOT NULL AND food_commodity != '' THEN food_commodity
                WHEN horti_commodity IS NOT NULL AND horti_commodity != '' THEN horti_commodity  
                WHEN plantation_commodity IS NOT NULL AND plantation_commodity != '' THEN plantation_commodity
                ELSE 'UNKNOWN'
            END as commodity,
            CASE 
                WHEN food_commodity IS NOT NULL AND food_commodity != '' THEN 'FOOD'
                WHEN horti_commodity IS NOT NULL AND horti_commodity != '' THEN 'HORTICULTURE'
                WHEN plantation_commodity IS NOT NULL AND plantation_commodity != '' THEN 'PLANTATION'
                ELSE 'UNKNOWN'
            END as commodity_type,
            COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0) as land_area
        FROM agriculture_reports
        WHERE latitude IS NOT NULL AND longitude IS NOT NULL
        AND (food_commodity IS NOT NULL OR horti_commodity IS NOT NULL OR plantation_commodity IS NOT NULL)
    `
    
    err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetCommodityCountBySector(ctx context.Context) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    // Food crops
    var foodCrops []map[string]interface{}
    r.db.WithContext(ctx).Raw(`
        SELECT food_commodity as name, COUNT(*) as count
        FROM agriculture_reports
        WHERE food_commodity IS NOT NULL AND food_commodity != ''
        GROUP BY food_commodity
        ORDER BY count DESC
    `).Scan(&foodCrops)
    result["food_crops"] = foodCrops
    
    // Horticulture
    var horticulture []map[string]interface{}
    r.db.WithContext(ctx).Raw(`
        SELECT horti_commodity as name, COUNT(*) as count
        FROM agriculture_reports
        WHERE horti_commodity IS NOT NULL AND horti_commodity != ''
        GROUP BY horti_commodity
        ORDER BY count DESC
    `).Scan(&horticulture)
    result["horticulture"] = horticulture
    
    // Plantation
    var plantation []map[string]interface{}
    r.db.WithContext(ctx).Raw(`
        SELECT plantation_commodity as name, COUNT(*) as count
        FROM agriculture_reports
        WHERE plantation_commodity IS NOT NULL AND plantation_commodity != ''
        GROUP BY plantation_commodity
        ORDER BY count DESC
    `).Scan(&plantation)
    result["plantation"] = plantation
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetLandStatusDistribution(ctx context.Context) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            land_status as status,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM (
            SELECT food_land_status as land_status FROM agriculture_reports WHERE food_land_status IS NOT NULL AND food_land_status != ''
            UNION ALL
            SELECT horti_land_status as land_status FROM agriculture_reports WHERE horti_land_status IS NOT NULL AND horti_land_status != ''
            UNION ALL
            SELECT plantation_land_status as land_status FROM agriculture_reports WHERE plantation_land_status IS NOT NULL AND plantation_land_status != ''
        ) as combined_status
        GROUP BY land_status
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetMainConstraintsDistribution(ctx context.Context) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            main_constraint as constraint,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE main_constraint IS NOT NULL AND main_constraint != ''
        GROUP BY main_constraint
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetFarmerHopesAndNeeds(ctx context.Context) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    // Farmer hopes
    var hopes []map[string]interface{}
    r.db.WithContext(ctx).Raw(`
        SELECT 
            farmer_hope as hope,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE farmer_hope IS NOT NULL AND farmer_hope != ''
        GROUP BY farmer_hope
        ORDER BY count DESC
    `).Scan(&hopes)
    result["hopes"] = hopes
    
    // Training needs
    var trainingNeeds []map[string]interface{}
    r.db.WithContext(ctx).Raw(`
        SELECT 
            training_needed as training,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE training_needed IS NOT NULL AND training_needed != ''
        GROUP BY training_needed
        ORDER BY count DESC
    `).Scan(&trainingNeeds)
    result["training_needs"] = trainingNeeds
    
    // Urgent needs
    var urgentNeeds []map[string]interface{}
    r.db.WithContext(ctx).Raw(`
        SELECT 
            urgent_needs as need,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE urgent_needs IS NOT NULL AND urgent_needs != ''
        GROUP BY urgent_needs
        ORDER BY count DESC
    `).Scan(&urgentNeeds)
    result["urgent_needs"] = urgentNeeds
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetCommodityAnalysis(ctx context.Context, startDate, endDate time.Time, commodityName string) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    // Build flexible WHERE clauses for commodity matching
    commodityWhereClause := `(
        (food_commodity IS NOT NULL AND food_commodity != '' AND UPPER(food_commodity) LIKE UPPER($3)) OR
        (horti_commodity IS NOT NULL AND horti_commodity != '' AND UPPER(horti_commodity) LIKE UPPER($3)) OR
        (plantation_commodity IS NOT NULL AND plantation_commodity != '' AND UPPER(plantation_commodity) LIKE UPPER($3))
    )`
    
    if commodityName == "" {
        // If no specific commodity, get all
        commodityWhereClause = `(
            food_commodity IS NOT NULL AND food_commodity != '' OR
            horti_commodity IS NOT NULL AND horti_commodity != '' OR
            plantation_commodity IS NOT NULL AND plantation_commodity != ''
        )`
    }
    
    // Get current year data
    currentYearQuery := fmt.Sprintf(`
        SELECT 
            COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0) as total_area,
            COUNT(*) as report_count
        FROM agriculture_reports
        WHERE visit_date BETWEEN $1 AND $2
        AND %s
    `, commodityWhereClause)
    
    type yearData struct {
        TotalArea   float64 `json:"total_area"`
        ReportCount int64   `json:"report_count"`
    }
    
    var currentYear yearData
    var args []interface{}
    
    if commodityName == "" {
        args = []interface{}{startDate, endDate}
    } else {
        commodityPattern := "%" + commodityName + "%"
        args = []interface{}{startDate, endDate, commodityPattern}
    }
    
    err := r.db.WithContext(ctx).Raw(currentYearQuery, args...).Scan(&currentYear).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get current year data: %w", err)
    }
    
    // Get previous year data for comparison
    prevYearStart := startDate.AddDate(-1, 0, 0)
    prevYearEnd := endDate.AddDate(-1, 0, 0)
    
    var previousYear yearData
    prevArgs := []interface{}{prevYearStart, prevYearEnd}
    if commodityName != "" {
        commodityPattern := "%" + commodityName + "%"
        prevArgs = append(prevArgs, commodityPattern)
    }
    
    err = r.db.WithContext(ctx).Raw(currentYearQuery, prevArgs...).Scan(&previousYear).Error
    if err != nil {
        return nil, fmt.Errorf("failed to get previous year data: %w", err)
    }
    
    // Calculate estimated production with realistic productivity values
    // Use different productivity rates for different crop types
    var currentProductivity, previousProductivity float64 = 3.0, 3.0 // Default values
    
    // Adjust productivity based on commodity type
    commodityUpper := strings.ToUpper(commodityName)
    switch {
    case strings.Contains(commodityUpper, "PADI"):
        currentProductivity, previousProductivity = 5.2, 5.0 // Rice productivity
    case strings.Contains(commodityUpper, "JAGUNG"):
        currentProductivity, previousProductivity = 4.8, 4.6 // Corn productivity
    case strings.Contains(commodityUpper, "KEDELAI"):
        currentProductivity, previousProductivity = 1.5, 1.4 // Soybean productivity
    case strings.Contains(commodityUpper, "CABAI"):
        currentProductivity, previousProductivity = 8.0, 7.5 // Chili productivity
    case strings.Contains(commodityUpper, "TOMAT"):
        currentProductivity, previousProductivity = 12.0, 11.5 // Tomato productivity
    case strings.Contains(commodityUpper, "KOPI"):
        currentProductivity, previousProductivity = 0.8, 0.75 // Coffee productivity
    case strings.Contains(commodityUpper, "KAKAO"):
        currentProductivity, previousProductivity = 0.6, 0.55 // Cocoa productivity
    }
    
    currentProduction := currentYear.TotalArea * currentProductivity
    previousProduction := previousYear.TotalArea * previousProductivity
    
    // Calculate growth percentages
    var productionGrowth, areaGrowth, productivityGrowth float64
    
    if previousProduction > 0 {
        productionGrowth = ((currentProduction - previousProduction) / previousProduction) * 100
    }
    
    if previousYear.TotalArea > 0 {
        areaGrowth = ((currentYear.TotalArea - previousYear.TotalArea) / previousYear.TotalArea) * 100
    }
    
    if previousProductivity > 0 {
        productivityGrowth = ((currentProductivity - previousProductivity) / previousProductivity) * 100
    }
    
    result["total_production"] = currentProduction
    result["production_growth"] = productionGrowth
    result["total_harvested_area"] = currentYear.TotalArea
    result["harvested_area_growth"] = areaGrowth
    result["productivity"] = currentProductivity
    result["productivity_growth"] = productivityGrowth
    
    return result, nil
}

// Fixed GetProductionByDistrict method
func (r *agricultureRepositoryImpl) GetProductionByDistrict(ctx context.Context, startDate, endDate time.Time, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    commodityWhereClause := `(
        (food_commodity IS NOT NULL AND food_commodity != '' AND UPPER(food_commodity) LIKE UPPER($4)) OR
        (horti_commodity IS NOT NULL AND horti_commodity != '' AND UPPER(horti_commodity) LIKE UPPER($4)) OR
        (plantation_commodity IS NOT NULL AND plantation_commodity != '' AND UPPER(plantation_commodity) LIKE UPPER($4))
    )`
    
    if commodityName == "" {
        commodityWhereClause = `(
            food_commodity IS NOT NULL AND food_commodity != '' OR
            horti_commodity IS NOT NULL AND horti_commodity != '' OR
            plantation_commodity IS NOT NULL AND plantation_commodity != ''
        )`
    }
    
    query := fmt.Sprintf(`
        SELECT 
            district,
            COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0) as harvested_area,
            COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) * 3.0, 0) as production,
            COUNT(DISTINCT farmer_name) as farmer_count
        FROM agriculture_reports
        WHERE visit_date BETWEEN $1 AND $2
        AND %s
        GROUP BY district
        HAVING SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) > 0
        ORDER BY production DESC
    `, commodityWhereClause)
    
    var args []interface{}
    if commodityName == "" {
        args = []interface{}{startDate, endDate}
    } else {
        commodityPattern := "%" + commodityName + "%"
        args = []interface{}{startDate, endDate, commodityPattern}
    }
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetProductivityTrend(ctx context.Context, commodityName string, years []int) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    for _, year := range years {
        startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
        endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
        
        // Use fmt.Sprintf to inject the year as a literal value instead of a parameter
        query := fmt.Sprintf(`
            SELECT 
                %d as year,
                COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0) as area,
                COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) * 3.0, 0) as production,
                CASE 
                    WHEN SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) > 0
                    THEN 3.0
                    ELSE 0
                END as productivity
            FROM agriculture_reports
            WHERE visit_date BETWEEN $1 AND $2
            AND (food_commodity = $3 OR horti_commodity = $4 OR plantation_commodity = $5)
        `, year)
        
        // Use a more flexible approach with interface{} and type assertions
        rows, err := r.db.WithContext(ctx).Raw(query, startDate, endDate, commodityName, commodityName, commodityName).Rows()
        if err != nil {
            return nil, fmt.Errorf("failed to execute query for year %d: %w", year, err)
        }
        defer rows.Close()
        
        if rows.Next() {
            var yearValue, areaValue, productionValue, productivityValue interface{}
            
            err := rows.Scan(&yearValue, &areaValue, &productionValue, &productivityValue)
            if err != nil {
                return nil, fmt.Errorf("failed to scan results for year %d: %w", year, err)
            }
            
            // Convert values with proper type handling
            result := map[string]interface{}{
                "year":         convertToInt64(yearValue),
                "area":         convertToFloat64(areaValue),
                "production":   convertToFloat64(productionValue),
                "productivity": convertToFloat64(productivityValue),
            }
            
            results = append(results, result)
        } else {
            // If no data for this year, add zero values
            result := map[string]interface{}{
                "year":         int64(year),
                "area":         float64(0),
                "production":   float64(0),
                "productivity": float64(0),
            }
            results = append(results, result)
        }
    }
    
    return results, nil
}

// Helper functions for type conversion
func convertToInt64(value interface{}) int64 {
    switch v := value.(type) {
    case int:
        return int64(v)
    case int32:
        return int64(v)
    case int64:
        return v
    case float64:
        return int64(v)
    case float32:
        return int64(v)
    default:
        return 0
    }
}

func convertToFloat64(value interface{}) float64 {
    switch v := value.(type) {
    case float64:
        return v
    case float32:
        return float64(v)
    case int:
        return float64(v)
    case int32:
        return float64(v)
    case int64:
        return float64(v)
    default:
        return 0.0
    }
}

func (r *agricultureRepositoryImpl) GetFoodCropStats(ctx context.Context, commodityName string) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    whereClause := "food_commodity IS NOT NULL AND food_commodity != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND food_commodity = ?"
        args = append(args, commodityName)
    }
    
    // Land area
    var landArea float64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where(whereClause, args...).
        Select("COALESCE(SUM(food_land_area), 0)").
        Scan(&landArea)
    result["land_area"] = landArea
    
    // Estimated production (assuming 3 tons/hectare for food crops)
    result["estimated_production"] = landArea * 3.0
    
    // Pest affected area
    var pestAffectedArea float64
    pestQuery := whereClause + " AND has_pest_disease = true"
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where(pestQuery, args...).
        Select("COALESCE(SUM(food_land_area), 0)").
        Scan(&pestAffectedArea)
    result["pest_affected_area"] = pestAffectedArea
    
    // Pest report count
    var pestReportCount int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where(pestQuery, args...).
        Count(&pestReportCount)
    result["pest_report_count"] = pestReportCount
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetFoodCropDistribution(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "food_commodity IS NOT NULL AND food_commodity != '' AND latitude IS NOT NULL AND longitude IS NOT NULL"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND food_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            latitude, longitude, village, district, 
            food_commodity as commodity,
            food_land_area as land_area
        FROM agriculture_reports
        WHERE ` + whereClause
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetFoodCropGrowthPhases(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "food_commodity IS NOT NULL AND food_commodity != '' AND food_growth_phase IS NOT NULL AND food_growth_phase != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND food_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            food_growth_phase as phase,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY food_growth_phase
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetFoodCropTechnology(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "food_commodity IS NOT NULL AND food_commodity != '' AND food_technology IS NOT NULL AND food_technology != '' AND food_technology != 'TIDAK_ADA'"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND food_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            food_technology as technology,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY food_technology
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetFoodCropPestDominance(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "food_commodity IS NOT NULL AND food_commodity != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND food_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            CASE 
                WHEN has_pest_disease = false OR pest_disease_type IS NULL OR pest_disease_type = '' THEN 'TIDAK_ADA'
                ELSE pest_disease_type
            END as pest_type,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY 
            CASE 
                WHEN has_pest_disease = false OR pest_disease_type IS NULL OR pest_disease_type = '' THEN 'TIDAK_ADA'
                ELSE pest_disease_type
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetFoodCropHarvestSchedule(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "food_commodity IS NOT NULL AND food_commodity != '' AND food_harvest_date IS NOT NULL"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND food_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            food_commodity as commodity_detail,
            food_harvest_date as harvest_date,
            farmer_name,
            village,
            food_land_area as land_area
        FROM agriculture_reports
        WHERE ` + whereClause + `
        AND food_harvest_date >= CURRENT_DATE
        ORDER BY food_harvest_date ASC
        LIMIT 20
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetHorticultureStats(ctx context.Context, commodityName string) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    baseQuery := r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("horti_commodity IS NOT NULL AND horti_commodity != ''")
    
    if commodityName != "" {
        baseQuery = baseQuery.Where("horti_commodity = ?", commodityName)
    }
    
    // Land area
    var landArea float64
    baseQuery.Select("COALESCE(SUM(horti_land_area), 0)").Scan(&landArea)
    result["land_area"] = landArea
    
    // Estimated production (assuming 5 tons/hectare for horticulture)
    result["estimated_production"] = landArea * 5.0
    
    // Pest affected area
    var pestAffectedArea float64
    pestQuery := r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("horti_commodity IS NOT NULL AND horti_commodity != '' AND has_pest_disease = true")
    
    if commodityName != "" {
        pestQuery = pestQuery.Where("horti_commodity = ?", commodityName)
    }
    
    pestQuery.Select("COALESCE(SUM(horti_land_area), 0)").Scan(&pestAffectedArea)
    result["pest_affected_area"] = pestAffectedArea
    
    // Pest report count
    var pestReportCount int64
    pestQuery.Count(&pestReportCount)
    result["pest_report_count"] = pestReportCount
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetHorticultureDistribution(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            latitude, longitude, village, district, 
            horti_commodity as commodity,
            horti_land_area as land_area
        FROM agriculture_reports
        WHERE horti_commodity IS NOT NULL AND horti_commodity != '' 
        AND latitude IS NOT NULL AND longitude IS NOT NULL`
    
    args := []interface{}{}
    
    if commodityName != "" {
        query += " AND horti_commodity = ?"
        args = append(args, commodityName)
    }
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetHorticultureGrowthPhases(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            horti_growth_phase as phase,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE horti_commodity IS NOT NULL AND horti_commodity != '' 
        AND horti_growth_phase IS NOT NULL AND horti_growth_phase != ''`
    
    args := []interface{}{}
    
    if commodityName != "" {
        query += " AND horti_commodity = ?"
        args = append(args, commodityName)
    }
    
    query += `
        GROUP BY horti_growth_phase
        ORDER BY count DESC`
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetHorticultureTechnology(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "horti_commodity IS NOT NULL AND horti_commodity != '' AND horti_technology IS NOT NULL AND horti_technology != '' AND horti_technology != 'TIDAK_ADA'"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND horti_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            horti_technology as technology,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY horti_technology
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetHorticulturePestDominance(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "horti_commodity IS NOT NULL AND horti_commodity != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND horti_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            CASE 
                WHEN has_pest_disease = false OR pest_disease_type IS NULL OR pest_disease_type = '' THEN 'TIDAK_ADA'
                ELSE pest_disease_type
            END as pest_type,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY 
            CASE 
                WHEN has_pest_disease = false OR pest_disease_type IS NULL OR pest_disease_type = '' THEN 'TIDAK_ADA'
                ELSE pest_disease_type
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetHorticultureHarvestSchedule(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "horti_commodity IS NOT NULL AND horti_commodity != '' AND horti_harvest_date IS NOT NULL"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND horti_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
    SELECT 
        horti_commodity as commodity_detail,
        DATE(horti_harvest_date) as harvest_date,  
        farmer_name,
        village,
        horti_land_area as land_area
    FROM agriculture_reports
    WHERE ` + whereClause + `
    AND horti_harvest_date >= CURRENT_DATE
    ORDER BY horti_harvest_date ASC
    LIMIT 20
`
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

// Similar implementations for Plantation methods
func (r *agricultureRepositoryImpl) GetPlantationStats(ctx context.Context, commodityName string) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    whereClause := "plantation_commodity IS NOT NULL AND plantation_commodity != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND plantation_commodity = ?"
        args = append(args, commodityName)
    }
    
    // Land area
    var landArea float64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where(whereClause, args...).
        Select("COALESCE(SUM(plantation_land_area), 0)").
        Scan(&landArea)
    result["land_area"] = landArea
    
    // Estimated production (assuming 2 tons/hectare for plantation)
    result["estimated_production"] = landArea * 2.0
    
    // Pest affected area
    var pestAffectedArea float64
    pestQuery := whereClause + " AND has_pest_disease = true"
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where(pestQuery, args...).
        Select("COALESCE(SUM(plantation_land_area), 0)").
        Scan(&pestAffectedArea)
    result["pest_affected_area"] = pestAffectedArea
    
    // Pest report count
    var pestReportCount int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where(pestQuery, args...).
        Count(&pestReportCount)
    result["pest_report_count"] = pestReportCount
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetPlantationDistribution(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "plantation_commodity IS NOT NULL AND plantation_commodity != '' AND latitude IS NOT NULL AND longitude IS NOT NULL"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND plantation_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            latitude, longitude, village, district, 
            plantation_commodity as commodity,
            plantation_land_area as land_area
        FROM agriculture_reports
        WHERE ` + whereClause
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetPlantationGrowthPhases(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "plantation_commodity IS NOT NULL AND plantation_commodity != '' AND plantation_growth_phase IS NOT NULL AND plantation_growth_phase != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND plantation_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            plantation_growth_phase as phase,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY plantation_growth_phase
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetPlantationTechnology(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "plantation_commodity IS NOT NULL AND plantation_commodity != '' AND plantation_technology IS NOT NULL AND plantation_technology != '' AND plantation_technology != 'TIDAK_ADA'"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND plantation_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            plantation_technology as technology,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY plantation_technology
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetPlantationPestDominance(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "plantation_commodity IS NOT NULL AND plantation_commodity != ''"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND plantation_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            CASE 
                WHEN has_pest_disease = false OR pest_disease_type IS NULL OR pest_disease_type = '' THEN 'TIDAK_ADA'
                ELSE pest_disease_type
            END as pest_type,
            COUNT(*) as count,
            ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
        FROM agriculture_reports
        WHERE ` + whereClause + `
        GROUP BY 
            CASE 
                WHEN has_pest_disease = false OR pest_disease_type IS NULL OR pest_disease_type = '' THEN 'TIDAK_ADA'
                ELSE pest_disease_type
            END
        ORDER BY count DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetPlantationHarvestSchedule(ctx context.Context, commodityName string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    whereClause := "plantation_commodity IS NOT NULL AND plantation_commodity != '' AND plantation_harvest_date IS NOT NULL"
    args := []interface{}{}
    
    if commodityName != "" {
        whereClause += " AND plantation_commodity = ?"
        args = append(args, commodityName)
    }
    
    query := `
        SELECT 
            plantation_commodity as commodity_detail,
            plantation_harvest_date as harvest_date,
            farmer_name,
            village,
            plantation_land_area as land_area
        FROM agriculture_reports
        WHERE ` + whereClause + `
        AND plantation_harvest_date >= CURRENT_DATE
        ORDER BY plantation_harvest_date ASC
        LIMIT 20
    `
    
    err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error
    return results, err
}

func (r *agricultureRepositoryImpl) GetAgriculturalEquipmentStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    // Get current year reports with technology
    var currentYearReports int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", startDate, endDate).
        Where("food_technology IS NOT NULL OR horti_technology IS NOT NULL OR plantation_technology IS NOT NULL").
        Count(&currentYearReports)
    
    // Get previous year for comparison
    prevYearStart := startDate.AddDate(-1, 0, 0)
    prevYearEnd := endDate.AddDate(-1, 0, 0)
    
    var prevYearReports int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", prevYearStart, prevYearEnd).
        Where("food_technology IS NOT NULL OR horti_technology IS NOT NULL OR plantation_technology IS NOT NULL").
        Count(&prevYearReports)
    
    // Calculate growth
    var growth float64 = 0
    if prevYearReports > 0 {
        growth = ((float64(currentYearReports) - float64(prevYearReports)) / float64(prevYearReports)) * 100
    }
    
    // Estimated equipment counts based on technology adoption
    result["grain_processor_count"] = int64(float64(currentYearReports) * 0.3)
    result["grain_processor_growth"] = growth
    
    result["thresher_count"] = int64(float64(currentYearReports) * 0.25)
    result["thresher_growth"] = growth
    
    result["machinery_count"] = int64(float64(currentYearReports) * 0.4)
    result["machinery_growth"] = growth
    
    result["water_pump_count"] = int64(float64(currentYearReports) * 0.6)
    result["water_pump_growth"] = growth
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetEquipmentDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    // Use a different approach - scan into a struct first, then convert to map
    type equipmentRow struct {
        District       string `db:"district"`
        GrainProcessor int64  `db:"grain_processor"`
        Thresher       int64  `db:"thresher"`
        FarmMachinery  int64  `db:"farm_machinery"`
        WaterPump      int64  `db:"water_pump"`
    }
    
    var rows []equipmentRow
    
    query := `
        SELECT 
            district,
            CAST(FLOOR(COUNT(*) * 0.3) AS BIGINT) as grain_processor,
            CAST(FLOOR(COUNT(*) * 0.25) AS BIGINT) as thresher,
            CAST(FLOOR(COUNT(*) * 0.4) AS BIGINT) as farm_machinery,
            CAST(FLOOR(COUNT(*) * 0.6) AS BIGINT) as water_pump
        FROM agriculture_reports
        WHERE visit_date BETWEEN $1 AND $2
        AND (food_technology IS NOT NULL OR horti_technology IS NOT NULL OR plantation_technology IS NOT NULL)
        GROUP BY district
        ORDER BY district
    `
    
    err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&rows).Error
    if err != nil {
        return nil, err
    }
    
    // Convert struct rows to map[string]interface{}
    for _, row := range rows {
        results = append(results, map[string]interface{}{
            "district":        row.District,
            "grain_processor": row.GrainProcessor,
            "thresher":        row.Thresher,
            "farm_machinery":  row.FarmMachinery,
            "water_pump":      row.WaterPump,
        })
    }
    
    return results, err
}

func (r *agricultureRepositoryImpl) GetEquipmentTrend(ctx context.Context, equipmentType string, years []int) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    for _, year := range years {
        startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
        endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
        
        var count int64
        r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
            Where("visit_date BETWEEN ? AND ?", startDate, endDate).
            Where("food_technology IS NOT NULL OR horti_technology IS NOT NULL OR plantation_technology IS NOT NULL").
            Count(&count)
        
        // Adjust count based on equipment type
        var adjustedCount int64
        switch equipmentType {
        case "water_pump":
            adjustedCount = int64(float64(count) * 0.6)
        case "grain_processor":
            adjustedCount = int64(float64(count) * 0.3)
        case "thresher":
            adjustedCount = int64(float64(count) * 0.25)
        default:
            adjustedCount = int64(float64(count) * 0.4)
        }
        
        results = append(results, map[string]interface{}{
            "year":  year,
            "count": adjustedCount,
        })
    }
    
    return results, nil
}

func (r *agricultureRepositoryImpl) GetLandAndIrrigationStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
    result := make(map[string]interface{})
    
    // Total land area current year
    var currentTotalArea float64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", startDate, endDate).
        Select("COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0)").
        Scan(&currentTotalArea)
    
    // Previous year for comparison
    prevYearStart := startDate.AddDate(-1, 0, 0)
    prevYearEnd := endDate.AddDate(-1, 0, 0)
    
    var prevTotalArea float64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", prevYearStart, prevYearEnd).
        Select("COALESCE(SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)), 0)").
        Scan(&prevTotalArea)
    
    // Calculate growth
    var totalGrowth float64 = 0
    if prevTotalArea > 0 {
        totalGrowth = ((currentTotalArea - prevTotalArea) / prevTotalArea) * 100
    }
    
    // Estimate irrigated vs non-irrigated based on water access
    var goodWaterAccess int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", startDate, endDate).
        Where("water_access IN ('MUDAH_TERSEDIA', 'TERSEDIA_BERBAYAR')").
        Count(&goodWaterAccess)
    
    var totalReports int64
    r.db.WithContext(ctx).Model(&entity.AgricultureReport{}).
        Where("visit_date BETWEEN ? AND ?", startDate, endDate).
        Count(&totalReports)
    
    var irrigationRatio float64 = 0.7 // Default 70%
    if totalReports > 0 {
        irrigationRatio = float64(goodWaterAccess) / float64(totalReports)
    }
    
    irrigatedArea := currentTotalArea * irrigationRatio
    nonIrrigatedArea := currentTotalArea * (1 - irrigationRatio)
    
    result["total_land_area"] = currentTotalArea
    result["total_land_growth"] = totalGrowth
    result["irrigated_land_area"] = irrigatedArea
    result["irrigated_land_growth"] = totalGrowth * 1.1 // Slightly higher growth for irrigated
    result["non_irrigated_land_area"] = nonIrrigatedArea
    result["non_irrigated_land_growth"] = totalGrowth * 0.9 // Slightly lower growth for non-irrigated
    
    return result, nil
}

func (r *agricultureRepositoryImpl) GetLandDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    var results []map[string]interface{}
    
    query := `
        SELECT 
            district,
            SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) as total_area,
            SUM(COALESCE(food_land_area, 0) + COALESCE(horti_land_area, 0) + COALESCE(plantation_land_area, 0)) * 0.7 as irrigated_area,
            SUM(COALESCE(food_land_area, 0)) as food_crop_area,
            SUM(COALESCE(horti_land_area, 0)) as horti_area,
            SUM(COALESCE(plantation_land_area, 0)) as plantation_area,
            COUNT(DISTINCT farmer_name) as farmer_count
        FROM agriculture_reports
        WHERE visit_date BETWEEN ? AND ?
        GROUP BY district
        ORDER BY total_area DESC
    `
    
    err := r.db.WithContext(ctx).Raw(query, startDate, endDate).Scan(&results).Error
    return results, err
}