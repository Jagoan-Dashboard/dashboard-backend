package dto

import "building-report-backend/pkg/utils"


func (r *CreateAgricultureRequest) Normalize() {
    r.ExtensionOfficer = utils.NormalizeLocation(r.ExtensionOfficer)
    r.FarmerName = utils.NormalizeLocation(r.FarmerName)
    r.FarmerGroup = utils.NormalizeLocation(r.FarmerGroup)
    r.Village = utils.NormalizeLocation(r.Village)
    r.District = utils.NormalizeLocation(r.District)
    
    r.FarmerGroupType = utils.NormalizeEnum(r.FarmerGroupType)
    
    r.FoodCommodity = utils.NormalizeEnum(r.FoodCommodity)
    r.FoodLandStatus = utils.NormalizeEnum(r.FoodLandStatus)
    r.FoodGrowthPhase = utils.NormalizeEnum(r.FoodGrowthPhase)
    r.FoodDelayReason = utils.NormalizeEnum(r.FoodDelayReason)
    r.FoodTechnology = utils.NormalizeEnum(r.FoodTechnology)
    
    r.HortiCommodity = utils.NormalizeEnum(r.HortiCommodity)
    r.HortiSubCommodity = utils.NormalizeEnum(r.HortiSubCommodity)
    r.HortiLandStatus = utils.NormalizeEnum(r.HortiLandStatus)
    r.HortiGrowthPhase = utils.NormalizeEnum(r.HortiGrowthPhase)
    r.HortiDelayReason = utils.NormalizeEnum(r.HortiDelayReason)
    r.HortiTechnology = utils.NormalizeEnum(r.HortiTechnology)
    r.PostHarvestProblems = utils.NormalizeEnum(r.PostHarvestProblems)
    
    r.PlantationCommodity = utils.NormalizeEnum(r.PlantationCommodity)
    r.PlantationLandStatus = utils.NormalizeEnum(r.PlantationLandStatus)
    r.PlantationGrowthPhase = utils.NormalizeEnum(r.PlantationGrowthPhase)
    r.PlantationDelayReason = utils.NormalizeEnum(r.PlantationDelayReason)
    r.PlantationTechnology = utils.NormalizeEnum(r.PlantationTechnology)
    r.ProductionProblems = utils.NormalizeEnum(r.ProductionProblems)
    
    r.PestDiseaseType = utils.NormalizeEnum(r.PestDiseaseType)
    r.PestDiseaseCommodity = utils.NormalizeEnum(r.PestDiseaseCommodity)
    r.AffectedArea = utils.NormalizeEnum(r.AffectedArea)
    r.ControlAction = utils.NormalizeEnum(r.ControlAction)
    
    r.WeatherCondition = utils.NormalizeEnum(r.WeatherCondition)
    r.WeatherImpact = utils.NormalizeEnum(r.WeatherImpact)
    r.MainConstraint = utils.NormalizeEnum(r.MainConstraint)
    
    r.FarmerHope = utils.NormalizeEnum(r.FarmerHope)
    r.TrainingNeeded = utils.NormalizeEnum(r.TrainingNeeded)
    r.UrgentNeeds = utils.NormalizeEnum(r.UrgentNeeds)
    r.WaterAccess = utils.NormalizeEnum(r.WaterAccess)
}

func (r *UpdateAgricultureRequest) Normalize() {
	r.ExtensionOfficer = utils.NormalizeLocation(r.ExtensionOfficer)
	r.FarmerName = utils.NormalizeLocation(r.FarmerName)
	r.FarmerGroup = utils.NormalizeLocation(r.FarmerGroup)
	r.Village = utils.NormalizeLocation(r.Village)
	r.District = utils.NormalizeLocation(r.District)
	
	r.FarmerGroupType = utils.NormalizeEnum(r.FarmerGroupType)
	r.FoodCommodity = utils.NormalizeEnum(r.FoodCommodity)
	r.FoodLandStatus = utils.NormalizeEnum(r.FoodLandStatus)
	r.FoodGrowthPhase = utils.NormalizeEnum(r.FoodGrowthPhase)
	r.FoodDelayReason = utils.NormalizeEnum(r.FoodDelayReason)
	r.FoodTechnology = utils.NormalizeEnum(r.FoodTechnology)
	
	r.WeatherCondition = utils.NormalizeEnum(r.WeatherCondition)
	r.WeatherImpact = utils.NormalizeEnum(r.WeatherImpact)
	r.MainConstraint = utils.NormalizeEnum(r.MainConstraint)
	r.FarmerHope = utils.NormalizeEnum(r.FarmerHope)
	r.TrainingNeeded = utils.NormalizeEnum(r.TrainingNeeded)
	r.UrgentNeeds = utils.NormalizeEnum(r.UrgentNeeds)
	r.WaterAccess = utils.NormalizeEnum(r.WaterAccess)
}


