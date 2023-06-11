package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RESTHost string `envconfig:"REST_HOST" default:"localhost"`
	RESTPort int    `envconfig:"REST_PORT" default:"8080"`

	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconfig:"DB_PORT" default:"5433"`
	DBName     string `envconfig:"DB_NAME" default:"omnichannel"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"postgres"`

	KafkaHost         string `envconfig:"KAFKA_HOST" default:"localhost"`
	KafkaPort         string `envconfig:"KAFKA_PORT" default:"9092"`
	KafkaProductTopic string `envconfig:"KAFKA_PRODUCT_TOPIC" default:"product"`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	return cfg
}
