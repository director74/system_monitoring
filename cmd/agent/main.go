package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/director74/system_monitoring/internal/cfg"
	internalgrpc "github.com/director74/system_monitoring/internal/server/grpc"
)

var (
	port       string
	configFile string
)

func init() {
	flag.StringVar(&port, "port", "50051", "GRPC server port")
	flag.StringVar(&configFile, "config", "./measure.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := cfg.NewConfig()
	err := config.Parse(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	grpcServer := internalgrpc.NewServer(config)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
		log.Println("grpc server stopped")
	}()

	if err := grpcServer.Start(); err != nil {
		log.Println("failed to start grpc server: " + err.Error())
		cancel()
	}

	<-ctx.Done()
}
