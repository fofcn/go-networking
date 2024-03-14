package log

import (
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
)

var glogger zerolog.Logger

func InitLogger() error {
	logger, err := createLogger(zerolog.DebugLevel, "go-networking.log", true)
	if err != nil {
		return err
	}

	glogger = logger
	return nil
}

func Debug(msg string) {
	glogger.Debug().Msg(msg)
}

func Info(msg string) {
	glogger.Info().Msg(msg)
}

func Warn(msg string) {
	glogger.Warn().Msg(msg)
}

func Error(msg string) {
	glogger.Error().Msg(msg)
}

func ErrorErr(err error) {
	glogger.Error().Err(err).Msg("")
}

func ErrorErrMsg(err error, msg string) {
	glogger.Error().Err(err).Msg(msg)
}

func Fatal(msg string) {
	glogger.Fatal().Msg(msg)
}

func Panic(msg string) {
	glogger.Panic().Msg(msg)
}

func createLogger(level zerolog.Level, fileName string, console bool) (zerolog.Logger, error) {

	buildInfo, _ := debug.ReadBuildInfo()

	file, err := os.OpenFile(
		fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(file).
		Level(level).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	if console {
		consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: "2006-01-02 15:04:05"}
		logger = logger.Output(consoleOutput)
	}

	return logger, nil
}
