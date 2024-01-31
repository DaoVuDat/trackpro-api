package cstlogger

import (
	"fmt"
	"github.com/fatih/color"
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
		TimeFormat: time.RFC1123Z,
		FormatLevel: func(i interface{}) string {
			levelText := strings.ToUpper(fmt.Sprintf("%s", i))
			levelTextFormat := fmt.Sprintf("\t[%s]\t", levelText)
			switch levelText {
			case "TRACE":
				return color.CyanString(levelTextFormat)
			case "DEBUG":
				return color.HiWhiteString(levelTextFormat)
			case "INFO":
				return color.GreenString(levelTextFormat)
			case "WARN":
				return color.YellowString(levelTextFormat)
			case "ERROR":
				return color.RedString(levelTextFormat)
			case "FATAL":
				return color.HiBlueString(levelTextFormat)
			case "PANIC":
				return color.MagentaString(levelTextFormat)
			default:
				return color.WhiteString(levelText)
			}
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
