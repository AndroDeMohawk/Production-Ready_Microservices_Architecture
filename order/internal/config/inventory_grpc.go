package config

import (
	"net"
)

type inventoryGRPCConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

func (c *inventoryGRPCConfig) Address() string {
	return net.JoinHostPort(c.Addr, c.Port)
}
