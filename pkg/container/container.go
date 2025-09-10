// pkg/container/container.go
package container

import (
	"building-report-backend/internal/application/usecase"
	"building-report-backend/internal/domain/repository"
	"building-report-backend/internal/infrastructure/auth"
	"building-report-backend/internal/infrastructure/persistence/postgres"
	"building-report-backend/internal/infrastructure/storage"
	"building-report-backend/internal/interfaces/http/handler"
	"building-report-backend/pkg/config"
    "github.com/redis/go-redis/v9"
    redisPkg "building-report-backend/internal/infrastructure/persistence/redis"
	
	"github.com/minio/minio-go/v7"
	
	"gorm.io/gorm"
)

type Container struct {
    Config         *config.Config
    DB             *gorm.DB
    Redis          *redis.Client
    MinioClient    *minio.Client
    
    // Repositories
    UserRepo       repository.UserRepository
    ReportRepo     repository.ReportRepository
    CacheRepo      repository.CacheRepository
    
    // Services
    StorageService storage.StorageService
    AuthService    auth.JWTService
    
    // Use Cases
    AuthUseCase    *usecase.AuthUseCase
    ReportUseCase  *usecase.ReportUseCase
    
    // Handlers
    AuthHandler    *handler.AuthHandler
    ReportHandler  *handler.ReportHandler
}

func NewContainer(cfg *config.Config, db *gorm.DB, redisClient *redis.Client, minioClient *minio.Client) *Container {
    container := &Container{
        Config:      cfg,
        DB:          db,
        Redis:       redisClient,
        MinioClient: minioClient,
    }

    // Initialize repositories
    container.UserRepo = postgres.NewUserRepository(db)
    container.ReportRepo = postgres.NewReportRepository(db)
    container.CacheRepo = redisPkg.NewCacheRepository(redisClient)

    // Initialize services
    container.StorageService = storage.NewMinioStorage(
        minioClient,
        cfg.Minio.BucketName,
        cfg.Minio.PublicURL,
    )
    container.AuthService = auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

    // Initialize use cases
    container.AuthUseCase = usecase.NewAuthUseCase(
        container.UserRepo,
        container.AuthService,
        container.CacheRepo,
    )
    container.ReportUseCase = usecase.NewReportUseCase(
        container.ReportRepo,
        container.StorageService,
        container.CacheRepo,
    )

    // Initialize handlers
    container.AuthHandler = handler.NewAuthHandler(container.AuthUseCase)
    container.ReportHandler = handler.NewReportHandler(container.ReportUseCase)

    return container
}