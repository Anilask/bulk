package repositories

import (
	"context"

	"bulk/models/request"
	disbursementpb "bulk/pb/disbursement"
	filepb "bulk/pb/file"
)

//go:generate mockgen -source=./bulk_file.go -destination ./bulk_file_mock.go -package=repositories
type BulkFileRepository interface {
	CreateBulkFile(ctx context.Context, req *CreateBulkFileRequest) (*BulkFileModel, error)
	UpdateBulkFile(ctx context.Context, req *UpdateBulkFileRequest) error
	GetBulkFile(ctx context.Context, req *GetBulkFileRequest) (*BulkFileModel, error)
	GetBulkData(ctx context.Context, req *GetBulkFileRequest) ([]request.BulkDisbursementDetails, error)
	UpdateBulkDisbursementData(ctx context.Context, req *request.UpdateDisbursement) (int64, error)
	UpdateBulkFileStatus(ctx context.Context, req *filepb.UpdateBulkStatusRequest) (int64, error)
	GetBulkList(ctx context.Context, req *filepb.BulkFilesRequest) ([]*filepb.BulkFile, int, error)
	GetDisbursements(ctx context.Context, bulkID int, req *disbursementpb.DisbursementsRequest) ([]*disbursementpb.Disbursement, int, error)
	GetBulkDisbursementStatusCount(ctx context.Context, bulkID int) (*disbursementpb.BulkDisbursementStatusCount, error)
	GetBulkFileByBulkID(ctx context.Context, bulkID string) (*BulkFileModel, error)
}

type GetBulkFileRequest struct {
	Id     int64
	Bulkid string
	Status int32
}

type BulkFileModel struct {
	Id       int64
	BulkID   string
	BulkName string
	FileName string
	FilePath string
	FileSize int64
	UserId   string
	Status   int64
	Created  string
}

type CreateBulkFileRequest struct {
	MerchantCode string
	BulkId       string
	Name         string
	FileName     string
	FilePath     string
	FileSize     int64
	UserId       string
	Status       BulkStatus
}

type BulkStatus int64

func (s BulkStatus) String() string {
	return [...]string{"PENDING", "SUCCESS", "FAILED"}[s-1]
}

func (s BulkStatus) Int64() int64 {
	return int64(s)
}

const (
	PendingStatus BulkStatus = 1
	SuccessStatus BulkStatus = 2
	FailedStatus  BulkStatus = 9
)

type UpdateBulkFileRequest struct {
	Id       int64
	Status   BulkStatus
	FilePath string
}

type BulkDisbursementStatusCount struct {
	InquiryInitiated  int
	InquirySuccess    int
	TransferInitiated int
	TransferSuccess   int
	Failed            int
}
