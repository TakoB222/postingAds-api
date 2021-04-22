package config

import (
	"github.com/TakoB222/postingAds-api/pkg"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	defaultHttpPort               = "8000"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1

	defaultConfigPath = "./config/"
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
		URI      string
		User     string `mapstructure:"user"`
		Name     string `mapstructure:"name"`
		Password string
	}
)

//Set up config
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
	setFromEnv(&cfg)

}

func populateDefaults() {
	viper.SetDefault("http.host", "localhost:"+defaultHttpPort)
	viper.SetDefault("http.readTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHttpRWTimeout)
	viper.SetDefault("http.maxHeaderBytes", defaultHttpMaxHeaderMegabytes)
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0]) // folder
	viper.SetConfigName(path[1]) // config file name

	return viper.ReadInConfig()
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = viper.GetString("password")
	cfg.Postgres.URI = viper.GetString("uri")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", cfg.Http); err != nil {
		return err
	}
	return viper.UnmarshalKey("db.postgres", cfg.Postgres)
}

func parseEnv() error {
	if err := godotenv.Load(); err != nil {
		pkg.Error("error occurred with environment load")
	}

	return parsePostgresEnv()
}

func parsePostgresEnv() error {
	viper.SetEnvPrefix("postgres")
	if err := viper.BindEnv("password"); err != nil {
		return err
	}

	return viper.BindEnv("uri")
}
