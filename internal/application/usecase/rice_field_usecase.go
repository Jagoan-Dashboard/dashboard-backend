package usecase

import (
	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type RiceFieldUseCase struct {
	riceFieldRepo repository.RiceFieldRepository
	cacheRepo     repository.CacheRepository
}

func NewRiceFieldUseCase(
	riceFieldRepo repository.RiceFieldRepository,
	cacheRepo repository.CacheRepository,
) *RiceFieldUseCase {
	return &RiceFieldUseCase{
		riceFieldRepo: riceFieldRepo,
		cacheRepo:     cacheRepo,
	}
}

// Create new rice field record
func (uc *RiceFieldUseCase) Create(ctx context.Context, req *dto.CreateRiceFieldRequest) (*dto.RiceFieldResponse, error) {
	// Basic validation
	if req.District == "" {
		return nil, errors.New("district is required")
	}
	if req.Date.IsZero() {
		return nil, errors.New("date is required")
	}

	// Create entity
	riceField := &entity.RiceField{
		District:                req.District,
		Longitude:               req.Longitude,
		Latitude:                req.Latitude,
		Date:                    req.Date,
		Year:                    req.Year,
		RainfedRiceFields:       req.RainfedRiceFields,
		IrrigatedRiceFields:     req.IrrigatedRiceFields,
		TotalRiceFieldArea:      req.TotalRiceFieldArea,
		DryfieldArea:            req.DryfieldArea,
		ShiftingCultivationArea: req.ShiftingCultivationArea,
		UnusedLandArea:          req.UnusedLandArea,
		TotalNonRiceFieldArea:   req.TotalNonRiceFieldArea,
		TotalLandArea:           req.TotalLandArea,
		DataSource:              req.DataSource,
	}

	// Save to database
	if err := uc.riceFieldRepo.Create(ctx, riceField); err != nil {
		return nil, fmt.Errorf("failed to create rice field: %w", err)
	}

	// Invalidate cache
	uc.invalidateCache(ctx)

	return uc.toResponse(riceField), nil
}

// Update rice field record
func (uc *RiceFieldUseCase) Update(ctx context.Context, id string, req *dto.UpdateRiceFieldRequest) (*dto.RiceFieldResponse, error) {
	// Get existing record
	riceField, err := uc.riceFieldRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("rice field not found: %w", err)
	}

	// Update fields if provided
	if req.District != "" {
		riceField.District = req.District
	}
	if req.Longitude != 0 {
		riceField.Longitude = req.Longitude
	}
	if req.Latitude != 0 {
		riceField.Latitude = req.Latitude
	}
	if !req.Date.IsZero() {
		riceField.Date = req.Date
	}
	if req.Year != 0 {
		riceField.Year = req.Year
	}
	if req.RainfedRiceFields != 0 {
		riceField.RainfedRiceFields = req.RainfedRiceFields
	}
	if req.IrrigatedRiceFields != 0 {
		riceField.IrrigatedRiceFields = req.IrrigatedRiceFields
	}
	if req.DryfieldArea != 0 {
		riceField.DryfieldArea = req.DryfieldArea
	}
	if req.ShiftingCultivationArea != 0 {
		riceField.ShiftingCultivationArea = req.ShiftingCultivationArea
	}
	if req.UnusedLandArea != 0 {
		riceField.UnusedLandArea = req.UnusedLandArea
	}

	// Save changes
	if err := uc.riceFieldRepo.Update(ctx, riceField); err != nil {
		return nil, fmt.Errorf("failed to update rice field: %w", err)
	}

	// Invalidate cache
	uc.invalidateCache(ctx)

	return uc.toResponse(riceField), nil
}

// Delete rice field record
func (uc *RiceFieldUseCase) Delete(ctx context.Context, id string) error {
	// Check if exists
	_, err := uc.riceFieldRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("rice field not found: %w", err)
	}

	// Delete
	if err := uc.riceFieldRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete rice field: %w", err)
	}

	// Invalidate cache
	uc.invalidateCache(ctx)

	return nil
}

// GetByID retrieves a single rice field by ID
func (uc *RiceFieldUseCase) GetByID(ctx context.Context, id string) (*dto.RiceFieldResponse, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("rice_field:%s", id)
	var cached dto.RiceFieldResponse
	err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err == nil {
		return &cached, nil
	}

	// Get from database
	riceField, err := uc.riceFieldRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("rice field not found: %w", err)
	}

	response := uc.toResponse(riceField)

	// Cache the result
	uc.cacheRepo.Set(ctx, cacheKey, response, 15*time.Minute)

	return response, nil
}

