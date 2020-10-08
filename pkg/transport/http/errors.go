package http

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kit_log "github.com/go-kit/kit/log"
	kit_http "github.com/go-kit/kit/transport/http"
)

// ErrorResponse defines an error message returned via HTTP
type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    codes.Code `json:"code"`
	Message string     `json:"message"`
}

// StatusCode returns the HTTP status code to return
func (res ErrorResponse) StatusCode() int {
	switch res.Error.Code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// EncodeErrorResponse encodes an error response to JSON
func NewErrorEncoder(logger kit_log.Logger) kit_http.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		var e error
		_, ok := status.FromError(err)
		if ok {
			e = err
		} else {
			_ = logger.Log("err", err)
			e = status.Error(codes.Unknown, "Unknown error")
		}
		s, _ := status.FromError(e)
		_, err = fmt.Fprintf(w, "%s (code %d)", s.Message(), s.Code())
		if err != nil {
			_ = logger.Log("err", err)
		}
	}
}
