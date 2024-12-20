package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"go-redis/internal/config"
	"go-redis/internal/model"
	"go-redis/internal/repository"
	"go-redis/internal/scheduler"
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

	initApp(configPath)

	// Gracefully shutdown initialization
	ctx, cancel := context.WithCancel(context.Background())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go server.StartTcpServer(ctx, os.Args, &wg)
	if !config.GetConfigValueBool("read-only") {
		go scheduler.StartExpiryScheduler(ctx, &wg)
	}

	<-signalCh
	cancel()
	wg.Wait()
}

func initApp(configPath *string) {
	config.InitConfParser(*configPath)
	log.Init(config.GetConfigValueString("log-dir"))
	repository.InitMemoryRepository()
	initAppState()
}

func initAppState() {
	model.State = model.AppState{ReplicationOffset: 0, ReplicationId: uuid.Must(uuid.NewRandom()).String()}
}
