package dto

import (
    "time"
    "building-report-backend/internal/domain/entity"
)

type CreateAgricultureRequest struct {
    // Data Penyuluh
    ExtensionOfficer       string    `json:"extension_officer" validate:"omitempty"`
    VisitDate              time.Time `json:"visit_date" validate:"omitempty"`
    FarmerName             string    `json:"farmer_name" validate:"omitempty"`
    FarmerGroup            string    `json:"farmer_group,omitempty"`
    FarmerGroupType        string    `json:"farmer_group_type,omitempty"` // POKTAN, GAPOKTAN
    Village                string    `json:"village" validate:"omitempty"`
    District               string    `json:"district" validate:"omitempty"`
    Latitude               float64   `json:"latitude" validate:"omitempty,min=-90,max=90"`
    Longitude              float64   `json:"longitude" validate:"omitempty,min=-180,max=180"`
    
    // Pangan (Food Crops)
    FoodCommodity          string    `json:"food_commodity,omitempty"`
    FoodLandStatus         string    `json:"food_land_status,omitempty"`
    FoodLandArea           float64   `json:"food_land_area,omitempty"`
    FoodGrowthPhase        string    `json:"food_growth_phase,omitempty"`
    FoodPlantAge           int       `json:"food_plant_age,omitempty"`
    FoodPlantingDate       string    `json:"food_planting_date,omitempty"`
    FoodHarvestDate        string    `json:"food_harvest_date,omitempty"`
    FoodDelayReason        string    `json:"food_delay_reason,omitempty"`
    FoodTechnology         string    `json:"food_technology,omitempty"`
    
    // Hortikultura (Horticulture)
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
    
    // Perkebunan (Plantation)
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
    
    // Hama dan Penyakit (Pest and Disease)
    HasPestDisease         bool      `json:"has_pest_disease"`
    PestDiseaseType        string    `json:"pest_disease_type,omitempty"`
    PestDiseaseCommodity   string    `json:"pest_disease_commodity,omitempty"` // PANGAN, HORTIKULTURA, PERKEBUNAN
    AffectedArea           string    `json:"affected_area,omitempty"`
    ControlAction          string    `json:"control_action,omitempty"`
    
    // Cuaca dan Lingkungan (Weather and Environment)
    WeatherCondition       string    `json:"weather_condition" validate:"omitempty"`
    WeatherImpact          string    `json:"weather_impact" validate:"omitempty"`
    MainConstraint         string    `json:"main_constraint" validate:"omitempty"`
    
    // Harapan dan Kebutuhan Petani (Farmer Needs and Aspirations)
    FarmerHope             string    `json:"farmer_hope" validate:"omitempty"`
    TrainingNeeded         string    `json:"training_needed" validate:"omitempty"`
    UrgentNeeds            string    `json:"urgent_needs" validate:"omitempty"`
    WaterAccess            string    `json:"water_access" validate:"omitempty"`
    Suggestions            string    `json:"suggestions,omitempty"`
}

