package db

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func ConnectOptional(conf config.PostgresConfig) (*pgxpool.Pool, error) {
	if conf.Host == "" {
		return nil, nil
	}
	return Connect(conf)
}

func Connect(conf config.PostgresConfig) (*pgxpool.Pool, error) {
	connString := buildConnUrl(conf)

	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Error().Str("url", connString).Err(err).Msg("failed to parse database url")
		return nil, err
	}

	if conf.MaxOpenConns > 0 {
		connConfig.MaxConns = int32(conf.MaxOpenConns)
	}
	if conf.MaxIdleConns > 0 {
		connConfig.MinConns = int32(conf.MaxIdleConns)
	}
	if conf.ConnMaxLifetime > 0 {
		connConfig.MaxConnLifetime = time.Second * conf.ConnMaxLifetime
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to db")
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to ping db")
		return nil, err
	}

	log.Debug().Fields(map[string]any{
		"host": conf.Host,
		"port": conf.Port,
		"user": conf.User,
		"db":   conf.DbName,
	}).Msg("connected to db")

	return conn, nil
}

func buildConnUrl(conf config.PostgresConfig) string {
	host := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	dbHost, err := netip.ParseAddr(conf.Host)
	if err == nil && dbHost.IsValid() {
		host = netip.AddrPortFrom(dbHost, conf.Port).String()
	}

	ssl := "disable"
	if conf.SSLMode {
		ssl = "require"
	}

	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", conf.User, conf.Password, host, conf.DbName, ssl)
}
