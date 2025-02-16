package config

import (
	"time"

	"github.com/inhies/go-bytesize"
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

func (nc NetworkConfig) MessageSizeToSizeInBytes() (int, error) {
	b, err := bytesize.Parse(nc.MaxMessageSize)
	if err != nil {
		return 0, err
	}

	return int(b), nil
}