func (r *CreateBinaMargaRequest) Normalize() {
	r.ReporterName = utils.NormalizeLocation(r.ReporterName)
	r.RoadName = utils.NormalizeLocation(r.RoadName)
	r.BridgeName = utils.NormalizeLocation(r.BridgeName)
	
	r.InstitutionUnit = utils.NormalizeEnum(r.InstitutionUnit)
	// r.RoadType = utils.NormalizeEnum(r.RoadType)
	// r.RoadClass = utils.NormalizeEnum(r.RoadClass)
	r.PavementType = utils.NormalizeEnum(r.PavementType)
	r.DamageType = utils.NormalizeEnum(r.DamageType)
	r.DamageLevel = utils.NormalizeEnum(r.DamageLevel)
	r.BridgeStructureType = utils.NormalizeEnum(r.BridgeStructureType)
	r.BridgeDamageType = utils.NormalizeEnum(r.BridgeDamageType)
	r.BridgeDamageLevel = utils.NormalizeEnum(r.BridgeDamageLevel)
	r.TrafficCondition = utils.NormalizeEnum(r.TrafficCondition)
	r.TrafficImpact = utils.NormalizeEnum(r.TrafficImpact)
	r.UrgencyLevel = utils.NormalizeEnum(r.UrgencyLevel)
}

func (r *UpdateBinaMargaRequest) Normalize() {
	r.RoadName = utils.NormalizeLocation(r.RoadName)
	r.BridgeName = utils.NormalizeLocation(r.BridgeName)
	
	// r.RoadType = utils.NormalizeEnum(r.RoadType)
	// r.RoadClass = utils.NormalizeEnum(r.RoadClass)
	r.PavementType = utils.NormalizeEnum(r.PavementType)
	r.DamageType = utils.NormalizeEnum(r.DamageType)
	r.DamageLevel = utils.NormalizeEnum(r.DamageLevel)
	r.BridgeStructureType = utils.NormalizeEnum(r.BridgeStructureType)
	r.BridgeDamageType = utils.NormalizeEnum(r.BridgeDamageType)
	r.BridgeDamageLevel = utils.NormalizeEnum(r.BridgeDamageLevel)
	r.TrafficCondition = utils.NormalizeEnum(r.TrafficCondition)
	r.TrafficImpact = utils.NormalizeEnum(r.TrafficImpact)
	r.UrgencyLevel = utils.NormalizeEnum(r.UrgencyLevel)
}


func (r *CreateWaterResourcesRequest) Normalize() {
	r.ReporterName = utils.NormalizeLocation(r.ReporterName)
	r.IrrigationAreaName = utils.NormalizeLocation(r.IrrigationAreaName)
	
	r.InstitutionUnit = utils.NormalizeEnum(r.InstitutionUnit)
	r.IrrigationType = utils.NormalizeEnum(r.IrrigationType)
	r.DamageType = utils.NormalizeEnum(r.DamageType)
	r.DamageLevel = utils.NormalizeEnum(r.DamageLevel)
	r.UrgencyCategory = utils.NormalizeEnum(r.UrgencyCategory)
}

func (r *UpdateWaterResourcesRequest) Normalize() {
	r.IrrigationAreaName = utils.NormalizeLocation(r.IrrigationAreaName)
	
	r.IrrigationType = utils.NormalizeEnum(r.IrrigationType)
	r.DamageType = utils.NormalizeEnum(r.DamageType)
	r.DamageLevel = utils.NormalizeEnum(r.DamageLevel)
	r.UrgencyCategory = utils.NormalizeEnum(r.UrgencyCategory)
}


func (r *CreateSpatialPlanningRequest) Normalize() {
	r.ReporterName = utils.NormalizeLocation(r.ReporterName)
	
	r.Institution = utils.NormalizeEnum(r.Institution)
	r.AreaCategory = utils.NormalizeEnum(r.AreaCategory)
	r.ViolationType = utils.NormalizeEnum(r.ViolationType)
	r.ViolationLevel = utils.NormalizeEnum(r.ViolationLevel)
	r.EnvironmentalImpact = utils.NormalizeEnum(r.EnvironmentalImpact)
	r.UrgencyLevel = utils.NormalizeEnum(r.UrgencyLevel)
}

func (r *UpdateSpatialPlanningRequest) Normalize() {
	r.AreaCategory = utils.NormalizeEnum(r.AreaCategory)
	r.ViolationType = utils.NormalizeEnum(r.ViolationType)
	r.ViolationLevel = utils.NormalizeEnum(r.ViolationLevel)
	r.EnvironmentalImpact = utils.NormalizeEnum(r.EnvironmentalImpact)
	r.UrgencyLevel = utils.NormalizeEnum(r.UrgencyLevel)
	r.Status = utils.NormalizeEnum(r.Status)
}


func (r *CreateReportRequest) Normalize() {
	r.ReporterName = utils.NormalizeLocation(r.ReporterName)
	r.Village = utils.NormalizeLocation(r.Village)
	r.District = utils.NormalizeLocation(r.District)
	r.BuildingName = utils.NormalizeLocation(r.BuildingName)
	
	r.ReporterRole = utils.NormalizeEnum(r.ReporterRole)
	r.BuildingType = utils.NormalizeEnum(r.BuildingType)
	r.ReportStatus = utils.NormalizeEnum(r.ReportStatus)
	r.FundingSource = utils.NormalizeEnum(r.FundingSource)
	r.WorkType = utils.NormalizeEnum(r.WorkType)
	r.ConditionAfterRehab = utils.NormalizeEnum(r.ConditionAfterRehab)
}

func (r *UpdateReportRequest) Normalize() {
	r.BuildingName = utils.NormalizeLocation(r.BuildingName)
	r.BuildingType = utils.NormalizeEnum(r.BuildingType)
	r.ReportStatus = utils.NormalizeEnum(r.ReportStatus)
	r.FundingSource = utils.NormalizeEnum(r.FundingSource)
	r.WorkType = utils.NormalizeEnum(r.WorkType)
	r.ConditionAfterRehab = utils.NormalizeEnum(r.ConditionAfterRehab)
}