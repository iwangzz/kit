package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/consul/api"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//curl -X POST http://127.0.0.1:9878/kit/calc/add/10/20

func main() {
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.String("consul.port", "8500", "consul server port")
	)
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	consulConfig := api.DefaultConfig()
	consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	//创建反向代理
	proxy := NewReverseProxy(consulClient, logger)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", "9878")
		errc <- http.ListenAndServe(":9878", proxy)
	}()

	logger.Log("exit", <-errc)
}

// NewReverseProxy 创建反向代理处理方法
func NewReverseProxy(client *api.Client, logger log.Logger) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		reqPath := req.URL.Path
		if reqPath == "" {
			return
		}
		//按照分隔符'/'对路径进行分解，获取服务名称serviceName
		pathArray := strings.Split(reqPath, "/")
		serviceName := pathArray[1]

		//调用consul api查询serviceName的服务实例列表
		result, _, err := client.Catalog().Service(serviceName, "", nil)
		if err != nil {
			logger.Log("ReverseProxy failed", "query service instace error", err.Error())
			return
		}

		if len(result) == 0 {
			logger.Log("ReverseProxy failed", "no such service instance", serviceName)
			return
		}

		destPath := strings.Join(pathArray[2:], "/")
		tgt := result[rand.Int()%len(result)]
		logger.Log("service id", tgt.ServiceID)
		fmt.Println(tgt.ServiceID)

		//设置代理服务地址信息
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
		req.URL.Path = "/" + destPath
	}

	return &httputil.ReverseProxy{Director: director}
}