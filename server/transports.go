package main

import (
	"context"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-kit/kit/log"
	"strconv"
	"encoding/json"
		kithttp "github.com/go-kit/kit/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func DecodeCalcRequest(_ context.Context, r *http.Request) (interface{}, error){
	vars := mux.Vars(r)
	calcType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	aStr, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	bStr, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	a, _ := strconv.Atoi(aStr)
	b, _ := strconv.Atoi(bStr)

	return CalcRequest{CalcType:calcType, A:a, B:b}, nil
}

func EncodeCalcResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error{
	w.Header().Set("content-type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func DecodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error){
	return HealthRequest{}, nil
}

//MakeCalcHandler
func MakeCalcHandler(ctx context.Context, eps CalcEndpoints, log log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(log),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}
	//calc
	r.Methods("POST").Path("/calc/{type}/{a}/{b}").Handler(kithttp.NewServer(
		eps.CalcEndpoint,
		DecodeCalcRequest,
		EncodeCalcResponse,
		options...
	))
	//health
	r.Methods("GET").Handler(kithttp.NewServer(
		eps.HealthEndpoint,
		DecodeHealthRequest,
		EncodeCalcResponse,
		options...
	))
	//prometheus
	r.Path("/metrics").Path("/health").Handler(promhttp.Handler())
	return r
}


