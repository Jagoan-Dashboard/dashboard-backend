package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/storage"
	"building-report-backend/pkg/utils"
)

type AgricultureUseCase struct {
    agricultureRepo repository.AgricultureRepository
    storage         storage.StorageService
    cache           repository.CacheRepository
}

func NewAgricultureUseCase(
    agricultureRepo repository.AgricultureRepository,
    storage storage.StorageService,
    cache repository.CacheRepository,
) *AgricultureUseCase {
    return &AgricultureUseCase{
        agricultureRepo: agricultureRepo,
        storage:         storage,
        cache:           cache,
    }
}

func (uc *AgricultureUseCase) CreateReport(ctx context.Context, req *dto.CreateAgricultureRequest, photos []*multipart.FileHeader) (*entity.AgricultureReport, error) {
    report := &entity.AgricultureReport{
        ID:               utils.GenerateULID(),
        ExtensionOfficer: req.ExtensionOfficer,
        VisitDate:        req.VisitDate,
        FarmerName:       req.FarmerName,
        FarmerGroup:      req.FarmerGroup,
        Village:          req.Village,
        District:         req.District,
        Latitude:         req.Latitude,
        Longitude:        req.Longitude,
        WeatherCondition: entity.WeatherCondition(req.WeatherCondition),
        WeatherImpact:    entity.WeatherImpact(req.WeatherImpact),
        MainConstraint:   entity.MainConstraint(req.MainConstraint),
        FarmerHope:       entity.FarmerHope(req.FarmerHope),
        TrainingNeeded:   entity.TrainingNeeded(req.TrainingNeeded),
        UrgentNeeds:      entity.UrgentNeeds(req.UrgentNeeds),
        WaterAccess:      entity.WaterAccess(req.WaterAccess),
        Suggestions:      req.Suggestions,
        HasPestDisease:   req.HasPestDisease,
    }

    
    if req.FarmerGroupType != "" {
        report.FarmerGroupType = entity.FarmerGroupType(req.FarmerGroupType)
    }

    
    if req.FoodCommodity != "" {
        report.FoodCommodity = entity.FoodCommodity(req.FoodCommodity)
        report.FoodLandStatus = entity.LandStatus(req.FoodLandStatus)
        report.FoodLandArea = req.FoodLandArea
        report.FoodGrowthPhase = entity.GrowthPhase(req.FoodGrowthPhase)
        report.FoodPlantAge = req.FoodPlantAge
        report.FoodDelayReason = entity.DelayReason(req.FoodDelayReason)
        report.FoodTechnology = entity.TechnologyMethod(req.FoodTechnology)
        
        
        if req.FoodPlantingDate != "" {
            if plantingDate, err := time.Parse("2006-01-02", req.FoodPlantingDate); err == nil {
                report.FoodPlantingDate = &plantingDate
            }
        }
        if req.FoodHarvestDate != "" {
            if harvestDate, err := time.Parse("2006-01-02", req.FoodHarvestDate); err == nil {
                report.FoodHarvestDate = &harvestDate
            }
        }
    }

    
    if req.HortiCommodity != "" {
        report.HortiCommodity = entity.HorticultureCommodity(req.HortiCommodity)
        report.HortiSubCommodity = req.HortiSubCommodity
        report.HortiLandStatus = entity.LandStatus(req.HortiLandStatus)
        report.HortiLandArea = req.HortiLandArea
        report.HortiGrowthPhase = entity.HortiGrowthPhase(req.HortiGrowthPhase)
        report.HortiPlantAge = req.HortiPlantAge
        report.HortiDelayReason = entity.DelayReason(req.HortiDelayReason)
        report.HortiTechnology = entity.HortiTechnology(req.HortiTechnology)
        report.PostHarvestProblems = entity.PostHarvestProblem(req.PostHarvestProblems)
        
        
        if req.HortiPlantingDate != "" {
            if plantingDate, err := time.Parse("2006-01-02", req.HortiPlantingDate); err == nil {
                report.HortiPlantingDate = &plantingDate
            }
        }
        if req.HortiHarvestDate != "" {
            if harvestDate, err := time.Parse("2006-01-02", req.HortiHarvestDate); err == nil {
                report.HortiHarvestDate = &harvestDate
            }
        }
    }

    
    if req.PlantationCommodity != "" {
        report.PlantationCommodity = entity.PlantationCommodity(req.PlantationCommodity)
        report.PlantationLandStatus = entity.LandStatus(req.PlantationLandStatus)
        report.PlantationLandArea = req.PlantationLandArea
        report.PlantationGrowthPhase = entity.PlantationGrowthPhase(req.PlantationGrowthPhase)
        report.PlantationPlantAge = req.PlantationPlantAge
        report.PlantationDelayReason = entity.DelayReason(req.PlantationDelayReason)
        report.PlantationTechnology = entity.PlantationTechnology(req.PlantationTechnology)
        report.ProductionProblems = entity.ProductionProblem(req.ProductionProblems)
        
        
        if req.PlantationPlantingDate != "" {
            if plantingDate, err := time.Parse("2006-01-02", req.PlantationPlantingDate); err == nil {
                report.PlantationPlantingDate = &plantingDate
            }
        }
        if req.PlantationHarvestDate != "" {
            if harvestDate, err := time.Parse("2006-01-02", req.PlantationHarvestDate); err == nil {
                report.PlantationHarvestDate = &harvestDate
            }
        }
    }

    
    if req.HasPestDisease && req.PestDiseaseType != "" {
        report.PestDiseaseType = entity.PestDiseaseType(req.PestDiseaseType)
        report.PestDiseaseCommodity = req.PestDiseaseCommodity
        report.AffectedArea = entity.AffectedAreaLevel(req.AffectedArea)
        report.ControlAction = entity.ControlAction(req.ControlAction)
    }

    
    photoTypes := []string{"field", "crop", "general", "pest_disease"}
    for i, photo := range photos {
        photoType := "general"
        if i < len(photoTypes) {
            photoType = photoTypes[i]
        }

        photoURL, err := uc.storage.UploadFile(ctx, photo, "agriculture")
        if err != nil {
            return nil, fmt.Errorf("failed to upload photo: %w", err)
        }

        caption := fmt.Sprintf("%s - %s (%s)", photoType, report.FarmerName, report.Village)
        report.Photos = append(report.Photos, entity.AgriculturePhoto{
            ID:       utils.GenerateULID(),
            PhotoURL:  photoURL,
            PhotoType: photoType,
            Caption:   caption,
        })
    }

    if err := uc.agricultureRepo.Create(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "agriculture:list")
    uc.cache.Delete(ctx, "agriculture:stats")

    return report, nil
}

