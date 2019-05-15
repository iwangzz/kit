package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"os"
	"strconv"
	"github.com/pborman/uuid"
)

func Register (consulHost, consulPort, svcHost, svcPort string, logger log.Logger) (register sd.Registrar) {
	var client consul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = consulHost + ":" + consulPort
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}

		client = consul.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		HTTP: "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout: "1s",
		Notes: "consul check service health status",
	}

	port, _ := strconv.Atoi(svcPort)
	reg := api.AgentServiceRegistration{
		ID: "kit" + uuid.New(),
		Name: "kit",
		Address: svcHost,
		Port: port,
		Tags: []string{"kit", "kit_service"},
		Check: &check,
	}

	//注册
	register = consul.NewRegistrar(client, &reg, logger)

	return
}
