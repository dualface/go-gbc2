package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

var (
	L *zap.SugaredLogger

	backend *zap.Logger
)

func init() {
	mode := strings.TrimSpace(strings.ToLower(os.Getenv("GOGBC_RUN_MODE")))
	if mode == "prod" {
		backend, _ = zap.NewProduction()
	} else {
		backend, _ = zap.NewDevelopment()
	}
	L = backend.Sugar()
}

func Sync() {
	backend.Sync() // flushes buffer, if any
}
