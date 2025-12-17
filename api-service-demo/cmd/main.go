package main

import (
	"log/slog"

	"github.com/braunkc/todo-app/api-service-demo/config"
	client "github.com/braunkc/todo-app/api-service-demo/internal/grpc"
	server "github.com/braunkc/todo-app/api-service-demo/internal/http"
	"github.com/braunkc/todo-app/api-service-demo/internal/token"
	"github.com/braunkc/todo-app/api-service-demo/pkg/log"
	pb "github.com/braunkc/todo-app/api-service-demo/proto/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logCfg := log.Config{
		Service:    "api-demo",
		OutputType: log.Console,
		Level:      slog.LevelDebug,
	}

	loggerHandler, err := log.NewHandler(&logCfg)
	if err != nil {
	}
	l := slog.New(loggerHandler)

	cfg, err := config.New()
	if err != nil {
	}
	l.Debug("config inited", slog.Any("cfg", cfg))

	conn, err := grpc.NewClient(cfg.DatabaseService.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
	}
	defer conn.Close()
	c := pb.NewDataBaseServiceClient(conn)

	jwtService := token.NewJWTService([]byte(cfg.SecretKey))
	dbService := client.New(c)
	r := server.New(jwtService, dbService)

	r.Run(cfg.HTTPServer.Port)
}
