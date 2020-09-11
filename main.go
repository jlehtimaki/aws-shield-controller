package main

import (
	"context"
	"github.com/jlehtimaki/aws-shield-controller/controller"
	"github.com/jlehtimaki/aws-shield-controller/pkg/apis/awsshield"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := awsshield.NewConfig()
	if err := cfg.ParseFlags(os.Args[1:]); err != nil {
		log.Fatalf("flag parsing error: %v", err)
	}
	log.Infof("Check interval: %s", cfg.Interval)

	ctrl := controller.Controller{
		Interval: cfg.Interval,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleSigterm(cancel)
	log.Info("Running")
	ctrl.Run(ctx)
}

func handleSigterm(cancel func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Info("Received SIGTERM. Terminating...")
	cancel()
}
