package log

import (
	"github.com/rs/zerolog"
	"os"
)

var Logger *zerolog.Logger

func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	Logger = &logger
}
