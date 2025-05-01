package main

import (
	"context"
	"flag"
	"github.com/CT1403-2/Code-Judgement/judge/internal/runner"
	"github.com/CT1403-2/Code-Judgement/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CT1403-2/Code-Judgement/judge/config"
	"github.com/CT1403-2/Code-Judgement/judge/internal/controller"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	runnerInstance := runner.New(cfg)
	ctrl := controller.New(cfg, func(c *config.Config) (proto.ManagerClient, error) {
		conn, err := grpc.NewClient(c.Manager.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		return proto.NewManagerClient(conn), nil
	}, runnerInstance)
	ctrl.Run(ctx)
}
