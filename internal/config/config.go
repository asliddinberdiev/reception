package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Auth     Auth     `mapstructure:"auth"`
	Postgres Postgres `mapstructure:"postgres"`
}

type App struct {
	Environment string `envconfig:"APP_ENVIRONMENT" default:"dev" required:"true" mapstructure:"environment"`
	ServiceName string `envconfig:"APP_NAME" default:"reception" required:"true" mapstructure:"app_name"`
	Host        string `envconfig:"APP_HOST" default:"localhost" required:"true" mapstructure:"host"`
	Port        int    `envconfig:"APP_PORT" default:"8000" required:"true" mapstructure:"port"`
	LogLevel    string `envconfig:"APP_LOG_LEVEL" default:"info" required:"true" mapstructure:"log_level"`
}

type Auth struct {
	AccessTTL  time.Duration `envconfig:"AUTH_ACCESS_TTL" default:"15m" required:"true" mapstructure:"access_ttl"`
	RefreshTTL time.Duration `envconfig:"AUTH_REFRESH_TTL" default:"1h" required:"true" mapstructure:"refresh_ttl"`
	CodeLength int           `envconfig:"AUTH_CODE_LENGTH" default:"8" required:"true" mapstructure:"code_length"`
	SecretKey  string        `envconfig:"AUTH_KEY" default:"auth_secret_key" required:"true" mapstructure:"secret_key"`
}

type Postgres struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"localhost" required:"true" mapstructure:"host"`
	Port     int    `envconfig:"POSTGRES_PORT" default:"5432" required:"true" mapstructure:"port"`
	User     string `envconfig:"POSTGRES_USER" default:"postgres" required:"true" mapstructure:"user"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"password" required:"true" mapstructure:"password"`
	Database string `envconfig:"POSTGRES_DATABASE" default:"reception" required:"true" mapstructure:"database"`
	SslMode  string `envconfig:"POSTGRES_SSLMODE" default:"disable" required:"true" mapstructure:"sslmode"`
}

func Init(confDir string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if confDir != "" {
		v.AddConfigPath(confDir)
	}
	v.AddConfigPath("./conf")
	v.AddConfigPath(".")

	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "dev"
	}

	v.SetConfigName(env)

	envFile := fmt.Sprintf("%s/%s.env", confDir, env)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("WARNING: No %s.env file found in %s: %v", env, confDir, err)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	if v.ConfigFileUsed() != "" {
		if err := v.Unmarshal(&cfg); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Database,
		c.Postgres.SslMode,
	)
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}
