package logger

import (
	"os"
	"path"

	"github.com/pagient/pagient-server/internal/config"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logFile *os.File

// Init sets up global logger using configuration
func Init() error {
	level, err := zerolog.ParseLevel(config.Log.Level)
	if err != nil {
		return errors.New("parse log level failed")
	}
	zerolog.SetGlobalLevel(level)

	logFile, err = os.OpenFile(path.Join(config.General.Root, "pagient.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return errors.New("logfile could not be opened")
	}

	if config.Log.Pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     logFile,
				NoColor: !config.Log.Colored,
			},
		)

		return nil
	}

	log.Logger = log.Output(logFile)
	return nil
}

// Close closes the log file
func Close() error {
	return logFile.Close()
}
