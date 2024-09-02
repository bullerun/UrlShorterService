package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	ENV        string `migrate.yaml:"env" default:"local"`
	Postgres   `migrate.yaml:"postgres"`
	HTTPServer `migrate.yaml:"http_server"`
}
type Postgres struct {
	User     string `migrate.yaml:"user" default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `migrate.yaml:"host" default:"localhost"`
	Port     string `migrate.yaml:"port" default:"5432"`
	Database string `migrate.yaml:"database" default:"postgres"`
}
type HTTPServer struct {
	Addr        string        `migrate.yaml:"address" default:"localhost:8080"`
	TimeOut     time.Duration `migrate.yaml:"timeout" default:"10s"`
	IdleTimeout time.Duration `migrate.yaml:"idle_timeout" default:"10s"`
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
	if config.Postgres.Password == "" {
		config.Postgres.Password = getPGPassword()
	}
	return &config
}
func getConfigPath() string {
	return os.Getenv("CONFIG_PATH")
}
func getPGPassword() string {
	return os.Getenv("POSTGRES_PASSWORD")
}
