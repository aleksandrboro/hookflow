package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string    `yaml:"log_level"`
	Server   ServerCfg `yaml:"server"`
	DB       DBCfg     `yaml:"db"`
}

type ServerCfg struct {
	Addr         string        `yaml:"addr"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}

type DBCfg struct {
	User string `yaml:"user"`
}

func MustLoad() *Config {
	cfgPath := fetchCfgPath()

	if cfgPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("invalid config path: " + cfgPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(cfgPath, &config); err != nil {
		panic(err)
	}

	return &config
}

func fetchCfgPath() string {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "", "path to config file")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = os.Getenv("CONFIG_PATH")
	}

	return cfgPath
}
