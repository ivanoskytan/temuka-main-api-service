package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/temuka-api-service/internal/constant"
	fileStorage "github.com/temuka-api-service/util/file_storage"
)

type FileService interface {
	UploadFile(ctx context.Context, fileName string, fileData any) (string, error)
}

type FileServiceImpl struct {
	storage           fileStorage.S3Wrapper
	allowedExtensions map[string]bool
}

func NewFileService(storage fileStorage.S3Wrapper) FileService {
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

	s3Key := fmt.Sprintf("uploads/%s", fileName)

	reader, ok := fileData.(io.Reader)
	if !ok {
		return "", fmt.Errorf("invalid file reader")
	}

	if err := s.storage.UploadStream(ctx, s3Key, reader); err != nil {
		return "", err
	}

	var url string
	fmt.Printf("S3 Endpoint: %s", constant.EnvS3Endpoint)
	if constant.EnvS3Endpoint != "" {
		url = fmt.Sprintf("%s/%s/%s", constant.EnvS3Endpoint, constant.EnvS3Bucket, s3Key)
	} else {
		url = fmt.Sprintf(
			"https://%s.s3.%s.amazonaws.com/%s",
			constant.EnvS3Bucket,
			constant.EnvAWSRegion,
			s3Key,
		)
	}

	return url, nil
}
