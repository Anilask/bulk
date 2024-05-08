package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"strconv"
	"time"

	errors "bulk/grpc-errors"

	"bulk/api/v1/repositories"
	"bulk/constant"
	"bulk/models/request"
	filepb "bulk/pb/file"
	"bulk/tracer"
	"bulk/utils"

	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"
)

var (
	ErrWrongBulkFileFormat = "wrong bulk file format"
)

const (
	BenficiaryAccNo       = "Beneficiary Account Number"
	BeneficiaryBankCode   = "Beneficiary Bank Code"
	MobileNumber          = "Mobile Number"
	Amount                = "Amount"
	TransferNotes         = "Transfer Notes"
	MissingRequiredColumn = "missing required column"
)

func (s *Service) UploadBulkFile(
	ctx context.Context,
	req *UploadBulkFileRequest,
) (*filepb.UploadBulkFileResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "CreateBulkFile")
	defer span.End()
	var err error
	var link string
	bulkId := utils.GenerateNewUniqueNumber()
	// create job record
	f, err := s.bulkFileRepo.CreateBulkFile(childCtx, &repositories.CreateBulkFileRequest{
		MerchantCode: req.MerchantCode,
		BulkId:       bulkId,
		Name:         req.BulkName,
		FileName:     req.FileName,
		FilePath:     "N/A",
		FileSize:     req.FileSize,
		UserId:       req.UserId,
		Status:       repositories.PendingStatus,
	})
	if err != nil {
		span.RecordError(err)
		// create custom error
		s.logger.Errorf("failed to update bulk file record: %s", err.Error())
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	defer func(err *error, link *string) {
		if *err != nil {
			// update job record
			if err := s.bulkFileRepo.UpdateBulkFile(childCtx, &repositories.UpdateBulkFileRequest{
				Id:       f.Id,
				Status:   repositories.FailedStatus,
				FilePath: *link,
			}); err != nil {
				span.RecordError(err)
				s.logger.Errorf("failed to update bulk file record: %s", err.Error())
			}
		}
	}(&err, &link)

	// parse raw data to struct
	// validate data structure
	var bulkData []request.BulkData
	s.logger.Debugf("Before parseBulkFileData| Data: %v | MimeType:", req.Data, req.MimeType)
	bulkData, err = s.parseBulkFileData(childCtx, req.Data, req.MimeType)
	s.logger.Debugf("parseBulkFileData %v", bulkData)
	if err != nil {
		span.RecordError(err)
		// upload to gcs
		var uploadErr error
		link, uploadErr = s.uploadBulkFile(
			childCtx,
			req.FileName,
			req.MimeType,
			req.Data,
			req.UploadAt,
			repositories.FailedStatus,
		)
		if uploadErr != nil {
			span.RecordError(uploadErr)
			s.logger.Errorf("failed to parse data: %s", err.Error())
			return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
		}
		return nil, err
	}
	inquiry := request.Inquiry{
		TransactionId:   childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber: childCtx.Value(constant.CtxReferenceNumber).(string),
		CorrelationId:   childCtx.Value(constant.HeaderCorrelationID).(string),
		MerchantCode:    req.MerchantCode,
		BulkId:          f.Id,
		Data:            bulkData,
	}
	inquiryObj, ErrMarshalinquiry := json.Marshal(inquiry)
	if ErrMarshalinquiry != nil {
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}
	// // Create event to validate, insert data into db and do inquiries
	s.logger.Debugf("Bulk File Data %v", bulkData)
	s.pubsubClient.PublishData(ctx, s.cfg.Pubsub.InquiryTopic, inquiryObj, nil)

	// upload to gcs
	link, err = s.uploadBulkFile(
		childCtx,
		req.FileName,
		req.MimeType,
		req.Data,
		req.UploadAt,
		repositories.SuccessStatus,
	)
	if err != nil {
		span.RecordError(err)
		s.logger.Errorf("failed to upload to gcs: %s", err.Error())
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	// update job status to success
	err = s.bulkFileRepo.UpdateBulkFile(childCtx, &repositories.UpdateBulkFileRequest{
		Id:       f.Id,
		Status:   repositories.SuccessStatus,
		FilePath: link,
	})
	if err != nil {
		span.RecordError(err)
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	return &filepb.UploadBulkFileResponse{
		ResponseCode:    "200",
		ResponseMessage: "ok",
		ResponseTime:    utils.GetTime().String(),
		TransactionId:   childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber: childCtx.Value(constant.CtxReferenceNumber).(string),
		File: &filepb.File{
			BulkId:         bulkId,
			UploadAt:       req.UploadAt.String(),
			FileName:       req.FileName,
			FileSize:       req.FileSize,
			MerchantUserId: req.UserId,
			Status:         2,
		},
	}, err
}

func (s *Service) UpdateBulkStatus(ctx context.Context, req *filepb.UpdateBulkStatusRequest) (*filepb.UpdateBulkStatusResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "UpdateBulkStatus")
	defer span.End()
	if _, err := s.bulkFileRepo.UpdateBulkFileStatus(childCtx, req); err != nil {
		span.RecordError(err)
		s.logger.Errorf("failed to update bulk status: %s", err.Error())
		return nil, err
	}
	return &filepb.UpdateBulkStatusResponse{
		ResponseCode:    "200",
		ResponseMessage: "Success",
		ResponseTime:    utils.GetGMT7DateTime(constant.TimestampFormat),
		TransactionId:   req.TransactionId,
		ReferenceNumber: childCtx.Value(constant.CtxReferenceNumber).(string),
		BulkId:          req.BulkId,
		Status:          req.Status,
	}, nil
}

func (s *Service) parseBulkFileData(childCtx context.Context, data []byte, mimeType string) ([]request.BulkData, error) {
	switch mimeType {
	case utils.CSVMimeType:
		return s.parseBulkFileCSV(childCtx, data)
	case utils.XLSXMimeType:
		return s.parseBulkFileXLSX(childCtx, data)
	case utils.XLSMimeType:
		return s.parseBulkFileXLS(childCtx, data)
	default:
		return nil, errors2.New("unsupported file type")
	}
}

func (s *Service) parseBulkFileCSV(childCtx context.Context, raw []byte) ([]request.BulkData, error) {
	s.logger.Debugf("Inside parseBulkFileCSV")
	reader := csv.NewReader(bytes.NewReader(raw))
	s.logger.Debugf("parseBulkFileCSV | reader: %v", reader)
	columns, err := reader.Read()
	if err != nil {
		s.logger.Errorf("Error while reading headers %v", err)
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "BadHeader0302")
	}

	s.logger.Debugf("parseBulkFileCSV | columns: %v", columns)
	if len(columns) == 0 {
		s.logger.Error("No columns found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	_, err = s.validateHeaders(childCtx, columns)
	if err != nil {
		return nil, err
	}
	data, err := reader.ReadAll()
	if err != nil {
		s.logger.Errorf("failed to read csv data: %v", err)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}
	s.logger.Debugf("parseBulkFileCSV | data: %v", data)
	resp := ReadDatafromCSV(data)
	if resp == nil {
		s.logger.Errorf("failed to append csv data to struct: %v", err)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}
	s.logger.Debugf("parseBulkFileCSV | resp: %v", resp)
	return resp, nil
}
func ReadDatafromCSV(data [][]string) []request.BulkData {
	var resp []request.BulkData
	for _, line := range data {
		var rec request.BulkData
		for j, field := range line {
			if j == 0 {
				rec.BeneficiaryBankCode = field
			} else if j == 1 {
				rec.BeneficiaryAccountNumber = field
			} else if j == 2 {
				rec.MobileNumber = field
			} else if j == 3 {
				f, err := strconv.ParseFloat(field, 8)
				if err != nil {
					return nil
					break
				}
				rec.Amount = f
			} else if j == 4 {
				rec.PaymentInfo = field
			}
			rec.InquiryCorrelationId = utils.GenerateNewUniqueNumber()
		}
		resp = append(resp, rec)

	}
	return resp
}
func (s *Service) uploadBulkFile(
	ctx context.Context,
	name string,
	mimeType string,
	data []byte,
	now time.Time,
	status repositories.BulkStatus,
) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	name = fmt.Sprintf("%d-%s", now.Unix(), name)
	resp, err := s.bulkFileUploadRepo.UploadFile(ctx, &repositories.UploadFileRequest{
		Data:         data,
		MimeType:     mimeType,
		FileName:     name,
		ParentFolder: status.String(),
	})
	if err != nil {
		s.logger.Errorf("failed to upload bulk file: %v", err)
		return "", errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	return resp.Path, nil
}
func (s *Service) parseBulkFileXLSX(childCtx context.Context, data []byte) ([]request.BulkData, error) {
	resp := make([]request.BulkData, 0, 100)
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		s.logger.Errorf("failed to open excel file: %v", err)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	listSheets := f.GetSheetList()
	if len(listSheets) == 0 {
		s.logger.Error("No sheet found")
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoFileDoesNotExist0209")
	}

	rows, err := f.GetRows(listSheets[0])
	if err != nil {
		s.logger.Errorf("Failed to get sheet rows: %v", err)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "DashIncorrectFileFormat0513")
	}

	if len(rows) == 0 {
		s.logger.Error("No rows found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	_, err = s.validateHeaders(childCtx, rows[0])
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) parseBulkFileXLS(childCtx context.Context, data []byte) ([]request.BulkData, error) {
	xl, err := xls.OpenReader(bytes.NewReader(data), "utf-8")
	if err != nil {
		s.logger.Errorf("failed to open excel file: %v", err)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	if xl.NumSheets() == 0 {
		s.logger.Error("No sheet found")
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoFileDoesNotExist0209")
	}

	sheet := xl.GetSheet(0)

	if sheet.MaxRow < 2 {
		s.logger.Error("parseBulkFileXLS: MaxRow is less than 2")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	row0 := sheet.Row(0)
	if row0 == nil {
		s.logger.Error("parseBulkFileXLS: no rows found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	fields := make([]string, 0, 20)

	if row0.LastCol()-row0.FirstCol() < 4 {
		s.logger.Error("parseBulkFileXLS: Wrong file format")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	for i := row0.FirstCol(); i <= row0.LastCol(); i++ {
		fields = append(fields, row0.Col(i))
	}

	_, err = s.validateHeaders(childCtx, fields)
	if err != nil {
		return nil, err
	}

	return make([]request.BulkData, 0), nil
}

type validateHeaderResp struct {
	BenficiaryAccNoIndex     int
	BeneficiaryBankCodeIndex int
	MobileNumberIndex        int
	AmountIndex              int
	TransferNotesIndex       int
}

func (s *Service) validateHeaders(childCtx context.Context, columns []string) (*validateHeaderResp, error) {

	s.logger.Debugf("validateHeaders | columns: %v", columns)
	benficiaryAccNoIndex := slices.Index(columns, BenficiaryAccNo)
	if benficiaryAccNoIndex == -1 {
		s.logger.Error("benficiaryAccNoIndex is not found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}
	beneficiaryBankCodeIndex := slices.Index(columns, BeneficiaryBankCode)
	if beneficiaryBankCodeIndex == -1 {
		s.logger.Error("beneficiaryBankCodeIndex is not found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")

	}

	mobileNumberIndex := slices.Index(columns, MobileNumber)
	if mobileNumberIndex == -1 {
		s.logger.Error("mobileNumberIndex is not found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")

	}
	amountIndex := slices.Index(columns, Amount)
	if amountIndex == -1 {
		s.logger.Error("amountIndex is not found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	transferNotesIndex := slices.Index(columns, TransferNotes)
	if amountIndex == -1 {
		s.logger.Error("transferNotesIndex is not found")
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "DashIncorrectFileFormat0513")
	}

	return &validateHeaderResp{
		BenficiaryAccNoIndex:     benficiaryAccNoIndex,
		BeneficiaryBankCodeIndex: beneficiaryBankCodeIndex,
		MobileNumberIndex:        mobileNumberIndex,
		AmountIndex:              amountIndex,
		TransferNotesIndex:       transferNotesIndex,
	}, nil
}

func (s *Service) GetSignedURLFromBulkFile(
	ctx context.Context,
	req *GetSignedURLFromBulkFileRequest,
) (*GetSignedURLFromBulkFileResponse, error) {
	ctx, span := tracer.StartSpan(ctx, "GetSignedURLFromBulkFile")
	defer span.End()

	model, err := s.bulkFileRepo.GetBulkFile(ctx, &repositories.GetBulkFileRequest{
		Id: req.Id,
	})
	if err != nil {
		s.logger.Errorf("failed to get bulk file: %v", err)
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	resp, err := s.bulkFileUploadRepo.GetSignedURL(ctx, &repositories.GetSignedURLRequest{
		Path: model.FilePath,
	})
	if err != nil {
		s.logger.Errorf("failed to get signed url: %v", err)
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	return &GetSignedURLFromBulkFileResponse{
		URL: resp.URL,
	}, nil
}
