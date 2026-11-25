package config

type loggerConfig struct {
	LogLevel string `yaml:"level" env:"LOG_LEVEL"`
}

func (l *loggerConfig) Level() string {
	return l.LogLevel
}
