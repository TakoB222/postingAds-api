package config

import (
	"errors"
	"fmt"
	"github.com/TakoB222/postingAds-api/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultHttpPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1

	defaultConfigPath = "../configs/config.yml"
	//envBase = "../"
)

type (
	Config struct {
		Http HttpServer
		Postgres
		Auth Auth
	}

	HttpServer struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes"`
	}

	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		DBName   string `mapstructure:"dbName"`
		Password string
	}

	Auth struct {
		PasswordSalt    string
		TokenSigningKey string
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	}
)

func Init(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigPath
	}
	populateDefaults()

	if err := parseEnv(); err != nil {
		return nil, err
	}

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func populateDefaults() {
	viper.SetDefault("http.host", "localhost:"+defaultHttpPort)
	viper.SetDefault("http.readTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.maxHeaderBytes", defaultHttpMaxHeaderMegabytes)
}

func parseConfigFile(filePath string) error {
	//rawPath := strings.Split(filePath, "/")

	viper.AddConfigPath("configs") // folder
	viper.SetConfigName("config")  // config file name

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = viper.GetString("password")
	cfg.Auth.PasswordSalt = viper.GetString("password_salt")
	cfg.Auth.TokenSigningKey = viper.GetString("signing_key")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.Http); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	return viper.UnmarshalKey("db.postgres", &cfg.Postgres)
}

func grabDirectory(dataPath string) ([]string, error) {
	//fmt.Printf("Scan from dir - %s\n", dataPath)

	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		fmt.Printf("error occurred with a ReadDir: %v", err.Error())
	}

	var filesArray []string
	for _, file := range files {
		filePath := filepath.Join(dataPath, file.Name())
		if file.IsDir() {
			files, err := grabDirectory(filePath)
			if err != nil {
				return nil, err
			}
			filesArray = append(filesArray, files...)
		}
		if filepath.Ext(strings.TrimSpace(filePath)) == ".env" {
			filesArray = append(filesArray, filePath)
		}
	}

	return filesArray, nil
}

func parseEnv() error {
	//files, err := grabDirectory(envBase)
	//if err != nil {
	//	logger.Error("error occurred with grabbing directory")
	//}
	if err := godotenv.Load(); err != nil {
		logger.Error("error occurred with environment load")
	}

	if err := parseAuthEnv(); err != nil {
		logger.Error(err.Error())
	}

	return parsePostgresEnv()
}

func parsePostgresEnv() error {
	viper.SetEnvPrefix("postgres")
	return viper.BindEnv("password")
}

func parseAuthEnv() error {
	viper.SetEnvPrefix("auth")
	if err := viper.BindEnv("password_salt"); err != nil {
		return errors.New("error occurred with password last environment binding")
	}
	return viper.BindEnv("signing_key")
}
