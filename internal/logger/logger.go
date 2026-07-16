// Package logger provides a single, shared structured logger for the whole
// service. No package under internal/ should ever call fmt.Println/log.Println
// directly — everything goes through this logger so every line has consistent
// structure (timestamp, level, fields) and can be shipped to a log aggregator.
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New builds a zerolog.Logger configured from the given level string
// ("debug" | "info" | "warn" | "error"). In "development" it also writes
// human-readable console output; in any other environment it writes JSON,
// which is what log aggregators (Loki, CloudWatch, Datadog, etc.) expect.
func New(levelStr string, environment string) zerolog.Logger {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if environment == "development" {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}
		return zerolog.New(consoleWriter).With().Timestamp().Logger()
	}

	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}