package app

import (
	"L0_EVRONE/config"
	"L0_EVRONE/internal/controller/http"
	"L0_EVRONE/internal/controller/natsstreaming"
	"L0_EVRONE/internal/usecase"
	"L0_EVRONE/internal/usecase/memory"
	"L0_EVRONE/internal/usecase/postgresdb"
	"L0_EVRONE/pkg/httpserver"
	"L0_EVRONE/pkg/logger"
	"L0_EVRONE/pkg/postgres"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	var err error

	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()
	pgdb := postgresdb.New(pg)

	mem := memory.New(pgdb)

	// Use case
	orderUseCase := usecase.New(mem)

	// HTTP
	handler := gin.New()
	v1.NewRouter(handler, l, *orderUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Nats-streaming
	nats := natsstreaming.New(l, *orderUseCase)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = nats.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - nats.Shutdown: %w", err))
	}
}
