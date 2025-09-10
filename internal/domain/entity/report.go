// internal/domain/entity/report.go
package entity

import (
    "time"
    "github.com/google/uuid"
)

type Report struct {
    ID                     uuid.UUID              `json:"id" gorm:"type:uuid;primary_key"`
    ReporterName          string                 `json:"reporter_name" gorm:"not null"`
    ReporterRole          ReporterRole           `json:"reporter_role" gorm:"type:varchar(50)"`
    Village               string                 `json:"village" gorm:"not null"`
    District              string                 `json:"district" gorm:"not null"`
    BuildingName          string                 `json:"building_name" gorm:"not null"`
    BuildingType          BuildingType           `json:"building_type" gorm:"type:varchar(50)"`
    ReportStatus          ReportStatusType       `json:"report_status" gorm:"type:varchar(50)"`
    FundingSource         FundingSource          `json:"funding_source" gorm:"type:varchar(50)"`
    LastYearConstruction  int                    `json:"last_year_construction"`
    FullAddress           string                 `json:"full_address" gorm:"type:text"`
    Latitude              float64                `json:"latitude"`
    Longitude             float64                `json:"longitude"`
    FloorArea             float64                `json:"floor_area"`
    FloorCount            int                    `json:"floor_count"`
    WorkType              *WorkType              `json:"work_type,omitempty" gorm:"type:varchar(50)"`
    ConditionAfterRehab   *ConditionAfterRehab  `json:"condition_after_rehab,omitempty" gorm:"type:varchar(100)"`
    Photos                []ReportPhoto          `json:"photos" gorm:"foreignKey:ReportID"`
    CreatedBy             uuid.UUID              `json:"created_by" gorm:"type:uuid"`
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
}

type ReportPhoto struct {
    ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    ReportID   uuid.UUID `json:"report_id" gorm:"type:uuid;not null"`
    PhotoURL   string    `json:"photo_url" gorm:"not null"`
    PhotoType  string    `json:"photo_type" gorm:"type:varchar(50)"` // closeup / overall
    CreatedAt  time.Time `json:"created_at"`
}

func (r *Report) BeforeCreate() {
    if r.ID == uuid.Nil {
        r.ID = uuid.New()
    }
    r.CreatedAt = time.Now()
    r.UpdatedAt = time.Now()
}

func (rp *ReportPhoto) BeforeCreate() {
    if rp.ID == uuid.Nil {
        rp.ID = uuid.New()
    }
    rp.CreatedAt = time.Now()
}