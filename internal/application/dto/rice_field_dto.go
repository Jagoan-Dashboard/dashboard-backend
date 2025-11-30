package dto

import "time"

type CreateRiceFieldRequest struct {
	District            string    `json:"district" validate:"required"`
	Longitude           float64   `json:"longitude" validate:"required"`
	Latitude            float64   `json:"latitude" validate:"required"`
	Date                time.Time `json:"date" validate:"required"`
	RainfedRiceFields   float64   `json:"rainfed_rice_fields"`
	IrrigatedRiceFields float64   `json:"irrigated_rice_fields"`
}

type UpdateRiceFieldRequest struct {
	District            string    `json:"district"`
	Longitude           float64   `json:"longitude"`
	Latitude            float64   `json:"latitude"`
	Date                time.Time `json:"date"`
	RainfedRiceFields   float64   `json:"rainfed_rice_fields"`
	IrrigatedRiceFields float64   `json:"irrigated_rice_fields"`
}

type RiceFieldResponse struct {
	ID                  string    `json:"id"`
	District            string    `json:"district"`
	Longitude           float64   `json:"longitude"`
	Latitude            float64   `json:"latitude"`
	Date                time.Time `json:"date"`
	RainfedRiceFields   float64   `json:"rainfed_rice_fields"`
	IrrigatedRiceFields float64   `json:"irrigated_rice_fields"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type RiceFieldStatsResponse struct {
	District              string  `json:"district"`
	TotalRainfedArea      float64 `json:"total_rainfed_area"`
	TotalIrrigatedArea    float64 `json:"total_irrigated_area"`
	TotalRiceFieldArea    float64 `json:"total_rice_field_area"`
	AverageRainfedArea    float64 `json:"average_rainfed_area"`
	AverageIrrigatedArea  float64 `json:"average_irrigated_area"`
	RainfedRiceFieldsCount int64  `json:"rainfed_rice_fields_count"`
	IrrigatedRiceFieldsCount int64 `json:"irrigated_rice_fields_count"`
}

type RiceFieldMapPoint struct {
	ID                  string  `json:"id"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	District            string  `json:"district"`
	Date                string  `json:"date"`
	RainfedRiceFields   float64 `json:"rainfed_rice_fields"`
	IrrigatedRiceFields float64 `json:"irrigated_rice_fields"`
	TotalArea           float64 `json:"total_area"`
}

type PaginatedRiceFieldResponse struct {
	RiceFields []*RiceFieldResponse `json:"rice_fields"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	TotalPages int64                `json:"total_pages"`
}

type RiceFieldAnalysisResponse struct {
	TotalRainfedArea        float64               `json:"total_rainfed_area"`
	TotalIrrigatedArea      float64               `json:"total_irrigated_area"`
	TotalRiceFieldArea      float64               `json:"total_rice_field_area"`
	DistributionByDistrict  []RiceFieldStatsResponse `json:"distribution_by_district"`
	DistributionMap         []RiceFieldMapPoint      `json:"distribution_map"`
	AreaTrendByMonth        []map[string]interface{} `json:"area_trend_by_month"`
	IrrigationRatio         float64               `json:"irrigation_ratio"`
}