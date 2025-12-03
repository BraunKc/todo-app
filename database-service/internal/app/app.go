package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/braunkc/todo-app/database-service/config"
	"github.com/braunkc/todo-app/database-service/internal/application/usecases"
	database "github.com/braunkc/todo-app/database-service/internal/infra/database/postgres"
	grpcServer "github.com/braunkc/todo-app/database-service/internal/interfaces/grpc"
	"github.com/braunkc/todo-app/database-service/pkg/log"
)

func Run() error {
	logCfg := log.Config{
		Service:    "database",
		OutputType: log.Console,
		Level:      slog.LevelDebug,
	}

	loggerHandler, err := log.NewHandler(&logCfg)
	if err != nil {
		return fmt.Errorf("failed to create logger handler: %w", err)
	}
	l := slog.New(loggerHandler)

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}
	l.Debug("config inited", slog.Any("cfg", cfg))

	mapper := database.NewMapper()
	db, err := database.NewDatabaseService(cfg, mapper)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	l.Info("successful connected to DB")

	usecasesService := usecases.NewUsecasesService(db)

	server := grpcServer.New(usecasesService)

	listener, err := net.Listen("tcp", cfg.GRPCServer.Port)
	if err != nil {
		return fmt.Errorf("failed to create tcp listener: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		l.Info("server running")
		if err := server.Serve(listener); err != nil {
			l.Error("failed to serve server", slog.String("err", err.Error()))
			cancel()
		}
	}()

	<-ctx.Done()
	l.Info("graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	done := make(chan any)
	go func() {
		server.GracefulStop()
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		server.Stop()
		return errors.New("shutdown context is done, forced shutdown")
	case <-done:
		l.Info("service gracefully stopped")
		return nil
	}
}
