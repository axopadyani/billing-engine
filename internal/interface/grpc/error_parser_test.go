package grpc

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/axopadyani/billing-engine/internal/common/businesserror"
)

func TestToGrpcError(t *testing.T) {
	testCases := []struct {
		name     string
		inputErr error
		wantErr  *status.Status
	}{
		{
			name:     "unknown error",
			inputErr: errors.New("unknown error"),
			wantErr:  status.New(codes.Unknown, "unknown error"),
		},
		{
			name:     "internal error",
			inputErr: businesserror.New("internal error", businesserror.KindInternal),
			wantErr:  status.New(codes.Internal, "internal error"),
		},
		{
			name:     "bad request error",
			inputErr: businesserror.New("bad request", businesserror.KindBadRequest),
			wantErr:  status.New(codes.InvalidArgument, "bad request"),
		},
		{
			name:     "unprocessable entity error",
			inputErr: businesserror.New("unprocessable entity", businesserror.KindUnprocessableEntity),
			wantErr:  status.New(codes.FailedPrecondition, "unprocessable entity"),
		},
		{
			name:     "not found error",
			inputErr: businesserror.New("not found", businesserror.KindNotFound),
			wantErr:  status.New(codes.NotFound, "not found"),
		},
		{
			name:     "already exists error",
			inputErr: businesserror.New("already exists", businesserror.KindAlreadyExists),
			wantErr:  status.New(codes.AlreadyExists, "already exists"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grpcErr := toGrpcError(tc.inputErr)
			st, ok := status.FromError(grpcErr)
			if !ok {
				t.Fatalf("expected gRPC status error, got %v", grpcErr)
			}

			if st.Code() != tc.wantErr.Code() {
				t.Fatalf("expected code %v, got %v", tc.wantErr.Code(), st.Code())
			}

			if st.Message() != tc.wantErr.Message() {
				t.Fatalf("expected message %s, got %s", tc.wantErr.Message(), st.Message())
			}
		})
	}
}
