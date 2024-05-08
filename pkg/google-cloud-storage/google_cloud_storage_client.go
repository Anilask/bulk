package google_cloud_storage

import (
	"bulk/config"
	"bulk/logger"
	"context"

	"cloud.google.com/go/storage"
)

//go:generate mockgen -source google_cloud_storage_client.go -destination google_cloud_storage_client_mock.go -package google_cloud_storage
type (
	Repository interface {
		SignedUrl(ctx context.Context, filePath string) (string, error)
	}

	GoogleCloudStorage struct {
		client *storage.Client
		cfg    config.Config
		log    logger.ILogger
	}
)

func (gcs *GoogleCloudStorage) WithLogger(logger logger.ILogger) *GoogleCloudStorage {
	gcs.log = logger
	return gcs
}

func (gcs *GoogleCloudStorage) WithConfig(cfg config.Config) *GoogleCloudStorage {
	gcs.cfg = cfg
	return gcs
}

func (gcs *GoogleCloudStorage) WithClient(client *storage.Client) *GoogleCloudStorage {
	gcs.client = client
	return gcs
}

func New() *GoogleCloudStorage {
	return &GoogleCloudStorage{}
}
