package errors_test

import (
	"context"
	"testing"

	"bulk/constant"
	err "bulk/grpc-errors"
	pb "bulk/pbs/error"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetError(t *testing.T) {
	ctx := context.WithValue(context.Background(), constant.CtxTransactionId, "123456")
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, "789012")

	// Test case for BadRequestErrMsg
	errMsg := err.BadRequestErrMsg
	messageCode := "missing_parameter"
	expectedErrorCode := err.ProductCodeBadRequest
	expectedHTTPCode := codes.Code(expectedErrorCode)
	expectedMessage := "bad.request"

	err := err.GetError(ctx, errMsg, messageCode)
	s := status.Convert(err)
	if s.Code() != expectedHTTPCode {
		t.Errorf("Expected HTTP code %v, got %v", expectedHTTPCode, s.Code())
	}
	if s.Message() != messageCode {
		t.Errorf("Expected message code %s, got %s", messageCode, s.Message())
	}
	if len(s.Details()) != 1 {
		t.Errorf("Expected 1 detail, got %d", len(s.Details()))
	}
	detail := s.Details()[0]
	re, ok := detail.(*pb.ResponseError)
	if !ok {
		t.Errorf("Expected ResponseError detail type, got %T", detail)
	}
	if re.ResponseMessage != expectedMessage {
		t.Errorf("Expected response message %s, got %s", expectedMessage, re.ResponseMessage)
	}
}

func TestGenerateCustomError(t *testing.T) {
	// Simulate a gRPC error with details
	errMsg := "sample error message"
	s := status.New(codes.InvalidArgument, errMsg)
	se, _ := s.WithDetails(&pb.ResponseError{ResponseMessage: errMsg})

	// Create a context with transaction ID and reference number
	ctx := context.WithValue(context.Background(), constant.CtxTransactionId, "123456")
	ctx = context.WithValue(ctx, constant.CtxReferenceNumber, "789012")

	customErr, ok := err.GenerateCustomError(ctx, se.Err())
	if !ok {
		t.Errorf("Expected custom error, got nil")
	}
	if customErr.GetHTTPCode() != int(codes.InvalidArgument) {
		t.Errorf("Expected HTTP code %d, got %d", int(codes.InvalidArgument), customErr.GetHTTPCode())
	}
	if customErr.GetErrorCode() != errMsg {
		t.Errorf("Expected error code %s, got %s", errMsg, customErr.GetErrorCode())
	}
	if len(customErr.GetErrors()) != 1 {
		t.Errorf("Expected 1 error item, got %d", len(customErr.GetErrors()))
	}
	// Add more assertions as needed
}
