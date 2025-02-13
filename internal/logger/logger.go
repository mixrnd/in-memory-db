package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"in-memory-db/internal/config"
)

func CreateLogger(cfg config.LogConfig) (*zap.Logger, error) {
	zapCfg := zap.NewProductionConfig()

	defaultLevel := zapcore.DebugLevel
	if cfg.Level != "" {
		if err := defaultLevel.Set(cfg.Level); err != nil {
			fmt.Printf("wrong log level %s", cfg.Level)
			return nil, err
		}
	}

	zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	zapCfg.OutputPaths = []string{cfg.Output}

	return zapCfg.Build()
}
