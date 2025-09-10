// internal/domain/entity/user.go
package entity

import (
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    Role      UserRole  `json:"role" gorm:"type:varchar(20)"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserRole string

const (
    RoleAdmin     UserRole = "ADMIN"
    RoleOperator  UserRole = "OPERATOR"
    RoleViewer    UserRole = "VIEWER"
)

func (u *User) BeforeCreate() error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    
    u.CreatedAt = time.Now()
    u.UpdatedAt = time.Now()
    return nil
}

func (u *User) ComparePassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}