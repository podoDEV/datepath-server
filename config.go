package datepath_server

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/PODO/datepath-server/pkg/logrotater"
)

const (
	defaultServerPort              = 80
	defaultHttpsOn                 = false
	defaultSSLKeyPath              = "config/ssl/key.pem"
	defaultSSLCrtPath              = "config/ssl/cert.pem"
	defaultShutdownTimeoutInSecond = 60
)

type ServerConfig struct {
	Port                    int    `toml:"port"`
	HttpsOn                 bool   `toml:"https_on"`
	SSLKeyPath              string `toml:"ssl_key_path"`
	SSLCrtPath              string `toml:"ssl_crt_path"`
	ShutdownTimeoutInSecond int    `toml:"shutdown_timeout_in_second"`
}

type Config struct {
	Server     ServerConfig      `toml:"server"`
	LogRotater logrotater.Config `toml:"logrotater"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                    defaultServerPort,
			HttpsOn:                 defaultHttpsOn,
			SSLKeyPath:              defaultSSLKeyPath,
			SSLCrtPath:              defaultSSLCrtPath,
			ShutdownTimeoutInSecond: defaultShutdownTimeoutInSecond,
		},
		LogRotater: logrotater.NewConfig(),
	}
}

func LoadConfig(path string) (*Config, error) {
	conf := NewDefaultConfig()
	if path != "" {
		if _, err := toml.DecodeFile(path, &conf); err != nil {
			return nil, fmt.Errorf("toml decode file fail: %w", err)
		}
	}
	return conf, nil
}
