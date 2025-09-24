
package entity

import (
    "time"
    "building-report-backend/pkg/utils"
)

type WaterResourcesReport struct {
    ID                      string                   `json:"id" gorm:"type:varchar(26);primary_key"`
    ReporterName           string                   `json:"reporter_name" gorm:"not null"`
    InstitutionUnit        InstitutionUnitType      `json:"institution_unit" gorm:"type:varchar(50)"`
    PhoneNumber            string                   `json:"phone_number" gorm:"type:varchar(20)"`
    ReportDateTime         time.Time                `json:"report_datetime" gorm:"not null"`
    IrrigationAreaName     string                   `json:"irrigation_area_name" gorm:"type:varchar(255)"`
    IrrigationType         IrrigationType           `json:"irrigation_type" gorm:"type:varchar(50)"`
    Latitude               float64                  `json:"latitude"`
    Longitude              float64                  `json:"longitude"`
    DamageType             DamageType               `json:"damage_type" gorm:"type:varchar(100)"`
    DamageLevel            DamageLevel              `json:"damage_level" gorm:"type:varchar(50)"`
    EstimatedLength        float64                  `json:"estimated_length" gorm:"comment:'in meters'"`
    EstimatedWidth         float64                  `json:"estimated_width" gorm:"comment:'in meters'"`
    EstimatedVolume        float64                  `json:"estimated_volume" gorm:"comment:'in m² or ha'"`
    AffectedRiceFieldArea  float64                  `json:"affected_rice_field_area" gorm:"comment:'in hectares'"`
    AffectedFarmersCount   int                      `json:"affected_farmers_count"`
    UrgencyCategory        UrgencyCategory          `json:"urgency_category" gorm:"type:varchar(50)"`
    Photos                 []WaterResourcesPhoto    `json:"photos" gorm:"foreignKey:ReportID"`
    Status                 WaterResourceStatus      `json:"status" gorm:"type:varchar(50);default:'PENDING'"`
    Notes                  string                   `json:"notes" gorm:"type:text"`
    HandlingRecommendation string                   `json:"handling_recommendation" gorm:"type:text"`
    EstimatedBudget        float64                  `json:"estimated_budget"`
    CreatedBy              string                   `json:"created_by" gorm:"type:varchar(26);not null"`
    CreatedAt              time.Time                `json:"created_at"`
    UpdatedAt              time.Time                `json:"updated_at"`
}

type WaterResourcesPhoto struct {
    ID         string    `json:"id" gorm:"type:varchar(26);primary_key"`
    ReportID   string    `json:"report_id" gorm:"type:varchar(26);not null"`
    PhotoURL   string    `json:"photo_url" gorm:"not null"`
    PhotoAngle string    `json:"photo_angle" gorm:"type:varchar(50)"` 
    Caption    string    `json:"caption" gorm:"type:varchar(255)"`
    CreatedAt  time.Time `json:"created_at"`
}


func (WaterResourcesReport) TableName() string {
    return "water_resources_reports"
}

func (WaterResourcesPhoto) TableName() string {
    return "water_resources_photos"
}


func (r *WaterResourcesReport) BeforeCreate() {
    if r.ID == "" {
        r.ID = utils.GenerateULID()
    }
    r.CreatedAt = time.Now()
    r.UpdatedAt = time.Now()
    if r.Status == "" {
        r.Status = WaterResourceStatusPending
    }
}

func (r *WaterResourcesReport) BeforeUpdate() {
    r.UpdatedAt = time.Now()
}

func (rp *WaterResourcesPhoto) BeforeCreate() {
    if rp.ID == "" {
        rp.ID = utils.GenerateULID()
    }
    rp.CreatedAt = time.Now()
}


func (r *WaterResourcesReport) CalculatePriority() int {
    priority := 0
    
    
    if r.UrgencyCategory == UrgencyCategoryMendesak {
        priority += 100
    }
    
    
    switch r.DamageLevel {
    case DamageLevelBerat:
        priority += 50
    case DamageLevelSedang:
        priority += 25
    case DamageLevelRingan:
        priority += 10
    }
    
    
    if r.AffectedRiceFieldArea > 10 {
        priority += 30
    }
    if r.AffectedFarmersCount > 50 {
        priority += 20
    }
    
    return priority
}
