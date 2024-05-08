package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"bulk/api/v1/repositories"
	"bulk/constant"
	"bulk/models/request"
	"bulk/pb/common"
	disbursementpb "bulk/pb/disbursement"
	"bulk/utils"
	errors "bulk/grpc-errors"
	"bulk/tracer"
)

func (s *Service) Disburse(ctx context.Context, req *Disbursement) (*disbursementpb.DisbursementResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "Disburse")
	defer span.End()
	s.logger.Debugf("%v | BulkDisbursementId: %v", "Service.Disburse", req.BulkDisbursementId)
	bulkDataList, errGetBulkData := s.bulkFileRepo.GetBulkData(childCtx, &repositories.GetBulkFileRequest{
		Id:     0,
		Bulkid: req.BulkDisbursementId,
		Status: 2,
	})
	if errGetBulkData != nil {
		if errGetBulkData == sql.ErrNoRows {
			return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
		}
		return nil, errGetBulkData
	}
	s.logger.Debugf("BulkData List: %v", bulkDataList)
	disbursementDetails := request.Dibursement{
		TransactionId:   childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber: childCtx.Value(constant.CtxReferenceNumber).(string),
		CorrelationId:   childCtx.Value(constant.HeaderCorrelationID).(string),
		MerchantUserId:  req.MerchantUserId,
		MerchantCode:    req.MerchantCode,
		BulkId:          req.BulkDisbursementId,
		Data:            bulkDataList,
	}
	disbursementDetailsObj, ErrMarshalDisbursement := json.Marshal(disbursementDetails)
	if ErrMarshalDisbursement != nil {
		return nil, errors.GetError(ctx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}
	// // Create event to validate, insert data into db and do inquiries
	s.logger.Debugf("Bulk Disbursement Details %v", disbursementDetailsObj)
	s.pubsubClient.PublishData(ctx, s.cfg.Pubsub.TransferTopic, disbursementDetailsObj, nil)

	return &disbursementpb.DisbursementResponse{
		ResponseCode:    "200",
		ResponseMessage: "Success",
		ResponseTime:    utils.GetTime().String(),
		TransactionId:   childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber: childCtx.Value(constant.CtxReferenceNumber).(string),
		Status:          6,
	}, nil
}

func (s *Service) UpdateDisburse(ctx context.Context, req *request.UpdateDisbursement) (*disbursementpb.UpdateDisbursementResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "Disburse")
	defer span.End()

	_, errUpdateBulkData := s.bulkFileRepo.UpdateBulkDisbursementData(childCtx, req)
	if errUpdateBulkData != nil {
		if errUpdateBulkData == sql.ErrNoRows {
			return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
		}
	}

	return &disbursementpb.UpdateDisbursementResponse{
		ResponseCode:    "200",
		ResponseMessage: "Success",
		ResponseTime:    utils.GetTime().String(),
		TransactionId:   childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber: childCtx.Value(constant.CtxReferenceNumber).(string),
	}, nil
}

