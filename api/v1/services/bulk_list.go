package services

import (
	"context"

	"bulk/constant"
	"bulk/pb/common"
	filepb "bulk/pb/file"
	"bulk/utils"
	"bulk/grpc-errors"
	"bulk/tracer"
)

func (s *Service) GetBulkList(ctx context.Context, req *filepb.BulkFilesRequest) (*filepb.BulkFilesResponse, error) {
	getBulkListChildCtx, getBulkListSpan := tracer.StartSpan(ctx, "GetBulkList")
	defer getBulkListSpan.End()

	bulkList, totalCount, err := s.bulkFileRepo.GetBulkList(getBulkListChildCtx, req)
	if err != nil {
		return nil, errors.GetError(getBulkListChildCtx, errors.BadRequestErrMsg, "AyoRecordNotFound0206")
	}

	return &filepb.BulkFilesResponse{
		ResponseCode:    "200",
		ResponseMessage: "Success",
		ResponseTime:    utils.GetTime().String(),
		TransactionId:   getBulkListChildCtx.Value(constant.CtxTransactionID).(string),
		ReferenceNumber: getBulkListChildCtx.Value(constant.CtxReferenceNumber).(string),
		BulkList:        bulkList,
		Pagination: &common.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalCount: int32(totalCount),
		},
	}, nil
}
