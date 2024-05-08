package errors

import (
	"bulk/constant"
	cusErr "bulk/errors"
	pb "bulk/pbs/error"
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	BadRequestErrMsg         = "bad.request"
	PreconditionErrMsg       = "precondition.failed"
	ServiceUnavailableErrMsg = "service.unavailable"
)

const (
	ProductCodeBadRequest         = 400
	ProductCodePreconditionFailed = 412
	ProductCodeInternal           = 503
)

const ResponseTimeLayout = "20060102150405"

func GetError(ctx context.Context, errMsg string, messageCode string) error {
	err := cusErr.GetHttpCode(errMsg)
	se := status.New(codes.Code(err.HttpCode), messageCode)
	ds, e := se.WithDetails(
		&pb.ResponseError{
			ResponseMessage: err.Message,
		},
	)
	if e != nil {
		log.Printf("[GetError] error while creating gRPC status: %v", e)
		return errors.New("error while creating gRPC status")
	}
	return ds.Err()
}

// GenerateCustomError generates a *errors.CustomError from err
// It is assumed that err passed is a grpc error of type *status.Status
func GenerateCustomError(ctx context.Context, err error) (*cusErr.CustomError, bool) {
	if err == nil {
		return nil, false
	}
	var ce *cusErr.CustomError
	s := status.Convert(err)
	if len(s.Details()) <= 0 {
		return nil, false
	}
	for _, d := range s.Details() {
		switch re := d.(type) {
		case *pb.ResponseError:
			ce = &cusErr.CustomError{
				ErrorVal:        cusErr.GetHttpCode(re.ResponseMessage),
				ResponseTime:    time.Now().Format(ResponseTimeLayout),
				TransactionId:   ctx.Value(constant.CtxTransactionId).(string),
				ReferenceNumber: ctx.Value(constant.CtxReferenceNumber).(string),
				Errors:          cusErr.GetErrorItems(s.Message()),
			}
		default:
			log.Printf("Unexpected type: %s", re)
			return nil, false
		}
	}
	return ce, true
}

func getHTTPCode(code codes.Code) string {
	return strconv.Itoa(int(code))
}
