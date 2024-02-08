package servicediscovery

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	port      = 8080
	serviceId = "user-service"
)

func RegisterService() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	addr := "localhost"
	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    "user-server",
		Port:    port,
		Address: addr,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", addr, port, serviceId),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	regiErr := consul.Agent().ServiceRegister(registration)
	if regiErr != nil {
		log.Fatalf("failed to register service")
	} else {
		log.Printf("successfully registered services %s:%v", addr, port)
	}
}
