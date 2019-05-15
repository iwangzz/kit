package main

import (
	"net/http"
	"context"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"os"
	"fmt"
	"os/signal"
	"syscall"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"flag"
)

//启动 prometheus
//prometheus --config.file="/wangzz/go/src/kit/config/prometheus.yml"
//启动grafana
///usr/local/opt/grafana/bin/grafana-server --config /usr/local/etc/grafana/grafana.ini --homepath /usr/local/opt/grafana/share/grafana cfg:default.paths.logs=/usr/local/var/log/grafana cfg:default.paths.data=/usr/local/var/lib/grafana cfg:default.paths.plugins=/usr/local/var/lib/grafana/plugins
//启动consul
//consul agent -dev
//启动服务
//./server -consul.host=127.0.0.1 -consul.port=8500 -service.host=127.0.0.1 -service.port=9876
//consul admin
//http://127.0.0.1:8500/ui/#/dc1/services
//测试
//curl -X POST http://127.0.0.1:9876/calc/add/2/3

func main() {
	//解析参数
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
		consulPort = flag.String("consul.port", "8500", "consul host")
		svcHost = flag.String("service.host", "127.0.0.1", "service host")
		svcPort = flag.String("service.port", "9876", "service host")
	)
	flag.Parse()

	ctx := context.Background()
	errChan :=  make(chan error)

	var svc Service
	svc = &CalcService{}
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestamp)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	filedKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "kit",
		Subsystem: "kit_service",
		Name: "request_count",
		Help: "qps",
	}, filedKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "kit",
		Subsystem: "kit_service",
		Name: "request_latency",
		Help: "request time",
	}, filedKeys)

	//log middleware
	svc = LoggingMiddleware(logger) (svc)
	//prometheus middleware
	svc = MetricMiddleware(requestCount, requestLatency) (svc)

	//ep :=  MakeCalcEndpoint(svc)
	eps := CalcEndpoints{
		CalcEndpoint:	MakeCalcEndpoint(svc),
		HealthEndpoint:	MakeHealthEndpoint(svc),
	}
	r := MakeCalcHandler(ctx, eps, logger)

	//consul register
	register :=  Register(*consulHost, *consulPort, *svcHost, *svcPort, logger)
	go func() {
		fmt.Println("server start at port:" + *svcPort)
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(*svcHost + ":" + *svcPort , handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1,syscall.SIGUSR2)
		errChan <- fmt.Errorf("sig %s", <-c)
	}()

	fmt.Println(<-errChan)
	//解绑服务
	register.Deregister()
}