func (uc *AgricultureUseCase) GetReport(ctx context.Context, id string) (*entity.AgricultureReport, error) {
    cacheKey := "agriculture:" + id
    
    report, err := uc.agricultureRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, report, 3600)

    return report, nil
}

func (uc *AgricultureUseCase) ListReports(ctx context.Context, page, limit int, filters map[string]interface{}) (*dto.PaginatedAgricultureResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.agricultureRepo.FindAll(ctx, limit, offset, filters)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedAgricultureResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *AgricultureUseCase) UpdateReport(ctx context.Context, id string, req *dto.UpdateAgricultureRequest, userID string) (*entity.AgricultureReport, error) {
    report, err := uc.agricultureRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    
    if report.CreatedBy != userID {
        return nil, ErrUnauthorized
    }

    
    if req.ExtensionOfficer != "" {
        report.ExtensionOfficer = req.ExtensionOfficer
    }
    if req.FarmerName != "" {
        report.FarmerName = req.FarmerName
    }
    if req.FarmerGroup != "" {
        report.FarmerGroup = req.FarmerGroup
    }
    if req.FarmerGroupType != "" {
        report.FarmerGroupType = entity.FarmerGroupType(req.FarmerGroupType)
    }
    if req.Village != "" {
        report.Village = req.Village
    }
    if req.District != "" {
        report.District = req.District
    }

    
    if req.FoodCommodity != "" {
        report.FoodCommodity = entity.FoodCommodity(req.FoodCommodity)
    }
    if req.FoodLandStatus != "" {
        report.FoodLandStatus = entity.LandStatus(req.FoodLandStatus)
    }
    if req.FoodLandArea > 0 {
        report.FoodLandArea = req.FoodLandArea
    }
    if req.FoodGrowthPhase != "" {
        report.FoodGrowthPhase = entity.GrowthPhase(req.FoodGrowthPhase)
    }
    if req.FoodPlantAge > 0 {
        report.FoodPlantAge = req.FoodPlantAge
    }

    
    

    if req.WeatherCondition != "" {
        report.WeatherCondition = entity.WeatherCondition(req.WeatherCondition)
    }
    if req.WeatherImpact != "" {
        report.WeatherImpact = entity.WeatherImpact(req.WeatherImpact)
    }
    if req.MainConstraint != "" {
        report.MainConstraint = entity.MainConstraint(req.MainConstraint)
    }
    if req.FarmerHope != "" {
        report.FarmerHope = entity.FarmerHope(req.FarmerHope)
    }
    if req.TrainingNeeded != "" {
        report.TrainingNeeded = entity.TrainingNeeded(req.TrainingNeeded)
    }
    if req.UrgentNeeds != "" {
        report.UrgentNeeds = entity.UrgentNeeds(req.UrgentNeeds)
    }
    if req.WaterAccess != "" {
        report.WaterAccess = entity.WaterAccess(req.WaterAccess)
    }
    if req.Suggestions != "" {
        report.Suggestions = req.Suggestions
    }

    if err := uc.agricultureRepo.Update(ctx, report); err != nil {
        return nil, err
    }

    
    uc.cache.Delete(ctx, "agriculture:"+id)
    uc.cache.Delete(ctx, "agriculture:list")
    uc.cache.Delete(ctx, "agriculture:stats")

    return report, nil
}

