package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Log struct {
	*slog.Logger
}

func New(opts ...Option) *Log {
	cfg := config(opts...)

	log := slog.New(slog.NewJSONHandler(os.Stdout, &cfg))

	slog.SetDefault(log)

	return &Log{log}
}

func (l *Log) Info(format string, msg ...any) {
	l.Logger.Info(fmt.Sprintf(format, msg...))
}

func (l *Log) Debug(err error, format string, msg ...any) {
	l.Logger.Debug(
		fmt.Sprintf(format, msg...),
		handleErr(err),
	)
}

func (l *Log) Error(err error, format string, msg ...any) {
	l.Logger.Error(
		fmt.Sprintf(format, msg...),
		handleErr(err),
	)
}

func handleErr(err error) slog.Attr {
	if err != nil {
		return slog.String("err", err.Error())
	}

	return slog.Attr{}
}
