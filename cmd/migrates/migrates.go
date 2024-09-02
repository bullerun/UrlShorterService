package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type migrateConfig struct {
	SourceURL string `yaml:"sourceURL"`
	Storage   `yaml:"storage"`
}
type Storage struct {
	User     string `yaml:"user" default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `yaml:"host" default:"localhost"`
	Port     string `yaml:"port" default:"5432"`
	Database string `yaml:"database" default:"postgres"`
}

func main() {
	var migrateConfigPath string
	flag.StringVar(&migrateConfigPath, "migrate-config", "", "path to migrate config file")
	flag.Parse()
	if migrateConfigPath == "" {
		log.Fatal("Please specify migrate-config")
	}
	cfg := parseMigrate(migrateConfigPath)
	m, err := migrate.New(
		"file://"+cfg.SourceURL,
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			cfg.Storage.User,
			cfg.Storage.Password,
			cfg.Storage.Host,
			cfg.Storage.Port,
			cfg.Storage.Database))

	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
func parseMigrate(migratePath string) *migrateConfig {
	if migratePath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(migratePath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", migratePath)
	}

	var config migrateConfig
	if err := cleanenv.ReadConfig(migratePath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	if config.Password == "" {
		log.Println("password not set")
	}
	return &config
}
