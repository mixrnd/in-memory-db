package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"in-memory-db/internal"
	"in-memory-db/internal/compute"
	"in-memory-db/internal/config"
	intlogger "in-memory-db/internal/logger"
	"in-memory-db/internal/network"
	inmemory "in-memory-db/internal/storage/in-memory"
	"in-memory-db/internal/storage/wal"
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

	maxSegmentSize, err := cfg.Wal.MaxSegmentSizeToSizeInBytes()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	e := inmemory.NewEngine()
	p := compute.NewParser()
	segment := wal.NewSegment(maxSegmentSize, cfg.Wal.DataDirectory)
	walInst := wal.NewWal(ctx, cfg.Wal.FlushingBatchSize, cfg.Wal.FlushingBatchTimeout, segment, logger)
	db := internal.NewDatabase(e, p, logger, walInst)
	db.Init()

	server := network.NewServer(ctx, cfg.Network.Address, db, logger,
		network.WithServerIdleTimeout(cfg.Network.IdleTimeout),
		network.WithServerBufferSize(maxMessageSize),
		network.WithServerMaxConnectionsNumber(cfg.Network.MaxConnections),
	)

	doneSignal := make(chan os.Signal, 1)
	signal.Notify(doneSignal, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("server ready to run")
		if err := server.Run(); err != nil {
			logger.Error("server run error", zap.Error(err))
			return
		}

		cancel()
	}()

	<-doneSignal

	cancel()

	walInst.WaitWrite()
}
