package app

import (
	"L0_EVRONE/config"
	"L0_EVRONE/internal/controller/http"
	"L0_EVRONE/internal/controller/natsstreaming"
	"L0_EVRONE/internal/usecase"
	"L0_EVRONE/internal/usecase/memory"
	"L0_EVRONE/pkg/httpserver"
	"L0_EVRONE/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	var err error
	
	l := logger.New(cfg.Log.Level)


	mem := memory.New()

	// Use case
	orderUseCase := usecase.New(mem)


	handler := gin.New()
	v1.NewRouter(handler, l, *orderUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	natsstreaming.New(l,*orderUseCase)

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
}
