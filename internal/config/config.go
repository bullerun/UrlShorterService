package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	ENV        string `yaml:"env" default:"local"`
	Postgres   `yaml:"postgres"`
	HTTPServer `yaml:"http_server"`
}
type Postgres struct {
	User     string `yaml:"user" default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `yaml:"host" default:"localhost"`
	Port     string `yaml:"port" default:"5432"`
	Database string `yaml:"database" default:"postgres"`
}
type HTTPServer struct {
	Addr        string        `yaml:"address" default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" default:"10s"`
}
type Path struct {
	configPath string
}

func New() *Config {
	configPath := getConfigPath()
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &config
}
func getConfigPath() string {
	return os.Getenv("CONFIG_PATH")
}
