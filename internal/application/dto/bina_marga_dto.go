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
    
    
    RoadName            string    `json:"road_name" validate:"omitempty"`
    RoadType            string    `json:"road_type" validate:"omitempty,oneof=JALAN_NASIONAL JALAN_PROVINSI JALAN_KABUPATEN JALAN_DESA"`
    RoadClass           string    `json:"road_class" validate:"omitempty,oneof=ARTERI KOLEKTOR LOKAL LINGKUNGAN"`
    SegmentLength       float64   `json:"segment_length" validate:"min=0"` 
    Latitude            float64   `json:"latitude" validate:"required,min=-90,max=90"`
    Longitude           float64   `json:"longitude" validate:"required,min=-180,max=180"`
    
    
    PavementType        string    `json:"pavement_type" validate:"omitempty,oneof=ASPAL_FLEXIBLE BETON_RIGID PAVING JALAN_TANAH"`
    DamageType          string    `json:"damage_type" validate:"omitempty"`
    DamageLevel         string    `json:"damage_level" validate:"omitempty,oneof=RINGAN SEDANG BERAT"`
    DamagedLength       float64   `json:"damaged_length" validate:"min=0"` 
    DamagedWidth        float64   `json:"damaged_width" validate:"min=0"`  
    TotalDamagedArea    float64   `json:"total_damaged_area" validate:"min=0"` 
    
    
    BridgeName          string    `json:"bridge_name,omitempty"`
    BridgeStructureType string    `json:"bridge_structure_type,omitempty"` 
    BridgeDamageType    string    `json:"bridge_damage_type,omitempty"`
    BridgeDamageLevel   string    `json:"bridge_damage_level,omitempty"` 
    
    
    TrafficCondition    string    `json:"traffic_condition" validate:"required"` 
    TrafficImpact       string    `json:"traffic_impact"`
    DailyTrafficVolume  int       `json:"daily_traffic_volume" validate:"min=0"`
    UrgencyLevel        string    `json:"urgency_level" validate:"required,oneof=DARURAT CEPAT RUTIN"`
    
    
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
    SegmentLength          float64 `json:"segment_length,omitempty"`
    
    
    PavementType           string  `json:"pavement_type,omitempty"`
    DamageType             string  `json:"damage_type,omitempty"`
    DamageLevel            string  `json:"damage_level,omitempty"`
    DamagedLength          float64 `json:"damaged_length,omitempty"`
    DamagedWidth           float64 `json:"damaged_width,omitempty"`
    TotalDamagedArea       float64 `json:"total_damaged_area,omitempty"`
    
    
    BridgeName             string  `json:"bridge_name,omitempty"`
    BridgeStructureType    string  `json:"bridge_structure_type,omitempty"`
    BridgeDamageType       string  `json:"bridge_damage_type,omitempty"`
    BridgeDamageLevel      string  `json:"bridge_damage_level,omitempty"`
    
    
    TrafficCondition       string  `json:"traffic_condition,omitempty"`
    TrafficImpact          string  `json:"traffic_impact,omitempty"`
    DailyTrafficVolume     int     `json:"daily_traffic_volume,omitempty"`
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
    PavementTypeDistribution []map[string]interface{} `json:"pavement_type_distribution"`
    BridgeReports          int64                    `json:"bridge_reports"`
    EstimatedTotalBudget   float64                  `json:"estimated_total_budget"`
    AverageRepairTime      float64                  `json:"average_repair_time_days"`
}

type RoadDamageByTypeResponse struct {
    RoadType             string  `json:"road_type"`
    RoadClass            string  `json:"road_class"`
    PavementType         string  `json:"pavement_type"`
    ReportCount          int     `json:"report_count"`
    TotalDamagedArea     float64 `json:"total_damaged_area"`
    TotalDamagedLength   float64 `json:"total_damaged_length"`
    TotalEstimatedBudget float64 `json:"total_estimated_budget"`
    AvgRepairTime        float64 `json:"avg_repair_time_days"`
    EmergencyCount       int     `json:"emergency_count"`
}


