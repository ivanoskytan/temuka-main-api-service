package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/temuka-api-service/internal/constant"
	fileStorage "github.com/temuka-api-service/util/file_storage"
)

type FileService interface {
	UploadFile(ctx context.Context, fileName string, fileData any) (string, error)
}

type FileServiceImpl struct {
	storage           fileStorage.AzureBlobWrapper
	allowedExtensions map[string]bool
}

func NewFileService(storage fileStorage.AzureBlobWrapper) FileService {
	return &FileServiceImpl{
		storage: storage,
		allowedExtensions: map[string]bool{
			".jpg":  true,
			".png":  true,
			".mp4":  true,
			".mkv":  true,
			".jpeg": true,
		},
	}
}

func (s *FileServiceImpl) UploadFile(ctx context.Context, fileName string, fileData any) (string, error) {
	ext := filepath.Ext(fileName)
	if !s.allowedExtensions[ext] {
		return "", fmt.Errorf("file type not allowed")
	}

	blobName := fmt.Sprintf("uploads/%s", fileName)

	reader, ok := fileData.(io.Reader)
	if !ok {
		return "", fmt.Errorf("invalid file reader")
	}

	if err := s.storage.UploadStream(ctx, blobName, reader); err != nil {
		return "", err
	}

	var url string
	if constant.EnvAzureStorageAccountName == "" {
		url = fmt.Sprintf("%s/%s/%s", os.Getenv(constant.EnvS3Endpoint), os.Getenv(constant.EnvS3Bucket), blobName)
	} else {
		url = fmt.Sprintf(
			"https://%s.blob.core.windows.net/%s/%s",
			os.Getenv(constant.EnvAzureStorageAccountName),
			os.Getenv(constant.EnvAzureStorageContainerName),
			blobName,
		)
	}

	return url, nil
}
