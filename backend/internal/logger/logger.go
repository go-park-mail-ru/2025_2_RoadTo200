package logger

import (
	"io"
	"os"
	"sync"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
)

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[Level]string{
	TRACE: "TRACE",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger - основная структура логгера
type Logger struct {
	mu        sync.Mutex
	out       io.Writer
	level     Level
	prefix    string
	color     bool
	timestamp bool
}

func New(cfg *config.LoggerConfig) *Logger {
	lvl := INFO
	for l, name := range levelNames {
		if name == cfg.Level {
			lvl = l
		}
	}
	return &Logger{
		out:   os.Stdout,
		level: lvl,
		// prefix:    cfg.Prefix,
		// color:     cfg.Color,
		// timestamp: cfg.Timestamp,
	}
}