func (r *CreateAgricultureRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateAgricultureRequest struct {
    // Data Penyuluh
    ExtensionOfficer       string  `json:"extension_officer,omitempty"`
    FarmerName             string  `json:"farmer_name,omitempty"`
    FarmerGroup            string  `json:"farmer_group,omitempty"`
    FarmerGroupType        string  `json:"farmer_group_type,omitempty"`
    Village                string  `json:"village,omitempty"`
    District               string  `json:"district,omitempty"`
    
    // Pangan (Food Crops)
    FoodCommodity          string  `json:"food_commodity,omitempty"`
    FoodLandStatus         string  `json:"food_land_status,omitempty"`
    FoodLandArea           float64 `json:"food_land_area,omitempty"`
    FoodGrowthPhase        string  `json:"food_growth_phase,omitempty"`
    FoodPlantAge           int     `json:"food_plant_age,omitempty"`
    FoodPlantingDate       string  `json:"food_planting_date,omitempty"`
    FoodHarvestDate        string  `json:"food_harvest_date,omitempty"`
    FoodDelayReason        string  `json:"food_delay_reason,omitempty"`
    FoodTechnology         string  `json:"food_technology,omitempty"`
    
    // Hortikultura (Horticulture)
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
    
    // Perkebunan (Plantation)
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
    
    // Hama dan Penyakit (Pest and Disease)
    HasPestDisease         *bool   `json:"has_pest_disease,omitempty"`
    PestDiseaseType        string  `json:"pest_disease_type,omitempty"`
    PestDiseaseCommodity   string  `json:"pest_disease_commodity,omitempty"`
    AffectedArea           string  `json:"affected_area,omitempty"`
    ControlAction          string  `json:"control_action,omitempty"`
    
    // Cuaca dan Lingkungan (Weather and Environment)
    WeatherCondition       string  `json:"weather_condition,omitempty"`
    WeatherImpact          string  `json:"weather_impact,omitempty"`
    MainConstraint         string  `json:"main_constraint,omitempty"`
    
    // Harapan dan Kebutuhan Petani (Farmer Needs and Aspirations)
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
    
    // Commodity Reports
    FoodCropReports           int64                    `json:"food_crop_reports"`
    HorticultureReports       int64                    `json:"horticulture_reports"`
    PlantationReports         int64                    `json:"plantation_reports"`
    
    // Pest and Disease
    ReportsWithPestDisease    int64                    `json:"reports_with_pest_disease"`
    PestDiseasePercentage     float64                  `json:"pest_disease_percentage"`
    
    // Problems
    PostHarvestProblemReports int64                    `json:"post_harvest_problem_reports"`
    ProductionProblemReports  int64                    `json:"production_problem_reports"`
    
    // Distributions
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
    WaterAccessDistribution   []map[string]interface{} `json:"water_access_distribution"`
    PestDiseaseByType         []map[string]interface{} `json:"pest_disease_by_type"`
    LandStatusDistribution    []map[string]interface{} `json:"land_status_distribution"`
}

type CommodityProductionResponse struct {
    Commodity     string  `json:"commodity"`
    CommodityType string  `json:"commodity_type"` // FOOD, HORTICULTURE, PLANTATION
    ReportCount   int     `json:"report_count"`
    TotalArea     float64 `json:"total_area_ha"`
    AverageArea   float64 `json:"average_area_ha"`
    FarmerCount   int     `json:"farmer_count"`
    VillageCount  int     `json:"village_count"`
}

type ExtensionOfficerPerformanceResponse struct {
    ExtensionOfficer      string    `json:"extension_officer"`
    TotalVisits           int       `json:"total_visits"`
    FarmersVisited        int       `json:"farmers_visited"`
    VillagesCovered       int       `json:"villages_covered"`
    LastVisit             time.Time `json:"last_visit"`
    AverageVisitsPerMonth float64   `json:"average_visits_per_month"`
    CommodityTypes        []string  `json:"commodity_types"`
    TotalLandAreaCovered  float64   `json:"total_land_area_covered"`
}

type VillageAgricultureSummary struct {
    Village               string  `json:"village"`
    District              string  `json:"district"`
    TotalReports          int     `json:"total_reports"`
    TotalFarmers          int     `json:"total_farmers"`
    TotalLandArea         float64 `json:"total_land_area"`
    MainCommodities       []string `json:"main_commodities"`
    PestDiseaseReports    int     `json:"pest_disease_reports"`
    MainConstraints       []string `json:"main_constraints"`
    ExtensionOfficers     []string `json:"extension_officers"`
}

type FarmerNeedsAnalysisResponse struct {
    MainConstraints       []map[string]interface{} `json:"main_constraints"`
    FarmerHopes           []map[string]interface{} `json:"farmer_hopes"`
    TrainingNeeds         []map[string]interface{} `json:"training_needs"`
    UrgentNeeds           []map[string]interface{} `json:"urgent_needs"`
    WaterAccessIssues     []map[string]interface{} `json:"water_access_issues"`
    TechnologyAdoption    []map[string]interface{} `json:"technology_adoption"`
    CommonSuggestions     []map[string]interface{} `json:"common_suggestions"`
}

type PestDiseaseAnalysisResponse struct {
    TotalReportsWithPestDisease int64                    `json:"total_reports_with_pest_disease"`
    PestDiseasePercentage       float64                  `json:"pest_disease_percentage"`
    PestDiseaseByType           []map[string]interface{} `json:"pest_disease_by_type"`
    PestDiseaseByCommodity      []map[string]interface{} `json:"pest_disease_by_commodity"`
    AffectedAreaDistribution    []map[string]interface{} `json:"affected_area_distribution"`
    ControlActionsUsed          []map[string]interface{} `json:"control_actions_used"`
    SeasonalTrends              []map[string]interface{} `json:"seasonal_trends"`
}

// Executive Dashboard Response
type AgricultureExecutiveResponse struct {
    TotalLandArea         float64                 `json:"total_land_area"`
    PestDiseaseReports    int64                   `json:"pest_disease_reports"`
    TotalExtensionReports int64                   `json:"total_extension_reports"`
    CommodityMap          []CommodityMapPoint     `json:"commodity_map"`
    CommodityBySector     CommodityBySectorData   `json:"commodity_by_sector"`
    LandStatusDistrib     []LandStatusCount       `json:"land_status_distribution"`
    MainConstraints       []ConstraintCount       `json:"main_constraints"`
    FarmerHopesNeeds      FarmerHopesNeedsData    `json:"farmer_hopes_needs"`
}

type CommodityMapPoint struct {
    Latitude      float64 `json:"latitude"`
    Longitude     float64 `json:"longitude"`
    Village       string  `json:"village"`
    District      string  `json:"district"`
    Commodity     string  `json:"commodity"`
    CommodityType string  `json:"commodity_type"`
    LandArea      float64 `json:"land_area"`
}

type CommodityBySectorData struct {
    FoodCrops    []CommodityCount `json:"food_crops"`
    Horticulture []CommodityCount `json:"horticulture"`
    Plantation   []CommodityCount `json:"plantation"`
}

type CommodityCount struct {
    Name  string `json:"name"`
    Count int64  `json:"count"`
}

type LandStatusCount struct {
    Status     string  `json:"status"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

type ConstraintCount struct {
    Constraint string  `json:"constraint"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

type FarmerHopesNeedsData struct {
    Hopes []HopeCount `json:"hopes"`
    Needs []NeedCount `json:"needs"`
}

type HopeCount struct {
    Hope       string  `json:"hope"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

type NeedCount struct {
    Need       string  `json:"need"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

// Commodity Analysis Response
type CommodityAnalysisResponse struct {
    TotalProduction       float64              `json:"total_production"`
    ProductionGrowth      float64              `json:"production_growth"`
    TotalHarvestedArea    float64              `json:"total_harvested_area"`
    HarvestedAreaGrowth   float64              `json:"harvested_area_growth"`
    Productivity          float64              `json:"productivity"`
    ProductivityGrowth    float64              `json:"productivity_growth"`
    ProductionByDistrict  []ProductionDistrict `json:"production_by_district"`
    ProductivityTrend     []ProductivityTrend  `json:"productivity_trend"`
}

type ProductionDistrict struct {
    District      string  `json:"district"`
    Production    float64 `json:"production"`
    HarvestedArea float64 `json:"harvested_area"`
    FarmerCount   int     `json:"farmer_count"`
}

type ProductivityTrend struct {
    Year         int     `json:"year"`
    Productivity float64 `json:"productivity"`
    Production   float64 `json:"production"`
    Area         float64 `json:"area"`
}

// Food Crop Response
type FoodCropResponse struct {
    LandArea           float64               `json:"land_area"`
    EstimatedProduction float64              `json:"estimated_production"`
    PestAffectedArea   float64               `json:"pest_affected_area"`
    PestReportCount    int64                 `json:"pest_report_count"`
    DistributionMap    []CommodityMapPoint   `json:"distribution_map"`
    GrowthPhases       []GrowthPhaseCount    `json:"growth_phases"`
    TechnologyUsed     []TechnologyCount     `json:"technology_used"`
    PestDominance      []PestDominanceCount  `json:"pest_dominance"`
    HarvestSchedule    []HarvestScheduleItem `json:"harvest_schedule"`
}

type GrowthPhaseCount struct {
    Phase      string  `json:"phase"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

type TechnologyCount struct {
    Technology string  `json:"technology"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

type PestDominanceCount struct {
    PestType   string  `json:"pest_type"`
    Count      int64   `json:"count"`
    Percentage float64 `json:"percentage"`
}

type HarvestScheduleItem struct {
    CommodityDetail string    `json:"commodity_detail"`
    HarvestDate     time.Time `json:"harvest_date"`
    FarmerName      string    `json:"farmer_name"`
    Village         string    `json:"village"`
    LandArea        float64   `json:"land_area"`
}

// Horticulture Response
type HorticultureResponse struct {
    LandArea           float64               `json:"land_area"`
    EstimatedProduction float64              `json:"estimated_production"`
    PestAffectedArea   float64               `json:"pest_affected_area"`
    PestReportCount    int64                 `json:"pest_report_count"`
    DistributionMap    []CommodityMapPoint   `json:"distribution_map"`
    GrowthPhases       []GrowthPhaseCount    `json:"growth_phases"`
    TechnologyUsed     []TechnologyCount     `json:"technology_used"`
    PestDominance      []PestDominanceCount  `json:"pest_dominance"`
    HarvestSchedule    []HarvestScheduleItem `json:"harvest_schedule"`
}

// Plantation Response
type PlantationResponse struct {
    LandArea           float64               `json:"land_area"`
    EstimatedProduction float64              `json:"estimated_production"`
    PestAffectedArea   float64               `json:"pest_affected_area"`
    PestReportCount    int64                 `json:"pest_report_count"`
    DistributionMap    []CommodityMapPoint   `json:"distribution_map"`
    GrowthPhases       []GrowthPhaseCount    `json:"growth_phases"`
    TechnologyUsed     []TechnologyCount     `json:"technology_used"`
    PestDominance      []PestDominanceCount  `json:"pest_dominance"`
    HarvestSchedule    []HarvestScheduleItem `json:"harvest_schedule"`
}

// Agricultural Equipment Response
type AgriculturalEquipmentResponse struct {
    GrainProcessor        EquipmentCount      `json:"grain_processor"`
    MultipurposeThresher  EquipmentCount      `json:"multipurpose_thresher"`
    FarmMachinery         EquipmentCount      `json:"farm_machinery"`
    WaterPump             EquipmentCount      `json:"water_pump"`
    DistributionByDistrict []EquipmentDistrict `json:"distribution_by_district"`
    WaterPumpTrend        []EquipmentTrend    `json:"water_pump_trend"`
}

type EquipmentCount struct {
    Count         int64   `json:"count"`
    GrowthPercent float64 `json:"growth_percent"`
}

type EquipmentDistrict struct {
    District       string `json:"district"`
    GrainProcessor int64  `json:"grain_processor"`
    Thresher       int64  `json:"thresher"`
    FarmMachinery  int64  `json:"farm_machinery"`
    WaterPump      int64  `json:"water_pump"`
}

type EquipmentTrend struct {
    Year  int   `json:"year"`
    Count int64 `json:"count"`
}

// Land and Irrigation Response
type LandIrrigationResponse struct {
    TotalLandArea        LandAreaCount          `json:"total_land_area"`
    IrrigatedLandArea    LandAreaCount          `json:"irrigated_land_area"`
    NonIrrigatedLandArea LandAreaCount          `json:"non_irrigated_land_area"`
    IrrigatedByDistrict  []LandDistrict         `json:"irrigated_by_district"`
    LandDistribution     []LandDistributionItem `json:"land_distribution"`
}

type LandAreaCount struct {
    Area          float64 `json:"area"`
    GrowthPercent float64 `json:"growth_percent"`
}

type LandDistrict struct {
    District      string  `json:"district"`
    IrrigatedArea float64 `json:"irrigated_area"`
    TotalArea     float64 `json:"total_area"`
    FarmerCount   int     `json:"farmer_count"`
}

type LandDistributionItem struct {
    District       string  `json:"district"`
    TotalArea      float64 `json:"total_area"`
    IrrigatedArea  float64 `json:"irrigated_area"`
    FoodCropArea   float64 `json:"food_crop_area"`
    HortiArea      float64 `json:"horti_area"`
    PlantationArea float64 `json:"plantation_area"`
}