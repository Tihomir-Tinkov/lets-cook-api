package config

import (
	"reflect"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Logger    LoggerConfig   `envPrefix:"LOG_"`
	Env       string         `env:"ENV"`
	Server    ServerConfig   `envPrefix:"API_"`
	Routes    RoutesConfig   `envPrefix:"ROUTES_"`
	DB        PostgresConfig `envPrefix:"DB_"`
	StorePath string         `envPrefix:"STORE_PATH_" envDefault:"storage"`
}

type ServerConfig struct {
	Cors         CorsConfig `envPrefix:"CORS_"`
	InternalPort uint16     `env:"PORT"`
}

type RoutesConfig struct {
	Prefix           string `env:"PREFIX" envDefault:"api"`
	LoggerMiddleware bool   `env:"LOGGER_MIDDLEWARE" envDefault:"true"`
}

type LoggerConfig struct {
	Encoding string `env:"ENCODING" envDefault:"json"`
	Level    string `env:"LEVEL" envDefault:"info"`
}

type PostgresConfig struct {
	Host            string        `env:"HOST"`
	User            string        `env:"USER"`
	Password        string        `env:"PASS"`
	DbName          string        `env:"NAME"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"10"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"0s"`
	Port            uint16        `env:"PORT"`
	SSLMode         bool          `env:"SSL_MODE"`
}

type CorsConfig struct {
	AllowOrigins string `env:"ALLOW_ORIGINS"`
}

func (cfg *Config) Parse() error {
	if err := env.ParseWithOptions(cfg, env.Options{
		Prefix: "APP_",
	}); err != nil {
		return err
	}

	trimStrings(reflect.ValueOf(cfg).Elem())
	return nil
}

func trimStrings(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(strings.TrimSpace(field.String()))
		case reflect.Struct:
			trimStrings(field)
		}
	}
}
