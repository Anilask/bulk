package handlers

import (
	"context"
	"fmt"
	"strings"
	"bulk/constant"
	filepb "bulk/pb/file"
	"bulk/utils"
	"bulk/grpc-errors"
	"bulk/tracer"
	"bulk/validator"
	cusErr "bulk/errors"
	"google.golang.org/grpc/metadata"
)

func (handler *BulkDisbursementHandler) BulkFiles(ctx context.Context, req *filepb.BulkFilesRequest) (*filepb.BulkFilesResponse, error) {
	bulkFilesSpanName := "BulkDisbursementHandler.BulkFiles"
	ctx, bulkFilesSpan := tracer.StartSpan(ctx, bulkFilesSpanName)

	defer utils.PanicRecover(handler.log)
	defer bulkFilesSpan.End()
	bulkFilesMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}
	correlationId := bulkFilesMD[strings.ToLower(constant.HeaderCorrelationID)]

	if err := validator.ValidateRequest(req); err != nil {
		reqValidatorErr := cusErr.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		bulkFilesSpan.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}

	bulkFilesReferenceNumber := bulkFilesMD[strings.ToLower(constant.CtxReferenceNumber)]
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, bulkFilesReferenceNumber[0])

	bulkFilesSpanAttributes := map[string]string{
		constant.CtxTransactionID:    req.TransactionId,
		constant.HeaderCorrelationID: correlationId[0],
		constant.MerchantCodeString:  req.MerchantCode,
		constant.BulkDisbursmentID:   *req.BulkDisbursementId,
		constant.BulkDisbursmentName: *req.BulkDisbursementName,
		constant.Uploader:            *req.Uploader,
		constant.Status:              fmt.Sprintf("%d", req.Status),
		constant.Page:                fmt.Sprintf("%d", req.Page),
		constant.LimitString:         fmt.Sprintf("%d", req.Limit),
	}

	tracer.SetSpanAttributes(bulkFilesSpan, bulkFilesSpanAttributes)
	resp, bulkFileserr := handler.service.GetBulkList(ctx, req)
	if bulkFileserr != nil {
		handler.log.Errorf("Failed to get the bulkfile list: %v", bulkFileserr)
		bulkFilesSpan.RecordError(bulkFileserr)
		return nil, bulkFileserr
	}

	return resp, nil
}
