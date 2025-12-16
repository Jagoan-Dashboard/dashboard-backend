
package storage

import (
    "context"
    "log"
    "building-report-backend/pkg/config"
    
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg config.MinioConfig) (*minio.Client, error) {
    client, err := minio.New(cfg.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
        Secure: cfg.UseSSL,
    })
    if err != nil {
        return nil, err
    }

    ctx := context.Background()
    exists, err := client.BucketExists(ctx, cfg.BucketName)
    if err != nil {
        return nil, err
    }

    if !exists {
        if err := client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{}); err != nil {
            return nil, err
        }
        log.Printf("Bucket %s created", cfg.BucketName)
    }

    return client, nil
}