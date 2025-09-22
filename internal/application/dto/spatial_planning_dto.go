
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


type TataRuangBasicStatistics struct {
    TotalReports            int64   `json:"total_reports"`
    EstimatedTotalLengthM   float64 `json:"estimated_total_length_m"`
    EstimatedTotalAreaM2    float64 `json:"estimated_total_area_m2"`
    UrgentReportsCount      int64   `json:"urgent_reports_count"`
}

type TataRuangLocationDistribution struct {
    District       string  `json:"district"`
    Village        string  `json:"village"`
    ViolationCount int     `json:"violation_count"`
    AvgLatitude    float64 `json:"avg_latitude"`
    AvgLongitude   float64 `json:"avg_longitude"`
    UrgentCount    int     `json:"urgent_count"`
    SevereCount    int     `json:"severe_count"`
}

type TataRuangUrgencyStatistics struct {
    UrgencyLevel string  `json:"urgency_level"`
    Count        int64   `json:"count"`
    Percentage   float64 `json:"percentage"`
}

type TataRuangViolationTypeStatistics struct {
    ViolationType string  `json:"violation_type"`
    Count         int64   `json:"count"`
    Percentage    float64 `json:"percentage"`
    SevereCount   int     `json:"severe_count"`
    UrgentCount   int     `json:"urgent_count"`
}

type TataRuangViolationLevelStatistics struct {
    ViolationLevel string  `json:"violation_level"`
    Count          int64   `json:"count"`
    Percentage     float64 `json:"percentage"`
    UrgentCount    int     `json:"urgent_count"`
}

type TataRuangAreaCategoryDistribution struct {
    AreaCategory string  `json:"area_category"`
    Count        int64   `json:"count"`
    Percentage   float64 `json:"percentage"`
    UrgentCount  int     `json:"urgent_count"`
    SevereCount  int     `json:"severe_count"`
}

type TataRuangEnvironmentalImpactStatistics struct {
    EnvironmentalImpact string  `json:"environmental_impact"`
    Count               int64   `json:"count"`
    Percentage          float64 `json:"percentage"`
    SevereCount         int     `json:"severe_count"`
}

type TataRuangOverviewResponse struct {
    // Baris pertama - Basic statistics
    BasicStats TataRuangBasicStatistics `json:"basic_stats"`
    
    // Baris kedua - Location and urgency
    LocationDistribution []TataRuangLocationDistribution `json:"location_distribution"`
    UrgencyStatistics    []TataRuangUrgencyStatistics    `json:"urgency_statistics"`
    
    // Baris ketiga - Violation details  
    ViolationTypeStatistics  []TataRuangViolationTypeStatistics  `json:"violation_type_statistics"`
    ViolationLevelStatistics []TataRuangViolationLevelStatistics `json:"violation_level_statistics"`
    
    // Additional insights
    AreaCategoryDistribution       []TataRuangAreaCategoryDistribution       `json:"area_category_distribution"`
    EnvironmentalImpactStatistics  []TataRuangEnvironmentalImpactStatistics  `json:"environmental_impact_statistics"`
}