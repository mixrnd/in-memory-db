package config

import (
	"time"

	"github.com/inhies/go-bytesize"
)

const (
	defaultEngineType = "in-memory"

	defaultNetworkAddress     = "127.0.0.1:8080"
	defaultMaxConnection      = 2
	defaultMaxMessageSize     = "4KB"
	defaultIdleTimeoutMinutes = 5

	defaultLogLevel  = "INFO"
	defaultLogOutput = "/tmp/in-mem-db.log"
)

type Config struct {
	Engine  EngineConfig  `yaml:"engine"`
	Network NetworkConfig `yaml:"network"`
	Log     LogConfig     `yaml:"logging"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

func (c *Config) SetDefaults() {
	c.Engine = EngineConfig{Type: defaultEngineType}
	c.Network = NetworkConfig{
		Address:        defaultNetworkAddress,
		MaxConnections: defaultMaxConnection,
		MaxMessageSize: defaultMaxMessageSize,
		IdleTimeout:    defaultIdleTimeoutMinutes * time.Minute,
	}
	c.Log = LogConfig{
		Level:  defaultLogLevel,
		Output: defaultLogOutput,
	}
}

func (nc NetworkConfig) MessageSizeToSizeInBytes() (int, error) {
	b, err := bytesize.Parse(nc.MaxMessageSize)
	if err != nil {
		return 0, err
	}

	return int(b), nil
}
