package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RiceField struct {
	ID                  string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	District            string    `json:"district" gorm:"not null"`
	Longitude           float64   `json:"longitude" gorm:"type:decimal(10,8)"`
	Latitude            float64   `json:"latitude" gorm:"type:decimal(11,8)"`
	Date                time.Time `json:"date" gorm:"not null"`
	RainfedRiceFields   float64   `json:"rainfed_rice_fields" gorm:"type:decimal(15,2);comment:'Area of rainfed rice fields in hectares'"`
	IrrigatedRiceFields float64   `json:"irrigated_rice_fields" gorm:"type:decimal(15,2);comment:'Area of irrigated rice fields in hectares'"`
	CreatedAt           time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"not null"`
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
	return nil
}

func (rf *RiceField) BeforeUpdate(tx *gorm.DB) error {
	rf.UpdatedAt = time.Now()
	return nil
}
