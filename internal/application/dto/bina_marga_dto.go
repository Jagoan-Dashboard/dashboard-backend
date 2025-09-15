
package dto

import (
    "time"
    "building-report-backend/internal/domain/entity"
)

type CreateBinaMargaRequest struct {
    ReporterName        string    `json:"reporter_name" validate:"required"`
    InstitutionUnit     string    `json:"institution_unit" validate:"required,oneof=DINAS DESA KECAMATAN"`
    PhoneNumber         string    `json:"phone_number" validate:"required"`
    ReportDateTime      time.Time `json:"report_datetime" validate:"required"`
    RoadName            string    `json:"road_name" validate:"required"`
    RoadType            string    `json:"road_type" validate:"required"`
    RoadClass           string    `json:"road_class" validate:"required"`
    Latitude            float64   `json:"latitude" validate:"required,min=-90,max=90"`
    Longitude           float64   `json:"longitude" validate:"required,min=-180,max=180"`
    DamageType          string    `json:"damage_type" validate:"required"`
    DamageLevel         string    `json:"damage_level" validate:"required,oneof=RINGAN SEDANG BERAT"`
    DamagedLength       float64   `json:"damaged_length" validate:"min=0"`
    DamagedWidth        float64   `json:"damaged_width" validate:"min=0"`
    TrafficImpact       string    `json:"traffic_impact" validate:"required"`
    UrgencyLevel        string    `json:"urgency_level" validate:"required,oneof=RENDAH SEDANG TINGGI DARURAT"`
    CauseOfDamage       string    `json:"cause_of_damage,omitempty"`
    Notes               string    `json:"notes,omitempty"`
}

func (r *CreateBinaMargaRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateBinaMargaRequest struct {
    RoadName               string  `json:"road_name,omitempty"`
    RoadType               string  `json:"road_type,omitempty"`
    RoadClass              string  `json:"road_class,omitempty"`
    DamageType             string  `json:"damage_type,omitempty"`
    DamageLevel            string  `json:"damage_level,omitempty"`
    DamagedLength          float64 `json:"damaged_length,omitempty"`
    DamagedWidth           float64 `json:"damaged_width,omitempty"`
    TrafficImpact          string  `json:"traffic_impact,omitempty"`
    UrgencyLevel           string  `json:"urgency_level,omitempty"`
    CauseOfDamage          string  `json:"cause_of_damage,omitempty"`
    Notes                  string  `json:"notes,omitempty"`
    HandlingRecommendation string  `json:"handling_recommendation,omitempty"`
    EstimatedBudget        float64 `json:"estimated_budget,omitempty"`
    EstimatedRepairTime    int     `json:"estimated_repair_time,omitempty"`
}

func (r *UpdateBinaMargaRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateBinaMargaStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=PENDING VERIFIED PLANNED IN_PROGRESS COMPLETED POSTPONED REJECTED"`
    Notes  string `json:"notes,omitempty"`
}

func (r *UpdateBinaMargaStatusRequest) Validate() error {
    return validate.Struct(r)
}

type PaginatedBinaMargaResponse struct {
    Reports    []*entity.BinaMargaReport `json:"reports"`
    Total      int64                     `json:"total"`
    Page       int                       `json:"page"`
    PerPage    int                       `json:"per_page"`
    TotalPages int64                     `json:"total_pages"`
}

type BinaMargaStatisticsResponse struct {
    TotalReports           int64                    `json:"total_reports"`
    EmergencyReports       int64                    `json:"emergency_reports"`
    BlockedRoads           int64                    `json:"blocked_roads"`
    TotalDamagedArea       float64                  `json:"total_damaged_area_sqm"`
    TotalDamagedLength     float64                  `json:"total_damaged_length_m"`
    RoadTypeDistribution   []map[string]interface{} `json:"road_type_distribution"`
    DamageTypeDistribution []map[string]interface{} `json:"damage_type_distribution"`
    DamageLevelCounts      []map[string]interface{} `json:"damage_level_counts"`
    UrgencyLevelCounts     []map[string]interface{} `json:"urgency_level_counts"`
    StatusDistribution     []map[string]interface{} `json:"status_distribution"`
    TrafficImpactCounts    []map[string]interface{} `json:"traffic_impact_counts"`
    EstimatedTotalBudget   float64                  `json:"estimated_total_budget"`
    AverageRepairTime      float64                  `json:"average_repair_time_days"`
}

type RoadDamageByTypeResponse struct {
    RoadType             string  `json:"road_type"`
    RoadClass            string  `json:"road_class"`
    ReportCount          int     `json:"report_count"`
    TotalDamagedArea     float64 `json:"total_damaged_area"`
    TotalDamagedLength   float64 `json:"total_damaged_length"`
    TotalEstimatedBudget float64 `json:"total_estimated_budget"`
    AvgRepairTime        float64 `json:"avg_repair_time_days"`
    EmergencyCount       int     `json:"emergency_count"`
}