// GetAll retrieves paginated rice fields with filters
func (uc *RiceFieldUseCase) GetAll(ctx context.Context, page, perPage int, filters map[string]interface{}) (*dto.PaginatedRiceFieldResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Try cache
	cacheKey := fmt.Sprintf("rice_fields:page:%d:per_page:%d:filters:%v", page, perPage, filters)
	var cached dto.PaginatedRiceFieldResponse
	err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err == nil {
		return &cached, nil
	}

	// Get from database
	riceFields, total, err := uc.riceFieldRepo.FindAll(ctx, perPage, offset, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get rice fields: %w", err)
	}

	// Convert to response
	responses := make([]*dto.RiceFieldResponse, len(riceFields))
	for i, rf := range riceFields {
		responses[i] = uc.toResponse(rf)
	}

	totalPages := (total + int64(perPage) - 1) / int64(perPage)

	result := &dto.PaginatedRiceFieldResponse{
		RiceFields: responses,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}

	// Cache result
	uc.cacheRepo.Set(ctx, cacheKey, result, 10*time.Minute)

	return result, nil
}

// GetStatistics retrieves rice field statistics
func (uc *RiceFieldUseCase) GetStatistics(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before end date")
	}

	// Try cache
	cacheKey := fmt.Sprintf("rice_field_stats:%s:%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	var cached map[string]interface{}
	err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err == nil {
		return cached, nil
	}

	// Get from repository
	stats, err := uc.riceFieldRepo.GetRiceFieldStatistics(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	// Cache result
	uc.cacheRepo.Set(ctx, cacheKey, stats, 30*time.Minute)

	return stats, nil
}

// GetDistributionByDistrict retrieves rice field distribution by district
func (uc *RiceFieldUseCase) GetDistributionByDistrict(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	// Try cache
	cacheKey := fmt.Sprintf("rice_field_dist:%s:%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	var cached []map[string]interface{}
	err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err == nil {
		return cached, nil
	}

	// Get from repository
	dist, err := uc.riceFieldRepo.GetRiceFieldDistributionByDistrict(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get distribution: %w", err)
	}

	// Cache result
	uc.cacheRepo.Set(ctx, cacheKey, dist, 30*time.Minute)

	return dist, nil
}

// GetMapData retrieves rice field data for map visualization
func (uc *RiceFieldUseCase) GetMapData(ctx context.Context, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	// Try cache
	cacheKey := fmt.Sprintf("rice_field_map:%s:%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	var cached []map[string]interface{}
	err := uc.cacheRepo.Get(ctx, cacheKey, &cached)
	if err == nil {
		return cached, nil
	}

	// Get from repository
	mapData, err := uc.riceFieldRepo.GetIndividualRiceFieldDistribution(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get map data: %w", err)
	}

	// Cache result
	uc.cacheRepo.Set(ctx, cacheKey, mapData, 20*time.Minute)

	return mapData, nil
}

// Helper: Convert entity to response DTO
func (uc *RiceFieldUseCase) toResponse(rf *entity.RiceField) *dto.RiceFieldResponse {
	return &dto.RiceFieldResponse{
		ID:                      rf.ID,
		District:                rf.District,
		Longitude:               rf.Longitude,
		Latitude:                rf.Latitude,
		Date:                    rf.Date,
		Year:                    rf.Year,
		RainfedRiceFields:       rf.RainfedRiceFields,
		IrrigatedRiceFields:     rf.IrrigatedRiceFields,
		TotalRiceFieldArea:      rf.TotalRiceFieldArea,
		DryfieldArea:            rf.DryfieldArea,
		ShiftingCultivationArea: rf.ShiftingCultivationArea,
		UnusedLandArea:          rf.UnusedLandArea,
		TotalNonRiceFieldArea:   rf.TotalNonRiceFieldArea,
		TotalLandArea:           rf.TotalLandArea,
		DataSource:              rf.DataSource,
		CreatedAt:               rf.CreatedAt,
		UpdatedAt:               rf.UpdatedAt,
	}
}

// Helper: Invalidate related caches
func (uc *RiceFieldUseCase) invalidateCache(ctx context.Context) {
	// Invalidate all rice field related caches
	patterns := []string{
		"rice_field:*",
		"rice_fields:*",
		"rice_field_stats:*",
		"rice_field_dist:*",
		"rice_field_map:*",
	}

	for _, pattern := range patterns {
		uc.cacheRepo.Delete(ctx, pattern)
	}
}