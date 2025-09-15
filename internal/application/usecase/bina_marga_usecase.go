
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
    
    "github.com/google/uuid"
)

type BinaMargaUseCase struct {
    binaMargaRepo repository.BinaMargaRepository
    storage       storage.StorageService
    cache         repository.CacheRepository
}

func NewBinaMargaUseCase(
    binaMargaRepo repository.BinaMargaRepository,
    storage storage.StorageService,
    cache repository.CacheRepository,
) *BinaMargaUseCase {
    return &BinaMargaUseCase{
        binaMargaRepo: binaMargaRepo,
        storage:       storage,
        cache:         cache,
    }
}

func (uc *BinaMargaUseCase) CreateReport(ctx context.Context, req *dto.CreateBinaMargaRequest, photos []*multipart.FileHeader, userID uuid.UUID) (*entity.BinaMargaReport, error) {
    
    damagedArea := req.DamagedLength * req.DamagedWidth
    
    report := &entity.BinaMargaReport{
        ReporterName:     req.ReporterName,
        InstitutionUnit:  entity.InstitutionUnitType(req.InstitutionUnit),
        PhoneNumber:      req.PhoneNumber,
        ReportDateTime:   req.ReportDateTime,
        RoadName:         req.RoadName,
        RoadType:         entity.RoadType(req.RoadType),
        RoadClass:        entity.RoadClass(req.RoadClass),
        Latitude:         req.Latitude,
        Longitude:        req.Longitude,
        DamageType:       entity.RoadDamageType(req.DamageType),
        DamageLevel:      entity.RoadDamageLevel(req.DamageLevel),
        DamagedLength:    req.DamagedLength,
        DamagedWidth:     req.DamagedWidth,
        DamagedArea:      damagedArea,
        TrafficImpact:    entity.TrafficImpact(req.TrafficImpact),
        UrgencyLevel:     entity.RoadUrgencyLevel(req.UrgencyLevel),
        CauseOfDamage:    req.CauseOfDamage,
        Notes:            req.Notes,
        Status:           entity.BinaMargaStatusPending,
        CreatedBy:        userID,
    }

    
    report.EstimatedBudget = uc.calculateEstimatedBudget(report)
    report.EstimatedRepairTime = uc.calculateEstimatedRepairTime(report)

    
    photoAngles := []string{"before", "damage_detail", "traffic_impact", "aerial", "surrounding"}
    for i, photo := range photos {
        angle := "general"
        if i < len(photoAngles) {
            angle = photoAngles[i]
        }

        photoURL, err := uc.storage.UploadFile(ctx, photo, "bina-marga")
        if err != nil {
            return nil, fmt.Errorf("failed to upload photo: %w", err)
        }

        caption := fmt.Sprintf("%s view - %s (%s)", angle, report.RoadName, report.DamageType)
        report.Photos = append(report.Photos, entity.BinaMargaPhoto{
            PhotoURL:   photoURL,
            PhotoAngle: angle,
            Caption:    caption,
        })
    }

    if err := uc.binaMargaRepo.Create(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "bina_marga:list")
    uc.cache.Delete(ctx, "bina_marga:stats")
    uc.cache.Delete(ctx, "bina_marga:emergency")

    
    if report.UrgencyLevel == entity.RoadUrgencyEmergency || report.TrafficImpact == entity.TrafficImpactBlocked {
        uc.sendUrgentNotification(ctx, report)
    }

    return report, nil
}

func (uc *BinaMargaUseCase) GetReport(ctx context.Context, id uuid.UUID) (*entity.BinaMargaReport, error) {
    cacheKey := "bina_marga:" + id.String()
    
    report, err := uc.binaMargaRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, report, 3600)

    return report, nil
}

