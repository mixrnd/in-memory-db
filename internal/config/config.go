package config

import (
	"time"

	"github.com/inhies/go-bytesize"
)

type Config struct {
	Engine  EngineConfig  `yaml:"engine"`
	Network NetworkConfig `yaml:"network"`
	Log     LogConfig     `yaml:"logging"`
	Wal     WalConfig     `yaml:"wal"`
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

type WalConfig struct {
	FlushingBatchSize    int           `yaml:"flushing_batch_size"`
	FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string        `yaml:"max_segment_size"`
	DataDirectory        string        `yaml:"data_directory"`
}

func (nc NetworkConfig) MessageSizeToSizeInBytes() (int, error) {
	return sizeInStringToBytes(nc.MaxMessageSize)
}

func (wc WalConfig) MaxSegmentSizeToSizeInBytes() (int, error) {
	return sizeInStringToBytes(wc.MaxSegmentSize)
}

func sizeInStringToBytes(st string) (int, error) {
	b, err := bytesize.Parse(st)
	if err != nil {
		return 0, err
	}

	return int(b), nil
}
