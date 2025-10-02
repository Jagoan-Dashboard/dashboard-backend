package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/storage"
	"building-report-backend/pkg/utils"

	
)

type WaterResourcesUseCase struct {
    waterRepo repository.WaterResourcesRepository
    storage   storage.StorageService
    cache     repository.CacheRepository
}

func NewWaterResourcesUseCase(
    waterRepo repository.WaterResourcesRepository,
    storage storage.StorageService,
    cache repository.CacheRepository,
) *WaterResourcesUseCase {
    return &WaterResourcesUseCase{
        waterRepo: waterRepo,
        storage:   storage,
        cache:     cache,
    }
}

func (uc *WaterResourcesUseCase) CreateReport(ctx context.Context, req *dto.CreateWaterResourcesRequest, photos []*multipart.FileHeader) (*entity.WaterResourcesReport, error) {
    report := &entity.WaterResourcesReport{
        ID:                    utils.GenerateULID(),
        ReporterName:          req.ReporterName,
        InstitutionUnit:       entity.InstitutionUnitType(req.InstitutionUnit),
        PhoneNumber:           req.PhoneNumber,
        ReportDateTime:        req.ReportDateTime,
        IrrigationAreaName:    req.IrrigationAreaName,
        IrrigationType:        entity.IrrigationType(req.IrrigationType),
        Latitude:              req.Latitude,
        Longitude:             req.Longitude,
        DamageType:            entity.DamageType(req.DamageType),
        DamageLevel:           entity.DamageLevel(req.DamageLevel),
        EstimatedLength:       req.EstimatedLength,
        EstimatedWidth:        req.EstimatedWidth,
        EstimatedVolume:       req.EstimatedVolume,
        AffectedRiceFieldArea: req.AffectedRiceFieldArea,
        AffectedFarmersCount:  req.AffectedFarmersCount,
        UrgencyCategory:       entity.UrgencyCategory(req.UrgencyCategory),
        Notes:                 req.Notes,
        Status:                entity.WaterResourceStatusPending,
    }

    
    report.EstimatedBudget = uc.calculateEstimatedBudget(report)

    
    photoAngles := []string{"front", "side", "damage_detail", "aerial"}
    for i, photo := range photos {
        if i >= len(photoAngles) {
            break
        }

        photoURL, err := uc.storage.UploadFile(ctx, photo, "water-resources")
        if err != nil {
            return nil, fmt.Errorf("failed to upload photo: %w", err)
        }

        caption := fmt.Sprintf("%s view - %s", photoAngles[i], report.IrrigationAreaName)
        report.Photos = append(report.Photos, entity.WaterResourcesPhoto{
            ID:         utils.GenerateULID(),
            PhotoURL:   photoURL,
            PhotoAngle: photoAngles[i],
            Caption:    caption,
        })
    }

    if err := uc.waterRepo.Create(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "water:list")
    uc.cache.Delete(ctx, "water:stats")
    uc.cache.Delete(ctx, "water:urgent")

    
    if report.UrgencyCategory == entity.UrgencyCategoryMendesak {
        uc.sendUrgentNotification(ctx, report)
    }

    return report, nil
}

func (uc *WaterResourcesUseCase) GetReport(ctx context.Context, id string) (*entity.WaterResourcesReport, error) {
    cacheKey := "water:" + id
    
    report, err := uc.waterRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, report, 3600)

    return report, nil
}

func (uc *WaterResourcesUseCase) ListReports(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.PaginatedWaterResourcesResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.waterRepo.FindAll(ctx, limit, offset, filters)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedWaterResourcesResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *WaterResourcesUseCase) ListByPriority(ctx context.Context, page, limit int) (*dto.PaginatedWaterResourcesResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.waterRepo.FindByPriority(ctx, limit, offset)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedWaterResourcesResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *WaterResourcesUseCase) UpdateReport(ctx context.Context, id string, req *dto.UpdateWaterResourcesRequest, userID string) (*entity.WaterResourcesReport, error) {
    report, err := uc.waterRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    if report.CreatedBy != userID {
        return nil, ErrUnauthorized
    }

    
    if req.IrrigationAreaName != "" {
        report.IrrigationAreaName = req.IrrigationAreaName
    }
    if req.IrrigationType != "" {
        report.IrrigationType = entity.IrrigationType(req.IrrigationType)
    }
    if req.DamageType != "" {
        report.DamageType = entity.DamageType(req.DamageType)
    }
    if req.DamageLevel != "" {
        report.DamageLevel = entity.DamageLevel(req.DamageLevel)
    }
    if req.EstimatedLength > 0 {
        report.EstimatedLength = req.EstimatedLength
    }
    if req.EstimatedWidth > 0 {
        report.EstimatedWidth = req.EstimatedWidth
    }
    if req.EstimatedVolume > 0 {
        report.EstimatedVolume = req.EstimatedVolume
    }
    if req.AffectedRiceFieldArea > 0 {
        report.AffectedRiceFieldArea = req.AffectedRiceFieldArea
    }
    if req.AffectedFarmersCount > 0 {
        report.AffectedFarmersCount = req.AffectedFarmersCount
    }
    if req.UrgencyCategory != "" {
        report.UrgencyCategory = entity.UrgencyCategory(req.UrgencyCategory)
    }
    if req.Notes != "" {
        report.Notes = req.Notes
    }
    if req.HandlingRecommendation != "" {
        report.HandlingRecommendation = req.HandlingRecommendation
    }
    if req.EstimatedBudget > 0 {
        report.EstimatedBudget = req.EstimatedBudget
    }

    
    if req.EstimatedBudget == 0 {
        report.EstimatedBudget = uc.calculateEstimatedBudget(report)
    }

    if err := uc.waterRepo.Update(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "water:"+id)
    uc.cache.Delete(ctx, "water:list")
    uc.cache.Delete(ctx, "water:stats")

    return report, nil
}

