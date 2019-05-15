package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"time"
)

func MakeDiscoverEndpoint(ctx context.Context, client consul.Client, logger log.Logger) endpoint.Endpoint {
	serviceName := "kit"
	tags := []string{"kit", "kit_service"}
	duration := 500 * time.Millisecond

	instancer := consul.NewInstancer(client, logger, serviceName, tags, true)

	factory := arithmeticFactory(ctx, "POST", "calc")
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(1, duration, balancer)

	return retry
}