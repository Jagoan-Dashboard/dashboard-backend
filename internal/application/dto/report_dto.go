
package dto

import (
    "building-report-backend/internal/domain/entity"
)

type CreateReportRequest struct {
    ReporterName         string  `json:"reporter_name" form:"reporter_name" validate:"required"`
    ReporterRole         string  `json:"reporter_role" form:"reporter_role" validate:"required"`
    Village              string  `json:"village" form:"village" validate:"required"`
    District             string  `json:"district" form:"district" validate:"required"`
    BuildingName         string  `json:"building_name" form:"building_name" validate:"required"`
    BuildingType         string  `json:"building_type" form:"building_type" validate:"required"`
    ReportStatus         string  `json:"report_status" form:"report_status" validate:"required"`
    FundingSource        string  `json:"funding_source" form:"funding_source" validate:"required"`
    LastYearConstruction int     `json:"last_year_construction" form:"last_year_construction" validate:"required,min=1900,max=2100"`
    FullAddress          string  `json:"full_address" form:"full_address" validate:"required"`
    Latitude             float64 `json:"latitude" form:"latitude" validate:"required,min=-90,max=90"`
    Longitude            float64 `json:"longitude" form:"longitude" validate:"required,min=-180,max=180"`
    FloorArea            float64 `json:"floor_area" form:"floor_area" validate:"required,min=0"`
    FloorCount           int     `json:"floor_count" form:"floor_count" validate:"required,min=1"`
    WorkType             string  `json:"work_type,omitempty" form:"work_type"`
    ConditionAfterRehab  string  `json:"condition_after_rehab,omitempty" form:"condition_after_rehab"`
}

func (r *CreateReportRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateReportRequest struct {
    BuildingName         string  `json:"building_name,omitempty"`
    BuildingType         string  `json:"building_type,omitempty"`
    ReportStatus         string  `json:"report_status,omitempty"`
    FundingSource        string  `json:"funding_source,omitempty"`
    LastYearConstruction int     `json:"last_year_construction,omitempty"`
    FullAddress          string  `json:"full_address,omitempty"`
    Latitude             float64 `json:"latitude,omitempty"`
    Longitude            float64 `json:"longitude,omitempty"`
    FloorArea            float64 `json:"floor_area,omitempty"`
    FloorCount           int     `json:"floor_count,omitempty"`
    WorkType             string  `json:"work_type,omitempty"`
    ConditionAfterRehab  string  `json:"condition_after_rehab,omitempty"`
}

func (r *UpdateReportRequest) Validate() error {
    return validate.Struct(r)
}

type PaginatedReportsResponse struct {
    Reports    []*entity.Report `json:"reports"`
    Total      int64            `json:"total"`
    Page       int              `json:"page"`
    PerPage    int              `json:"per_page"`
    TotalPages int64            `json:"total_pages"`
}