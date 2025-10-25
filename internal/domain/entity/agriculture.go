package entity

import (
	"building-report-backend/pkg/utils"
	"time"
)

type AgricultureReport struct {
	ID               string             `json:"id" gorm:"type:varchar(26);primary_key"`
	ExtensionOfficer string             `json:"extension_officer" gorm:"not null"`
	VisitDate        time.Time          `json:"visit_date" gorm:"not null"`
	FarmerName       string             `json:"farmer_name" gorm:"not null"`
	FarmerGroup      string             `json:"farmer_group"`
	FarmerGroupType  FarmerGroupType    `json:"farmer_group_type" gorm:"type:varchar(50)"`
	Village          string             `json:"village" gorm:"not null"`
	District         string             `json:"district" gorm:"not null"`
	Latitude         float64            `json:"latitude"`
	Longitude        float64            `json:"longitude"`
	Photos           []AgriculturePhoto `json:"photos" gorm:"foreignKey:ReportID"`

	FoodCommodity    FoodCommodity    `json:"food_commodity,omitempty" gorm:"type:varchar(50)"`
	FoodLandStatus   LandStatus       `json:"food_land_status,omitempty" gorm:"type:varchar(50)"`
	FoodLandArea     float64          `json:"food_land_area,omitempty" gorm:"comment:'in hectares'"`
	FoodGrowthPhase  GrowthPhase      `json:"food_growth_phase,omitempty" gorm:"type:varchar(50)"`
	FoodPlantAge     int              `json:"food_plant_age,omitempty" gorm:"comment:'in days'"`
	FoodPlantingDate *time.Time       `json:"food_planting_date,omitempty"`
	FoodHarvestDate  *time.Time       `json:"food_harvest_date,omitempty"`
	FoodDelayReason  DelayReason      `json:"food_delay_reason,omitempty" gorm:"type:varchar(100)"`
	FoodTechnology   TechnologyMethod `json:"food_technology,omitempty" gorm:"type:varchar(100)"`

	HortiCommodity      HorticultureCommodity `json:"horti_commodity,omitempty" gorm:"type:varchar(50)"`
	HortiSubCommodity   string                `json:"horti_sub_commodity,omitempty" gorm:"type:varchar(100)"`
	HortiLandStatus     LandStatus            `json:"horti_land_status,omitempty" gorm:"type:varchar(50)"`
	HortiLandArea       float64               `json:"horti_land_area,omitempty" gorm:"comment:'in hectares'"`
	HortiGrowthPhase    HortiGrowthPhase      `json:"horti_growth_phase,omitempty" gorm:"type:varchar(50)"`
	HortiPlantAge       int                   `json:"horti_plant_age,omitempty" gorm:"comment:'in days'"`
	HortiPlantingDate   *time.Time            `json:"horti_planting_date,omitempty"`
	HortiHarvestDate    *time.Time            `json:"horti_harvest_date,omitempty"`
	HortiDelayReason    DelayReason           `json:"horti_delay_reason,omitempty" gorm:"type:varchar(100)"`
	HortiTechnology     HortiTechnology       `json:"horti_technology,omitempty" gorm:"type:varchar(100)"`
	PostHarvestProblems PostHarvestProblem    `json:"post_harvest_problems,omitempty" gorm:"type:varchar(100)"`

	PlantationCommodity    PlantationCommodity   `json:"plantation_commodity,omitempty" gorm:"type:varchar(50)"`
	PlantationLandStatus   LandStatus            `json:"plantation_land_status,omitempty" gorm:"type:varchar(50)"`
	PlantationLandArea     float64               `json:"plantation_land_area,omitempty" gorm:"comment:'in hectares'"`
	PlantationGrowthPhase  PlantationGrowthPhase `json:"plantation_growth_phase,omitempty" gorm:"type:varchar(50)"`
	PlantationPlantAge     int                   `json:"plantation_plant_age,omitempty" gorm:"comment:'in days'"`
	PlantationPlantingDate *time.Time            `json:"plantation_planting_date,omitempty"`
	PlantationHarvestDate  *time.Time            `json:"plantation_harvest_date,omitempty"`
	PlantationDelayReason  DelayReason           `json:"plantation_delay_reason,omitempty" gorm:"type:varchar(100)"`
	PlantationTechnology   PlantationTechnology  `json:"plantation_technology,omitempty" gorm:"type:varchar(100)"`
	ProductionProblems     ProductionProblem     `json:"production_problems,omitempty" gorm:"type:varchar(100)"`

	HasPestDisease       bool              `json:"has_pest_disease"`
	PestDiseaseType      PestDiseaseType   `json:"pest_disease_type,omitempty" gorm:"type:varchar(100)"`
	PestDiseaseCommodity string            `json:"pest_disease_commodity,omitempty" gorm:"type:varchar(50)"`
	AffectedArea         AffectedAreaLevel `json:"affected_area,omitempty" gorm:"type:varchar(50)"`
	ControlAction        ControlAction     `json:"control_action,omitempty" gorm:"type:varchar(100)"`

	WeatherCondition WeatherCondition `json:"weather_condition" gorm:"type:varchar(50)"`
	WeatherImpact    WeatherImpact    `json:"weather_impact" gorm:"type:varchar(50)"`
	MainConstraint   MainConstraint   `json:"main_constraint" gorm:"type:varchar(50)"`

	FarmerHope     FarmerHope     `json:"farmer_hope" gorm:"type:varchar(100)"`
	TrainingNeeded TrainingNeeded `json:"training_needed" gorm:"type:varchar(100)"`
	UrgentNeeds    UrgentNeeds    `json:"urgent_needs" gorm:"type:varchar(100)"`
	WaterAccess    WaterAccess    `json:"water_access" gorm:"type:varchar(50)"`
	Suggestions    string         `json:"suggestions" gorm:"type:text"`

	// CreatedBy              string              `json:"created_by" gorm:"type:varchar(26);not null"
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type AgriculturePhoto struct {
	ID        string    `json:"id" gorm:"type:varchar(26);primary_key"`
	ReportID  string    `json:"report_id" gorm:"type:varchar(26);not null"`
	PhotoURL  string    `json:"photo_url" gorm:"not null;size:500"`
	PhotoType string    `json:"photo_type" gorm:"type:varchar(50)"`
	Caption   string    `json:"caption" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

func (AgricultureReport) TableName() string {
	return "agriculture_reports"
}

func (AgriculturePhoto) TableName() string {
	return "agriculture_photos"
}

func (r *AgricultureReport) BeforeCreate() {
	if r.ID == "" {
		r.ID = utils.GenerateULID()
	}
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
}

func (r *AgricultureReport) BeforeUpdate() {
	r.UpdatedAt = time.Now()
}

func (ap *AgriculturePhoto) BeforeCreate() {
	if ap.ID == "" {
		ap.ID = utils.GenerateULID()
	}
	ap.CreatedAt = time.Now()
}
