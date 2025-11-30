package main

import (
	"log/slog"
	"os"

	"gihtub.com/braunkc/todo-db/config"
	database "gihtub.com/braunkc/todo-db/internal/infra/database/postgres"
	"gihtub.com/braunkc/todo-db/pkg/log"
)

func main() {
	logCfg := log.Config{
		Service:    "database",
		OutputType: log.Console,
		Level:      slog.LevelDebug,
	}

	loggerHandler, err := log.NewHandler(&logCfg)
	if err != nil {
		slog.Error("failed to create logger handler", slog.String("err", err.Error()))
		os.Exit(1)
	}
	l := slog.New(loggerHandler)

	cfg, err := config.New()
	if err != nil {
		l.Error("failed to init config", slog.String("err", err.Error()))
		os.Exit(1)
	}
	l.Info("config inited", slog.Any("cfg", cfg))

	_, err = database.NewDatabaseService(cfg)
	if err != nil {
		l.Error("failed to connect to DB", slog.String("err", err.Error()))
		os.Exit(1)
	}
	l.Info("successful connected to DB")
}
