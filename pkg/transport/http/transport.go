package http

import (
	"bufio"
	"context"
	"net/http"

	kit_log "github.com/go-kit/kit/log"
	kit_http "github.com/go-kit/kit/transport/http"

	weasel "github.com/revas-hq/weasel/pkg"
)

func DecodeGetObjectRequest(_ context.Context, request *http.Request) (interface{}, error) {
	d := weasel.GetObjectRequest{
		Host: request.Host,
		Path: request.URL.Path,
	}
	return &d, nil
}

func EncodeGetObjectResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(*weasel.GetObjectResponse)
	h := w.Header()
	h.Set("Content-Type", res.Object.ContentType)
	h.Set("Cache-Control", res.Object.CacheControl)
	h.Set("Content-Disposition", res.Object.ContentDisposition)
	h.Set("Etag", res.Object.Etag)
	_, err := bufio.NewReader(res.Object.Body).WriteTo(w)
	return err
}

func NewHTTPHandler(logger kit_log.Logger, endpoints *weasel.Endpoints) http.Handler {
	options := []kit_http.ServerOption{
		kit_http.ServerErrorLogger(logger),
		kit_http.ServerErrorEncoder(NewErrorEncoder(logger)),
	}

	GetObjectHandler := kit_http.NewServer(
		endpoints.GetObject,
		DecodeGetObjectRequest,
		EncodeGetObjectResponse,
		options...,
	)

	return GetObjectHandler
}
