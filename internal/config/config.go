package config

type ServerConfig struct {
	ServiceName string `env:"SERVICE_NAME,required"`
}
