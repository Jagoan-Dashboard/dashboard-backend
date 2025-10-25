package dto

import (
    "building-report-backend/internal/domain/entity"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type RegisterRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func (r *RegisterRequest) Validate() error {
    return validate.Struct(r)
}

type LoginRequest struct {
    Identifier string `json:"identifier" validate:"required"` // bisa username atau email
    Password   string `json:"password" validate:"required"`
}

func (l *LoginRequest) Validate() error {
    return validate.Struct(l)
}

type AuthResponse struct {
    Token     string      `json:"token"`
    User      interface{} `json:"user"`
    ExpiresIn int         `json:"expires_in"`
}

type CreateUserRequest struct {
    Username string           `json:"username" validate:"required,min=3,max=50"`
    Email    string           `json:"email" validate:"required,email"`
    Password string           `json:"password" validate:"required,min=6"`
    Role     entity.UserRole  `json:"role" validate:"required,oneof=SUPERADMIN USER"`
}

func (r *CreateUserRequest) Validate() error {
    return validate.Struct(r)
}

type UpdateUserRequest struct {
    Role entity.UserRole `json:"role" validate:"required,oneof=SUPERADMIN USER"`
}

func (r *UpdateUserRequest) Validate() error {
    return validate.Struct(r)
}

type UserResponse struct {
    ID        string          `json:"id"`
    Username  string          `json:"username"`
    Email     string          `json:"email"`
    Role      entity.UserRole `json:"role"`
    IsActive  bool            `json:"is_active"`
    CreatedAt string          `json:"created_at"`
    UpdatedAt string          `json:"updated_at"`
}

type UserListResponse struct {
    Users []*UserResponse `json:"users"`
    Total int64           `json:"total"`
}