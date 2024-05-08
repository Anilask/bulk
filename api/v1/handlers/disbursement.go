package handlers

import (
	"context"
	"strings"

	"bulk/api/v1/services"
	"bulk/constant"
	"bulk/models/request"
	disbursementpb "bulk/pb/disbursement"
	"bulk/utils"
	"bulk/grpc-errors"
	cusErr "bulk/errors"
	"bulk/tracer"
	"bulk/validator"
	"google.golang.org/grpc/metadata"
)

func (handler *BulkDisbursementHandler) Disburse(ctx context.Context, req *disbursementpb.DisbursementRequest) (*disbursementpb.DisbursementResponse, error) {
	ctx, disburseSpan := tracer.StartSpan(ctx, "BulkDisbursementHandler.Disburse")

	defer utils.PanicRecover(handler.log)
	defer disburseSpan.End()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}
	correlationId := md[strings.ToLower(constant.HeaderCorrelationID)]

	if err := validator.ValidateRequest(req); err != nil {
		reqValidatorErr := cusErr.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		disburseSpan.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}
	disburseReferenceNumber := md[strings.ToLower(constant.CtxReferenceNumber)]
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, disburseReferenceNumber[0])

	disburseSpanAttributes := map[string]string{
		constant.CtxTransactionID:    req.TransactionId,
		constant.CtxReferenceNumber:  disburseReferenceNumber[0],
		constant.HeaderCorrelationID: correlationId[0],
	}

	tracer.SetSpanAttributes(disburseSpan, disburseSpanAttributes)
	resp, errDisburse := handler.service.Disburse(ctx, &services.Disbursement{
		MerchantUserId:     req.MerchantUserId,
		MerchantCode:       req.MerchantCode,
		BulkDisbursementId: req.BulkDisbursementId,
	})
	if errDisburse != nil {
		handler.log.Errorf("Failed to proceed disburse: %v", errDisburse)
		disburseSpan.RecordError(errDisburse)
		return nil, errDisburse
	}

	return resp, nil
}

func (handler *BulkDisbursementHandler) UpdateDisburseDetails(ctx context.Context, req *disbursementpb.UpdateDisbursementRequest) (*disbursementpb.UpdateDisbursementResponse, error) {
	ctx, updateDisburseDetailsSpan := tracer.StartSpan(ctx, "BulkDisbursementHandler.UpdateDisburseDetails")

	defer utils.PanicRecover(handler.log)
	defer updateDisburseDetailsSpan.End()
	udateDisburseDetailsMD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}
	correlationId := udateDisburseDetailsMD[strings.ToLower(constant.HeaderCorrelationID)]

	if err := validator.ValidateRequest(req); err != nil {
		reqValidatorErr := cusErr.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		updateDisburseDetailsSpan.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, req.ReferenceNumber)

	updateDisburseDetailsSpanAttributes := map[string]string{
		constant.CtxTransactionID:    req.TransactionId,
		constant.CtxReferenceNumber:  req.ReferenceNumber,
		constant.HeaderCorrelationID: correlationId[0],
	}

	tracer.SetSpanAttributes(updateDisburseDetailsSpan, updateDisburseDetailsSpanAttributes)
	resp, errUpdateDisburse := handler.service.UpdateDisburse(ctx, &request.UpdateDisbursement{
		TransactionId:               req.TransactionId,
		ReferenceNumber:             req.ReferenceNumber,
		CustomerId:                  req.CustomerId,
		BeneficiaryCorrelationId:    req.BeneficiaryCorrelationId,
		BeneficiaryId:               req.BeneficiaryId,
		BeneficiaryStatus:           req.BeneficiaryStatus,
		DisbursementReferenceNumber: req.DisbursementReferenceNumber,
		DisbursementStatus:          req.DisbursementStatus,
		Status:                      req.Status,
		FailedReason:                req.FailedReason,
		BulkId:                      req.BulkId,
		Id:                          req.Id,
		BenficiaryName:              req.BeneficiaryName,
		BeneficiaryBankName:         req.BeneficiaryBankName,
		Type:                        req.Type,
	})
	if errUpdateDisburse != nil {
		handler.log.Errorf("Failed to update disbursement details: %v", errUpdateDisburse)
		updateDisburseDetailsSpan.RecordError(errUpdateDisburse)
		return nil, errUpdateDisburse
	}

	return resp, nil
}

func (handler *BulkDisbursementHandler) DownloadDisburse(ctx context.Context, req *disbursementpb.DownloadDisbursementDataRequest) (*disbursementpb.DownloadDisbursementDataResponse, error) {
	ctx, downloadSpan := tracer.StartSpan(ctx, "BulkDisbursementHandler.DownloadDisburse")

	defer utils.PanicRecover(handler.log)
	defer downloadSpan.End()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}
	correlationId := md[strings.ToLower(constant.HeaderCorrelationID)]

	if err := validator.ValidateRequest(req); err != nil {
		reqValidatorErr := cusErr.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		downloadSpan.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}
	downloadReferenceNumber := md[strings.ToLower(constant.CtxReferenceNumber)]
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, downloadReferenceNumber[0])

	downloadSpanAttributes := map[string]string{
		constant.CtxTransactionID:    req.TransactionId,
		constant.CtxReferenceNumber:  downloadReferenceNumber[0],
		constant.HeaderCorrelationID: correlationId[0],
	}

	tracer.SetSpanAttributes(downloadSpan, downloadSpanAttributes)
	resp, errDownload := handler.service.Download(ctx, req)
	if errDownload != nil {
		handler.log.Errorf("Failed to download disbursement file: %v", errDownload)
		downloadSpan.RecordError(errDownload)
		return nil, errDownload
	}

	return resp, nil
}