func (uc *AgricultureUseCase) DeleteReport(ctx context.Context, id string, userID string) error {
    report, err := uc.agricultureRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    
    if report.CreatedBy != userID {
        return ErrUnauthorized
    }

    
    for _, photo := range report.Photos {
        uc.storage.DeleteFile(ctx, photo.PhotoURL)
    }

    if err := uc.agricultureRepo.Delete(ctx, id); err != nil {
        return err
    }

    
    uc.cache.Delete(ctx, "agriculture:"+id)
    uc.cache.Delete(ctx, "agriculture:list")
    uc.cache.Delete(ctx, "agriculture:stats")

    return nil
}

func (uc *AgricultureUseCase) GetStatistics(ctx context.Context) (*dto.AgricultureStatisticsResponse, error) {
    
    cacheKey := "agriculture:stats"
    var stats dto.AgricultureStatisticsResponse
    
    err := uc.cache.Get(ctx, cacheKey, &stats)
    if err == nil {
        return &stats, nil
    }

    
    rawStats, err := uc.agricultureRepo.GetStatistics(ctx)
    if err != nil {
        return nil, err
    }

    
    response := &dto.AgricultureStatisticsResponse{
        TotalReports:               rawStats["total_reports"].(int64),
        TotalFarmers:               rawStats["total_farmers"].(int64),
        TotalLandArea:              rawStats["total_land_area_ha"].(float64),
        FoodCropReports:            rawStats["food_crop_reports"].(int64),
        HorticultureReports:        rawStats["horticulture_reports"].(int64),
        PlantationReports:          rawStats["plantation_reports"].(int64),
        ReportsWithPestDisease:     rawStats["reports_with_pest_disease"].(int64),
        PestDiseasePercentage:      rawStats["pest_disease_percentage"].(float64),
        PostHarvestProblemReports:  rawStats["post_harvest_problem_reports"].(int64),
        ProductionProblemReports:   rawStats["production_problem_reports"].(int64),
    }

    
    if villageStats, ok := rawStats["village_distribution"].([]interface{}); ok {
        for _, v := range villageStats {
            if m, ok := v.(map[string]interface{}); ok {
                response.VillageDistribution = append(response.VillageDistribution, m)
            }
        }
    }

    if extensionStats, ok := rawStats["extension_officer_stats"].([]interface{}); ok {
        for _, v := range extensionStats {
            if m, ok := v.(map[string]interface{}); ok {
                response.ExtensionOfficerStats = append(response.ExtensionOfficerStats, m)
            }
        }
    }

    
    uc.cache.Set(ctx, cacheKey, response, 300)

    return response, nil
}

func (uc *AgricultureUseCase) GetByExtensionOfficer(ctx context.Context, extensionOfficer string, page, limit int) (*dto.PaginatedAgricultureResponse, error) {
    offset := (page - 1) * limit
    
    reports, total, err := uc.agricultureRepo.FindByExtensionOfficer(ctx, extensionOfficer, limit, offset)
    if err != nil {
        return nil, err
    }

    return &dto.PaginatedAgricultureResponse{
        Reports:     reports,
        Total:       total,
        Page:        page,
        PerPage:     limit,
        TotalPages:  (total + int64(limit) - 1) / int64(limit),
    }, nil
}

func (uc *AgricultureUseCase) GetPestDiseaseReports(ctx context.Context, limit int) ([]*entity.AgricultureReport, error) {
    cacheKey := fmt.Sprintf("agriculture:pest_disease:%d", limit)
    
    reports, err := uc.agricultureRepo.GetPestDiseaseReports(ctx, limit)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, reports, 300)

    return reports, nil
}

func (uc *AgricultureUseCase) GetCommodityProduction(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    return uc.agricultureRepo.GetCommodityProduction(ctx, startDate, endDate)
}

func (uc *AgricultureUseCase) GetExtensionOfficerPerformance(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
    return uc.agricultureRepo.GetExtensionOfficerPerformance(ctx, startDate, endDate)
}

func (uc *AgricultureUseCase) GetTechnologyAdoptionStats(ctx context.Context) (map[string]interface{}, error) {
    cacheKey := "agriculture:technology_adoption"
    
    var stats map[string]interface{}
    err := uc.cache.Get(ctx, cacheKey, &stats)
    if err == nil {
        return stats, nil
    }

    stats, err = uc.agricultureRepo.GetTechnologyAdoptionStats(ctx)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, stats, 600)

    return stats, nil
}

func (uc *AgricultureUseCase) GetFarmerNeedsAnalysis(ctx context.Context) (map[string]interface{}, error) {
    cacheKey := "agriculture:farmer_needs"
    
    var analysis map[string]interface{}
    err := uc.cache.Get(ctx, cacheKey, &analysis)
    if err == nil {
        return analysis, nil
    }

    analysis, err = uc.agricultureRepo.GetFarmerNeedsAnalysis(ctx)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, analysis, 600)

    return analysis, nil
}

