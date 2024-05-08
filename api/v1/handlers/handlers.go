package handlers

import (
	"context"

	"bulk/api/v1/services"
	"bulk/config"
	disbursementpb "bulk/pb/disbursement"
	filepb "bulk/pb/file"
	"bulk/logger"
)

type BulkDisbursementHandler struct {
	cfg     *config.Config
	log     logger.ILogger
	service services.Repository
	filepb.UnimplementedBulkFileHandlerServer
	disbursementpb.UnimplementedBulkDisbursementHandlerServer
}

type BulkDisbursementInterface interface {
	UploadBulkFile(ctx context.Context, req *filepb.UploadBulkFileRequest) (*filepb.UploadBulkFileResponse, error)
	Disburse(ctx context.Context, req *disbursementpb.DisbursementRequest) (*disbursementpb.DisbursementResponse, error)
	UpdateDisburseDetails(ctx context.Context, req *disbursementpb.UpdateDisbursementRequest) (*disbursementpb.UpdateDisbursementResponse, error)
	UpdateBulkFileStatus(ctx context.Context, req *filepb.UpdateBulkStatusRequest) (*filepb.UpdateBulkStatusResponse, error)
	GetBulkList(ctx context.Context, req *filepb.BulkFilesRequest) (*filepb.BulkFilesResponse, error)
	GetDisbursements(ctx context.Context, req *disbursementpb.DisbursementsRequest) (*disbursementpb.DisbursementsResponse, error)
	DownloadDisburse(ctx context.Context, req *disbursementpb.DownloadDisbursementDataRequest) (*disbursementpb.DownloadDisbursementDataResponse, error)
}

func (handler *BulkDisbursementHandler) WithLogger(l logger.ILogger) *BulkDisbursementHandler {
	handler.log = l

	return handler
}

func (handler *BulkDisbursementHandler) WithService(s services.Repository) *BulkDisbursementHandler {
	handler.service = s

	return handler
}
func (handler *BulkDisbursementHandler) WithConfig(c *config.Config) *BulkDisbursementHandler {
	handler.cfg = c
	return handler
}

func New() *BulkDisbursementHandler {
	return &BulkDisbursementHandler{}
}
