
package dto

import (
    "time"
    "building-report-backend/internal/domain/entity"
)

type CreateAgricultureRequest struct {
    ExtensionOfficer       string    `json:"extension_officer" validate:"required"`
    VisitDate              time.Time `json:"visit_date" validate:"required"`
    FarmerName             string    `json:"farmer_name" validate:"required"`
    FarmerGroup            string    `json:"farmer_group,omitempty"`
    FarmerGroupType        string    `json:"farmer_group_type,omitempty"`
    Village                string    `json:"village" validate:"required"`
    District               string    `json:"district" validate:"required"`
    Latitude               float64   `json:"latitude" validate:"required,min=-90,max=90"`
    Longitude              float64   `json:"longitude" validate:"required,min=-180,max=180"`
    
    
    FoodCommodity          string    `json:"food_commodity,omitempty"`
    FoodLandStatus         string    `json:"food_land_status,omitempty"`
    FoodLandArea           float64   `json:"food_land_area,omitempty"`
    FoodGrowthPhase        string    `json:"food_growth_phase,omitempty"`
    FoodPlantAge           int       `json:"food_plant_age,omitempty"`
    FoodPlantingDate       string    `json:"food_planting_date,omitempty"`
    FoodHarvestDate        string    `json:"food_harvest_date,omitempty"`
    FoodDelayReason        string    `json:"food_delay_reason,omitempty"`
    FoodTechnology         string    `json:"food_technology,omitempty"`
    
    
    HortiCommodity         string    `json:"horti_commodity,omitempty"`
    HortiSubCommodity      string    `json:"horti_sub_commodity,omitempty"`
    HortiLandStatus        string    `json:"horti_land_status,omitempty"`
    HortiLandArea          float64   `json:"horti_land_area,omitempty"`
    HortiGrowthPhase       string    `json:"horti_growth_phase,omitempty"`
    HortiPlantAge          int       `json:"horti_plant_age,omitempty"`
    HortiPlantingDate      string    `json:"horti_planting_date,omitempty"`
    HortiHarvestDate       string    `json:"horti_harvest_date,omitempty"`
    HortiDelayReason       string    `json:"horti_delay_reason,omitempty"`
    HortiTechnology        string    `json:"horti_technology,omitempty"`
    PostHarvestProblems    string    `json:"post_harvest_problems,omitempty"`
    
    
    PlantationCommodity    string    `json:"plantation_commodity,omitempty"`
    PlantationLandStatus   string    `json:"plantation_land_status,omitempty"`
    PlantationLandArea     float64   `json:"plantation_land_area,omitempty"`
    PlantationGrowthPhase  string    `json:"plantation_growth_phase,omitempty"`
    PlantationPlantAge     int       `json:"plantation_plant_age,omitempty"`
    PlantationPlantingDate string    `json:"plantation_planting_date,omitempty"`
    PlantationHarvestDate  string    `json:"plantation_harvest_date,omitempty"`
    PlantationDelayReason  string    `json:"plantation_delay_reason,omitempty"`
    PlantationTechnology   string    `json:"plantation_technology,omitempty"`
    ProductionProblems     string    `json:"production_problems,omitempty"`
    
    
    HasPestDisease         bool      `json:"has_pest_disease"`
    PestDiseaseType        string    `json:"pest_disease_type,omitempty"`
    PestDiseaseCommodity   string    `json:"pest_disease_commodity,omitempty"`
    AffectedArea           string    `json:"affected_area,omitempty"`
    ControlAction          string    `json:"control_action,omitempty"`
    
    
    WeatherCondition       string    `json:"weather_condition" validate:"required"`
    WeatherImpact          string    `json:"weather_impact" validate:"required"`
    MainConstraint         string    `json:"main_constraint" validate:"required"`
    
    
    FarmerHope             string    `json:"farmer_hope" validate:"required"`
    TrainingNeeded         string    `json:"training_needed" validate:"required"`
    UrgentNeeds            string    `json:"urgent_needs" validate:"required"`
    WaterAccess            string    `json:"water_access" validate:"required"`
    Suggestions            string    `json:"suggestions,omitempty"`
}

