
package usecase

import (
	"context"
	"errors"
	"time"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/auth"

	"github.com/google/uuid"
)

var (
    ErrUserExists      = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUnauthorized    = errors.New("unauthorized")
)

type AuthUseCase struct {
    userRepo    repository.UserRepository
    authService auth.JWTService
    cache       repository.CacheRepository
}

func NewAuthUseCase(
    userRepo repository.UserRepository,
    authService auth.JWTService,
    cache repository.CacheRepository,
) *AuthUseCase {
    return &AuthUseCase{
        userRepo:    userRepo,
        authService: authService,
        cache:       cache,
    }
}

func (uc *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    
    existingUser, _ := uc.userRepo.FindByUsername(ctx, req.Username)
    if existingUser != nil {
        return nil, ErrUserExists
    }

    existingUser, _ = uc.userRepo.FindByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, ErrUserExists
    }

    
    user := &entity.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
        Role:     entity.RoleOperator,
        IsActive: true,
    }

    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    
    token, err := uc.authService.GenerateToken(user.ID, user.Username, string(user.Role))
    if err != nil {
        return nil, err
    }

    
    cacheKey := "user:" + user.ID.String()
    uc.cache.Set(ctx, cacheKey, user, 24*time.Hour)

    return &dto.AuthResponse{
        Token:     token,
        User:      user,
        ExpiresIn: 24 * 3600, 
    }, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := uc.userRepo.FindByUsername(ctx, req.Username)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if !user.ComparePassword(req.Password) {
        return nil, ErrInvalidCredentials
    }

    if !user.IsActive {
        return nil, errors.New("user account is inactive")
    }

    
    token, err := uc.authService.GenerateToken(user.ID, user.Username, string(user.Role))
    if err != nil {
        return nil, err
    }

    
    cacheKey := "user:" + user.ID.String()
    uc.cache.Set(ctx, cacheKey, user, 24*time.Hour)

    return &dto.AuthResponse{
        Token:     token,
        User:      user,
        ExpiresIn: 24 * 3600,
    }, nil
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
    
    cacheKey := "user:" + userID
    var user entity.User
    
    err := uc.cache.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }

    
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, err
    }

    dbUser, err := uc.userRepo.FindByID(ctx, userUUID)
    if err != nil {
        return nil, err
    }

    
    uc.cache.Set(ctx, cacheKey, dbUser, 24*time.Hour)

    return dbUser, nil
}