func (uc *AgricultureUseCase) GetExecutiveSummary(ctx context.Context, commodityType string) (*dto.AgricultureExecutiveResponse, error) {
    
    commodityType = strings.ToUpper(strings.TrimSpace(commodityType))
    
    cacheKey := fmt.Sprintf("agriculture:executive_summary:%s", commodityType)
    var response dto.AgricultureExecutiveResponse
    
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    summary, err := uc.agricultureRepo.GetExecutiveSummary(ctx, commodityType)
    if err != nil {
        return nil, err
    }

    
    commodityMap, err := uc.agricultureRepo.GetCommodityDistributionByDistrict(ctx, commodityType)
    if err != nil {
        return nil, err
    }

    
    sectorData, err := uc.agricultureRepo.GetCommodityCountBySector(ctx, commodityType)
    if err != nil {
        return nil, err
    }

    
    landStatus, err := uc.agricultureRepo.GetLandStatusDistribution(ctx, commodityType)
    if err != nil {
        return nil, err
    }

    
    constraints, err := uc.agricultureRepo.GetMainConstraintsDistribution(ctx, commodityType)
    if err != nil {
        return nil, err
    }

    farmerNeeds, err := uc.agricultureRepo.GetFarmerHopesAndNeeds(ctx, commodityType)
    if err != nil {
        return nil, err
    }

    
    response.TotalLandArea = summary["total_land_area"].(float64)
    response.PestDiseaseReports = summary["pest_disease_reports"].(int64)
    response.TotalExtensionReports = summary["total_extension_reports"].(int64)

    for _, mapItem := range commodityMap {
        response.CommodityMap = append(response.CommodityMap, dto.CommodityMapPoint{
            Latitude:      mapItem["latitude"].(float64),
            Longitude:     mapItem["longitude"].(float64),
            Village:       mapItem["village"].(string),
            District:      mapItem["district"].(string),
            Commodity:     mapItem["commodity"].(string),
            CommodityType: mapItem["commodity_type"].(string),
            LandArea:      mapItem["land_area"].(float64),
        })
    }

    // Convert sector data
    if foodCrops, ok := sectorData["food_crops"].([]map[string]interface{}); ok {
        for _, item := range foodCrops {
            response.CommodityBySector.FoodCrops = append(response.CommodityBySector.FoodCrops,
                dto.CommodityCount{
                    Name:  item["name"].(string),
                    Count: item["count"].(int64),
                })
        }
    }

    if horticulture, ok := sectorData["horticulture"].([]map[string]interface{}); ok {
        for _, item := range horticulture {
            response.CommodityBySector.Horticulture = append(response.CommodityBySector.Horticulture,
                dto.CommodityCount{
                    Name:  item["name"].(string),
                    Count: item["count"].(int64),
                })
        }
    }

    if plantation, ok := sectorData["plantation"].([]map[string]interface{}); ok {
        for _, item := range plantation {
            response.CommodityBySector.Plantation = append(response.CommodityBySector.Plantation,
                dto.CommodityCount{
                    Name:  item["name"].(string),
                    Count: item["count"].(int64),
                })
        }
    }

    // Convert land status
    for _, item := range landStatus {
        response.LandStatusDistrib = append(response.LandStatusDistrib, dto.LandStatusCount{
            Status:     item["status"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    // Convert constraints
    for _, item := range constraints {
        response.MainConstraints = append(response.MainConstraints, dto.ConstraintCount{
            Constraint: item["constraint"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    // Convert farmer needs
    if hopes, ok := farmerNeeds["hopes"].([]map[string]interface{}); ok {
        for _, item := range hopes {
            response.FarmerHopesNeeds.Hopes = append(response.FarmerHopesNeeds.Hopes,
                dto.HopeCount{
                    Hope:       item["hope"].(string),
                    Count:      item["count"].(int64),
                    Percentage: item["percentage"].(float64),
                })
        }
    }

    if needs, ok := farmerNeeds["needs"].([]map[string]interface{}); ok {
        for _, item := range needs {
            response.FarmerHopesNeeds.Needs = append(response.FarmerHopesNeeds.Needs,
                dto.NeedCount{
                    Need:       item["need"].(string),
                    Count:      item["count"].(int64),
                    Percentage: item["percentage"].(float64),
                })
        }
    }


    
    uc.cache.Set(ctx, cacheKey, &response, 600*time.Second)

    return &response, nil
}


func (uc *AgricultureUseCase) GetFoodCropStats(ctx context.Context, commodityName string) (*dto.FoodCropResponse, error) {
    cacheKey := fmt.Sprintf("agriculture:food_crop_stats:%s", commodityName)
    
    var response dto.FoodCropResponse
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    stats, err := uc.agricultureRepo.GetFoodCropStats(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    distribution, err := uc.agricultureRepo.GetFoodCropDistribution(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    growthPhases, err := uc.agricultureRepo.GetFoodCropGrowthPhases(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    technology, err := uc.agricultureRepo.GetFoodCropTechnology(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    pestDominance, err := uc.agricultureRepo.GetFoodCropPestDominance(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    harvestSchedule, err := uc.agricultureRepo.GetFoodCropHarvestSchedule(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    response.LandArea = stats["land_area"].(float64)
    response.EstimatedProduction = stats["estimated_production"].(float64)
    response.PestAffectedArea = stats["pest_affected_area"].(float64)
    response.PestReportCount = stats["pest_report_count"].(int64)

    
    for _, item := range distribution {
        response.DistributionMap = append(response.DistributionMap, dto.CommodityMapPoint{
            Latitude:      item["latitude"].(float64),
            Longitude:     item["longitude"].(float64),
            Village:       item["village"].(string),
            District:      item["district"].(string),
            Commodity:     item["commodity"].(string),
            CommodityType: "FOOD",
            LandArea:      item["land_area"].(float64),
        })
    }

    
    for _, item := range growthPhases {
        response.GrowthPhases = append(response.GrowthPhases, dto.GrowthPhaseCount{
            Phase:      item["phase"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range technology {
        response.TechnologyUsed = append(response.TechnologyUsed, dto.TechnologyCount{
            Technology: item["technology"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range pestDominance {
        response.PestDominance = append(response.PestDominance, dto.PestDominanceCount{
            PestType:   item["pest_type"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range harvestSchedule {
    var harvestDate time.Time
    
    switch v := item["harvest_date"].(type) {
    case time.Time:
        harvestDate = v
    case string:
        if parsedDate, err := time.Parse("2006-01-02", v); err == nil {
            harvestDate = parsedDate
        }
    case nil:
        continue
    }
    
    response.HarvestSchedule = append(response.HarvestSchedule, dto.HarvestScheduleItem{
        CommodityDetail: item["commodity_detail"].(string),
        HarvestDate:     harvestDate,
        FarmerName:      item["farmer_name"].(string),
        Village:         item["village"].(string),
        LandArea:        item["land_area"].(float64),
    })
}

    
    uc.cache.Set(ctx, cacheKey, &response, 900*time.Second)

    return &response, nil
}

func (uc *AgricultureUseCase) GetHorticultureStats(ctx context.Context, commodityName string) (*dto.HorticultureResponse, error) {
    cacheKey := fmt.Sprintf("agriculture:horticulture_stats:%s", commodityName)
    
    var response dto.HorticultureResponse
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    stats, err := uc.agricultureRepo.GetHorticultureStats(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    distribution, err := uc.agricultureRepo.GetHorticultureDistribution(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    growthPhases, err := uc.agricultureRepo.GetHorticultureGrowthPhases(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    technology, err := uc.agricultureRepo.GetHorticultureTechnology(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    pestDominance, err := uc.agricultureRepo.GetHorticulturePestDominance(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    harvestSchedule, err := uc.agricultureRepo.GetHorticultureHarvestSchedule(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    response.LandArea = stats["land_area"].(float64)
    response.EstimatedProduction = stats["estimated_production"].(float64)
    response.PestAffectedArea = stats["pest_affected_area"].(float64)
    response.PestReportCount = stats["pest_report_count"].(int64)

    
    for _, item := range distribution {
        response.DistributionMap = append(response.DistributionMap, dto.CommodityMapPoint{
            Latitude:      item["latitude"].(float64),
            Longitude:     item["longitude"].(float64),
            Village:       item["village"].(string),
            District:      item["district"].(string),
            Commodity:     item["commodity"].(string),
            CommodityType: "HORTICULTURE",
            LandArea:      item["land_area"].(float64),
        })
    }

    
    for _, item := range growthPhases {
        response.GrowthPhases = append(response.GrowthPhases, dto.GrowthPhaseCount{
            Phase:      item["phase"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range technology {
        response.TechnologyUsed = append(response.TechnologyUsed, dto.TechnologyCount{
            Technology: item["technology"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range pestDominance {
        response.PestDominance = append(response.PestDominance, dto.PestDominanceCount{
            PestType:   item["pest_type"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range harvestSchedule {
    var harvestDate time.Time
    
    
    switch v := item["harvest_date"].(type) {
    case time.Time:
        harvestDate = v
    case string:
        if parsedDate, err := time.Parse("2006-01-02", v); err == nil {
            harvestDate = parsedDate
        }
    case nil:
        
        continue
    }
    
    response.HarvestSchedule = append(response.HarvestSchedule, dto.HarvestScheduleItem{
        CommodityDetail: item["commodity_detail"].(string),
        HarvestDate:     harvestDate,
        FarmerName:      item["farmer_name"].(string),
        Village:         item["village"].(string),
        LandArea:        item["land_area"].(float64),
    })
}

    
    uc.cache.Set(ctx, cacheKey, &response, 900*time.Second)

    return &response, nil
}

func (uc *AgricultureUseCase) GetPlantationStats(ctx context.Context, commodityName string) (*dto.PlantationResponse, error) {
    cacheKey := fmt.Sprintf("agriculture:plantation_stats:%s", commodityName)
    
    var response dto.PlantationResponse
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    stats, err := uc.agricultureRepo.GetPlantationStats(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    distribution, err := uc.agricultureRepo.GetPlantationDistribution(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    growthPhases, err := uc.agricultureRepo.GetPlantationGrowthPhases(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    technology, err := uc.agricultureRepo.GetPlantationTechnology(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    pestDominance, err := uc.agricultureRepo.GetPlantationPestDominance(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    harvestSchedule, err := uc.agricultureRepo.GetPlantationHarvestSchedule(ctx, commodityName)
    if err != nil {
        return nil, err
    }

    
    response.LandArea = stats["land_area"].(float64)
    response.EstimatedProduction = stats["estimated_production"].(float64)
    response.PestAffectedArea = stats["pest_affected_area"].(float64)
    response.PestReportCount = stats["pest_report_count"].(int64)

    
    for _, item := range distribution {
        response.DistributionMap = append(response.DistributionMap, dto.CommodityMapPoint{
            Latitude:      item["latitude"].(float64),
            Longitude:     item["longitude"].(float64),
            Village:       item["village"].(string),
            District:      item["district"].(string),
            Commodity:     item["commodity"].(string),
            CommodityType: "PLANTATION",
            LandArea:      item["land_area"].(float64),
        })
    }

    
    for _, item := range growthPhases {
        response.GrowthPhases = append(response.GrowthPhases, dto.GrowthPhaseCount{
            Phase:      item["phase"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range technology {
        response.TechnologyUsed = append(response.TechnologyUsed, dto.TechnologyCount{
            Technology: item["technology"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range pestDominance {
        response.PestDominance = append(response.PestDominance, dto.PestDominanceCount{
            PestType:   item["pest_type"].(string),
            Count:      item["count"].(int64),
            Percentage: item["percentage"].(float64),
        })
    }

    
    for _, item := range harvestSchedule {
        var harvestDate time.Time
        
        
        switch v := item["harvest_date"].(type) {
        case time.Time:
            harvestDate = v
        case string:
            if parsedDate, err := time.Parse("2006-01-02", v); err == nil {
                harvestDate = parsedDate
            }
        case nil:
            
            continue
        default:
            
            continue
        }
        
        response.HarvestSchedule = append(response.HarvestSchedule, dto.HarvestScheduleItem{
            CommodityDetail: item["commodity_detail"].(string),
            HarvestDate:     harvestDate,
            FarmerName:      item["farmer_name"].(string),
            Village:         item["village"].(string),
            LandArea:        item["land_area"].(float64),
        })
    }

    
    uc.cache.Set(ctx, cacheKey, &response, 900*time.Second)

    return &response, nil
}

func (uc *AgricultureUseCase) GetAgriculturalEquipmentStats(ctx context.Context, startDate, endDate time.Time) (*dto.AgriculturalEquipmentResponse, error) {
    cacheKey := fmt.Sprintf("agriculture:equipment_stats:%s:%s", 
        startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    
    var response dto.AgriculturalEquipmentResponse
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    stats, err := uc.agricultureRepo.GetAgriculturalEquipmentStats(ctx, startDate, endDate)
    if err != nil {
        return nil, err
    }

    
    distribution, err := uc.agricultureRepo.GetEquipmentDistributionByDistrict(ctx, startDate, endDate)
    if err != nil {
        return nil, err
    }

    
    years := []int{2018, 2019, 2020, 2021, 2022, 2023, 2024}
    waterPumpTrend, err := uc.agricultureRepo.GetEquipmentTrend(ctx, "water_pump", years)
    if err != nil {
        return nil, err
    }

    
    response.GrainProcessor = dto.EquipmentCount{
        Count:         convertToInt64(stats["grain_processor_count"]),
        GrowthPercent: convertToFloat64(stats["grain_processor_growth"]),
    }
    
    response.MultipurposeThresher = dto.EquipmentCount{
        Count:         convertToInt64(stats["thresher_count"]),
        GrowthPercent: convertToFloat64(stats["thresher_growth"]),
    }
    
    response.FarmMachinery = dto.EquipmentCount{
        Count:         convertToInt64(stats["machinery_count"]),
        GrowthPercent: convertToFloat64(stats["machinery_growth"]),
    }
    
    response.WaterPump = dto.EquipmentCount{
        Count:         convertToInt64(stats["water_pump_count"]),
        GrowthPercent: convertToFloat64(stats["water_pump_growth"]),
    }

    
    for _, item := range distribution {
        response.DistributionByDistrict = append(response.DistributionByDistrict, dto.EquipmentDistrict{
            District:       convertToString(item["district"]),
            GrainProcessor: convertToInt64(item["grain_processor"]),
            Thresher:       convertToInt64(item["thresher"]),
            FarmMachinery:  convertToInt64(item["farm_machinery"]),
            WaterPump:      convertToInt64(item["water_pump"]),
        })
    }

    
    for _, item := range waterPumpTrend {
        response.WaterPumpTrend = append(response.WaterPumpTrend, dto.EquipmentTrend{
            Year:  int(convertToInt64(item["year"])),
            Count: convertToInt64(item["count"]),
        })
    }

    
    uc.cache.Set(ctx, cacheKey, &response, 1200*time.Second)

    return &response, nil
}

func (uc *AgricultureUseCase) GetLandAndIrrigationStats(ctx context.Context, startDate, endDate time.Time) (*dto.LandIrrigationResponse, error) {
    cacheKey := fmt.Sprintf("agriculture:land_irrigation:%s:%s", 
        startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    
    var response dto.LandIrrigationResponse
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    stats, err := uc.agricultureRepo.GetLandAndIrrigationStats(ctx, startDate, endDate)
    if err != nil {
        return nil, err
    }

    
    distribution, err := uc.agricultureRepo.GetLandDistributionByDistrict(ctx, startDate, endDate)
    if err != nil {
        return nil, err
    }

    
    response.TotalLandArea = dto.LandAreaCount{
        Area:          stats["total_land_area"].(float64),
        GrowthPercent: stats["total_land_growth"].(float64),
    }
    
    response.IrrigatedLandArea = dto.LandAreaCount{
        Area:          stats["irrigated_land_area"].(float64),
        GrowthPercent: stats["irrigated_land_growth"].(float64),
    }
    
    response.NonIrrigatedLandArea = dto.LandAreaCount{
        Area:          stats["non_irrigated_land_area"].(float64),
        GrowthPercent: stats["non_irrigated_land_growth"].(float64),
    }

    
    for _, item := range distribution {
        districtItem := dto.LandDistrict{
            District:      item["district"].(string),
            IrrigatedArea: item["irrigated_area"].(float64),
            TotalArea:     item["total_area"].(float64),
            FarmerCount:   int(item["farmer_count"].(int64)),
        }
        response.IrrigatedByDistrict = append(response.IrrigatedByDistrict, districtItem)
        
        response.LandDistribution = append(response.LandDistribution, dto.LandDistributionItem{
            District:       item["district"].(string),
            TotalArea:      item["total_area"].(float64),
            IrrigatedArea:  item["irrigated_area"].(float64),
            FoodCropArea:   item["food_crop_area"].(float64),
            HortiArea:      item["horti_area"].(float64),
            PlantationArea: item["plantation_area"].(float64),
        })
    }

    
    uc.cache.Set(ctx, cacheKey, &response, 1200*time.Second)

    return &response, nil
}
func (uc *AgricultureUseCase) GetCommodityAnalysis(ctx context.Context, startDate, endDate time.Time, commodityName string) (*dto.CommodityAnalysisResponse, error) {
    cacheKey := fmt.Sprintf("agriculture:commodity_analysis:%s:%s:%s", 
        commodityName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
    
    var response dto.CommodityAnalysisResponse
    err := uc.cache.Get(ctx, cacheKey, &response)
    if err == nil {
        return &response, nil
    }

    
    analysis, err := uc.agricultureRepo.GetCommodityAnalysis(ctx, startDate, endDate, commodityName)
    if err != nil {
        return nil, err
    }

    
    districtProduction, err := uc.agricultureRepo.GetProductionByDistrict(ctx, startDate, endDate, commodityName)
    if err != nil {
        return nil, err
    }

    
    years := []int{2018, 2019, 2020, 2021, 2022, 2023, 2024}
    productivityTrend, err := uc.agricultureRepo.GetProductivityTrend(ctx, commodityName, years)
    if err != nil {
        return nil, err
    }

    
    response.TotalProduction = analysis["total_production"].(float64)
    response.ProductionGrowth = analysis["production_growth"].(float64)
    response.TotalHarvestedArea = analysis["total_harvested_area"].(float64)
    response.HarvestedAreaGrowth = analysis["harvested_area_growth"].(float64)
    response.Productivity = analysis["productivity"].(float64)
    response.ProductivityGrowth = analysis["productivity_growth"].(float64)

    
    for _, item := range districtProduction {
        response.ProductionByDistrict = append(response.ProductionByDistrict, dto.ProductionDistrict{
            District:      item["district"].(string),
            Production:    item["production"].(float64),
            HarvestedArea: item["harvested_area"].(float64),
            FarmerCount:   int(item["farmer_count"].(int64)),
        })
    }

    
    for _, item := range productivityTrend {
        response.ProductivityTrend = append(response.ProductivityTrend, dto.ProductivityTrend{
            Year:         int(item["year"].(int64)),
            Productivity: item["productivity"].(float64),
            Production:   item["production"].(float64),
            Area:         item["area"].(float64),
        })
    }

    
    uc.cache.Set(ctx, cacheKey, &response, 1800*time.Second)

    return &response, nil
}



func (uc *AgricultureUseCase) GetDashboardWithFilters(ctx context.Context, filters map[string]interface{}) (*dto.AgricultureExecutiveResponse, error) {


    return uc.GetExecutiveSummary(ctx,  "")
}

func (uc *AgricultureUseCase) GetCommodityTrendAnalysis(ctx context.Context, commodityName string, years []int) ([]dto.ProductivityTrend, error) {
    cacheKey := fmt.Sprintf("agriculture:commodity_trend:%s", commodityName)
    
    var results []dto.ProductivityTrend
    err := uc.cache.Get(ctx, cacheKey, &results)
    if err == nil {
        return results, nil
    }

    
    trendData, err := uc.agricultureRepo.GetProductivityTrend(ctx, commodityName, years)
    if err != nil {
        return nil, err
    }

    
    for _, item := range trendData {
        results = append(results, dto.ProductivityTrend{
            Year:         int(item["year"].(int64)),
            Productivity: item["productivity"].(float64),
            Production:   item["production"].(float64),
            Area:         item["area"].(float64),
        })
    }

    
    uc.cache.Set(ctx, cacheKey, results, 3600*time.Second)

    return results, nil
}

func (uc *AgricultureUseCase) GetFoodCropsByDistrict(ctx context.Context, district string) (*dto.FoodCropResponse, error) {
    
    
    return uc.GetFoodCropStats(ctx, "")
}

func (uc *AgricultureUseCase) GetHortticultureByDistrict(ctx context.Context, district string) (*dto.HorticultureResponse, error) {
    
    return uc.GetHorticultureStats(ctx, "")
}

func (uc *AgricultureUseCase) GetPlantationByDistrict(ctx context.Context, district string) (*dto.PlantationResponse, error) {
    
    return uc.GetPlantationStats(ctx, "")
}


func (uc *AgricultureUseCase) GetCommodityComparison(ctx context.Context, startDate, endDate time.Time, commodities []string) (map[string]*dto.CommodityAnalysisResponse, error) {
    results := make(map[string]*dto.CommodityAnalysisResponse)
    
    for _, commodity := range commodities {
        analysis, err := uc.GetCommodityAnalysis(ctx, startDate, endDate, commodity)
        if err != nil {
            
            continue
        }
        results[commodity] = analysis
    }
    
    return results, nil
}


func (uc *AgricultureUseCase) ValidateCommodityName(commodityName, commodityType string) bool {
    validFoodCrops := map[string]bool{
        "PADI_SAWAH": true, "PADI_LADANG": true, "JAGUNG": true, 
        "KEDELAI": true, "KACANG_TANAH": true, "UBI_KAYU": true, "UBI_JALAR": true,
    }
    
    validHorticulture := map[string]bool{
        "SAYURAN": true, "BUAH": true, "FLORIKULTURA": true, "TANAMAN_OBAT_TRADISIONAL": true,
    }
    
    validPlantation := map[string]bool{
        "KOPI": true, "KAKAO": true, "KELAPA": true, "KELAPA_SAWIT": true, 
        "CENGKEH": true, "TEBU": true, "KARET": true, "TEMBAKAU": true, 
        "VANILI": true, "LADA": true, "PALA": true,
    }
    
    switch commodityType {
    case "FOOD":
        return validFoodCrops[commodityName]
    case "HORTICULTURE":
        return validHorticulture[commodityName]
    case "PLANTATION":
        return validPlantation[commodityName]
    default:
        return validFoodCrops[commodityName] || validHorticulture[commodityName] || validPlantation[commodityName]
    }
}


func (uc *AgricultureUseCase) GetAvailableCommodities(ctx context.Context) (map[string][]string, error) {
    cacheKey := "agriculture:available_commodities"
    
    var result map[string][]string
    err := uc.cache.Get(ctx, cacheKey, &result)
    if err == nil {
        return result, nil
    }

    result = map[string][]string{
        "food_crops": {
            "PADI_SAWAH", "PADI_LADANG", "JAGUNG", "KEDELAI", 
            "KACANG_TANAH", "UBI_KAYU", "UBI_JALAR",
        },
        "horticulture": {
            "SAYURAN", "BUAH", "FLORIKULTURA", "TANAMAN_OBAT_TRADISIONAL",
        },
        "plantation": {
            "KOPI", "KAKAO", "KELAPA", "KELAPA_SAWIT", "CENGKEH", 
            "TEBU", "KARET", "TEMBAKAU", "VANILI", "LADA", "PALA",
        },
    }
    
    
    uc.cache.Set(ctx, cacheKey, result, 24*3600*time.Second)
    
    return result, nil
}


func (uc *AgricultureUseCase) CalculateProductionEstimate(landArea float64, commodityType string, technologyUsed bool, irrigated bool) float64 {
    var baseProductivity float64
    
    
    switch commodityType {
    case "FOOD":
        baseProductivity = 3.0
    case "HORTICULTURE":
        baseProductivity = 5.0
    case "PLANTATION":
        baseProductivity = 2.0
    default:
        baseProductivity = 3.0
    }
    
    
    if technologyUsed {
        baseProductivity *= 1.2 
    }
    
    
    if irrigated {
        baseProductivity *= 1.15 
    }
    
    return landArea * baseProductivity
}


func (uc *AgricultureUseCase) GetPlantingCalendar(ctx context.Context, commodityType string) (map[string][]string, error) {
    
    calendar := map[string][]string{
        "PADI_SAWAH": {"Oktober-November", "Februari-Maret", "Juni-Juli"},
        "JAGUNG": {"September-Oktober", "Januari-Februari", "Mei-Juni"},
        "KEDELAI": {"April-Mei", "Juli-Agustus"},
        "SAYURAN": {"Sepanjang Tahun (tergantung jenis)"},
        "KOPI": {"Oktober-Desember"},
        "KAKAO": {"Sepanjang Tahun"},
    }
    
    if commodityType != "" {
        if seasons, exists := calendar[commodityType]; exists {
            return map[string][]string{commodityType: seasons}, nil
        }
    }
    
    return calendar, nil
}

func convertToString(value interface{}) string {
    switch v := value.(type) {
    case string:
        return v
    case []byte:
        return string(v)
    default:
        return fmt.Sprintf("%v", v)
    }
}


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
    case nil:
        return 0
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
    case nil:
        return 0.0
    default:
        return 0.0
    }
}