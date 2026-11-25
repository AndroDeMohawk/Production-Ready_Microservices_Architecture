package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PGConfig     pgConfig     `yaml:"pg"`
	GRPCConfig   grpcConfig   `yaml:"grpc"`
	LoggerConfig loggerConfig `yaml:"logger"`
}

var appConfig *Config

func AppConfig() *Config {
	return appConfig
}

func Load(path string) (*Config, error) {
	var cfg Config
	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("ошибка загрузки конфига из %s : %w", path, err)
		}
		appConfig = &cfg
		return appConfig, nil

	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфига из env файла %w", err)
	}
	appConfig = &cfg
	return appConfig, nil
}

const defCfgPath = "config/dev.config.yaml"

func FetchConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу (например, stage.config.yaml)")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}
	return defCfgPath

}
