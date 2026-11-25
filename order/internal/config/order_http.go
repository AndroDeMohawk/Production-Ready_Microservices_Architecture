package config

import "net"

type orderHTTPConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

func (c *orderHTTPConfig) Address() string {
	return net.JoinHostPort(c.Addr, c.Port)
}
