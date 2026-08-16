package file_storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureBlobWrapper struct {
	Client    *azblob.Client
	Container string
}

func NewAzureBlob(accountName, accountKey, containerName string) (*AzureBlobWrapper, error) {
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %w", err)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azure blob client: %w", err)
	}

	log.Println("Azure Blob Storage client successfully created")

	return &AzureBlobWrapper{
		Client:    client,
		Container: containerName,
	}, nil
}

func (a *AzureBlobWrapper) UploadByte(ctx context.Context, blobName string, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := a.Client.UploadStream(ctx, a.Container, blobName, reader, nil)
	if err != nil {
		return fmt.Errorf("failed to upload byte buffer to Azure Blob: %w", err)
	}

	return nil
}

func (a *AzureBlobWrapper) UploadStream(ctx context.Context, blobName string, data io.Reader) error {
	_, err := a.Client.UploadStream(ctx, a.Container, blobName, data, nil)
	if err != nil {
		return fmt.Errorf("failed to upload stream to Azure Blob: %w", err)
	}

	return nil
}

func (a *AzureBlobWrapper) Download(ctx context.Context, blobName string) ([]byte, error) {
	get, err := a.Client.DownloadStream(ctx, a.Container, blobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob from Azure Blob: %w", err)
	}
	defer get.Body.Close()

	data, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read downloaded blob data: %w", err)
	}

	return data, nil
}

func (a *AzureBlobWrapper) Delete(ctx context.Context, blobName string) error {
	_, err := a.Client.DeleteBlob(ctx, a.Container, blobName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob from Azure Blob: %w", err)
	}

	return nil
}
