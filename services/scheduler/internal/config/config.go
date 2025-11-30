package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB        DBConfig        `mapstructure:"db"`
	Scheduler SchedulerConfig `mapstructure:"scheduler"`
	Server    ServerConfig    `mapstructure:"server"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ServerConfig struct {
	Port                   int `mapstructure:"port"`
	ShutdownTimeoutSeconds int `mapstructure:"shutdown_timeout_seconds"`
}

type SchedulerConfig struct {
	// Use string in YAML, then parse to time.Duration automatically
	Interval  time.Duration `mapstructure:"interval"`
	BatchSize int           `mapstructure:batch_size`
}

// Load loads the configuration based on the environment
func Load(env string) (*Config, error) {
	v := viper.New()

	// Base config file
	v.SetConfigName("default") // config.yaml
	v.AddConfigPath("./services/scheduler/config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// Optional: environment-specific config (e.g., dev.yaml, prod.yaml)
	if env != "" {
		v.SetConfigName(env)
		_ = v.MergeInConfig() // ignore if env file does not exist
	}

	// Environment variables override everything
	v.SetEnvPrefix("APP_SCHEDULER")
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // FIX: cannot be nil

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
