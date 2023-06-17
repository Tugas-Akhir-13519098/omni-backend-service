package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RESTHost string `envconfig:"REST_HOST" default:"localhost"`
	RESTPort int    `envconfig:"REST_PORT" default:"8080"`
	ClientHost string`envconfig:"CLIENT_HOST" default:"localhost"`
	ClientPort int    `envconfig:"REST_PORT" default:"3000"`

	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconfig:"DB_PORT" default:"5433"`
	DBName     string `envconfig:"DB_NAME" default:"omni"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"postgres"`

	KafkaHost         string `envconfig:"KAFKA_HOST" default:"localhost"`
	KafkaPort         string `envconfig:"KAFKA_PORT" default:"9092"`
	KafkaProductTopic string `envconfig:"KAFKA_PRODUCT_TOPIC" default:"product"`

	FirebaseKeyPath           string `envconfig:"FIREBASE_KEY_PATH" default:"config/firebase-key.json"`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	return cfg
}
