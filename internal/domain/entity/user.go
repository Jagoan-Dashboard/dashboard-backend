package entity

import (
	"building-report-backend/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        string    `json:"id" gorm:"type:varchar(26);primary_key"`
    Username  string    `json:"username" gorm:"unique;not null;size:50"`
    Email     string    `json:"email" gorm:"unique;not null;size:100"`
    Password  string    `json:"-" gorm:"not null;size:255"`
    Role      UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'USER'"`
    IsActive  bool      `json:"is_active" gorm:"default:true;not null"`
    CreatedAt time.Time `json:"created_at" gorm:"not null"`
    UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type UserRole string

const (
    RoleSuperAdmin UserRole = "SUPERADMIN"
    RoleUser       UserRole = "USER"
)

func (u *User) BeforeCreate() error {
    if u.ID == "" {
        u.ID = utils.GenerateULID()
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)

    now := time.Now()
    u.CreatedAt = now
    u.UpdatedAt = now
    return nil
}

func (u *User) BeforeUpdate() error {
    u.UpdatedAt = time.Now()
    return nil
}

func (u *User) ComparePassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}

func (u *User) IsSuperAdmin() bool {
    return u.Role == RoleSuperAdmin
}