
package entity

import (
    "time"
    "building-report-backend/pkg/utils"
)

type Report struct {
    ID                     string                 `json:"id" gorm:"type:varchar(26);primary_key"`
    ReporterName          string                 `json:"reporter_name" gorm:"not null;size:255"`
    ReporterRole          ReporterRole           `json:"reporter_role" gorm:"type:varchar(50);not null"`
    Village               string                 `json:"village" gorm:"not null;size:255"`
    District              string                 `json:"district" gorm:"not null;size:255"`
    BuildingName          string                 `json:"building_name" gorm:"not null;size:255"`
    BuildingType          BuildingType           `json:"building_type" gorm:"type:varchar(50);not null"`
    ReportStatus          ReportStatusType       `json:"report_status" gorm:"type:varchar(50);not null"`
    FundingSource         FundingSource          `json:"funding_source" gorm:"type:varchar(50);not null"`
    LastYearConstruction  int                    `json:"last_year_construction"`
    FullAddress           string                 `json:"full_address" gorm:"type:text"`
    Latitude              float64                `json:"latitude" gorm:"type:decimal(10,8)"`
    Longitude             float64                `json:"longitude" gorm:"type:decimal(11,8)"`
    FloorArea             float64                `json:"floor_area" gorm:"type:decimal(10,2)"`
    FloorCount            int                    `json:"floor_count"`
    WorkType              *WorkType              `json:"work_type,omitempty" gorm:"type:varchar(50)"`
    ConditionAfterRehab   *ConditionAfterRehab  `json:"condition_after_rehab,omitempty" gorm:"type:varchar(100)"`
    Photos                []ReportPhoto          `json:"photos" gorm:"foreignKey:ReportID"`
    // CreatedBy             string                 `json:"created_by" gorm:"type:varchar(26);not null"`
    CreatedAt             time.Time              `json:"created_at" gorm:"not null"`
    UpdatedAt             time.Time              `json:"updated_at" gorm:"not null"`
}

type ReportPhoto struct {
    ID         string    `json:"id" gorm:"type:varchar(26);primary_key"`
    ReportID   string    `json:"report_id" gorm:"type:varchar(26);not null"`
    PhotoURL   string    `json:"photo_url" gorm:"not null;size:500"`
    PhotoType  string    `json:"photo_type" gorm:"type:varchar(50)"`
    CreatedAt  time.Time `json:"created_at" gorm:"not null"`
}

func (r *Report) BeforeCreate() {
    if r.ID == "" {
        r.ID = utils.GenerateULID()
    }
    now := time.Now()
    r.CreatedAt = now
    r.UpdatedAt = now
}

func (r *Report) BeforeUpdate() {
    r.UpdatedAt = time.Now()
}

func (rp *ReportPhoto) BeforeCreate() {
    if rp.ID == "" {
        rp.ID = utils.GenerateULID()
    }
    rp.CreatedAt = time.Now()
}