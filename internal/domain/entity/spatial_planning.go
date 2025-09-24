
package entity

import (
    "time"
    "building-report-backend/pkg/utils"
)

type SpatialPlanningReport struct {
    ID                      string                       `json:"id" gorm:"type:varchar(26);primary_key"`
    ReporterName           string                       `json:"reporter_name" gorm:"not null"`
    Institution            InstitutionType              `json:"institution" gorm:"type:varchar(50)"`
    PhoneNumber            string                       `json:"phone_number" gorm:"type:varchar(20)"`
    ReportDateTime         time.Time                    `json:"report_datetime" gorm:"not null"`
    AreaDescription        string                       `json:"area_description" gorm:"type:text"`
    AreaCategory           AreaCategory                 `json:"area_category" gorm:"type:varchar(100)"`
    ViolationType          SpatialViolationType         `json:"violation_type" gorm:"type:varchar(150)"`
    ViolationLevel         ViolationLevel               `json:"violation_level" gorm:"type:varchar(50)"`
    EnvironmentalImpact    EnvironmentalImpact         `json:"environmental_impact" gorm:"type:varchar(100)"`
    UrgencyLevel           UrgencyLevel                 `json:"urgency_level" gorm:"type:varchar(20)"`
    Latitude               float64                      `json:"latitude"`
    Longitude              float64                      `json:"longitude"`
    Address                string                       `json:"address" gorm:"type:text"`
    Photos                 []SpatialPlanningPhoto       `json:"photos" gorm:"foreignKey:ReportID"`
    Status                 SpatialReportStatus          `json:"status" gorm:"type:varchar(50);default:'PENDING'"`
    Notes                  string                       `json:"notes" gorm:"type:text"`
    CreatedBy              string                       `json:"created_by" gorm:"type:varchar(26);not null"`
    CreatedAt              time.Time                    `json:"created_at"`
    UpdatedAt              time.Time                    `json:"updated_at"`
}

type SpatialPlanningPhoto struct {
    ID         string    `json:"id" gorm:"type:varchar(26);primary_key"`
    ReportID   string    `json:"report_id" gorm:"type:varchar(26);not null"`
    PhotoURL   string    `json:"photo_url" gorm:"not null"`
    Caption    string    `json:"caption" gorm:"type:varchar(255)"`
    CreatedAt  time.Time `json:"created_at"`
}


func (SpatialPlanningReport) TableName() string {
    return "spatial_planning_reports"
}

func (SpatialPlanningPhoto) TableName() string {
    return "spatial_planning_photos"
}


func (r *SpatialPlanningReport) BeforeCreate() {
    if r.ID == "" {
        r.ID = utils.GenerateULID()
    }
    r.CreatedAt = time.Now()
    r.UpdatedAt = time.Now()
}

func (r *SpatialPlanningReport) BeforeUpdate() {
    r.UpdatedAt = time.Now()
}

func (rp *SpatialPlanningPhoto) BeforeCreate() {
    if rp.ID == "" {
        rp.ID = utils.GenerateULID()
    }
    rp.CreatedAt = time.Now()
}