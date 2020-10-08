package pkg

import (
	"context"

	kit_endpoint "github.com/go-kit/kit/endpoint"
)

type GetObjectRequest struct {
	Host string
	Path string
}

type GetObjectResponse struct {
	Object *Object
}

type Endpoints struct {
	GetObject kit_endpoint.Endpoint
}

func NewGetObjectEndpoint(svc Service) kit_endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*GetObjectRequest)
		o := Object{}
		err := svc.GetObject(ctx, req.Host, req.Path, &o)
		if err != nil {
			return nil, err
		}
		return &GetObjectResponse{&o}, nil
	}
}
