package storage

import (
    "context"
    "fmt"
    "mime/multipart"
    "path"
    
    "github.com/google/uuid"
    "github.com/minio/minio-go/v7"
)

type StorageService interface {
    UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
    DeleteFile(ctx context.Context, objectKey string) error
    GetPublicURL(objectKey string) string 
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

    return objectKey, nil
}

func (s *minioStorage) DeleteFile(ctx context.Context, objectKey string) error {
    return s.client.RemoveObject(
        ctx,
        s.bucketName,
        objectKey,
        minio.RemoveObjectOptions{},
    )
}

func (s *minioStorage) GetPublicURL(objectKey string) string {
    return fmt.Sprintf("%s/%s", s.publicURL, objectKey)
}