package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bulk/tracer"
	"cloud.google.com/go/storage"

	"bulk/api/v1/repositories"
	"bulk/config"
	"bulk/utils"
)

type FileUploadRepository struct {
	objectPrefix   string
	bucketName     string
	googleAccessID string
	cl             *storage.Client
}

func NewFileUploadRepository(cfg *config.GCS) (*FileUploadRepository, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create gcs client: %w", err)
	}
	return &FileUploadRepository{
		objectPrefix:   "files",
		bucketName:     cfg.Bucket,
		googleAccessID: cfg.GoogleAccessId,
		cl:             client,
	}, nil
}

func (s *FileUploadRepository) UploadFile(ctx context.Context, req *repositories.UploadFileRequest) (*repositories.UploadFileResponse, error) {
	ctx, span := tracer.StartSpan(ctx, "UploadFile")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	path := s.getObjectName(req.FileName, req.ParentFolder)
	wc := s.cl.Bucket(s.bucketName).Object(path).NewWriter(ctx)
	if req.MimeType != "" {
		wc.ContentType = req.MimeType
	}

	if _, err := wc.Write(req.Data); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to write file to gcs: %w", err)
		
	}

	if err := wc.Close(); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to close gcs writer: %w", err)
	}

	return &repositories.UploadFileResponse{
		Path: path,
	}, nil
}

func (s *FileUploadRepository) getObjectName(name string, folder string) string {
	if folder != "" {
		return fmt.Sprintf("%s/%s/%s", s.objectPrefix, folder, name)
	}
	return fmt.Sprintf("%s/%s", s.objectPrefix, name)
}

func (s *FileUploadRepository) GetSignedURL(
	ctx context.Context,
	req *repositories.GetSignedURLRequest,
) (*repositories.GetSignedURLResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "GetSignedURL")
	defer span.End()

	_, cancel := context.WithTimeout(childCtx, 10*time.Second)
	defer cancel()

	path := req.Path
	// create signed url for the object
	url, err := s.cl.Bucket(s.bucketName).SignedURL(path, &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		Expires:        utils.GetTime().Add(15 * time.Minute),
		GoogleAccessID: s.googleAccessID,
	})
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get signed url: %w", err)
	}

	return &repositories.GetSignedURLResponse{
		URL: url,
	}, nil
}
