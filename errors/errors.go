package errors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bulk/constant"
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

type CustomError struct {
	ErrorVal
	ResponseTime    string      `json:"responseTime"`
	TransactionId   string      `json:"transactionId"`
	ReferenceNumber string      `json:"referenceNumber"`
	Errors          []ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ErrorVal struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	HttpCode int    `json:"-"`
}

func GetErrorItems(msgCode ...string) []ErrorItem {
	var listErrItem []ErrorItem
	for _, code := range msgCode {
		listErrItem = append(listErrItem, ErrorItem{
			Code:    code,
			Message: GetMessageByCode(code),
			Details: GetDescriptionByCode(code),
		})
	}

	return listErrItem
}

func GetMessageByCode(code string) string {
	switch code {
	case PreconditionErrMsg:
		return "Precondition failed"
	case BadRequestErrMsg:
		return "Bad request"
	default:
		return "Service unavailable"
	}
}

func GetDescriptionByCode(code string) string {
	// Implement logic to get description based on error code
	return "Description placeholder"
}

func GetHttpCode(errContext string) ErrorVal {
	switch errContext {
	case PreconditionErrMsg:
		return ErrorVal{
			Code:     ProductCodePreconditionFailed,
			Message:  PreconditionErrMsg,
			HttpCode: http.StatusPreconditionFailed,
		}

	case BadRequestErrMsg:
		return ErrorVal{
			Code:     ProductCodeBadRequest,
			Message:  BadRequestErrMsg,
			HttpCode: http.StatusBadRequest,
		}
	default:
		return ErrorVal{
			Code:     ProductCodeInternal,
			Message:  ServiceUnavailableErrMsg,
			HttpCode: http.StatusServiceUnavailable,
		}
	}
}

func GetError(ctx context.Context, errMsg string, messageCode ...string) *CustomError {
	txID, _ := ctx.Value(constant.CtxTransactionId).(string)
	refNum, _ := ctx.Value(constant.CtxReferenceNumber).(string)
	return &CustomError{
		ErrorVal:        GetHttpCode(errMsg),
		ResponseTime:    time.Now().Format(ResponseTimeLayout),
		TransactionId:   txID,
		ReferenceNumber: refNum,
		Errors:          GetErrorItems(messageCode...),
	}
}

func (o CustomError) Error() string {
	return fmt.Sprintf("CustomError code = %v desc - %v errors = %v", o.Code, o.Message, o.Errors)
}

func (o CustomError) GetHTTPCode() int {
	return o.HttpCode
}

func (o CustomError) GetErrors() []ErrorItem {
	return o.Errors
}

func (o CustomError) GetErrorCode() string {
	if len(o.Errors) > 0 {
		return o.Errors[0].Code
	}
	return ""
}

// GetCustomError creates a custom error according to args provided
func GetCustomError(ctx context.Context, errMsg string) *CustomError {
	txID, _ := ctx.Value(constant.CtxTransactionId).(string)
	refNum, _ := ctx.Value(constant.CtxReferenceNumber).(string)
	val := GetHttpCode(errMsg)
	return &CustomError{
		ErrorVal:        val,
		ResponseTime:    time.Now().Format(ResponseTimeLayout),
		TransactionId:   txID,
		ReferenceNumber: refNum,
		Errors:          []ErrorItem{},
	}
}
