package storage

import (
    "context"
    "fmt"
    "mime/multipart"
    "path"
    "strings"
    
    "github.com/google/uuid"
    "github.com/minio/minio-go/v7"
)

type StorageService interface {
    UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
    DeleteFile(ctx context.Context, objectKey string) error
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

func (s *minioStorage) UploadFile(
    ctx context.Context,
    file *multipart.FileHeader,
    folder string,
) (string, error) {
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    ext := path.Ext(file.Filename)
    objectKey := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

    _, err = s.client.PutObject(
        ctx,
        s.bucketName,
        objectKey,
        src,
        file.Size,
        minio.PutObjectOptions{
            ContentType: file.Header.Get("Content-Type"),
        },
    )
    if err != nil {
        return "", err
    }

    publicURL := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucketName, objectKey)
    return publicURL, nil
}

func (s *minioStorage) DeleteFile(ctx context.Context, fileURL string) error {
    objectKey := extractObjectKeyFromURL(fileURL, s.bucketName)
    if objectKey == "" {
        return fmt.Errorf("invalid file URL: %s", fileURL)
    }

    return s.client.RemoveObject(
        ctx,
        s.bucketName,
        objectKey,
        minio.RemoveObjectOptions{},
    )
}

func extractObjectKeyFromURL(fileURL, bucketName string) string {
    parts := strings.Split(fileURL, "/"+bucketName+"/")
    if len(parts) == 2 {
        return parts[1]
    }
    return fileURL
}