package main

import (
	"context"
	"github.com/TakoB222/postingAds-api/internal/config"
	"github.com/TakoB222/postingAds-api/internal/delivery/http"
	"github.com/TakoB222/postingAds-api/internal/server"
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/pkg/database"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/logger"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Init() {
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cfg, err := config.Init("../configs/config")
	if err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	db, err := database.NewPostgresDB(database.DBConfig{Host: cfg.Postgres.Host, Port: cfg.Postgres.Port, Username: cfg.Postgres.Username,
					Password: cfg.Postgres.Password, DBName: cfg.Postgres.DBName, SSLMode: "disable"})
	if err != nil {
		logrus.Fatalf("error initializing database: %s", err.Error())
	}

	tokenManager, err := auth.NewManager("dsfsdf")
	if err != nil {
		logrus.Fatalf("error with initializing token manager: %s", err.Error())
	}

	repos := repository.NewRepositories(db)
	service := service.NewServices(repos)
	handler := http.NewHandler(service, tokenManager)

	server := server.NewServer(server.Config{cfg.Http.Host, cfg.Http.Port, cfg.Http.MaxHeaderMegabytes,
					cfg.Http.ReadTimeout, cfg.Http.WriteTimeout}, handler.Init())
	go func() {
		if err := server.Run(); err != nil{
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logrus.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
