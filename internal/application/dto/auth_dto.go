
package dto

import (
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