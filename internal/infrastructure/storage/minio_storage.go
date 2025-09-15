
package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type StorageService interface {
    UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
    DeleteFile(ctx context.Context, fileURL string) error
    GetFileURL(ctx context.Context, objectName string) (string, error)
}

type minioStorage struct {
    client     *minio.Client
    bucketName string
    publicURL  string
}

func NewMinioStorage(client *minio.Client, bucketName, publicURL string) StorageService {
    return &minioStorage{
        client:     client,
        bucketName: bucketName,
        publicURL:  publicURL,
    }
}

func (s *minioStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    
    ext := path.Ext(file.Filename)
    objectName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

    
    _, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
        ContentType: file.Header.Get("Content-Type"),
    })
    if err != nil {
        return "", err
    }

    
    return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, objectName), nil
}

func (s *minioStorage) DeleteFile(ctx context.Context, fileURL string) error {
    
    objectName := extractObjectName(fileURL)
    
    return s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *minioStorage) GetFileURL(ctx context.Context, objectName string) (string, error) {
    
    url, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, 7*24*time.Hour, nil)
    if err != nil {
        return "", err
    }
    return url.String(), nil
}

func extractObjectName(fileURL string) string {
    
    parsedURL, err := url.Parse(fileURL)
    if err != nil {
        return ""
    }
    
    
    parts := strings.SplitN(parsedURL.Path, "/", 3)
    if len(parts) >= 3 {
        return parts[2]
    }
    
    return ""
}