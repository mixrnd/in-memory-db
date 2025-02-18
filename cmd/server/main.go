package main

import (
	"context"
	"flag"
	"fmt"

	"go.uber.org/zap"
	"in-memory-db/internal"
	"in-memory-db/internal/compute"
	"in-memory-db/internal/config"
	intlogger "in-memory-db/internal/logger"
	"in-memory-db/internal/network"
	inmemory "in-memory-db/internal/storage/in-memory"
)

var configPath = flag.String("config", "config.yml", "Path to config file")

func main() {
	flag.Parse()

	cfg, err := config.ParseConfig(*configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	maxMessageSize, err := cfg.Network.MessageSizeToSizeInBytes()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger, err := intlogger.CreateLogger(cfg.Log)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	e := inmemory.NewEngine()
	p := compute.NewParser()
	db := internal.NewDatabase(e, p, logger)

	server := network.NewServer(context.Background(), cfg.Network.Address, db, logger,
		network.WithServerIdleTimeout(cfg.Network.IdleTimeout),
		network.WithServerBufferSize(maxMessageSize),
		network.WithServerMaxConnectionsNumber(cfg.Network.MaxConnections),
	)

	logger.Info("server ready to run")
	if err := server.Run(); err != nil {
		logger.Error("server run error", zap.Error(err))
		return
	}
}
