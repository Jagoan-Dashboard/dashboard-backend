package usecase

import (
	"context"
	"errors"

	"building-report-backend/internal/application/dto"
	"building-report-backend/internal/domain/constants"
	"building-report-backend/internal/domain/entity"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/auth"
	"building-report-backend/pkg/utils"
	"building-report-backend/pkg/validation"

	"golang.org/x/crypto/bcrypt"
)

var (
    ErrUserExists         = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrUserNotFound       = errors.New("user not found")
    ErrInactiveUser       = errors.New("user account is inactive")
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
    // Validate input
    if err := uc.validateRegistrationInput(req); err != nil {
        return nil, err
    }

    // Check if user already exists
    if err := uc.validateUserNotExists(ctx, req.Username, req.Email); err != nil {
        return nil, err
    }

    password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &entity.User{
        ID:       utils.GenerateULID(),
        Username: req.Username,
        Email:    req.Email,
        Password: string(password), 
        Role:     entity.RoleOperator,
        IsActive: true,
    }

    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return uc.generateAuthResponse(ctx, user)
}

func (uc *AuthUseCase) validateRegistrationInput(req *dto.RegisterRequest) error {
    if err := validation.ValidateRequired(req.Username, "username"); err != nil {
        return err
    }

    if err := validation.ValidateUsername(req.Username); err != nil {
        return err
    }

    if err := validation.ValidateRequired(req.Email, "email"); err != nil {
        return err
    }

    if err := validation.ValidateEmail(req.Email); err != nil {
        return err
    }

    if err := validation.ValidateRequired(req.Password, "password"); err != nil {
        return err
    }

    if err := validation.ValidatePassword(req.Password); err != nil {
        return err
    }

    return nil
}

func (uc *AuthUseCase) validateUserNotExists(ctx context.Context, username, email string) error {
    if existingUser, _ := uc.userRepo.FindByUsername(ctx, username); existingUser != nil {
        return ErrUserExists
    }

    if existingUser, _ := uc.userRepo.FindByEmail(ctx, email); existingUser != nil {
        return ErrUserExists
    }

    return nil
}

func (uc *AuthUseCase) generateAuthResponse(ctx context.Context, user *entity.User) (*dto.AuthResponse, error) {
    token, err := uc.authService.GenerateToken(user.ID, user.Username, string(user.Role))
    if err != nil {
        return nil, err
    }

    // Cache user
    cacheKey := constants.UserCachePrefix + user.ID
    uc.cache.Set(ctx, cacheKey, user, constants.UserCacheDuration)

    return &dto.AuthResponse{
        Token:     token,
        User:      user,
        ExpiresIn: constants.TokenExpirationSecs,
    }, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {

    user, err := uc.userRepo.FindByUsernameOrEmail(ctx, req.Identifier)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if !user.ComparePassword(req.Password) {
        return nil, ErrInvalidCredentials
    }

    if !user.IsActive {
        return nil, ErrInactiveUser
    }

    return uc.generateAuthResponse(ctx, user)
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
    // Validate ULID format
    if !utils.IsValidULID(userID) {
        return nil, ErrInvalidCredentials
    }

    // Try to get from cache first
    cacheKey := constants.UserCachePrefix + userID
    var user entity.User

    err := uc.cache.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }

    // Get from database
    dbUser, err := uc.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, ErrUserNotFound
    }

    // Cache the result
    uc.cache.Set(ctx, cacheKey, dbUser, constants.UserCacheDuration)

    return dbUser, nil
}