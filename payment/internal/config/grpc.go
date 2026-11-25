package config

import "net"

type grpcConfig struct {
	Addr   string `yaml:"addr" env:"GRPC_ADDRESS" env-default:"localhost"`
	Port   string `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
	GwAddr string `yaml:"gw_addr" env:"GRPC_GATEWAY_ADDRESS" env-default:"40041"`
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.Addr, c.Port)
}

func (c *grpcConfig) GatewayAddress() string {
	return net.JoinHostPort("", c.GwAddr)
}
