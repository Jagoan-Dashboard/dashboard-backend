
package entity

import (
    "time"
    "github.com/google/uuid"
)

type IndikatorDemografi struct {
    ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Indikator string    `json:"indikator" gorm:"type:varchar(150);not null"`
    Tahun     int       `json:"tahun" gorm:"not null"`
    Nilai     float64   `json:"nilai" gorm:"type:decimal(12,2);not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (IndikatorDemografi) TableName() string {
    return "indikator_demografi"
}

func (i *IndikatorDemografi) BeforeCreate() error {
    if i.ID == "" {
        i.ID = uuid.New().String()
    }
    return nil
}