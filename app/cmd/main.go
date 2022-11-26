package main

import (
	"context"
	"github.com/ilkinabd/goods-manager/app/internal/app"
	"github.com/ilkinabd/goods-manager/app/internal/config"
	"github.com/ilkinabd/goods-manager/app/pkg/logging"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx.Done()

	logging.Info(ctx, "config initializing")
	cfg := config.GetConfig()

	ctx = logging.ContextWithLogger(ctx, logging.NewLogger())

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logging.Fatal(ctx, err)
	}

	logging.Info(ctx, "Running Application")
	err = a.Run(ctx)
	if err != nil {
		logging.Fatal(ctx, err)
	}
}
