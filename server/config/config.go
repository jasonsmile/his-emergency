package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	current *Config
	mu      sync.RWMutex
)

// SetConfig 设置当前进程使用的配置。
func SetConfig(cfg *Config) { mu.Lock(); current = cfg; mu.Unlock() }

// GetConfig 获取当前进程配置。未初始化时返回空配置，便于调用方给出明确错误。
func GetConfig() *Config { mu.RLock(); defer mu.RUnlock(); return current }

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
	HIS      HISConfig      `yaml:"his"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	Username               string `yaml:"username"`
	Password               string `yaml:"password"`
	DBName                 string `yaml:"dbname"`
	MaxOpenConns           int    `yaml:"max_open_conns"`
	MaxIdleConns           int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_seconds"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`
}

type HISConfig struct {
	EmergencyMode bool   `yaml:"emergency_mode"`
	InvoicePrefix string `yaml:"invoice_prefix"`
	BaseURL       string `yaml:"base_url"`
	APIKey        string `yaml:"api_key"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "release"
	}
	return &cfg, nil
}
