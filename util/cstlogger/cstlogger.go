package cstlogger

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

func NewLogger(env string) *zerolog.Logger {
	logLevel := zerolog.TraceLevel
	var output io.Writer = zerolog.ConsoleWriter{Out: os.Stdout,
		TimeFormat: time.RFC822Z,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		}}

	if strings.Contains(env, "prod") {
		logLevel = zerolog.ErrorLevel
		fileLogger := &lumberjack.Logger{
			Filename:   "track-pro-app.log",
			MaxSize:    10, //
			MaxBackups: 10,
			Compress:   true,
		}
		output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
	}

	logger := zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &logger
}
