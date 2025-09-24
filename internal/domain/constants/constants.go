package constants

import "time"

// Cache durations
const (
	UserCacheDuration      = 24 * time.Hour
	ReportCacheDuration    = 12 * time.Hour
	StatsticsCacheDuration = 30 * time.Minute
)

// Cache key prefixes
const (
	UserCachePrefix             = "user:"
	ReportCachePrefix           = "report:"
	AgricultureReportPrefix     = "agriculture_report:"
	SpatialPlanningReportPrefix = "spatial_planning_report:"
	WaterResourcesReportPrefix  = "water_resources_report:"
	BinaMargaReportPrefix       = "bina_marga_report:"
	StatisticsPrefix            = "stats:"
)

// JWT token expiration
const (
	TokenExpirationTime = 24 * time.Hour
	TokenExpirationSecs = 24 * 3600 // 24 hours in seconds
)

// Pagination defaults
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Database constraints
const (
	MaxTextLength       = 1000
	MaxURLLength        = 500
	MaxNameLength       = 255
	MaxDescLength       = 500
	MaxAddressLength    = 1000
	MaxSuggestionLength = 2000
)

// File upload constraints
const (
	MaxFileSize         = 10 * 1024 * 1024 // 10MB
	MaxPhotosPerReport  = 10
	AllowedImageFormats = "jpg,jpeg,png,webp"
)

// Priority score ranges
const (
	MinPriorityScore = 0
	MaxPriorityScore = 100
)

// Decimal precision
const (
	AreaDecimalPlaces       = 3  // For land area in hectares
	CoordinateDecimalPlaces = 8  // For latitude/longitude
	MoneyDecimalPlaces      = 2  // For currency amounts
)