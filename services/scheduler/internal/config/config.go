package config

import "time"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type SchedulerConfig struct {
	Interval *time.Duration
}
