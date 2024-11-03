package main

import (
	"context"
	"go-redis/internal/repository"
	"go-redis/internal/service/scheduler/expiryscheduler"
	"go-redis/internal/service/server"
	"go-redis/pkg/utils/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.InitLog("build/logs/server.log")
	repository.InitRepositories()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	go server.StartHttpServer(ctx, os.Args, &wg)
	go expiryscheduler.StartScheduler(ctx, &wg)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	cancel()

	wg.Wait()

}
