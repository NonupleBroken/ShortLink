package main

import (
	"ShortLink/config"
	"ShortLink/server"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func waitStopSignal(wg *sync.WaitGroup) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<- sig

	server.Stop()

	wg.Done()
}

func main() {
	configFileName := "config.toml"
	cfg, err := config.InitConfig(configFileName)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go waitStopSignal(&wg)

	server.Run(&cfg)

	server.WaitStop(&wg)
}