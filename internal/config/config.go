package config

import (
	"time"

	"github.com/kitti12911/lib-util/logger"
)

type Config struct {
	ServiceName       string        `mapstructure:"service_name"       env:"SERVICE_NAME"       validate:"required"`
	Port              int           `mapstructure:"port"               env:"PORT"               validate:"required,gte=1,lte=65535"`
	CollectorEndpoint string        `mapstructure:"collector_endpoint" env:"COLLECTOR_ENDPOINT" validate:"required_with=CollectorPort,omitempty,hostname|ip"`
	CollectorPort     int           `mapstructure:"collector_port"     env:"COLLECTOR_PORT"     validate:"required_with=CollectorEndpoint,omitempty,gte=1,lte=65535"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout"   env:"SHUTDOWN_TIMEOUT"`
	Logging           Logging       `mapstructure:"logging"`
	ExampleService    ServiceAddr   `mapstructure:"example_service"    validate:"required"`
}

type Logging struct {
	Level       logger.Level `mapstructure:"level"        env:"LOG_LEVEL"        validate:"oneof=debug info warn error"`
	AddSource   bool         `mapstructure:"add_source"   env:"LOG_ADD_SOURCE"`
	ServiceName string       `mapstructure:"service_name" env:"LOG_SERVICE_NAME"`
	EnableTrace bool         `mapstructure:"enable_trace" env:"LOG_ENABLE_TRACE"`
}

type ServiceAddr struct {
	Host string `mapstructure:"host" validate:"required"`
	Port int    `mapstructure:"port" validate:"required,gte=1,lte=65535"`
}
