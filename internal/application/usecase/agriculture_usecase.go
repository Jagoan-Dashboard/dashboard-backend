
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

func (uc *AgricultureUseCase) CreateReport(ctx context.Context, req *dto.CreateAgricultureRequest, photos []*multipart.FileHeader, userID uuid.UUID) (*entity.AgricultureReport, error) {
    report := &entity.AgricultureReport{
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
        CreatedBy:        userID,
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

func (uc *AgricultureUseCase) GetReport(ctx context.Context, id uuid.UUID) (*entity.AgricultureReport, error) {
    cacheKey := "agriculture:" + id.String()
    
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

func (uc *AgricultureUseCase) UpdateReport(ctx context.Context, id uuid.UUID, req *dto.UpdateAgricultureRequest, userID uuid.UUID) (*entity.AgricultureReport, error) {
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

    
    uc.cache.Delete(ctx, "agriculture:"+id.String())
    uc.cache.Delete(ctx, "agriculture:list")
    uc.cache.Delete(ctx, "agriculture:stats")

    return report, nil
}

func (uc *AgricultureUseCase) DeleteReport(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
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

    
    uc.cache.Delete(ctx, "agriculture:"+id.String())
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