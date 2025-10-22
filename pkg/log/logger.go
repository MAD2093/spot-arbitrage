package log

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

var (
	log *zerolog.Logger
)

// Config представляет конфигурацию логирования
type Config struct {
	Level  zerolog.Level `json:"level" yaml:"level"`
	Pretty bool          `json:"pretty" yaml:"pretty"`
}

func Init(cfg Config) *zerolog.Logger {
	wr := diode.NewWriter(os.Stdout, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})

	zerolog.SetGlobalLevel(cfg.Level)

	baseLogger := zerolog.New(wr).With().Caller().Logger()

	if cfg.Pretty {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC850,
		}
		baseLogger = zerolog.New(consoleWriter).
			Level(cfg.Level).
			With().
			Caller().
			Timestamp().
			Logger()
	}

	log = &baseLogger
	return log
}

func GetLogger() *zerolog.Logger {
	if log == nil {
		// panic("logger was not initialized")
	}
	return log
}
