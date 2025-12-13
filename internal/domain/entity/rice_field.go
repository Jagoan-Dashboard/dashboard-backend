package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RiceField struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	District  string    `json:"district" gorm:"not null;index"`
	Longitude float64   `json:"longitude" gorm:"type:decimal(10,8)"`
	Latitude  float64   `json:"latitude" gorm:"type:decimal(11,8)"`
	Date      time.Time `json:"date" gorm:"not null;index"`
	Year      int       `json:"year" gorm:"not null;index"`
	
	// Sawah (Rice Fields) - From template columns 4-6
	IrrigatedRiceFields float64 `json:"irrigated_rice_fields" gorm:"type:decimal(15,2);default:0;comment:'Luas Sawah Irigasi (ha)'"`
	RainfedRiceFields   float64 `json:"rainfed_rice_fields" gorm:"type:decimal(15,2);default:0;comment:'Luas Sawah Tadah Hujan (ha)'"`
	TotalRiceFieldArea  float64 `json:"total_rice_field_area" gorm:"type:decimal(15,2);default:0;comment:'Luas Lahan Sawah (ha)'"`
	
	// Non-Sawah (Non-Rice Fields) - From template columns 7-10
	DryfieldArea            float64 `json:"dryfield_area" gorm:"type:decimal(15,2);default:0;comment:'Luas Lahan Tegal/Kebun (ha)'"`
	ShiftingCultivationArea float64 `json:"shifting_cultivation_area" gorm:"type:decimal(15,2);default:0;comment:'Luas Lahan Ladang/Huma (ha)'"`
	UnusedLandArea          float64 `json:"unused_land_area" gorm:"type:decimal(15,2);default:0;comment:'Luas Lahan yang Sementara Tidak Diusahakan (ha)'"`
	TotalNonRiceFieldArea   float64 `json:"total_non_rice_field_area" gorm:"type:decimal(15,2);default:0;comment:'Luas Lahan Bukan Sawah (ha)'"`
	
	// Total - From template column 11
	TotalLandArea float64 `json:"total_land_area" gorm:"type:decimal(15,2);default:0;comment:'Total Luas Lahan (ha)'"`
	
	// Metadata
	DataSource string    `json:"data_source" gorm:"default:'manual';comment:'manual, import, api'"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null"`
}

func (RiceField) TableName() string {
	return "rice_fields"
}

func (rf *RiceField) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	rf.CreatedAt = now
	rf.UpdatedAt = now
	if rf.ID == "" {
		rf.ID = uuid.New().String()
	}
	
	// Auto-calculate year from date if not set
	if rf.Year == 0 && !rf.Date.IsZero() {
		rf.Year = rf.Date.Year()
	}
	
	// Auto-calculate totals if not provided
	rf.calculateTotals()
	
	return nil
}

func (rf *RiceField) BeforeUpdate(tx *gorm.DB) error {
	rf.UpdatedAt = time.Now()
	
	// Auto-calculate year from date if not set
	if rf.Year == 0 && !rf.Date.IsZero() {
		rf.Year = rf.Date.Year()
	}
	
	// Recalculate totals
	rf.calculateTotals()
	
	return nil
}

func (rf *RiceField) calculateTotals() {
	// Calculate total rice field area if not provided
	if rf.TotalRiceFieldArea == 0 {
		rf.TotalRiceFieldArea = rf.IrrigatedRiceFields + rf.RainfedRiceFields
	}
	
	// Calculate total non-rice field area if not provided
	if rf.TotalNonRiceFieldArea == 0 {
		rf.TotalNonRiceFieldArea = rf.DryfieldArea + rf.ShiftingCultivationArea + rf.UnusedLandArea
	}
	
	// Calculate total land area if not provided
	if rf.TotalLandArea == 0 {
		rf.TotalLandArea = rf.TotalRiceFieldArea + rf.TotalNonRiceFieldArea
	}
}