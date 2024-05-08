package repositories

import "context"

//go:generate mockgen -source=./file_upload.go -destination ./file_upload_mock.go -package=repositories
type FileUploadRepository interface {
	UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error)
	GetSignedURL(
		ctx context.Context,
		req *GetSignedURLRequest,
	) (*GetSignedURLResponse, error)
}

type UploadFileRequest struct {
	Data         []byte
	MimeType     string
	FileName     string
	ParentFolder string
}

type UploadFileResponse struct {
	Path string
}

type GetSignedURLRequest struct {
	Path string
}

type GetSignedURLResponse struct {
	URL string
}
