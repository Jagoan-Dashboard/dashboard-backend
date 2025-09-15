
package dto

import (
    "time"
    "building-report-backend/internal/domain/entity"
)

type CreateSpatialPlanningRequest struct {
    ReporterName        string    `json:"reporter_name" validate:"required"`
    Institution         string    `json:"institution" validate:"required,oneof=DINAS DESA KECAMATAN"`
    PhoneNumber         string    `json:"phone_number" validate:"required"`
    ReportDateTime      time.Time `json:"report_datetime" validate:"required"`
    AreaDescription     string    `json:"area_description" validate:"required"`
    AreaCategory        string    `json:"area_category" validate:"required"`
    ViolationType       string    `json:"violation_type" validate:"required"`
    ViolationLevel      string    `json:"violation_level" validate:"required,oneof=RINGAN SEDANG BERAT"`
    EnvironmentalImpact string    `json:"environmental_impact" validate:"required"`
    UrgencyLevel        string    `json:"urgency_level" validate:"required,oneof=MENDESAK BIASA"`
    Latitude            float64   `json:"latitude" validate:"required,min=-90,max=90"`
    Longitude           float64   `json:"longitude" validate:"required,min=-180,max=180"`
    Address             string    `json:"address" validate:"required"`
    Notes               string    `json:"notes,omitempty"`
}

func (r *CreateSpatialPlanningRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateSpatialPlanningRequest struct {
    AreaDescription     string    `json:"area_description,omitempty"`
    AreaCategory        string    `json:"area_category,omitempty"`
    ViolationType       string    `json:"violation_type,omitempty"`
    ViolationLevel      string    `json:"violation_level,omitempty"`
    EnvironmentalImpact string    `json:"environmental_impact,omitempty"`
    UrgencyLevel        string    `json:"urgency_level,omitempty"`
    Latitude            float64   `json:"latitude,omitempty"`
    Longitude           float64   `json:"longitude,omitempty"`
    Address             string    `json:"address,omitempty"`
    Notes               string    `json:"notes,omitempty"`
    Status              string    `json:"status,omitempty"`
}

func (r *UpdateSpatialPlanningRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateSpatialStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=PENDING REVIEWING PROCESSING RESOLVED REJECTED"`
    Notes  string `json:"notes,omitempty"`
}

func (r *UpdateSpatialStatusRequest) Validate() error {
    return validate.Struct(r)
}

type PaginatedSpatialReportsResponse struct {
    Reports    []*entity.SpatialPlanningReport `json:"reports"`
    Total      int64                           `json:"total"`
    Page       int                             `json:"page"`
    PerPage    int                             `json:"per_page"`
    TotalPages int64                           `json:"total_pages"`
}

type SpatialStatisticsResponse struct {
    TotalReports    int64                  `json:"total_reports"`
    UrgentReports   int64                  `json:"urgent_reports"`
    ViolationLevels []map[string]interface{} `json:"violation_levels"`
    StatusCounts    []map[string]interface{} `json:"status_counts"`
    AreaCategories  []map[string]interface{} `json:"area_categories"`
}