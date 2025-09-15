
package dto

import (
    "time"
    "building-report-backend/internal/domain/entity"
    
)

type CreateWaterResourcesRequest struct {
    ReporterName          string    `json:"reporter_name" validate:"required"`
    InstitutionUnit       string    `json:"institution_unit" validate:"required,oneof=DINAS DESA KECAMATAN"`
    PhoneNumber           string    `json:"phone_number" validate:"required"`
    ReportDateTime        time.Time `json:"report_datetime" validate:"required"`
    IrrigationAreaName    string    `json:"irrigation_area_name" validate:"required"`
    IrrigationType        string    `json:"irrigation_type" validate:"required"`
    Latitude              float64   `json:"latitude" validate:"required,min=-90,max=90"`
    Longitude             float64   `json:"longitude" validate:"required,min=-180,max=180"`
    DamageType            string    `json:"damage_type" validate:"required"`
    DamageLevel           string    `json:"damage_level" validate:"required,oneof=RINGAN SEDANG BERAT"`
    EstimatedLength       float64   `json:"estimated_length" validate:"min=0"`
    EstimatedWidth        float64   `json:"estimated_width" validate:"min=0"`
    EstimatedVolume       float64   `json:"estimated_volume" validate:"min=0"`
    AffectedRiceFieldArea float64   `json:"affected_rice_field_area" validate:"min=0"`
    AffectedFarmersCount  int       `json:"affected_farmers_count" validate:"min=0"`
    UrgencyCategory       string    `json:"urgency_category" validate:"required,oneof=MENDESAK RUTIN"`
    Notes                 string    `json:"notes,omitempty"`
}

func (r *CreateWaterResourcesRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateWaterResourcesRequest struct {
    IrrigationAreaName     string  `json:"irrigation_area_name,omitempty"`
    IrrigationType         string  `json:"irrigation_type,omitempty"`
    DamageType             string  `json:"damage_type,omitempty"`
    DamageLevel            string  `json:"damage_level,omitempty"`
    EstimatedLength        float64 `json:"estimated_length,omitempty"`
    EstimatedWidth         float64 `json:"estimated_width,omitempty"`
    EstimatedVolume        float64 `json:"estimated_volume,omitempty"`
    AffectedRiceFieldArea  float64 `json:"affected_rice_field_area,omitempty"`
    AffectedFarmersCount   int     `json:"affected_farmers_count,omitempty"`
    UrgencyCategory        string  `json:"urgency_category,omitempty"`
    Notes                  string  `json:"notes,omitempty"`
    HandlingRecommendation string  `json:"handling_recommendation,omitempty"`
    EstimatedBudget        float64 `json:"estimated_budget,omitempty"`
}

func (r *UpdateWaterResourcesRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateWaterStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=PENDING VERIFIED IN_PROGRESS COMPLETED POSTPONED REJECTED"`
    Notes  string `json:"notes,omitempty"`
}

func (r *UpdateWaterStatusRequest) Validate() error {
    return validate.Struct(r)
}

type PaginatedWaterResourcesResponse struct {
    Reports    []*entity.WaterResourcesReport `json:"reports"`
    Total      int64                          `json:"total"`
    Page       int                            `json:"page"`
    PerPage    int                            `json:"per_page"`
    TotalPages int64                          `json:"total_pages"`
}

type WaterResourcesStatisticsResponse struct {
    TotalReports          int64                    `json:"total_reports"`
    UrgentPending         int64                    `json:"urgent_pending"`
    TotalAffectedAreaHa   float64                  `json:"total_affected_area_ha"`
    TotalAffectedFarmers  int64                    `json:"total_affected_farmers"`
    DamageTypes           []map[string]interface{} `json:"damage_types"`
    IrrigationTypes       []map[string]interface{} `json:"irrigation_types"`
    StatusDistribution    []map[string]interface{} `json:"status_distribution"`
    EstimatedTotalBudget  float64                  `json:"estimated_total_budget"`
}

type DamageByAreaResponse struct {
    IrrigationAreaName    string  `json:"irrigation_area_name"`
    ReportCount          int     `json:"report_count"`
    TotalAffectedArea    float64 `json:"total_affected_area"`
    TotalAffectedFarmers int     `json:"total_affected_farmers"`
    TotalEstimatedBudget float64 `json:"total_estimated_budget"`
    AvgDamageArea        float64 `json:"avg_damage_area"`
}