package main

import (
	"context"
	"flag"
	"go-redis/internal/config"
	"go-redis/internal/repository"
	"go-redis/internal/scheduler/expiryscheduler"
	"go-redis/internal/server"
	"go-redis/pkg/utils/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	configPath := flag.String("config", "./go-redis.conf", "Config file path for this node")
	flag.Parse()

	config.Init(*configPath)
	log.Init(config.GetString("log-dir"))
	repository.Init()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	go server.StartHttpServer(ctx, os.Args, &wg)
	if !config.GetBool("read-only") {
		go expiryscheduler.StartScheduler(ctx, &wg)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	cancel()

	wg.Wait()

}