func (uc *WaterResourcesUseCase) UpdateStatus(ctx context.Context, id string, req *dto.UpdateWaterStatusRequest) error {
    err := uc.waterRepo.UpdateStatus(ctx, id, entity.WaterResourceStatus(req.Status), req.Notes)
    if err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "water:"+id)
    uc.cache.Delete(ctx, "water:stats")

    return nil
}

func (uc *WaterResourcesUseCase) DeleteReport(ctx context.Context, id string, userID string) error {
    report, err := uc.waterRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    
    if report.CreatedBy != userID {
        return ErrUnauthorized
    }

    
    for _, photo := range report.Photos {
        uc.storage.DeleteFile(ctx, photo.PhotoURL)
    }

    if err := uc.waterRepo.Delete(ctx, id); err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "water:"+id)
    uc.cache.Delete(ctx, "water:list")
    uc.cache.Delete(ctx, "water:stats")

    return nil
}

func (uc *WaterResourcesUseCase) GetStatistics(ctx context.Context) (*dto.WaterResourcesStatisticsResponse, error) {
    
    cacheKey := "water:stats"
    var stats dto.WaterResourcesStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &stats)
    if err == nil {
        return &stats, nil
    }

    
    rawStats, err := uc.waterRepo.GetStatistics(ctx)
    if err != nil {
        return nil, err
    }

    
    response := &dto.WaterResourcesStatisticsResponse{
        TotalReports:         rawStats["total_reports"].(int64),
        UrgentPending:        rawStats["urgent_pending"].(int64),
        TotalAffectedAreaHa:  rawStats["total_affected_area_ha"].(float64),
        TotalAffectedFarmers: rawStats["total_affected_farmers"].(int64),
        EstimatedTotalBudget: rawStats["estimated_total_budget"].(float64),
    }

    
    if damageTypes, ok := rawStats["damage_types"].([]interface{}); ok {
        for _, v := range damageTypes {
            if m, ok := v.(map[string]interface{}); ok {
                response.DamageTypes = append(response.DamageTypes, m)
            }
        }
    }

    if irrigationTypes, ok := rawStats["irrigation_types"].([]interface{}); ok {
        for _, v := range irrigationTypes {
            if m, ok := v.(map[string]interface{}); ok {
                response.IrrigationTypes = append(response.IrrigationTypes, m)
            }
        }
    }

    if statusDist, ok := rawStats["status_distribution"].([]interface{}); ok {
        for _, v := range statusDist {
            if m, ok := v.(map[string]interface{}); ok {
                response.StatusDistribution = append(response.StatusDistribution, m)
            }
        }
    }

    
    uc.cache.Set(ctx, cacheKey, response, 300)

    return response, nil
}

func (uc *WaterResourcesUseCase) GetUrgentReports(ctx context.Context, limit int) ([]*entity.WaterResourcesReport, error) {
    
    cacheKey := fmt.Sprintf("water:urgent:%d", limit)
    
    reports, err := uc.waterRepo.GetUrgentReports(ctx, limit)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, reports, 120)

    return reports, nil
}

func (uc *WaterResourcesUseCase) GetDamageByArea(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    return uc.waterRepo.GetDamageStatisticsByArea(ctx, startDate, endDate)
}


