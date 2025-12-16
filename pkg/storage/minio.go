package storage

import (
	"building-report-backend/pkg/config"
	"context"
	"fmt"
	"log"

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
        log.Printf("✅ Bucket %s created", cfg.BucketName)
    }

    policy := fmt.Sprintf(`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {"AWS": ["*"]},
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::%s/*"]
            }
        ]
    }`, cfg.BucketName)

    err = client.SetBucketPolicy(ctx, cfg.BucketName, policy)
    if err != nil {
        log.Printf("arning: Could not set bucket policy: %v", err)
    } else {
        log.Printf("Bucket %s is now publicly readable", cfg.BucketName)
    }

    return client, nil
}