package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type config struct {
	PGConfig            pgConfig            `yaml:"pg"`
	PaymentGRPCConfig   paymentGRPCConfig   `yaml:"payment_grpc"`
	InventoryGRPCConfig inventoryGRPCConfig `yaml:"inventory_grpc"`
	OrderHTTPConfig     orderHTTPConfig     `yaml:"order_http"`
	LoggerConfig        loggerConfig        `yaml:"logger"`
}

var appConfig *config

func AppConfig() *config {
	return appConfig
}

func Load(path string) (*config, error) {
	var cfg config
	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, fmt.Errorf("ошибка чтения конфига из %q:%w", path, err)

		}
		appConfig = &cfg
		return appConfig, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфиг из env: %w", err)
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
