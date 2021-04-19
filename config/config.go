package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel      string        `default:"info" split_words:"true"`
	ListenAddress string        `default:":8080" split_words:"true"`
	CacheSize     int64         `default:"5000" split_words:"true"`
	ClientTimeout time.Duration `default:"5s" split_words:"true"`
	ReadTimeout   time.Duration `default:"5s" split_words:"true"`
	WriteTimeout  time.Duration `default:"500s" split_words:"true"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("ethproxy", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
