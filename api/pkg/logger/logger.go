package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Setup sets up the logger
func Setup(level string) zerolog.Logger {
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if level == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Configure logger
	logger := zerolog.New(os.Stdout).
		With().Timestamp().Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return logger
}