func (uc *WaterResourcesUseCase) calculateEstimatedBudget(report *entity.WaterResourcesReport) float64 {
    baseCost := 1000000.0 
    
    
    damageArea := report.EstimatedLength * report.EstimatedWidth
    areaCost := damageArea * 500000.0 
    
    
    levelMultiplier := 1.0
    switch report.DamageLevel {
    case entity.DamageLevelSedang:
        levelMultiplier = 1.5
    case entity.DamageLevelBerat:
        levelMultiplier = 2.5
    }
    
    
    typeMultiplier := 1.0
    switch report.IrrigationType {
    case entity.IrrigationBendung:
        typeMultiplier = 2.0
    case entity.IrrigationEmbungDam:
        typeMultiplier = 2.5
    case entity.IrrigationPintuAir:
        typeMultiplier = 1.8
    }
    
    
    urgencyAdditional := 0.0
    if report.UrgencyCategory == entity.UrgencyCategoryMendesak {
        urgencyAdditional = baseCost * 0.3
    }
    
    totalBudget := (baseCost + areaCost) * levelMultiplier * typeMultiplier + urgencyAdditional
    
    return totalBudget
}

func (uc *WaterResourcesUseCase) sendUrgentNotification(ctx context.Context, report *entity.WaterResourcesReport) {
    
    
    
    fmt.Printf("URGENT: New water resource damage report at %s affecting %d farmers\n", 
        report.IrrigationAreaName, report.AffectedFarmersCount)
}

