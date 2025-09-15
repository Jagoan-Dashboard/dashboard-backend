
package entity

import (
    "time"
    "github.com/google/uuid"
)

type BinaMargaReport struct {
    ID                     uuid.UUID              `json:"id" gorm:"type:uuid;primary_key"`
    ReporterName          string                 `json:"reporter_name" gorm:"not null"`
    InstitutionUnit       InstitutionUnitType    `json:"institution_unit" gorm:"type:varchar(50)"`
    PhoneNumber           string                 `json:"phone_number" gorm:"type:varchar(20)"`
    ReportDateTime        time.Time              `json:"report_datetime" gorm:"not null"`
    RoadName              string                 `json:"road_name" gorm:"type:varchar(255)"`
    RoadType              RoadType               `json:"road_type" gorm:"type:varchar(50)"`
    RoadClass             RoadClass              `json:"road_class" gorm:"type:varchar(50)"`
    Latitude              float64                `json:"latitude"`
    Longitude             float64                `json:"longitude"`
    DamageType            RoadDamageType         `json:"damage_type" gorm:"type:varchar(100)"`
    DamageLevel           RoadDamageLevel        `json:"damage_level" gorm:"type:varchar(50)"`
    DamagedLength         float64                `json:"damaged_length" gorm:"comment:'in meters'"`
    DamagedWidth          float64                `json:"damaged_width" gorm:"comment:'in meters'"`
    DamagedArea           float64                `json:"damaged_area" gorm:"comment:'in square meters'"`
    TrafficImpact         TrafficImpact          `json:"traffic_impact" gorm:"type:varchar(50)"`
    UrgencyLevel          RoadUrgencyLevel       `json:"urgency_level" gorm:"type:varchar(50)"`
    CauseOfDamage         string                 `json:"cause_of_damage" gorm:"type:text"`
    Photos                []BinaMargaPhoto       `json:"photos" gorm:"foreignKey:ReportID"`
    Status                BinaMargaStatus        `json:"status" gorm:"type:varchar(50);default:'PENDING'"`
    Notes                 string                 `json:"notes" gorm:"type:text"`
    HandlingRecommendation string                `json:"handling_recommendation" gorm:"type:text"`
    EstimatedBudget       float64                `json:"estimated_budget"`
    EstimatedRepairTime   int                    `json:"estimated_repair_time" gorm:"comment:'in days'"`
    CreatedBy             uuid.UUID              `json:"created_by" gorm:"type:uuid"`
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
}

type BinaMargaPhoto struct {
    ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    ReportID   uuid.UUID `json:"report_id" gorm:"type:uuid;not null"`
    PhotoURL   string    `json:"photo_url" gorm:"not null"`
    PhotoAngle string    `json:"photo_angle" gorm:"type:varchar(50)"` 
    Caption    string    `json:"caption" gorm:"type:varchar(255)"`
    CreatedAt  time.Time `json:"created_at"`
}


func (BinaMargaReport) TableName() string {
    return "bina_marga_reports"
}

func (BinaMargaPhoto) TableName() string {
    return "bina_marga_photos"
}


func (r *BinaMargaReport) BeforeCreate() {
    if r.ID == uuid.Nil {
        r.ID = uuid.New()
    }
    r.CreatedAt = time.Now()
    r.UpdatedAt = time.Now()
    if r.Status == "" {
        r.Status = BinaMargaStatusPending
    }
}

func (rp *BinaMargaPhoto) BeforeCreate() {
    if rp.ID == uuid.Nil {
        rp.ID = uuid.New()
    }
    rp.CreatedAt = time.Now()
}


func (r *BinaMargaReport) CalculatePriority() int {
    priority := 0
    
    
    switch r.UrgencyLevel {
    case RoadUrgencyEmergency:
        priority += 100
    case RoadUrgencyHigh:
        priority += 75
    case RoadUrgencyMedium:
        priority += 50
    case RoadUrgencyLow:
        priority += 25
    }
    
    
    switch r.DamageLevel {
    case RoadDamageLevelSevere:
        priority += 50
    case RoadDamageLevelModerate:
        priority += 30
    case RoadDamageLevelMinor:
        priority += 15
    }
    
    
    switch r.RoadClass {
    case RoadClassArteri:
        priority += 40
    case RoadClassKolektor:
        priority += 30
    case RoadClassLokal:
        priority += 20
    case RoadClassLingkungan:
        priority += 10
    }
    
    
    switch r.TrafficImpact {
    case TrafficImpactBlocked:
        priority += 60
    case TrafficImpactSeverelyReduced:
        priority += 40
    case TrafficImpactReduced:
        priority += 20
    case TrafficImpactMinimal:
        priority += 5
    }
    
    
    if r.DamagedArea > 100 {
        priority += 25
    } else if r.DamagedArea > 50 {
        priority += 15
    }
    
    return priority
}