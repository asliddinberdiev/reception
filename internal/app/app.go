package app

import (
	"context"
	externalLog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asliddinberdiev/reception/internal/config"
	"github.com/asliddinberdiev/reception/internal/server"
	"github.com/asliddinberdiev/reception/internal/service"
	"github.com/asliddinberdiev/reception/internal/storage"
	transport "github.com/asliddinberdiev/reception/internal/transport/http"
	"github.com/asliddinberdiev/reception/pkg/db"
	"github.com/asliddinberdiev/reception/pkg/logger"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		externalLog.Fatalf("failed to initialize configuration: %+v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	log := logger.NewLogger(cfg.App.LogLevel, cfg.App.ServiceName)

	defer logger.Cleanup(log)
	defer cancel()

	psqlConn, err := db.Connect(cfg, log, ctx)
	if err != nil {
		log.Fatal("failed to connect to postgres", logger.Error(err))
		return
	}
	defer psqlConn.Close()

	strgPg := storage.NewStoragePg(psqlConn, log)

	services := service.NewService(strgPg, cfg)

	handlers := transport.NewHandler(log, cfg, services)

	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); err != nil {
			log.Error("failed to run server", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	const timeout = 4 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Error("failed to stop server", logger.Error(err))
	}
}
