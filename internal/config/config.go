package config

import (
	"github.com/TakoB222/postingAds-api/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"time"
)

const (
	defaultHttpPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1

	defaultConfigPath = "./configs/config.yml"
)

type (
	Config struct {
		Http HttpServer
		Postgres
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
		Port 	 string `mapstructure:"port"`
		Username     string `mapstructure:"username"`
		DBName     string `mapstructure:"dbName"`
		Password string
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

func parseConfigFile(filepath string) error {
	//path := strings.Split(filepath, "/")

	viper.AddConfigPath("configs") // folder
	viper.SetConfigName("config") // config file name

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = viper.GetString("password")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.Http); err != nil {
		return err
	}
	return viper.UnmarshalKey("db.postgres", &cfg.Postgres)
}

func parseEnv() error {
	if err := godotenv.Load(); err != nil {
		logger.Error("error occurred with environment load")
	}

	return parsePostgresEnv()
}

func parsePostgresEnv() error {
	viper.SetEnvPrefix("postgres")
	return viper.BindEnv("password")
}
