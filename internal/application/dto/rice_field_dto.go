package dto

import "time"

type CreateRiceFieldRequest struct {
	District                string    `json:"district" validate:"required"`
	Longitude               float64   `json:"longitude"`
	Latitude                float64   `json:"latitude"`
	Date                    time.Time `json:"date" validate:"required"`
	Year                    int       `json:"year"`
	RainfedRiceFields       float64   `json:"rainfed_rice_fields"`
	IrrigatedRiceFields     float64   `json:"irrigated_rice_fields"`
	TotalRiceFieldArea      float64   `json:"total_rice_field_area"`
	DryfieldArea            float64   `json:"dryfield_area"`
	ShiftingCultivationArea float64   `json:"shifting_cultivation_area"`
	UnusedLandArea          float64   `json:"unused_land_area"`
	TotalNonRiceFieldArea   float64   `json:"total_non_rice_field_area"`
	TotalLandArea           float64   `json:"total_land_area"`
	DataSource              string    `json:"data_source"`
}

type UpdateRiceFieldRequest struct {
	District                string    `json:"district"`
	Longitude               float64   `json:"longitude"`
	Latitude                float64   `json:"latitude"`
	Date                    time.Time `json:"date"`
	Year                    int       `json:"year"`
	RainfedRiceFields       float64   `json:"rainfed_rice_fields"`
	IrrigatedRiceFields     float64   `json:"irrigated_rice_fields"`
	TotalRiceFieldArea      float64   `json:"total_rice_field_area"`
	DryfieldArea            float64   `json:"dryfield_area"`
	ShiftingCultivationArea float64   `json:"shifting_cultivation_area"`
	UnusedLandArea          float64   `json:"unused_land_area"`
	TotalNonRiceFieldArea   float64   `json:"total_non_rice_field_area"`
	TotalLandArea           float64   `json:"total_land_area"`
}

type RiceFieldResponse struct {
	ID                      string    `json:"id"`
	District                string    `json:"district"`
	Longitude               float64   `json:"longitude"`
	Latitude                float64   `json:"latitude"`
	Date                    time.Time `json:"date"`
	Year                    int       `json:"year"`
	RainfedRiceFields       float64   `json:"rainfed_rice_fields"`
	IrrigatedRiceFields     float64   `json:"irrigated_rice_fields"`
	TotalRiceFieldArea      float64   `json:"total_rice_field_area"`
	DryfieldArea            float64   `json:"dryfield_area"`
	ShiftingCultivationArea float64   `json:"shifting_cultivation_area"`
	UnusedLandArea          float64   `json:"unused_land_area"`
	TotalNonRiceFieldArea   float64   `json:"total_non_rice_field_area"`
	TotalLandArea           float64   `json:"total_land_area"`
	DataSource              string    `json:"data_source"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type RiceFieldStatsResponse struct {
	District                     string  `json:"district"`
	TotalRainfedArea             float64 `json:"total_rainfed_area"`
	TotalIrrigatedArea           float64 `json:"total_irrigated_area"`
	TotalRiceFieldArea           float64 `json:"total_rice_field_area"`
	TotalDryfieldArea            float64 `json:"total_dryfield_area"`
	TotalShiftingCultivationArea float64 `json:"total_shifting_cultivation_area"`
	TotalUnusedLandArea          float64 `json:"total_unused_land_area"`
	TotalNonRiceFieldArea        float64 `json:"total_non_rice_field_area"`
	TotalLandArea                float64 `json:"total_land_area"`
	AverageRainfedArea           float64 `json:"average_rainfed_area"`
	AverageIrrigatedArea         float64 `json:"average_irrigated_area"`
	RecordCount                  int64   `json:"record_count"`
}

type RiceFieldMapPoint struct {
	ID                      string  `json:"id"`
	Latitude                float64 `json:"latitude"`
	Longitude               float64 `json:"longitude"`
	District                string  `json:"district"`
	Date                    string  `json:"date"`
	Year                    int     `json:"year"`
	RainfedRiceFields       float64 `json:"rainfed_rice_fields"`
	IrrigatedRiceFields     float64 `json:"irrigated_rice_fields"`
	TotalRiceFieldArea      float64 `json:"total_rice_field_area"`
	DryfieldArea            float64 `json:"dryfield_area"`
	ShiftingCultivationArea float64 `json:"shifting_cultivation_area"`
	UnusedLandArea          float64 `json:"unused_land_area"`
	TotalNonRiceFieldArea   float64 `json:"total_non_rice_field_area"`
	TotalLandArea           float64 `json:"total_land_area"`
}

type PaginatedRiceFieldResponse struct {
	RiceFields []*RiceFieldResponse `json:"rice_fields"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	TotalPages int64                `json:"total_pages"`
}

type RiceFieldAnalysisResponse struct {
	TotalRainfedArea            float64                  `json:"total_rainfed_area"`
	TotalIrrigatedArea          float64                  `json:"total_irrigated_area"`
	TotalRiceFieldArea          float64                  `json:"total_rice_field_area"`
	TotalDryfieldArea           float64                  `json:"total_dryfield_area"`
	TotalShiftingCultivationArea float64                 `json:"total_shifting_cultivation_area"`
	TotalUnusedLandArea         float64                  `json:"total_unused_land_area"`
	TotalNonRiceFieldArea       float64                  `json:"total_non_rice_field_area"`
	TotalLandArea               float64                  `json:"total_land_area"`
	IrrigationRatio             float64                  `json:"irrigation_ratio"`
	DistributionByDistrict      []RiceFieldStatsResponse `json:"distribution_by_district"`
	DistributionMap             []RiceFieldMapPoint      `json:"distribution_map"`
	AreaTrendByMonth            []map[string]interface{} `json:"area_trend_by_month"`
}