package handlers

import (
	"context"
	"strings"
	"time"

	"bulk/api/v1/services"
	"bulk/constant"
	cusErr "bulk/errors"
	filepb "bulk/pb/file"
	"bulk/tracer"
	"bulk/utils"
	"bulk/validator"

	"bulk/errors"
	"github.com/gabriel-vasile/mimetype"
	"google.golang.org/grpc/metadata"
)

func (handler *BulkDisbursementHandler) UploadBulkFile(ctx context.Context, req *filepb.UploadBulkFileRequest) (*filepb.UploadBulkFileResponse, error) {
	ctx, uploadBulkFileSpan := tracer.StartSpan(ctx, "BulkDisbursementHandler.UploadBulkFile")

	defer utils.PanicRecover(handler.log)
	defer uploadBulkFileSpan.End()
	uploadBulkFileMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}

	uploadBulkFileReferenceNumber := uploadBulkFileMD[strings.ToLower(constant.CtxReferenceNumber)]
	correlationId := uploadBulkFileMD[strings.ToLower(constant.HeaderCorrelationID)]
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, uploadBulkFileReferenceNumber[0])

	if err := validator.ValidateRequest(req); err != nil {
		// Have to convert to gRPC status
		reqValidatorErr := cusErr.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		uploadBulkFileSpan.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}

	spanAttributes := map[string]string{
		constant.CtxTransactionID:    req.TransactionId,
		constant.CtxReferenceNumber:  uploadBulkFileReferenceNumber[0],
		constant.HeaderCorrelationID: correlationId[0],
	}

	tracer.SetSpanAttributes(uploadBulkFileSpan, spanAttributes)

	mtype := mimetype.Detect(req.FileContent)
	if !mtype.Is(utils.CSVMimeType) && !mtype.Is(utils.XLSMimeType) && !mtype.Is(utils.XLSXMimeType) {
		handler.log.Errorf("Incorrect mime type")
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "DashWrongMIMEType0512")
	}
	resp, errUploadFile := handler.service.UploadBulkFile(ctx, &services.UploadBulkFileRequest{
		MerchantCode: req.MerchantCode,
		Data:         req.FileContent,
		FileName:     req.FileName,
		MimeType:     mtype.String(),
		UserId:       req.MerchantUserId,
		FileSize:     req.FileSize,
		UploadAt:     time.Now(),
		BulkName:     req.BulkName,
	})
	if errUploadFile != nil {
		handler.log.Errorf("Failed to upload bulk file: %v", errUploadFile)
		uploadBulkFileSpan.RecordError(errUploadFile)
		return nil, errUploadFile
	}

	return resp, nil
}

func (handler *BulkDisbursementHandler) UpdateBulkFileStatus(ctx context.Context, req *filepb.UpdateBulkStatusRequest) (*filepb.UpdateBulkStatusResponse, error) {
	ctx, updateBulkFileStatusSpan := tracer.StartSpan(ctx, "BulkDisbursementHandler.UpdateBulkFileStatus")

	defer utils.PanicRecover(handler.log)
	defer updateBulkFileStatusSpan.End()
	updateBulkFileStatusMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}
	correlationId := updateBulkFileStatusMD[strings.ToLower(constant.HeaderCorrelationID)]

	if err := validator.ValidateRequest(req); err != nil {
		reqValidatorErr := cusErr.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		updateBulkFileStatusSpan.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}
	updateBulkFileStatusReferenceNumber := updateBulkFileStatusMD[strings.ToLower(constant.CtxReferenceNumber)]
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, updateBulkFileStatusReferenceNumber[0])

	spanAttributes := map[string]string{
		constant.CtxTransactionID:    req.TransactionId,
		constant.CtxReferenceNumber:  updateBulkFileStatusReferenceNumber[0],
		constant.HeaderCorrelationID: correlationId[0],
	}

	tracer.SetSpanAttributes(updateBulkFileStatusSpan, spanAttributes)

	resp, errUpdateStatus := handler.service.UpdateBulkStatus(ctx, req)
	if errUpdateStatus != nil {
		handler.log.Errorf("Failed to update bulk status: %v", errUpdateStatus)
		updateBulkFileStatusSpan.RecordError(errUpdateStatus)
		return nil, errUpdateStatus
	}

	return resp, nil
}
