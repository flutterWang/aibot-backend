package main

import (
	"aibot-backend/api"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("config", "configs/config_local.toml", "config file path")

func main() {
	flag.Parse()

	api.Start(*configPath)
	defer api.Stop()

	WaitExitSignal()
}

func WaitExitSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
