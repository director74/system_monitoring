package main

import (
	"flag"
	"fmt"

	"github.com/director74/system_monitoring/internal/cfg"
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
}
