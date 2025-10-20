
    package config

    import (
        "log"
        "os"
        "strconv"
        "github.com/joho/godotenv"
    )

    type Config struct {
        App      AppConfig
        Database DatabaseConfig
        Redis    RedisConfig
        Minio    MinioConfig
        JWT      JWTConfig
    }

    type AppConfig struct {
        Env            string
        Port           string
        AllowedOrigins string
    }

    type DatabaseConfig struct {
        Host     string
        Port     string
        User     string
        Password string
        DBName   string
        SSLMode  string
    }

    type RedisConfig struct {
        Host     string
        Port     string
        Password string
        DB       int
    }

    type MinioConfig struct {
        Endpoint   string
        AccessKey  string
        SecretKey  string
        UseSSL     bool
        BucketName string
        PublicURL  string
    }

    type JWTConfig struct {
        Secret      string
        ExpiryHours int
    }

    func Load() *Config {
        err := godotenv.Load()
        if err != nil {
            log.Println("No .env file found, using environment variables")
        }

        return &Config{
            App: AppConfig{
                Env:            getEnv("APP_ENV", "development"),
                Port:           getEnv("APP_PORT", "8081"),
                AllowedOrigins: getEnv("APP_ALLOWED_ORIGINS", "http://localhost:3000"),
            },
            Database: DatabaseConfig{
                Host:     getEnv("DB_HOST", "localhost"),
                Port:     getEnv("DB_PORT", "5432"),
                User:     getEnv("DB_USER", "postgres"),
                Password: getEnv("DB_PASSWORD", "password"),
                DBName:   getEnv("DB_NAME", "building_reports"),
                SSLMode:  getEnv("DB_SSL_MODE", "disable"),
            },
            Redis: RedisConfig{
                Host:     getEnv("REDIS_HOST", "localhost"),
                Port:     getEnv("REDIS_PORT", "6379"),
                Password: getEnv("REDIS_PASSWORD", ""),
                DB:       getEnvAsInt("REDIS_DB", 0),
            },
            Minio: MinioConfig{
                Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
                AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
                SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
                UseSSL:     getEnvAsBool("MINIO_USE_SSL", false),
                BucketName: getEnv("MINIO_BUCKET_NAME", "reports"),
                PublicURL:  getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
            },
            JWT: JWTConfig{
                Secret:      getEnv("JWT_SECRET", "your-secret-key-here"),
                ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
            },
        }
    }

    func getEnv(key, defaultValue string) string {
        if value := os.Getenv(key); value != "" {
            return value
        }
        return defaultValue
    }

    func getEnvAsInt(key string, defaultValue int) int {
        valueStr := getEnv(key, "")
        if value, err := strconv.Atoi(valueStr); err == nil {
            return value
        }
        return defaultValue
    }

    func getEnvAsBool(key string, defaultValue bool) bool {
        valueStr := getEnv(key, "")
        if value, err := strconv.ParseBool(valueStr); err == nil {
            return value
        }
        return defaultValue
    }