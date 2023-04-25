package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lastbyte32/gofemart/internal/api"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer ctxCancel()

	app, err := api.New()
	if err != nil {
		log.Fatalf("err on prepare api-server: %v", err.Error())
	}
	if err = app.Run(ctx); err != nil {
		log.Fatalf("err on execute api-server: %v", err.Error())
	}
	log.Println("app terminated")
}
