
package repository

import (
    "context"
    "time"
    "building-report-backend/internal/domain/entity"
)

type AgricultureRepository interface {
    Create(ctx context.Context, report *entity.AgricultureReport) error
    Update(ctx context.Context, report *entity.AgricultureReport) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.AgricultureReport, error)
    FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.AgricultureReport, int64, error)
    FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    FindByExtensionOfficer(ctx context.Context, extensionOfficer string, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    FindByVillage(ctx context.Context, village string, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    FindByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.AgricultureReport, int64, error)
    
    GetStatistics(ctx context.Context) (map[string]interface{}, error)
    GetCommodityProduction(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetExtensionOfficerPerformance(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetVillageProductionStats(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetPestDiseaseReports(ctx context.Context, limit int) ([]*entity.AgricultureReport, error)
    GetTechnologyAdoptionStats(ctx context.Context) (map[string]interface{}, error)
    GetFarmerNeedsAnalysis(ctx context.Context) (map[string]interface{}, error)
    CountTotalFarmers(ctx context.Context) (int64, error)
    CalculateTotalLandArea(ctx context.Context) (float64, error)
    CountReportsByCommodityType(ctx context.Context) (map[string]int64, error)
    CountReportsWithPestDisease(ctx context.Context) (int64, error)
    GetTopConstraints(ctx context.Context, limit int) ([]map[string]interface{}, error)
    GetTopFarmerHopes(ctx context.Context, limit int) ([]map[string]interface{}, error)
    
    // Executive Dashboard Methods
     GetExecutiveSummary(ctx context.Context, commodityType string) (map[string]interface{}, error)
    GetCommodityDistributionByDistrict(ctx context.Context, commodityType string) ([]map[string]interface{}, error)
    GetCommodityCountBySector(ctx context.Context, commodityType string) (map[string]interface{}, error)
    GetLandStatusDistribution(ctx context.Context, commodityType string) ([]map[string]interface{}, error)
    GetMainConstraintsDistribution(ctx context.Context, commodityType string) ([]map[string]interface{}, error)
    GetFarmerHopesAndNeeds(ctx context.Context, commodityType string) (map[string]interface{}, error)
    
    
    // Commodity Analysis Methods
    GetCommodityAnalysis(ctx context.Context, startDate, endDate time.Time, commodityName string) (map[string]interface{}, error)
    GetProductionByDistrict(ctx context.Context, startDate, endDate time.Time, commodityName string) ([]map[string]interface{}, error)
    GetProductivityTrend(ctx context.Context, commodityName string, years []int) ([]map[string]interface{}, error)
    
    // Food Crop Specific Methods
    GetFoodCropStats(ctx context.Context, commodityName string) (map[string]interface{}, error)
    GetFoodCropDistribution(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetFoodCropGrowthPhases(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetFoodCropTechnology(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetFoodCropPestDominance(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetFoodCropHarvestSchedule(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    
    // Horticulture Specific Methods
    GetHorticultureStats(ctx context.Context, commodityName string) (map[string]interface{}, error)
    GetHorticultureDistribution(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetHorticultureGrowthPhases(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetHorticultureTechnology(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetHorticulturePestDominance(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetHorticultureHarvestSchedule(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    
    // Plantation Specific Methods
    GetPlantationStats(ctx context.Context, commodityName string) (map[string]interface{}, error)
    GetPlantationDistribution(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetPlantationGrowthPhases(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetPlantationTechnology(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetPlantationPestDominance(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    GetPlantationHarvestSchedule(ctx context.Context, commodityName string) ([]map[string]interface{}, error)
    
    // Agricultural Equipment Methods
    GetAgriculturalEquipmentStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
    GetEquipmentDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
    GetEquipmentTrend(ctx context.Context, equipmentType string, years []int) ([]map[string]interface{}, error)
    
    // Land and Irrigation Methods
    GetLandAndIrrigationStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error)
    GetLandDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error)
}