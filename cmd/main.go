package main

import (
	"context"
	"intmax2-withdrawal/cmd/migrator"
	"intmax2-withdrawal/cmd/server"
	"intmax2-withdrawal/cmd/withdrawal"
	"intmax2-withdrawal/configs"
	"intmax2-withdrawal/internal/blockchain"
	"intmax2-withdrawal/internal/cli"
	"intmax2-withdrawal/internal/open_telemetry"
	"intmax2-withdrawal/pkg/logger"
	"intmax2-withdrawal/pkg/sql_db/db_app"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dimiro1/health"
)

func main() {
	cfg := configs.New()
	log := logger.New(cfg.LOG.Level, cfg.LOG.TimeFormat, cfg.LOG.JSON, cfg.LOG.IsLogLine)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	const int1 = 1
	done := make(chan os.Signal, int1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	go func() {
		<-done
		const msg = "SIGTERM detected"
		log.Errorf(msg)
		if cancel != nil {
			cancel()
		}
	}()

	err := open_telemetry.Init(cfg.OpenTelemetry.Enable)
	if err != nil {
		const msg = "open_telemetry init: %v"
		log.Errorf(msg, err)
		return
	}

	var dbApp db_app.SQLDb
	dbApp, err = db_app.New(ctx, log, &cfg.SQLDb)
	if err != nil {
		const msg = "db application init: %v"
		log.Errorf(msg, err)
		return
	}

	hc := health.NewHandler()
	bc := blockchain.New(ctx, cfg)

	wg := sync.WaitGroup{}

	err = cli.Run(
		ctx,
		migrator.NewMigratorCmd(ctx, log, dbApp),
		server.NewServerCmd(&server.Server{
			Context: ctx,
			Cancel:  cancel,
			Config:  cfg,
			Log:     log,
			DbApp:   dbApp,
			WG:      &wg,
			HC:      &hc,
		}),
		withdrawal.NewWithdrawCmd(&withdrawal.Withdrawal{
			Context: ctx,
			Config:  cfg,
			Log:     log,
			DbApp:   dbApp,
			SB:      bc,
		}),
	)
	if err != nil {
		const msg = "cli: %v"
		log.Errorf(msg, err)
		return
	}

	wg.Wait()
}
