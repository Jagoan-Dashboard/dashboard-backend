
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
        err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
        if err != nil {
            return nil, err
        }
        log.Printf("Bucket %s created successfully", cfg.BucketName)
        
        
        policy := `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": "*",
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::` + cfg.BucketName + `/*"]
                }
            ]
        }`
        
        err = client.SetBucketPolicy(ctx, cfg.BucketName, policy)
        if err != nil {
            log.Printf("Failed to set bucket policy: %v", err)
        }
    }

    return client, nil
}