package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var HistoryLogger zerolog.Logger
var Logger zerolog.Logger

func init() {
	historyLogFile := openLogFile("history")
	HistoryLogger = zerolog.New(historyLogFile).With().Timestamp().Logger()
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func openLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Msgf("error opening file: %v", err)
	}
	return f
}

func generateName(value string) string {
	// Hol die aktuelle Zeit in Millisekunden
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	return fmt.Sprintf("%s-%d", value, milliseconds)
}