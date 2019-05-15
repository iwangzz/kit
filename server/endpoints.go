package main

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type CalcEndpoints struct {
	CalcEndpoint   endpoint.Endpoint
	HealthEndpoint   endpoint.Endpoint
}

type CalcRequest struct {
	CalcType	string 		`json:"calc_type"`
	A 			int 		`json:"a"`
	B 			int 		`json:"b"`
}

type CalcResponse struct {
	Result  	int  		`json:"result"`
	Error       error 		`json:"error"`
}

type HealthRequest struct {
}

type HealthResponse struct {
	Status 		bool	`json:"result"`
}

func MakeCalcEndpoint (svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{})(response interface{}, err error) {
		req := request.(CalcRequest)
		var (
			result, a, b int
			calcError error
		)
		a = req.A
		b = req.B

		switch true {
			case req.CalcType == "add":
				result = svc.Add(a, b)
			case req.CalcType == "sub":
				result = svc.Subtract(a, b)
			case req.CalcType == "multi":
				result = svc.Multiply(a, b)
			case req.CalcType == "div":
				result, calcError = svc.Divide(a, b)
			default:
				return nil, ErrInvalidRequestType
		}

		return CalcResponse{Result:result, Error:calcError}, nil
	}
}

func MakeHealthEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{Status:status}, nil
	}
}