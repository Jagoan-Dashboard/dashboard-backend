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
    ErrForbidden          = errors.New("forbidden: insufficient permissions")
    ErrCannotDeleteSelf   = errors.New("cannot delete your own account")
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
    
    if err := uc.validateRegistrationInput(req); err != nil {
        return nil, err
    }

    
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
        Role:     entity.RoleUser, 
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
    
    if !utils.IsValidULID(userID) {
        return nil, ErrInvalidCredentials
    }

    
    cacheKey := constants.UserCachePrefix + userID
    var user entity.User

    err := uc.cache.Get(ctx, cacheKey, &user)
    if err == nil {
        return &user, nil
    }

    
    dbUser, err := uc.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, ErrUserNotFound
    }

    
    uc.cache.Set(ctx, cacheKey, dbUser, constants.UserCacheDuration)

    return dbUser, nil
}

func (uc *AuthUseCase) GetAllUsers(ctx context.Context, requesterID string) (*dto.UserListResponse, error) {
    
    requester, err := uc.GetUserByID(ctx, requesterID)
    if err != nil {
        return nil, err
    }

    if !requester.IsSuperAdmin() {
        return nil, ErrForbidden
    }

    
    users, total, err := uc.userRepo.FindAll(ctx, 0, 0)
    if err != nil {
        return nil, err
    }

    
    userResponses := make([]*dto.UserResponse, len(users))
    for i, user := range users {
        userResponses[i] = &dto.UserResponse{
            ID:        user.ID,
            Username:  user.Username,
            Email:     user.Email,
            Role:      user.Role,
            IsActive:  user.IsActive,
            CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
        }
    }

    return &dto.UserListResponse{
        Users: userResponses,
        Total: total,
    }, nil
}


func (uc *AuthUseCase) GetUserDetailByID(ctx context.Context, requesterID, targetUserID string) (*dto.UserResponse, error) {
    
    requester, err := uc.GetUserByID(ctx, requesterID)
    if err != nil {
        return nil, err
    }

    if !requester.IsSuperAdmin() {
        return nil, ErrForbidden
    }

    
    if !utils.IsValidULID(targetUserID) {
        return nil, ErrInvalidCredentials
    }

    
    targetUser, err := uc.GetUserByID(ctx, targetUserID)
    if err != nil {
        return nil, ErrUserNotFound
    }

    return &dto.UserResponse{
        ID:        targetUser.ID,
        Username:  targetUser.Username,
        Email:     targetUser.Email,
        Role:      targetUser.Role,
        IsActive:  targetUser.IsActive,
        CreatedAt: targetUser.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: targetUser.UpdatedAt.Format("2006-01-02 15:04:05"),
    }, nil
}


func (uc *AuthUseCase) CreateUser(ctx context.Context, requesterID string, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    
    requester, err := uc.GetUserByID(ctx, requesterID)
    if err != nil {
        return nil, err
    }

    if !requester.IsSuperAdmin() {
        return nil, ErrForbidden
    }

    
    if err := validation.ValidateRequired(req.Username, "username"); err != nil {
        return nil, err
    }

    if err := validation.ValidateUsername(req.Username); err != nil {
        return nil, err
    }

    if err := validation.ValidateRequired(req.Email, "email"); err != nil {
        return nil, err
    }

    if err := validation.ValidateEmail(req.Email); err != nil {
        return nil, err
    }

    if err := validation.ValidateRequired(req.Password, "password"); err != nil {
        return nil, err
    }

    if err := validation.ValidatePassword(req.Password); err != nil {
        return nil, err
    }

    
    if err := uc.validateUserNotExists(ctx, req.Username, req.Email); err != nil {
        return nil, err
    }

    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    
    newUser := &entity.User{
        ID:       utils.GenerateULID(),
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        Role:     req.Role,
        IsActive: true,
    }

    if err := uc.userRepo.Create(ctx, newUser); err != nil {
        return nil, err
    }

    
    cacheKey := constants.UserCachePrefix + newUser.ID
    uc.cache.Set(ctx, cacheKey, newUser, constants.UserCacheDuration)

    return &dto.UserResponse{
        ID:        newUser.ID,
        Username:  newUser.Username,
        Email:     newUser.Email,
        Role:      newUser.Role,
        IsActive:  newUser.IsActive,
        CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
    }, nil
}


func (uc *AuthUseCase) UpdateUserRole(ctx context.Context, requesterID, targetUserID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
    
    requester, err := uc.GetUserByID(ctx, requesterID)
    if err != nil {
        return nil, err
    }

    if !requester.IsSuperAdmin() {
        return nil, ErrForbidden
    }

    
    if !utils.IsValidULID(targetUserID) {
        return nil, ErrInvalidCredentials
    }

    
    targetUser, err := uc.GetUserByID(ctx, targetUserID)
    if err != nil {
        return nil, ErrUserNotFound
    }

    
    targetUser.Role = req.Role

    
    if err := uc.userRepo.Update(ctx, targetUser); err != nil {
        return nil, err
    }

    
    cacheKey := constants.UserCachePrefix + targetUser.ID
    uc.cache.Set(ctx, cacheKey, targetUser, constants.UserCacheDuration)

    return &dto.UserResponse{
        ID:        targetUser.ID,
        Username:  targetUser.Username,
        Email:     targetUser.Email,
        Role:      targetUser.Role,
        IsActive:  targetUser.IsActive,
        CreatedAt: targetUser.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: targetUser.UpdatedAt.Format("2006-01-02 15:04:05"),
    }, nil
}


func (uc *AuthUseCase) DeleteUser(ctx context.Context, requesterID, targetUserID string) error {
    
    requester, err := uc.GetUserByID(ctx, requesterID)
    if err != nil {
        return err
    }

    if !requester.IsSuperAdmin() {
        return ErrForbidden
    }

    
    if requesterID == targetUserID {
        return ErrCannotDeleteSelf
    }

    
    if !utils.IsValidULID(targetUserID) {
        return ErrInvalidCredentials
    }

    
    targetUser, err := uc.GetUserByID(ctx, targetUserID)
    if err != nil {
        return ErrUserNotFound
    }

    
    if err := uc.userRepo.Delete(ctx, targetUser.ID); err != nil {
        return err
    }

    
    cacheKey := constants.UserCachePrefix + targetUser.ID
    uc.cache.Delete(ctx, cacheKey)

    return nil
}