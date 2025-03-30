package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%-6s", i))
		},
	}

	// Set log level based on environment
	if os.Getenv("ENVIRONMENT") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Debug logs a debug message
func Debug(msg string, fields ...interface{}) {
	log.Debug().Fields(fieldsToMap(fields...)).Msg(msg)
}

// Info logs an info message
func Info(msg string, fields ...interface{}) {
	log.Info().Fields(fieldsToMap(fields...)).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...interface{}) {
	log.Warn().Fields(fieldsToMap(fields...)).Msg(msg)
}

// Error logs an error message
func Error(msg string, fields ...interface{}) {
	log.Error().Fields(fieldsToMap(fields...)).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...interface{}) {
	log.Fatal().Fields(fieldsToMap(fields...)).Msg(msg)
}

func fieldsToMap(fields ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			m[fields[i].(string)] = fields[i+1]
		}
	}
	return m
}
