package seeds

import (
	"building-report-backend/internal/domain/entity"
	"building-report-backend/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
func SeedUsers(db *gorm.DB) error {
	
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return nil 
	}

	users := []entity.User{
		{
			ID:        utils.GenerateULID(),
			Username:  "admin",
			Email:     "admin@ngawikab.go.id",
			Password:  hashPassword("Admin123!"),
			Role:      entity.RoleAdmin,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        utils.GenerateULID(),
			Username:  "supervisor",
			Email:     "supervisor@ngawikab.go.id",
			Password:  hashPassword("Super123!"),
			Role:      entity.RoleSupervisor,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        utils.GenerateULID(),
			Username:  "operator_pertanian",
			Email:     "pertanian@ngawikab.go.id",
			Password:  hashPassword("Oper123!"),
			Role:      entity.RoleOperator,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        utils.GenerateULID(),
			Username:  "operator_pu",
			Email:     "pu@ngawikab.go.id",
			Password:  hashPassword("Oper123!"),
			Role:      entity.RoleOperator,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        utils.GenerateULID(),
			Username:  "operator_tataruang",
			Email:     "tataruang@ngawikab.go.id",
			Password:  hashPassword("Oper123!"),
			Role:      entity.RoleOperator,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        utils.GenerateULID(),
			Username:  "viewer",
			Email:     "viewer@ngawikab.go.id",
			Password:  hashPassword("View123!"),
			Role:      entity.RoleViewer,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		
		if err := db.Session(&gorm.Session{SkipHooks: true}).Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}