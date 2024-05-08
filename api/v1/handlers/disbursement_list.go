package handlers

import (
	"context"
	"strings"

	"bulk/constant"
	disbursementpb "bulk/pb/disbursement"
	"bulk/utils"
	"bulk/errors"
	"bulk/tracer"
	"bulk/validator"
	"google.golang.org/grpc/metadata"
)

func (handler *BulkDisbursementHandler) Disbursements(ctx context.Context, req *disbursementpb.DisbursementsRequest) (*disbursementpb.DisbursementsResponse, error) {
	spanName := "BulkDisbursementHandler.Disbursements"
	ctx, span := tracer.StartSpan(ctx, spanName)

	defer utils.PanicRecover(handler.log)
	defer span.End()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.GetError(ctx, errors.BadRequestErrMsg, "AyoErrorBadRequest0204")
	}
	correlationId := md[strings.ToLower(constant.HeaderCorrelationID)]

	if err := validator.ValidateRequest(req); err != nil {
		reqValidatorErr := errors.GetCustomError(ctx, errors.PreconditionErrMsg)
		reqValidatorErr.Errors = append(reqValidatorErr.Errors, *err...)
		span.RecordError(reqValidatorErr)
		return nil, reqValidatorErr
	}

	referenceNumber := utils.GenerateNewUniqueNumber()
	ctx = context.WithValue(ctx, constant.HeaderCorrelationID, correlationId[0])
	ctx = context.WithValue(ctx, constant.CtxTransactionID, req.TransactionId)
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, referenceNumber)
	ctx = context.WithValue(ctx, constant.CtxMerchantCode, req.MerchantCode)

	spanAttributes := map[string]string{
		constant.CtxTransactionID:     req.TransactionId,
		constant.CtxBulkDisbursmentID: req.BulkDisbursmentId,
		constant.CtxMerchantCode:      req.MerchantCode,
		constant.HeaderCorrelationID:  correlationId[0],
	}

	tracer.SetSpanAttributes(span, spanAttributes)
	resp, errGetDisbursements := handler.service.Disbursements(ctx, req)
	if errGetDisbursements != nil {
		handler.log.Errorf("Failed to retrive Disbursements: %v", errGetDisbursements)
		span.RecordError(errGetDisbursements)
		return nil, errGetDisbursements
	}

	return resp, nil
}
