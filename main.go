package main

import (
	. "github.com/Floor-Gang/subwatch/internal"
	"os"
	"os/signal"
	"syscall"
)

const (
	configLocation = "config.yml"
)

func main() {
	config := GetConfig(configLocation)
	Start(config, configLocation)

	keepAlive()
}

func keepAlive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
