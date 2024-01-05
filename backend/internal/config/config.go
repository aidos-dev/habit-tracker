package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	DB
}

type HTTPServer struct {
	Port        string        `yaml:"port" env-default:"8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func MustLoad() *Config {
	// Load environment variables from .env file
	err := godotenv.Load("build/.env")
	if err != nil {
		log.Fatalf("error loading env variables from .env: %s", err.Error())
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return &cfg
}