func (uc *BinaMargaUseCase) ListReports(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.PaginatedBinaMargaResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.binaMargaRepo.FindAll(ctx, limit, offset, filters)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedBinaMargaResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *BinaMargaUseCase) ListByPriority(ctx context.Context, page, limit int) (*dto.PaginatedBinaMargaResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.binaMargaRepo.FindByPriority(ctx, limit, offset)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedBinaMargaResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *BinaMargaUseCase) UpdateReport(ctx context.Context, id uuid.UUID, req *dto.UpdateBinaMargaRequest, userID uuid.UUID) (*entity.BinaMargaReport, error) {
    report, err := uc.binaMargaRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    if report.CreatedBy != userID {
        return nil, ErrUnauthorized
    }

    
    if req.RoadName != "" {
        report.RoadName = req.RoadName
    }
    if req.RoadType != "" {
        report.RoadType = entity.RoadType(req.RoadType)
    }
    if req.RoadClass != "" {
        report.RoadClass = entity.RoadClass(req.RoadClass)
    }
    if req.DamageType != "" {
        report.DamageType = entity.RoadDamageType(req.DamageType)
    }
    if req.DamageLevel != "" {
        report.DamageLevel = entity.RoadDamageLevel(req.DamageLevel)
    }
    if req.DamagedLength > 0 {
        report.DamagedLength = req.DamagedLength
    }
    if req.DamagedWidth > 0 {
        report.DamagedWidth = req.DamagedWidth
    }
    if req.TrafficImpact != "" {
        report.TrafficImpact = entity.TrafficImpact(req.TrafficImpact)
    }
    if req.UrgencyLevel != "" {
        report.UrgencyLevel = entity.RoadUrgencyLevel(req.UrgencyLevel)
    }
    if req.CauseOfDamage != "" {
        report.CauseOfDamage = req.CauseOfDamage
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
    if req.EstimatedRepairTime > 0 {
        report.EstimatedRepairTime = req.EstimatedRepairTime
    }

    
    if req.DamagedLength > 0 || req.DamagedWidth > 0 {
        report.DamagedArea = report.DamagedLength * report.DamagedWidth
    }

    
    if req.EstimatedBudget == 0 {
        report.EstimatedBudget = uc.calculateEstimatedBudget(report)
    }
    if req.EstimatedRepairTime == 0 {
        report.EstimatedRepairTime = uc.calculateEstimatedRepairTime(report)
    }

    if err := uc.binaMargaRepo.Update(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "bina_marga:"+id.String())
    uc.cache.Delete(ctx, "bina_marga:list")
    uc.cache.Delete(ctx, "bina_marga:stats")

    return report, nil
}

func (uc *BinaMargaUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateBinaMargaStatusRequest) error {
    err := uc.binaMargaRepo.UpdateStatus(ctx, id, entity.BinaMargaStatus(req.Status), req.Notes)
    if err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "bina_marga:"+id.String())
    uc.cache.Delete(ctx, "bina_marga:stats")

    return nil
}

func (uc *BinaMargaUseCase) DeleteReport(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
    report, err := uc.binaMargaRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    
    if report.CreatedBy != userID {
        return ErrUnauthorized
    }

    
    for _, photo := range report.Photos {
        uc.storage.DeleteFile(ctx, photo.PhotoURL)
    }

    if err := uc.binaMargaRepo.Delete(ctx, id); err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "bina_marga:"+id.String())
    uc.cache.Delete(ctx, "bina_marga:list")
    uc.cache.Delete(ctx, "bina_marga:stats")

    return nil
}

func (uc *BinaMargaUseCase) GetStatistics(ctx context.Context) (*dto.BinaMargaStatisticsResponse, error) {
    
    cacheKey := "bina_marga:stats"
    var stats dto.BinaMargaStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &stats)
    if err == nil {
        return &stats, nil
    }

    
    rawStats, err := uc.binaMargaRepo.GetStatistics(ctx)
    if err != nil {
        return nil, err
    }

    
    response := &dto.BinaMargaStatisticsResponse{
        TotalReports:         rawStats["total_reports"].(int64),
        EmergencyReports:     rawStats["emergency_reports"].(int64),
        BlockedRoads:         rawStats["blocked_roads"].(int64),
        TotalDamagedArea:     rawStats["total_damaged_area_sqm"].(float64),
        TotalDamagedLength:   rawStats["total_damaged_length_m"].(float64),
        EstimatedTotalBudget: rawStats["estimated_total_budget"].(float64),
        AverageRepairTime:    rawStats["average_repair_time_days"].(float64),
    }

    
    if roadTypes, ok := rawStats["road_type_distribution"].([]interface{}); ok {
        for _, v := range roadTypes {
            if m, ok := v.(map[string]interface{}); ok {
                response.RoadTypeDistribution = append(response.RoadTypeDistribution, m)
            }
        }
    }

    if damageTypes, ok := rawStats["damage_type_distribution"].([]interface{}); ok {
        for _, v := range damageTypes {
            if m, ok := v.(map[string]interface{}); ok {
                response.DamageTypeDistribution = append(response.DamageTypeDistribution, m)
            }
        }
    }

    if damageLevels, ok := rawStats["damage_level_counts"].([]interface{}); ok {
        for _, v := range damageLevels {
            if m, ok := v.(map[string]interface{}); ok {
                response.DamageLevelCounts = append(response.DamageLevelCounts, m)
            }
        }
    }

    if urgencyLevels, ok := rawStats["urgency_level_counts"].([]interface{}); ok {
        for _, v := range urgencyLevels {
            if m, ok := v.(map[string]interface{}); ok {
                response.UrgencyLevelCounts = append(response.UrgencyLevelCounts, m)
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

    if trafficImpacts, ok := rawStats["traffic_impact_counts"].([]interface{}); ok {
        for _, v := range trafficImpacts {
            if m, ok := v.(map[string]interface{}); ok {
                response.TrafficImpactCounts = append(response.TrafficImpactCounts, m)
            }
        }
    }

    
    uc.cache.Set(ctx, cacheKey, response, 300)

    return response, nil
}

func (uc *BinaMargaUseCase) GetEmergencyReports(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error) {
    cacheKey := fmt.Sprintf("bina_marga:emergency:%d", limit)
    
    reports, err := uc.binaMargaRepo.FindEmergencyReports(ctx, limit)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, reports, 120)

    return reports, nil
}

func (uc *BinaMargaUseCase) GetBlockedRoads(ctx context.Context, limit int) ([]*entity.BinaMargaReport, error) {
    cacheKey := fmt.Sprintf("bina_marga:blocked:%d", limit)
    
    reports, err := uc.binaMargaRepo.FindBlockedRoads(ctx, limit)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, reports, 120)

    return reports, nil
}

