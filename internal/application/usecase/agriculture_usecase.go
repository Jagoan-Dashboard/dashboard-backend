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
			ID:        utils.GenerateULID(),
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
		Reports:    reports,
		Total:      total,
		Page:       page,
		PerPage:    limit,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}, nil
}

func (uc *AgricultureUseCase) UpdateReport(ctx context.Context, id string, req *dto.UpdateAgricultureRequest, userID string) (*entity.AgricultureReport, error) {
	report, err := uc.agricultureRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
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

	for _, item := range landStatus {
		response.LandStatusDistrib = append(response.LandStatusDistrib, dto.LandStatusCount{
			Status:     item["status"].(string),
			Count:      item["count"].(int64),
			Percentage: item["percentage"].(float64),
		})
	}

	for _, item := range constraints {
		response.MainConstraints = append(response.MainConstraints, dto.ConstraintCount{
			Constraint: item["constraint"].(string),
			Count:      item["count"].(int64),
			Percentage: item["percentage"].(float64),
		})
	}

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

	individualDistribution, err := uc.agricultureRepo.GetEquipmentIndividualDistribution(ctx, startDate, endDate)
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

	response.IndividualDistribution = []dto.EquipmentIndividualLocation{}
	for _, item := range individualDistribution {
		response.IndividualDistribution = append(response.IndividualDistribution, dto.EquipmentIndividualLocation{
			Latitude:       convertToFloat64(item["latitude"]),
			Longitude:      convertToFloat64(item["longitude"]),
			Village:        convertToString(item["village"]),
			District:       convertToString(item["district"]),
			FarmerName:     convertToString(item["farmer_name"]),
			TechnologyType: convertToString(item["technology_type"]),
			Commodity:      convertToString(item["commodity"]),
			VisitDate:      convertToString(item["visit_date"]),
		})
	}

	response.WaterPumpTrend = []dto.EquipmentTrend{}
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

	individualPoints, err := uc.agricultureRepo.GetLandIndividualDistribution(ctx, startDate, endDate)
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

	for _, point := range individualPoints {
		ricePoint := dto.LandIndividualPoint{
			Latitude:            convertToFloat64(point["latitude"]),
			Longitude:           convertToFloat64(point["longitude"]),
			District:            convertToString(point["district"]),
			RainfedRiceFields:   convertToFloat64(point["rainfed_rice_fields"]),
			IrrigatedRiceFields: convertToFloat64(point["irrigated_rice_fields"]),
			TotalRiceFieldArea:  convertToFloat64(point["total_rice_field_area"]),
			DataSource:          convertToString(point["data_source"]),
		}

		// Format date jika perlu
		if dateVal, ok := point["date"]; ok && dateVal != nil {
			var dateStr string
			switch v := dateVal.(type) {
			case time.Time:
				dateStr = v.Format("2006-01-02")
			case string:
				dateStr = v
			default:
				dateStr = ""
			}
			ricePoint.Date = dateStr
		}

		response.IndividualPoints = append(response.IndividualPoints, ricePoint)
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

	distribution, err := uc.agricultureRepo.GetProductionByDistrict(ctx, startDate, endDate, commodityName)
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

	response.ProductionDistribution = []dto.ProductionLocation{}
	if len(distribution) > 0 {
		for _, item := range distribution {
			response.ProductionDistribution = append(response.ProductionDistribution, dto.ProductionLocation{
				Latitude:            item["latitude"].(float64),
				Longitude:           item["longitude"].(float64),
				Village:             item["village"].(string),
				District:            item["district"].(string),
				Commodity:           item["commodity"].(string),
				LandArea:            item["land_area"].(float64),
				EstimatedProduction: item["estimated_production"].(float64),
				FarmerName:          item["farmer_name"].(string),
			})
		}
	}

	response.ProductivityTrend = []dto.ProductivityTrend{}
	if len(productivityTrend) > 0 {
		for _, item := range productivityTrend {
			response.ProductivityTrend = append(response.ProductivityTrend, dto.ProductivityTrend{
				Year:         int(item["year"].(int64)),
				Productivity: item["productivity"].(float64),
				Production:   item["production"].(float64),
				Area:         item["area"].(float64),
			})
		}
	}

	uc.cache.Set(ctx, cacheKey, &response, 1800*time.Second)

	return &response, nil
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
