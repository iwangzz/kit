package main

import (
	"net/http"
	"context"
	"github.com/go-kit/kit/log"
	"os"
	"fmt"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"os/signal"
	"syscall"
	"flag"
)

//启动服务
//./client -consul.host=127.0.0.1 -consul.port=8500
//测试
//curl -X POST -d '{"type":"add","a":10,"b":20}' http://127.0.0.1:9877/calculate

func main() {
	//解析参数
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
		consulPort = flag.String("consul.port", "8500", "consul host")
	)
	flag.Parse()

	ctx := context.Background()
	errChan :=  make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestamp)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//创建consul客户端对象
	var client consul.Client
	{
		consulConfig := api.DefaultConfig()

		consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
		consulClient, err := api.NewClient(consulConfig)

		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consul.NewClient(consulClient)
	}

	//创建Endpoint
	discoverEndpoint := MakeDiscoverEndpoint(ctx, client, logger)

	//创建传输层
	r := MakeHttpHandler(discoverEndpoint)

	go func() {
		fmt.Println("server start at port:9877")
		errChan <- http.ListenAndServe(":9877" , r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1,syscall.SIGUSR2)
		errChan <- fmt.Errorf("sig %s", <-c)
	}()

	fmt.Println(<-errChan)
}