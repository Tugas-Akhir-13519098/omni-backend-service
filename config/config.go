package config

import (
	"encoding/base64"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RESTHost string `envconfig:"REST_HOST" default:"localhost"`
	RESTPort int    `envconfig:"REST_PORT" default:"8080"`

	DBHost     string `envconfig:"DB_HOST" default:""`
	DBPort     int    `envconfig:"DB_PORT" default:""`
	DBName     string `envconfig:"DB_NAME" default:""`
	DBUser     string `envconfig:"DB_USER" default:""`
	DBPassword string `envconfig:"DB_PASSWORD" default:""`

	RedisHost     string `envconfig:"REDIS_HOST" default:""`
	RedisPort     int    `envconfig:"REDIS_PORT" default:""`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`

	JWTPrivateKey string `envconfig:"JWT_PRIVATE_KEY" default:""` // base64 format
	JWTPublicKey  string `envconfig:"JWT_PUBLIC_KEY" default:""`  // base64 format
	JWTDuration   int    `envconfig:"JWT_DURATION" default:""`

	SessionKey string `envconfig:"SESSION_KEY" default:""`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	// convert private key back to string
	decodedJWTPrivateKey, err := base64.StdEncoding.DecodeString(cfg.JWTPrivateKey)
	if err != nil {
		panic(err)
	}
	cfg.JWTPrivateKey = string(decodedJWTPrivateKey)

	// convert private key back to string
	decodedJWTPublicKey, err := base64.StdEncoding.DecodeString(cfg.JWTPublicKey)
	if err != nil {
		panic(err)
	}
	cfg.JWTPublicKey = string(decodedJWTPublicKey)

	// convert private key back to string
	decodedSessionKey, err := base64.StdEncoding.DecodeString(cfg.SessionKey)
	if err != nil {
		panic(err)
	}
	cfg.SessionKey = string(decodedSessionKey)

	return cfg
}
