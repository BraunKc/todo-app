package main

import (
	"log/slog"
	"os"

	"github.com/braunkc/todo-db/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app failed", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