func (uc *WaterResourcesUseCase) GetDashboard(
    ctx context.Context,
    irrigationType string,
    startDate, endDate time.Time,
) (*dto.WaterResourcesDashboardResponse, error) {

    // KPI
    totalArea, totalRice, totalReports, err := uc.waterRepo.GetSummaryKPIs(ctx, irrigationType, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // Distribusi (urgency, damage type, level)
    urgRows, err := uc.waterRepo.GroupCountBy(ctx, "urgency_category", irrigationType, startDate, endDate)
    if err != nil {
        return nil, err
    }
    dmgTypeRows, err := uc.waterRepo.GroupCountBy(ctx, "damage_type", irrigationType, startDate, endDate)
    if err != nil {
        return nil, err
    }
    dmgLevelRows, err := uc.waterRepo.GroupCountBy(ctx, "damage_level", irrigationType, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // Map points
    pts, err := uc.waterRepo.GetMapPoints(ctx, irrigationType, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // Build DTO
    res := &dto.WaterResourcesDashboardResponse{}
    res.KPIs.TotalDamageAreaM2 = totalArea
    res.KPIs.TotalRiceFieldHa = totalRice
    res.KPIs.TotalReports = totalReports

    res.UrgencyDistribution = make([]dto.KeyCount, len(urgRows))
    for i, r0 := range urgRows {
        res.UrgencyDistribution[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.TopDamageTypes = make([]dto.KeyCount, len(dmgTypeRows))
    for i, r0 := range dmgTypeRows {
        res.TopDamageTypes[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.TopDamageLevels = make([]dto.KeyCount, len(dmgLevelRows))
    for i, r0 := range dmgLevelRows {
        res.TopDamageLevels[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.MapPoints = make([]dto.DashboardMapPoint, len(pts))
    for i, p := range pts {
        res.MapPoints[i] = dto.DashboardMapPoint{
            Latitude:        p.Latitude,
            Longitude:       p.Longitude,
            IrrigationArea:  p.IrrigationArea,
            DamageType:      p.DamageType,
            DamageLevel:     p.DamageLevel,
            UrgencyCategory: p.UrgencyCategory,
        }
    }

    return res, nil
}

func (uc *WaterResourcesUseCase) GetWaterResourcesOverview(ctx context.Context, irrigationType string) (*dto.WaterResourcesOverviewResponse, error) {
    // Cache key based on irrigation type
    cacheKey := fmt.Sprintf("water_resources:overview:%s", irrigationType)
    var response dto.WaterResourcesOverviewResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    // Initialize empty arrays to avoid null returns
    response.LocationDistribution = []dto.WaterLocationStatsResponse{}
    response.UrgencyDistribution = []dto.WaterUrgencyStatsResponse{}
    response.DamageTypeDistribution = []dto.WaterDamageTypeStatsResponse{}
    response.DamageLevelDistribution = []dto.WaterDamageLevelStatsResponse{}

    // Get basic statistics
    basicStatsRaw, err := uc.waterRepo.GetWaterResourcesOverviewStats(ctx, irrigationType)
    if err != nil {
        return nil, fmt.Errorf("failed to get basic statistics: %w", err)
    }
    
    // Safe type assertions with proper conversion
    response.BasicStats.TotalDamageVolumeM2 = basicStatsRaw["total_damage_volume_m2"].(float64)
    response.BasicStats.TotalRiceFieldAreaHa = basicStatsRaw["total_rice_field_area_ha"].(float64)
    
    // Handle int64 conversion safely
    if totalReports, ok := basicStatsRaw["total_damaged_reports"].(int64); ok {
        response.BasicStats.TotalDamagedReports = totalReports
    } else if totalReportsFloat, ok := basicStatsRaw["total_damaged_reports"].(float64); ok {
        response.BasicStats.TotalDamagedReports = int64(totalReportsFloat)
    }

    // Get location distribution
    locationStats, err := uc.waterRepo.GetWaterLocationStats(ctx, irrigationType)
    if err != nil {
        return nil, fmt.Errorf("failed to get location statistics: %w", err)
    }
    
    for _, loc := range locationStats {
        locationStat := dto.WaterLocationStatsResponse{
            IrrigationAreaName: loc["irrigation_area_name"].(string),
            AvgLatitude:        loc["avg_latitude"].(float64),
            AvgLongitude:       loc["avg_longitude"].(float64),
            TotalAffectedArea:  loc["total_affected_area"].(float64),
        }
        
        // Safe conversion for integer fields
        if reportCount, ok := loc["report_count"].(int64); ok {
            locationStat.ReportCount = int(reportCount)
        } else if reportCountFloat, ok := loc["report_count"].(float64); ok {
            locationStat.ReportCount = int(reportCountFloat)
        }
        
        if farmersCount, ok := loc["total_affected_farmers"].(int64); ok {
            locationStat.TotalAffectedFarmers = int(farmersCount)
        } else if farmersCountFloat, ok := loc["total_affected_farmers"].(float64); ok {
            locationStat.TotalAffectedFarmers = int(farmersCountFloat)
        }
        
        response.LocationDistribution = append(response.LocationDistribution, locationStat)
    }

    // Get urgency distribution
    urgencyStats, err := uc.waterRepo.GetWaterUrgencyStats(ctx, irrigationType)
    if err != nil {
        return nil, fmt.Errorf("failed to get urgency statistics: %w", err)
    }
    
    for _, urgency := range urgencyStats {
        urgencyStat := dto.WaterUrgencyStatsResponse{
            UrgencyCategory: urgency["urgency_category"].(string),
        }
        
        // Safe conversion for count
        if count, ok := urgency["count"].(int64); ok {
            urgencyStat.Count = count
        } else if countFloat, ok := urgency["count"].(float64); ok {
            urgencyStat.Count = int64(countFloat)
        }
        
        response.UrgencyDistribution = append(response.UrgencyDistribution, urgencyStat)
    }

    // Get damage type distribution
    damageTypeStats, err := uc.waterRepo.GetWaterDamageTypeStats(ctx, irrigationType)
    if err != nil {
        return nil, fmt.Errorf("failed to get damage type statistics: %w", err)
    }
    
    for _, damageType := range damageTypeStats {
        damageTypeStat := dto.WaterDamageTypeStatsResponse{
            DamageType: damageType["damage_type"].(string),
        }
        
        // Safe conversion for count
        if count, ok := damageType["count"].(int64); ok {
            damageTypeStat.Count = count
        } else if countFloat, ok := damageType["count"].(float64); ok {
            damageTypeStat.Count = int64(countFloat)
        }
        
        response.DamageTypeDistribution = append(response.DamageTypeDistribution, damageTypeStat)
    }

    // Get damage level distribution
    damageLevelStats, err := uc.waterRepo.GetWaterDamageLevelStats(ctx, irrigationType)
    if err != nil {
        return nil, fmt.Errorf("failed to get damage level statistics: %w", err)
    }
    
    for _, damageLevel := range damageLevelStats {
        damageLevelStat := dto.WaterDamageLevelStatsResponse{
            DamageLevel: damageLevel["damage_level"].(string),
        }
        
        // Safe conversion for count
        if count, ok := damageLevel["count"].(int64); ok {
            damageLevelStat.Count = count
        } else if countFloat, ok := damageLevel["count"].(float64); ok {
            damageLevelStat.Count = int64(countFloat)
        }
        
        response.DamageLevelDistribution = append(response.DamageLevelDistribution, damageLevelStat)
    }

    // Cache the response for 5 minutes
    uc.cache.Set(ctx, cacheKey, &response, 300*time.Second)

    return &response, nil
}

// Alternative: Generic helper function untuk safe type conversion
func safeInt64(value interface{}) int64 {
    switch v := value.(type) {
    case int64:
        return v
    case float64:
        return int64(v)
    case int:
        return int64(v)
    case int32:
        return int64(v)
    default:
        return 0
    }
}

func safeInt(value interface{}) int {
    return int(safeInt64(value))
}

func safeFloat64(value interface{}) float64 {
    switch v := value.(type) {
    case float64:
        return v
    case float32:
        return float64(v)
    case int64:
        return float64(v)
    case int:
        return float64(v)
    default:
        return 0.0
    }
}