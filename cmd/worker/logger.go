package worker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	log zerolog.Logger
}

func NewWorkerLogger(role string) *Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatCaller: func(i interface{}) string {
			return "[worker]"
		},
	}

	log := zerolog.New(output).With().Timestamp().Str("role", role).Logger()
	return &Logger{log}
}

func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
	logger.log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}