type BinaMargaMapPoint struct {
    Latitude        float64 `json:"latitude"`
    Longitude       float64 `json:"longitude"`
    RoadName        string  `json:"road_name"`
    RoadType        string  `json:"road_type"`
    DamageType      string  `json:"damage_type,omitempty"`
    DamageLevel     string  `json:"damage_level,omitempty"`
    BridgeName      string  `json:"bridge_name,omitempty"`
    BridgeDamageType  string `json:"bridge_damage_type,omitempty"`
    BridgeDamageLevel string `json:"bridge_damage_level,omitempty"`
    UrgencyLevel    string  `json:"urgency_level"`
}

type BinaMargaDashboardResponse struct {
    KPIs struct {
        AvgSegmentLengthM     float64 `json:"avg_segment_length_m"`
        AvgDamageAreaM2       float64 `json:"avg_damage_area_m2"`
        AvgDailyTrafficVolume float64 `json:"avg_daily_traffic_volume"`
        TotalReports          int64   `json:"total_reports"`
    } `json:"kpis"`

    PriorityDistribution           []KeyCount          `json:"priority_distribution"`

    MapPoints                      []BinaMargaMapPoint `json:"map_points"`
    RoadDamageLevelDistribution    []KeyCount          `json:"road_damage_level_distribution"`
    BridgeDamageLevelDistribution  []KeyCount          `json:"bridge_damage_level_distribution"`

    TopRoadDamageTypes             []KeyCount          `json:"top_road_damage_types"`
    TopBridgeDamageTypes           []KeyCount          `json:"top_bridge_damage_types"`
}

type BinaMargaOverviewResponse struct {
    BasicStats struct {
        AvgSegmentLengthM         float64 `json:"avg_segment_length_m"`
        AvgDamageAreaM2           float64 `json:"avg_damage_area_m2"`
        AvgDailyTrafficVolume     float64 `json:"avg_daily_traffic_volume"`
        TotalInfrastructureReports int64   `json:"total_infrastructure_reports"`
    } `json:"basic_stats"`
    
    LocationDistribution         []BinaMargaLocationStatsResponse `json:"location_distribution"`
    PriorityDistribution         []BinaMargaPriorityStatsResponse `json:"priority_distribution"`
    RoadDamageLevelDistribution  []BinaMargaRoadDamageLevelStatsResponse `json:"road_damage_level_distribution"`
    BridgeDamageLevelDistribution []BinaMargaBridgeDamageLevelStatsResponse `json:"bridge_damage_level_distribution"`
    TopRoadDamageTypes           []BinaMargaRoadDamageTypeStatsResponse `json:"top_road_damage_types"`
    TopBridgeDamageTypes         []BinaMargaBridgeDamageTypeStatsResponse `json:"top_bridge_damage_types"`
}

type BinaMargaLocationStatsResponse struct {
    RoadName              string  `json:"road_name"`
    Latitude              float64 `json:"latitude"`
    Longitude             float64 `json:"longitude"`
    DamageType            string  `json:"damage_type"`
    DamageLevel           string  `json:"damage_level"`
    UrgencyLevel          string  `json:"urgency_level"`
    TrafficImpact         string  `json:"traffic_impact"`
    DamagedArea           float64 `json:"damaged_area"`
}

type BinaMargaPriorityStatsResponse struct {
    PriorityLevel string `json:"priority_level"`
    Count         int64  `json:"count"`
}

type BinaMargaRoadDamageLevelStatsResponse struct {
    DamageLevel string `json:"damage_level"`
    Count       int64  `json:"count"`
}

type BinaMargaBridgeDamageLevelStatsResponse struct {
    DamageLevel string `json:"damage_level"`
    Count       int64  `json:"count"`
}

type BinaMargaRoadDamageTypeStatsResponse struct {
    DamageType string `json:"damage_type"`
    Count      int64  `json:"count"`
}

type BinaMargaBridgeDamageTypeStatsResponse struct {
    DamageType string `json:"damage_type"`
    Count      int64  `json:"count"`
}