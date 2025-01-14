package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/axopadyani/billing-engine/internal/common/businesserror"
)

// toGrpcError converts a business error to a gRPC error.
//
// This function takes a generic error and attempts to map it to an appropriate
// gRPC error code. If the error is a BusinessError, it is mapped to a specific
// gRPC error code based on its Kind. Otherwise, it defaults to codes.Unknown.
//
// Parameters:
//   - err: The error to be converted. This can be any error type, but special
//     handling is applied if it's a BusinessError.
//
// Returns:
//
//	An appropriate gRPC status code. The error message is
//	preserved from the original error.
func toGrpcError(err error) error {
	code := codes.Unknown

	var businessErr *businesserror.BusinessError
	if errors.As(err, &businessErr) {
		switch businessErr.Kind() {
		case businesserror.KindInternal:
			code = codes.Internal
		case businesserror.KindBadRequest:
			code = codes.InvalidArgument
		case businesserror.KindUnprocessableEntity:
			code = codes.FailedPrecondition
		case businesserror.KindNotFound:
			code = codes.NotFound
		case businesserror.KindAlreadyExists:
			code = codes.AlreadyExists
		}
	}

	return status.Error(code, err.Error())
}
