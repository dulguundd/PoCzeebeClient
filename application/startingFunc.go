package application

import (
	"os"
)

type Environment struct {
	serviceConfig ServiceConfig
	zeebeConfig   ZeebeConfig
}

type ServiceConfig struct {
	address string
}

type ZeebeConfig struct {
	zeebeAddress string
}

func getEnvironment() *Environment {
	var environment Environment
	environment = Environment{
		serviceConfig: struct{ address string }{address: "0.0.0.0:8080"},
		zeebeConfig:   struct{ zeebeAddress string }{zeebeAddress: "10.12.7.29:26500"},
	}

	if os.Getenv("SERVER_ADDRESS") != "" {
		environment.serviceConfig.address = os.Getenv("SERVER_ADDRESS")
	}

	if os.Getenv("ZEEBE_ADDRESS") != "" {
		environment.zeebeConfig.zeebeAddress = os.Getenv("ZEEBE_ADDRESS")
	}

	return &environment
}
