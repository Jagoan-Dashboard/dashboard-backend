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

func (uc *BinaMargaUseCase) CreateReport(ctx context.Context, req *dto.CreateBinaMargaRequest, photos []*multipart.FileHeader) (*entity.BinaMargaReport, error) {
    
    damagedArea := req.DamagedLength * req.DamagedWidth
    
    
    totalDamagedArea := req.TotalDamagedArea
    if totalDamagedArea == 0 && damagedArea > 0 {
        totalDamagedArea = damagedArea
    }
    
    report := &entity.BinaMargaReport{
        ID:                  utils.GenerateULID(),
        ReporterName:        req.ReporterName,
        InstitutionUnit:     entity.InstitutionUnitType(req.InstitutionUnit),
        PhoneNumber:         req.PhoneNumber,
        ReportDateTime:      req.ReportDateTime,
        RoadName:            req.RoadName,
        RoadType:            entity.RoadType(req.RoadType),
        RoadClass:           entity.RoadClass(req.RoadClass),
        SegmentLength:       req.SegmentLength,
        Latitude:            req.Latitude,
        Longitude:           req.Longitude,
        PavementType:        entity.PavementType(req.PavementType),
        DamageType:          entity.RoadDamageType(req.DamageType),
        DamageLevel:         entity.RoadDamageLevel(req.DamageLevel),
        DamagedLength:       req.DamagedLength,
        DamagedWidth:        req.DamagedWidth,
        DamagedArea:         damagedArea,
        TotalDamagedArea:    totalDamagedArea,
        TrafficCondition:    entity.TrafficCondition(req.TrafficCondition),
        TrafficImpact:       entity.TrafficImpact(req.TrafficImpact),
        DailyTrafficVolume:  req.DailyTrafficVolume,
        UrgencyLevel:        entity.RoadUrgencyLevel(req.UrgencyLevel),
        CauseOfDamage:       req.CauseOfDamage,
        Notes:               req.Notes,
        Status:              entity.BinaMargaStatusPending,
    }

    
    if req.BridgeName != "" {
        report.BridgeName = req.BridgeName
        if req.BridgeStructureType != "" {
            report.BridgeStructureType = entity.BridgeStructureType(req.BridgeStructureType)
        }
        if req.BridgeDamageType != "" {
            report.BridgeDamageType = entity.BridgeDamageType(req.BridgeDamageType)
        }
        if req.BridgeDamageLevel != "" {
            report.BridgeDamageLevel = entity.BridgeDamageLevel(req.BridgeDamageLevel)
        }
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
        if report.BridgeName != "" {
            caption = fmt.Sprintf("%s view - Bridge %s (%s)", angle, report.BridgeName, report.BridgeDamageType)
        }
        
        report.Photos = append(report.Photos, entity.BinaMargaPhoto{
            ID:         utils.GenerateULID(),
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

    
    if report.UrgencyLevel == entity.RoadUrgencyEmergency || 
       report.TrafficImpact == entity.TrafficImpactBlocked ||
       report.TrafficCondition == entity.TrafficConditionBlocked {
        uc.sendUrgentNotification(ctx, report)
    }

    return report, nil
}

func (uc *BinaMargaUseCase) GetReport(ctx context.Context, id string) (*entity.BinaMargaReport, error) {
    cacheKey := "bina_marga:" + id

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

func (uc *BinaMargaUseCase) UpdateReport(ctx context.Context, id string, req *dto.UpdateBinaMargaRequest, userID string) (*entity.BinaMargaReport, error) {
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
    if req.SegmentLength > 0 {
        report.SegmentLength = req.SegmentLength
    }
    
    
    if req.PavementType != "" {
        report.PavementType = entity.PavementType(req.PavementType)
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
    if req.TotalDamagedArea > 0 {
        report.TotalDamagedArea = req.TotalDamagedArea
    }
    
    
    if req.BridgeName != "" {
        report.BridgeName = req.BridgeName
    }
    if req.BridgeStructureType != "" {
        report.BridgeStructureType = entity.BridgeStructureType(req.BridgeStructureType)
    }
    if req.BridgeDamageType != "" {
        report.BridgeDamageType = entity.BridgeDamageType(req.BridgeDamageType)
    }
    if req.BridgeDamageLevel != "" {
        report.BridgeDamageLevel = entity.BridgeDamageLevel(req.BridgeDamageLevel)
    }
    
    
    if req.TrafficCondition != "" {
        report.TrafficCondition = entity.TrafficCondition(req.TrafficCondition)
    }
    if req.TrafficImpact != "" {
        report.TrafficImpact = entity.TrafficImpact(req.TrafficImpact)
    }
    if req.DailyTrafficVolume > 0 {
        report.DailyTrafficVolume = req.DailyTrafficVolume
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
        if report.TotalDamagedArea == 0 {
            report.TotalDamagedArea = report.DamagedArea
        }
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

    
    uc.cache.Delete(ctx, "bina_marga:"+id)
    uc.cache.Delete(ctx, "bina_marga:list")
    uc.cache.Delete(ctx, "bina_marga:stats")

    return report, nil
}

func (uc *BinaMargaUseCase) UpdateStatus(ctx context.Context, id string, req *dto.UpdateBinaMargaStatusRequest) error {
    err := uc.binaMargaRepo.UpdateStatus(ctx, id, entity.BinaMargaStatus(req.Status), req.Notes)
    if err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "bina_marga:"+id)
    uc.cache.Delete(ctx, "bina_marga:stats")

    return nil
}

func (uc *BinaMargaUseCase) DeleteReport(ctx context.Context, id string, userID string) error {
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

    
    uc.cache.Delete(ctx, "bina_marga:"+id)
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

    
    if bridgeCount, ok := rawStats["bridge_reports"].(int64); ok {
        response.BridgeReports = bridgeCount
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
    
    
    damagedArea := report.TotalDamagedArea
    if damagedArea == 0 {
        damagedArea = report.DamagedArea
    }
    
    areaCost := damagedArea * baseCost
    
    
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
    
    
    pavementMultiplier := 1.0
    switch report.PavementType {
    case entity.PavementBetonRigid:
        pavementMultiplier = 1.8
    case entity.PavementAspalFlexible:
        pavementMultiplier = 1.0
    case entity.PavementPaving:
        pavementMultiplier = 1.3
    case entity.PavementJalanTanah:
        pavementMultiplier = 0.5
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
    case entity.RoadDamageGenaganDrainase:
        typeMultiplier = 1.7
    }
    
    
    bridgeAdditional := 0.0
    if report.BridgeName != "" {
        bridgeAdditional = areaCost * 2.0 
        if report.BridgeDamageLevel == entity.BridgeDamageLevelSevere {
            bridgeAdditional *= 3.0
        }
    }
    
    
    urgencyAdditional := 0.0
    if report.UrgencyLevel == entity.RoadUrgencyEmergency {
        urgencyAdditional = areaCost * 0.5 
    } else if report.UrgencyLevel == entity.RoadUrgencyHigh {
        urgencyAdditional = areaCost * 0.2 
    }
    
    totalBudget := (areaCost * levelMultiplier * classMultiplier * pavementMultiplier * typeMultiplier) + bridgeAdditional + urgencyAdditional
    
    return totalBudget
}

func (uc *BinaMargaUseCase) calculateEstimatedRepairTime(report *entity.BinaMargaReport) int {
    baseTimePerSqm := 0.1 
    
    
    damagedArea := report.TotalDamagedArea
    if damagedArea == 0 {
        damagedArea = report.DamagedArea
    }
    
    baseTime := damagedArea * baseTimePerSqm
    
    
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
    case entity.RoadDamageDrainase, entity.RoadDamageGenaganDrainase:
        typeMultiplier = 2.0
    case entity.RoadDamageRetakBuaya:
        typeMultiplier = 1.8
    case entity.RoadDamageLubang:
        typeMultiplier = 1.5
    }
    
    
    pavementMultiplier := 1.0
    switch report.PavementType {
    case entity.PavementBetonRigid:
        pavementMultiplier = 2.0
    case entity.PavementAspalFlexible:
        pavementMultiplier = 1.0
    case entity.PavementPaving:
        pavementMultiplier = 1.2
    case entity.PavementJalanTanah:
        pavementMultiplier = 0.3
    }
    
    
    classMultiplier := 1.0
    switch report.RoadClass {
    case entity.RoadClassArteri:
        classMultiplier = 1.5
    case entity.RoadClassKolektor:
        classMultiplier = 1.2
    }
    
    
    bridgeAdditional := 0.0
    if report.BridgeName != "" {
        bridgeAdditional = baseTime * 1.5
        if report.BridgeDamageLevel == entity.BridgeDamageLevelSevere {
            bridgeAdditional *= 2.0
        }
    }
    
    totalTime := (baseTime * levelMultiplier * typeMultiplier * pavementMultiplier * classMultiplier) + bridgeAdditional
    
    
    if totalTime < 1 {
        totalTime = 1
    } else if totalTime > 365 {
        totalTime = 365
    }
    
    return int(totalTime)
}

func (uc *BinaMargaUseCase) sendUrgentNotification(ctx context.Context, report *entity.BinaMargaReport) {
    
    
    
    notificationType := "URGENT"
    if report.TrafficImpact == entity.TrafficImpactBlocked || report.TrafficCondition == entity.TrafficConditionBlocked {
        notificationType = "CRITICAL"
    }
    
    fmt.Printf("%s: Road damage report - %s on %s (%s impact, %s condition)\n", 
        notificationType, report.DamageType, report.RoadName, report.TrafficImpact, report.TrafficCondition)
    
    if report.BridgeName != "" {
        fmt.Printf("Bridge affected: %s with %s damage level\n", 
            report.BridgeName, report.BridgeDamageLevel)
    }
}


func (uc *BinaMargaUseCase) GetDashboard(ctx context.Context, roadType string, startDate, endDate time.Time) (*dto.BinaMargaDashboardResponse, error) {
    avgSeg, avgArea, avgTraffic, total, err := uc.binaMargaRepo.GetKPIs(ctx, roadType, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // Priority (urgency_level)
    priorityRows, err := uc.binaMargaRepo.GroupCountBy(ctx, "urgency_level", roadType, startDate, endDate, false, false)
    if err != nil {
        return nil, err
    }

    // Level distribusi
    roadLevelRows, err := uc.binaMargaRepo.GroupCountBy(ctx, "damage_level", roadType, startDate, endDate, false, true) // onlyRoad
    if err != nil {
        return nil, err
    }
    bridgeLevelRows, err := uc.binaMargaRepo.GroupCountBy(ctx, "bridge_damage_level", roadType, startDate, endDate, true, false) // onlyBridge
    if err != nil {
        return nil, err
    }

    // Top types
    roadTypeRows, err := uc.binaMargaRepo.GroupCountBy(ctx, "damage_type", roadType, startDate, endDate, false, true)
    if err != nil {
        return nil, err
    }
    bridgeTypeRows, err := uc.binaMargaRepo.GroupCountBy(ctx, "bridge_damage_type", roadType, startDate, endDate, true, false)
    if err != nil {
        return nil, err
    }

    // Map points
    pts, err := uc.binaMargaRepo.GetMapPoints(ctx, roadType, startDate, endDate)
    if err != nil {
        return nil, err
    }

    // Build DTO
    res := &dto.BinaMargaDashboardResponse{}
    res.KPIs.AvgSegmentLengthM = avgSeg
    res.KPIs.AvgDamageAreaM2 = avgArea
    res.KPIs.AvgDailyTrafficVolume = avgTraffic
    res.KPIs.TotalReports = total

    res.PriorityDistribution = make([]dto.KeyCount, len(priorityRows))
    for i, r0 := range priorityRows {
        res.PriorityDistribution[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.RoadDamageLevelDistribution = make([]dto.KeyCount, len(roadLevelRows))
    for i, r0 := range roadLevelRows {
        res.RoadDamageLevelDistribution[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.BridgeDamageLevelDistribution = make([]dto.KeyCount, len(bridgeLevelRows))
    for i, r0 := range bridgeLevelRows {
        res.BridgeDamageLevelDistribution[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.TopRoadDamageTypes = make([]dto.KeyCount, len(roadTypeRows))
    for i, r0 := range roadTypeRows {
        res.TopRoadDamageTypes[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.TopBridgeDamageTypes = make([]dto.KeyCount, len(bridgeTypeRows))
    for i, r0 := range bridgeTypeRows {
        res.TopBridgeDamageTypes[i] = dto.KeyCount{Key: r0.Key, Count: r0.Count}
    }

    res.MapPoints = make([]dto.BinaMargaMapPoint, len(pts))
    for i, p := range pts {
        res.MapPoints[i] = dto.BinaMargaMapPoint{
            Latitude:          p.Latitude,
            Longitude:         p.Longitude,
            RoadName:          p.RoadName,
            RoadType:          p.RoadType,
            DamageType:        p.DamageType,
            DamageLevel:       p.DamageLevel,
            BridgeName:        deref(p.BridgeName),
            BridgeDamageType:  deref(p.BridgeDamageType),
            BridgeDamageLevel: deref(p.BridgeDamageLevel),
            UrgencyLevel:      p.UrgencyLevel,
        }
    }
    return res, nil
}

func deref(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

func (uc *BinaMargaUseCase) GetBinaMargaOverview(ctx context.Context, roadType string) (*dto.BinaMargaOverviewResponse, error) {
    // Cache key based on road type
    cacheKey := fmt.Sprintf("bina_marga:overview:%s", roadType)
    var response dto.BinaMargaOverviewResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    // Initialize empty arrays to avoid null returns
    response.LocationDistribution = []dto.BinaMargaLocationStatsResponse{}
    response.PriorityDistribution = []dto.BinaMargaPriorityStatsResponse{}
    response.RoadDamageLevelDistribution = []dto.BinaMargaRoadDamageLevelStatsResponse{}
    response.BridgeDamageLevelDistribution = []dto.BinaMargaBridgeDamageLevelStatsResponse{}
    response.TopRoadDamageTypes = []dto.BinaMargaRoadDamageTypeStatsResponse{}
    response.TopBridgeDamageTypes = []dto.BinaMargaBridgeDamageTypeStatsResponse{}

    // Get basic statistics
    basicStatsRaw, err := uc.binaMargaRepo.GetBinaMargaOverviewStats(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get basic statistics: %w", err)
    }
    
    response.BasicStats.AvgSegmentLengthM = safeFloat64(basicStatsRaw["avg_segment_length_m"])
    response.BasicStats.AvgDamageAreaM2 = safeFloat64(basicStatsRaw["avg_damage_area_m2"])
    response.BasicStats.AvgDailyTrafficVolume = safeFloat64(basicStatsRaw["avg_daily_traffic_volume"])
    response.BasicStats.TotalInfrastructureReports = safeInt64(basicStatsRaw["total_infrastructure_reports"])

    // Get location distribution for mapping
    locationStats, err := uc.binaMargaRepo.GetBinaMargaLocationStats(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get location statistics: %w", err)
    }
    
    for _, loc := range locationStats {
        response.LocationDistribution = append(response.LocationDistribution, dto.BinaMargaLocationStatsResponse{
            RoadName:      loc["road_name"].(string),
            Latitude:      safeFloat64(loc["latitude"]),
            Longitude:     safeFloat64(loc["longitude"]),
            DamageType:    loc["damage_type"].(string),
            DamageLevel:   loc["damage_level"].(string),
            UrgencyLevel:  loc["urgency_level"].(string),
            TrafficImpact: loc["traffic_impact"].(string),
            DamagedArea:   safeFloat64(loc["damaged_area"]),
        })
    }

    // Get priority distribution
    priorityStats, err := uc.binaMargaRepo.GetBinaMargaPriorityStats(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get priority statistics: %w", err)
    }
    
    for _, priority := range priorityStats {
        response.PriorityDistribution = append(response.PriorityDistribution, dto.BinaMargaPriorityStatsResponse{
            PriorityLevel: priority["priority_level"].(string),
            Count:         safeInt64(priority["count"]),
        })
    }

    // Get road damage level distribution
    roadDamageLevelStats, err := uc.binaMargaRepo.GetBinaMargaRoadDamageLevelStats(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get road damage level statistics: %w", err)
    }
    
    for _, level := range roadDamageLevelStats {
        response.RoadDamageLevelDistribution = append(response.RoadDamageLevelDistribution, dto.BinaMargaRoadDamageLevelStatsResponse{
            DamageLevel: level["damage_level"].(string),
            Count:       safeInt64(level["count"]),
        })
    }

    // Get bridge damage level distribution
    bridgeDamageLevelStats, err := uc.binaMargaRepo.GetBinaMargaBridgeDamageLevelStats(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get bridge damage level statistics: %w", err)
    }
    
    for _, level := range bridgeDamageLevelStats {
        response.BridgeDamageLevelDistribution = append(response.BridgeDamageLevelDistribution, dto.BinaMargaBridgeDamageLevelStatsResponse{
            DamageLevel: level["damage_level"].(string),
            Count:       safeInt64(level["count"]),
        })
    }

    // Get top road damage types
    topRoadDamageTypes, err := uc.binaMargaRepo.GetBinaMargaTopRoadDamageTypes(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get top road damage types: %w", err)
    }
    
    for _, damageType := range topRoadDamageTypes {
        response.TopRoadDamageTypes = append(response.TopRoadDamageTypes, dto.BinaMargaRoadDamageTypeStatsResponse{
            DamageType: damageType["damage_type"].(string),
            Count:      safeInt64(damageType["count"]),
        })
    }

    // Get top bridge damage types
    topBridgeDamageTypes, err := uc.binaMargaRepo.GetBinaMargaTopBridgeDamageTypes(ctx, roadType)
    if err != nil {
        return nil, fmt.Errorf("failed to get top bridge damage types: %w", err)
    }
    
    for _, damageType := range topBridgeDamageTypes {
        response.TopBridgeDamageTypes = append(response.TopBridgeDamageTypes, dto.BinaMargaBridgeDamageTypeStatsResponse{
            DamageType: damageType["damage_type"].(string),
            Count:      safeInt64(damageType["count"]),
        })
    }

    // Cache the response for 5 minutes
    uc.cache.Set(ctx, cacheKey, &response, 300*time.Second)

    return &response, nil
}