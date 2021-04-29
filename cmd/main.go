package main

import (
	"context"
	"github.com/TakoB222/postingAds-api/internal/config"
	"github.com/TakoB222/postingAds-api/internal/delivery/http"
	"github.com/TakoB222/postingAds-api/internal/repository"
	"github.com/TakoB222/postingAds-api/internal/server"
	"github.com/TakoB222/postingAds-api/internal/service"
	"github.com/TakoB222/postingAds-api/pkg/auth"
	"github.com/TakoB222/postingAds-api/pkg/database"
	"github.com/TakoB222/postingAds-api/pkg/hash"
	"github.com/TakoB222/postingAds-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
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

type dependecies struct {
	tokenManager *auth.Manager
	hasher       *hash.SHA1Hasher
}

func main() {
	cfg, err := initConfigManager()
	if err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	db, err := initPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("error initializing database: %s", err.Error())
	}

	dep := initDependencies(cfg)

	repos := repository.NewRepositories(db)
	service := service.NewServices(service.Dependencies{
		Repository:   repos,
		TokenManager: dep.tokenManager,
		Hasher:       dep.hasher,
		AccessTokenTTL: cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})
	handler := http.NewHandler(service, dep.tokenManager)

	server := server.NewServer(server.Config{Host: cfg.Http.Host, Port: cfg.Http.Port, MaxHeaderBytes: cfg.Http.MaxHeaderMegabytes,
		ReadTimeout: cfg.Http.ReadTimeout, WriteTimeout: cfg.Http.WriteTimeout}, handler.Init())
	go func() {
		if err := server.Run(); err != nil {
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

	if err := db.Close(); err != nil {
		logger.Errorf("failed to stop db: %v", err)
	}
}

func initConfigManager() (*config.Config, error) {
	return config.Init("../configs/config")
}

func initPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	return database.NewPostgresDB(database.DBConfig{Host: cfg.Postgres.Host, Port: cfg.Postgres.Port, Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password, DBName: cfg.Postgres.DBName, SSLMode: "disable"})
}

func initDependencies(cfg *config.Config) *dependecies {
	tokenManager, err := auth.NewManager(cfg.Auth.TokenSigningKey)
	if err != nil {
		logrus.Fatalf("error with initializing token manager: %s", err.Error())
	}
	hasher, err := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
	if err != nil {
		logrus.Fatalf("error with initializing password hasher: %s", err.Error())
	}

	return &dependecies{tokenManager: tokenManager, hasher: hasher}
}
