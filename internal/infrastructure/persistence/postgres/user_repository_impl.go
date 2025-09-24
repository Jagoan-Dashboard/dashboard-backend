
package postgres

import (
    "context"
    "building-report-backend/internal/domain/entity"
    "building-report-backend/internal/domain/repository"
    
    "gorm.io/gorm"
)

type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{}).Error
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, int64, error) {
    var users []*entity.User
    var total int64

    query := r.db.WithContext(ctx).Model(&entity.User{})
    query.Count(&total)

    err := query.Limit(limit).Offset(offset).Find(&users).Error
    return users, total, err
}