func (s *Service) Disbursements(ctx context.Context, req *disbursementpb.DisbursementsRequest) (*disbursementpb.DisbursementsResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "Disbursements")
	defer span.End()

	bulkFileDetials, err := s.bulkFileRepo.GetBulkFileByBulkID(childCtx, req.BulkDisbursmentId)
	if err != nil {
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	disbursements, totalCount, err := s.bulkFileRepo.GetDisbursements(childCtx, int(bulkFileDetials.Id), req)
	if err != nil {
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	bulkDisbursementStatusCount, err := s.bulkFileRepo.GetBulkDisbursementStatusCount(ctx, int(bulkFileDetials.Id))
	if err != nil {
		return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	return &disbursementpb.DisbursementsResponse{
		ResponseCode:                "200",
		ResponseMessage:             "Success",
		ResponseTime:                utils.GetTime().String(),
		TransactionId:               childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber:             childCtx.Value(constant.CtxReferenceNumber).(string),
		BulkDisbursmentId:           &req.BulkDisbursmentId,
		BulkDisbursmentName:         &bulkFileDetials.BulkName,
		BulkDisbursmentStatus:       &bulkFileDetials.Status,
		BulkDisbursementStatusCount: bulkDisbursementStatusCount,
		Pagination: &common.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalCount: int32(totalCount),
		},
		DisbursementList: disbursements,
	}, nil

}
func (s *Service) Download(ctx context.Context, req *disbursementpb.DownloadDisbursementDataRequest) (*disbursementpb.DownloadDisbursementDataResponse, error) {
	childCtx, span := tracer.StartSpan(ctx, "Downlaod")
	defer span.End()
	s.logger.Debugf("%v | BulkDisbursementId: %v", "Service.Downlaod", req.BulkDisbursementId)

	bulkDataList, errGetBulkData := s.bulkFileRepo.GetBulkData(childCtx, &repositories.GetBulkFileRequest{
		Id:     0,
		Bulkid: req.BulkDisbursementId,
		Status: req.Status,
	})
	if errGetBulkData != nil {
		if errGetBulkData == sql.ErrNoRows {
			return nil, errors.GetError(childCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
		}
		return nil, errGetBulkData
	}
	s.logger.Debugf("BulkData List: %v", bulkDataList)
	data, errCsv := s.createCSVFileFromBulkData(childCtx, bulkDataList)
	if errCsv != nil {
		s.logger.Errorf("%v | Error while creating csv file", "Downlaod", errCsv)
		return nil, errCsv
	}
	uploadFile, errUpload := s.bulkFileUploadRepo.UploadFile(childCtx, &repositories.UploadFileRequest{
		Data:         data,
		MimeType:     "text/csv",
		FileName:     fmt.Sprintf("%v.csv", req.BulkDisbursementId),
		ParentFolder: "Final",
	})
	if errUpload != nil {
		s.logger.Errorf("%v | Error while uploading csv file", "Downlaod", errUpload)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}
	signedResp, errGetSignedResp := s.bulkFileUploadRepo.GetSignedURL(childCtx, &repositories.GetSignedURLRequest{
		Path: uploadFile.Path,
	})
	if errGetSignedResp != nil {
		s.logger.Errorf("%v | Error while get signed url", "Downlaod", errGetSignedResp)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}

	return &disbursementpb.DownloadDisbursementDataResponse{
		ResponseCode:       "200",
		ResponseMessage:    "Success",
		ResponseTime:       utils.GetTime().String(),
		TransactionId:      childCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber:    childCtx.Value(constant.CtxReferenceNumber).(string),
		MerchantUserId:     req.MerchantUserId,
		MerchantCode:       req.MerchantCode,
		BulkDisbursementId: req.BulkDisbursementId,
		File: &disbursementpb.DownloadFile{
			Url: signedResp.URL,
		},
	}, nil
}
func (s *Service) createCSVFileFromBulkData(ctx context.Context, req []request.BulkDisbursementDetails) ([]byte, error) {
	childCtx, span := tracer.StartSpan(ctx, "createCSVFileFromBulkData")
	defer span.End()
	bf := new(bytes.Buffer)

	csvWriter := csv.NewWriter(bf)
	if err := csvWriter.Write([]string{"Beneficiary ID", "Bank Swift Code", "Beneficiary Bank Name", "Beneficiary Account Number", "Beneficiary Name", "Amount", "Transfer Notes", "Status", "Failure Reason", "Transfer Reference Number"}); err != nil {
		s.logger.Errorf("Failed to read csv headers %v", err)
		return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
	}
	for _, row := range req {
		status := "N/A"
		switch row.Status {
		case 1:
			status = "Verifying"
		case 2:
			status = "Verified"
		case 3:
			status = "Processing"
		case 4:
			status = "Success"
		case 5:
			status = "Failed"
		}

		if err := csvWriter.Write([]string{
			row.BenficiaryId,
			row.BankCode,
			row.BeneficiaryBankName,
			row.AccountNumber,
			row.BenficiaryName,
			fmt.Sprintf("%.2f", row.Amount),
			row.PaymentInfo,
			status,
			row.FailedReason,
			row.DisbursementReferenceNo,
		}); err != nil {
			s.logger.Errorf("Error While writing csv data %v")
			return nil, errors.GetError(childCtx, errors.ServiceUnavailableErrMsg, "AyoErrorInternalServerError0201")
		}
	}
	csvWriter.Flush()

	return bf.Bytes(), nil
}
