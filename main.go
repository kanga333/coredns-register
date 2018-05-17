package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/kanga333/coredns-register/register"
)

func main() {
	config := flag.String("config", "", "Path of config file.")
	flag.Parse()

	lg, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	c := &register.Config{}
	err = register.LoadFile(*config, c)
	if err != nil {
		lg.Error("failed to load file", zap.Error(err))
	}

	scheduler, err := c.CreateScheduler(lg)
	if err != nil {
		lg.Error("failed to create scheduler", zap.Error(err))
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	scheduler.Start(s)
}