func (uc *BinaMargaUseCase) GetDamageByRoadType(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    return uc.binaMargaRepo.GetDamageStatisticsByRoadType(ctx, startDate, endDate)
}

func (uc *BinaMargaUseCase) GetDamageByLocation(ctx context.Context, bounds map[string]float64) ([]map[string]interface{}, error) {
    return uc.binaMargaRepo.GetDamageStatisticsByLocation(ctx, bounds)
}



func (uc *BinaMargaUseCase) calculateEstimatedBudget(report *entity.BinaMargaReport) float64 {
    baseCost := 2000000.0 
    
    
    areaCost := report.DamagedArea * baseCost
    
    
    levelMultiplier := 1.0
    switch report.DamageLevel {
    case entity.RoadDamageLevelModerate:
        levelMultiplier = 1.5
    case entity.RoadDamageLevelSevere:
        levelMultiplier = 2.5
    }
    
    
    classMultiplier := 1.0
    switch report.RoadClass {
    case entity.RoadClassArteri:
        classMultiplier = 2.0
    case entity.RoadClassKolektor:
        classMultiplier = 1.5
    case entity.RoadClassLokal:
        classMultiplier = 1.2
    }
    
    
    typeMultiplier := 1.0
    switch report.DamageType {
    case entity.RoadDamageJembatan:
        typeMultiplier = 3.0
    case entity.RoadDamageAmblas:
        typeMultiplier = 2.5
    case entity.RoadDamageLubang:
        typeMultiplier = 2.0
    case entity.RoadDamageRetakBuaya:
        typeMultiplier = 1.8
    case entity.RoadDamageDrainase:
        typeMultiplier = 1.5
    }
    
    
    urgencyAdditional := 0.0
    if report.UrgencyLevel == entity.RoadUrgencyEmergency {
        urgencyAdditional = areaCost * 0.5 
    } else if report.UrgencyLevel == entity.RoadUrgencyHigh {
        urgencyAdditional = areaCost * 0.2 
    }
    
    totalBudget := areaCost * levelMultiplier * classMultiplier * typeMultiplier + urgencyAdditional
    
    return totalBudget
}

func (uc *BinaMargaUseCase) calculateEstimatedRepairTime(report *entity.BinaMargaReport) int {
    baseTimePerSqm := 0.1 
    
    
    baseTime := report.DamagedArea * baseTimePerSqm
    
    
    levelMultiplier := 1.0
    switch report.DamageLevel {
    case entity.RoadDamageLevelModerate:
        levelMultiplier = 1.5
    case entity.RoadDamageLevelSevere:
        levelMultiplier = 2.5
    }
    
    
    typeMultiplier := 1.0
    switch report.DamageType {
    case entity.RoadDamageJembatan:
        typeMultiplier = 4.0
    case entity.RoadDamageAmblas:
        typeMultiplier = 3.0
    case entity.RoadDamageDrainase:
        typeMultiplier = 2.0
    case entity.RoadDamageRetakBuaya:
        typeMultiplier = 1.8
    case entity.RoadDamageLubang:
        typeMultiplier = 1.5
    }
    
    
    classMultiplier := 1.0
    switch report.RoadClass {
    case entity.RoadClassArteri:
        classMultiplier = 1.5
    case entity.RoadClassKolektor:
        classMultiplier = 1.2
    }
    
    totalTime := baseTime * levelMultiplier * typeMultiplier * classMultiplier
    
    
    if totalTime < 1 {
        totalTime = 1
    } else if totalTime > 365 {
        totalTime = 365
    }
    
    return int(totalTime)
}

func (uc *BinaMargaUseCase) sendUrgentNotification(ctx context.Context, report *entity.BinaMargaReport) {
    
    
    fmt.Printf("URGENT: Road damage report - %s on %s (%s impact)\n", 
        report.DamageType, report.RoadName, report.TrafficImpact)
}