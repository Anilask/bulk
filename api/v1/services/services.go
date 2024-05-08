package services

import (
	"context"
	"time"

	"bulk/api/v1/repositories"
	"bulk/config"
	"bulk/models/request"
	disbursementpb "bulk/pb/disbursement"
	filepb "bulk/pb/file"
	gcs "bulk/pkg/google-cloud-storage"
	"bulk/pkg/pubsub"
	"bulk/logger"
)

//go:generate mockgen -source=./services.go -destination ./services_mock.go -package=services
type (
	Service struct {
		logger             logger.ILogger
		cfg                config.Config
		bulkFileRepo       repositories.BulkFileRepository
		bulkFileUploadRepo repositories.FileUploadRepository
		gcsRepo            gcs.Repository
		pubsubClient       pubsub.PubSubClient
	}

	Repository interface {
		UploadBulkFile(
			ctx context.Context,
			req *UploadBulkFileRequest,
		) (*filepb.UploadBulkFileResponse, error)

		GetSignedURLFromBulkFile(
			ctx context.Context,
			req *GetSignedURLFromBulkFileRequest,
		) (*GetSignedURLFromBulkFileResponse, error)

		Disburse(ctx context.Context, req *Disbursement) (*disbursementpb.DisbursementResponse, error)
		UpdateDisburse(ctx context.Context, req *request.UpdateDisbursement) (*disbursementpb.UpdateDisbursementResponse, error)
		UpdateBulkStatus(ctx context.Context, req *filepb.UpdateBulkStatusRequest) (*filepb.UpdateBulkStatusResponse, error)
		Disbursements(ctx context.Context, req *disbursementpb.DisbursementsRequest) (*disbursementpb.DisbursementsResponse, error)
		GetBulkList(ctx context.Context, req *filepb.BulkFilesRequest) (*filepb.BulkFilesResponse, error)
		Download(ctx context.Context, req *disbursementpb.DownloadDisbursementDataRequest) (*disbursementpb.DownloadDisbursementDataResponse, error)
	}

	GetSignedURLFromBulkFileRequest struct {
		Id int64
	}

	GetSignedURLFromBulkFileResponse struct {
		URL string
	}

	UploadBulkFileRequest struct {
		MerchantCode string
		Data         []byte
		FileName     string
		MimeType     string
		UserId       string
		FileSize     int64
		UploadAt     time.Time
		BulkName     string
	}
	Disbursement struct {
		MerchantUserId     string
		MerchantCode       string
		BulkDisbursementId string
	}
)

func (s *Service) WithLogger(l logger.ILogger) *Service {
	s.logger = l
	return s
}

func (s *Service) WithConfig(cfg config.Config) *Service {
	s.cfg = cfg
	return s
}

func (s *Service) WithBulkFileRepo(b repositories.BulkFileRepository) *Service {
	s.bulkFileRepo = b
	return s
}
func (s *Service) WithBulkFileUploadRepo(b repositories.FileUploadRepository) *Service {
	s.bulkFileUploadRepo = b
	return s
}
func (s *Service) WithGCSRepo(gcsRepo gcs.Repository) *Service {
	s.gcsRepo = gcsRepo
	return s
}
func (s *Service) WithPubSub(pubsubClient pubsub.PubSubClient) *Service {
	s.pubsubClient = pubsubClient
	return s
}
func New() *Service {
	return &Service{}
}
