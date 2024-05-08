package errors_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"bulk/constant"
	"bulk/errors"
)

func TestGetError(t *testing.T) {
	// Create a context with transaction ID and reference number
	ctx := context.WithValue(context.Background(), constant.CtxTransactionId, "123456")
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, "789012")

	// Test case for PreconditionErrMsg
	err := errors.GetError(ctx, errors.PreconditionErrMsg)
	if err.GetErrorCode() != errors.PreconditionErrMsg {
		t.Errorf("Expected error code %s, got %s", errors.PreconditionErrMsg, err.GetErrorCode())
	}
	if err.GetHTTPCode() != http.StatusPreconditionFailed {
		t.Errorf("Expected HTTP code %d, got %d", http.StatusPreconditionFailed, err.GetHTTPCode())
	}
	if len(err.GetErrors()) != 1 {
		t.Errorf("Expected 1 error item, got %d", len(err.GetErrors()))
	}

	// Test case for BadRequestErrMsg with additional error message
	err = errors.GetError(ctx, errors.BadRequestErrMsg, "missing_parameter")
	if err.GetErrorCode() != errors.BadRequestErrMsg {
		t.Errorf("Expected error code %s, got %s", errors.BadRequestErrMsg, err.GetErrorCode())
	}
	if err.GetHTTPCode() != http.StatusBadRequest {
		t.Errorf("Expected HTTP code %d, got %d", http.StatusBadRequest, err.GetHTTPCode())
	}
	if len(err.GetErrors()) != 1 {
		t.Errorf("Expected 1 error item, got %d", len(err.GetErrors()))
	}
	if err.GetErrors()[0].Code != "missing_parameter" {
		t.Errorf("Expected error code %s, got %s", "missing_parameter", err.GetErrors()[0].Code)
	}

	// Add more test cases as needed
}

func TestGetCustomError(t *testing.T) {
	// Create a context with transaction ID and reference number
	ctx := context.WithValue(context.Background(), constant.CtxTransactionId, "123456")
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, "789012")

	// Test case for ServiceUnavailableErrMsg
	err := errors.GetCustomError(ctx, errors.ServiceUnavailableErrMsg)
	if err.GetErrorCode() != errors.ServiceUnavailableErrMsg {
		t.Errorf("Expected error code %s, got %s", errors.ServiceUnavailableErrMsg, err.GetErrorCode())
	}
	if err.GetHTTPCode() != http.StatusServiceUnavailable {
		t.Errorf("Expected HTTP code %d, got %d", http.StatusServiceUnavailable, err.GetHTTPCode())
	}
	if len(err.GetErrors()) != 0 {
		t.Errorf("Expected 0 error items, got %d", len(err.GetErrors()))
	}

	// Add more test cases as needed
}

func TestCustomError_Error(t *testing.T) {
	// Create a custom error instance for testing
	err := &errors.CustomError{
		ErrorVal:        errors.GetHttpCode(errors.BadRequestErrMsg),
		ResponseTime:    time.Now().Format(errors.ResponseTimeLayout),
		TransactionId:   "123456",
		ReferenceNumber: "789012",
		Errors: []errors.ErrorItem{
			{Code: "missing_parameter", Message: "Missing required parameter", Details: "Parameter 'id' is required."},
		},
	}

	expectedErrorStr := "CustomError code = 400 desc - bad.request errors = [{missing_parameter Missing required parameter Parameter 'id' is required.}]"
	if err.Error() != expectedErrorStr {
		t.Errorf("Expected error string '%s', got '%s'", expectedErrorStr, err.Error())
	}
}
