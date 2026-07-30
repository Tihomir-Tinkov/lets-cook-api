package logger

import (
	"io"
	"os"
	"time"

	"github.com/Tihomir-Tinkov/lets-cook-api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type splitWriter struct {
	out io.Writer
	err io.Writer
}

func (s splitWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level <= zerolog.WarnLevel {
		return s.out.Write(p)
	}
	return s.err.Write(p)
}

func (s splitWriter) Write(p []byte) (n int, err error) {
	return s.out.Write(p)
}

func Init(conf config.LoggerConfig) {
	var writer zerolog.LevelWriter

	if conf.Encoding != "json" {
		writer = splitWriter{
			out: zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
			err: zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
		}
	} else {
		writer = splitWriter{
			out: os.Stdout,
			err: os.Stderr,
		}
	}

	log.Logger = zerolog.New(writer).With().Timestamp().Logger()

	switch conf.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
