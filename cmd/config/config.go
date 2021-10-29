package config

import "time"

type Config struct {
	Logger  Logger
	Server  Server
	Storage Storage
}

type Server struct {
	HTTP HTTP
	GRPC Grpc
}

type HTTP struct {
	Host string `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

type Grpc struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"8081"`
}

type Storage struct {
	Type  string `yaml:"type" env:"STORAGE_TYPE"`
	Redis Redis
}

type Redis struct {
	Host    string        `yaml:"host" env:"REDIS_HOST" env-default:"0.0.0.0"`
	Port    string        `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
	DB      int           `yaml:"db" env:"REDIS_DB" env-default:"0"`
	Expires time.Duration `yaml:"expires" env:"REDIS_EXPIRES" env-default:"10"`
}

type Logger struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	File  string `yaml:"file" env:"LOG_FILE"`
}

func NewConfig() *Config {
	return &Config{}
}