func (r *CreateAgricultureRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateAgricultureRequest struct {
    ExtensionOfficer       string  `json:"extension_officer,omitempty"`
    FarmerName             string  `json:"farmer_name,omitempty"`
    FarmerGroup            string  `json:"farmer_group,omitempty"`
    FarmerGroupType        string  `json:"farmer_group_type,omitempty"`
    Village                string  `json:"village,omitempty"`
    District               string  `json:"district,omitempty"`
    
    
    FoodCommodity          string  `json:"food_commodity,omitempty"`
    FoodLandStatus         string  `json:"food_land_status,omitempty"`
    FoodLandArea           float64 `json:"food_land_area,omitempty"`
    FoodGrowthPhase        string  `json:"food_growth_phase,omitempty"`
    FoodPlantAge           int     `json:"food_plant_age,omitempty"`
    FoodPlantingDate       string  `json:"food_planting_date,omitempty"`
    FoodHarvestDate        string  `json:"food_harvest_date,omitempty"`
    FoodDelayReason        string  `json:"food_delay_reason,omitempty"`
    FoodTechnology         string  `json:"food_technology,omitempty"`
    
    
    HortiCommodity         string  `json:"horti_commodity,omitempty"`
    HortiSubCommodity      string  `json:"horti_sub_commodity,omitempty"`
    HortiLandStatus        string  `json:"horti_land_status,omitempty"`
    HortiLandArea          float64 `json:"horti_land_area,omitempty"`
    HortiGrowthPhase       string  `json:"horti_growth_phase,omitempty"`
    HortiPlantAge          int     `json:"horti_plant_age,omitempty"`
    HortiPlantingDate      string  `json:"horti_planting_date,omitempty"`
    HortiHarvestDate       string  `json:"horti_harvest_date,omitempty"`
    HortiDelayReason       string  `json:"horti_delay_reason,omitempty"`
    HortiTechnology        string  `json:"horti_technology,omitempty"`
    PostHarvestProblems    string  `json:"post_harvest_problems,omitempty"`
    
    
    PlantationCommodity    string  `json:"plantation_commodity,omitempty"`
    PlantationLandStatus   string  `json:"plantation_land_status,omitempty"`
    PlantationLandArea     float64 `json:"plantation_land_area,omitempty"`
    PlantationGrowthPhase  string  `json:"plantation_growth_phase,omitempty"`
    PlantationPlantAge     int     `json:"plantation_plant_age,omitempty"`
    PlantationPlantingDate string  `json:"plantation_planting_date,omitempty"`
    PlantationHarvestDate  string  `json:"plantation_harvest_date,omitempty"`
    PlantationDelayReason  string  `json:"plantation_delay_reason,omitempty"`
    PlantationTechnology   string  `json:"plantation_technology,omitempty"`
    ProductionProblems     string  `json:"production_problems,omitempty"`
    
    
    HasPestDisease         *bool   `json:"has_pest_disease,omitempty"`
    PestDiseaseType        string  `json:"pest_disease_type,omitempty"`
    PestDiseaseCommodity   string  `json:"pest_disease_commodity,omitempty"`
    AffectedArea           string  `json:"affected_area,omitempty"`
    ControlAction          string  `json:"control_action,omitempty"`
    
    
    WeatherCondition       string  `json:"weather_condition,omitempty"`
    WeatherImpact          string  `json:"weather_impact,omitempty"`
    MainConstraint         string  `json:"main_constraint,omitempty"`
    
    
    FarmerHope             string  `json:"farmer_hope,omitempty"`
    TrainingNeeded         string  `json:"training_needed,omitempty"`
    UrgentNeeds            string  `json:"urgent_needs,omitempty"`
    WaterAccess            string  `json:"water_access,omitempty"`
    Suggestions            string  `json:"suggestions,omitempty"`
}

func (r *UpdateAgricultureRequest) Validate() error {
    return validate.Struct(r)
}

type PaginatedAgricultureResponse struct {
    Reports    []*entity.AgricultureReport `json:"reports"`
    Total      int64                       `json:"total"`
    Page       int                         `json:"page"`
    PerPage    int                         `json:"per_page"`
    TotalPages int64                       `json:"total_pages"`
}

type AgricultureStatisticsResponse struct {
    TotalReports              int64                    `json:"total_reports"`
    TotalFarmers              int64                    `json:"total_farmers"`
    TotalLandArea             float64                  `json:"total_land_area_ha"`
    
    
    FoodCropReports           int64                    `json:"food_crop_reports"`
    HorticultureReports       int64                    `json:"horticulture_reports"`
    PlantationReports         int64                    `json:"plantation_reports"`
    
    
    ReportsWithPestDisease    int64                    `json:"reports_with_pest_disease"`
    PestDiseasePercentage     float64                  `json:"pest_disease_percentage"`
    
    
    PostHarvestProblemReports int64                    `json:"post_harvest_problem_reports"`
    ProductionProblemReports  int64                    `json:"production_problem_reports"`
    
    
    CommodityDistribution     []map[string]interface{} `json:"commodity_distribution"`
    VillageDistribution       []map[string]interface{} `json:"village_distribution"`
    ExtensionOfficerStats     []map[string]interface{} `json:"extension_officer_stats"`
    GrowthPhaseDistribution   []map[string]interface{} `json:"growth_phase_distribution"`
    TechnologyAdoption        []map[string]interface{} `json:"technology_adoption"`
    MainConstraints           []map[string]interface{} `json:"main_constraints"`
    FarmerHopes               []map[string]interface{} `json:"farmer_hopes"`
    TrainingNeeds             []map[string]interface{} `json:"training_needs"`
    UrgentNeedsDistribution   []map[string]interface{} `json:"urgent_needs_distribution"`
    WeatherImpactCounts       []map[string]interface{} `json:"weather_impact_counts"`
}

type CommodityProductionResponse struct {
    Commodity     string  `json:"commodity"`
    ReportCount   int     `json:"report_count"`
    TotalArea     float64 `json:"total_area_ha"`
    AverageArea   float64 `json:"average_area_ha"`
    FarmerCount   int     `json:"farmer_count"`
    VillageCount  int     `json:"village_count"`
}

type ExtensionOfficerPerformanceResponse struct {
    ExtensionOfficer string    `json:"extension_officer"`
    TotalVisits      int       `json:"total_visits"`
    FarmersVisited   int       `json:"farmers_visited"`
    VillagesCovered  int       `json:"villages_covered"`
    LastVisit        time.Time `json:"last_visit"`
    AverageVisitsPerMonth float64 `json:"average_visits_per_month"`
}