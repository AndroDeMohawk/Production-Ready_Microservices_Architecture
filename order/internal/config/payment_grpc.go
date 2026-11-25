package config

import "net"

type paymentGRPCConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

func (c *paymentGRPCConfig) Address() string {
	return net.JoinHostPort(c.Addr, c.Port)
}
