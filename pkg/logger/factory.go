package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/google/uuid"
)

var _ LogFactory = (*Factory)(nil)

type Factory struct {
	mu        *sync.Mutex
	out       io.Writer
	level     Level
	prefix    string
	color     bool
	timestamp bool
}

func NewFactory(cfg *config.LoggerConfig) *Factory {
	lvl := INFO
	for l, name := range levelNames {
		if name == cfg.Level {
			lvl = l
		}
	}
	return &Factory{
		mu:        &sync.Mutex{},
		out:       os.Stdout,
		level:     lvl,
		prefix:    cfg.Prefix,
		color:     cfg.Color,
		timestamp: cfg.Timestamp,
	}
}

//func (f *Factory) SetLevel(lvl string) *Factory {
//	f.level = levelNames[]
//	return f
//}

func (f *Factory) Fork(path string) *Logger {
	id := uuid.New().String()
	return &Logger{
		mu:        *f.mu,
		out:       f.out,
		level:     f.level,
		prefix:    fmt.Sprintf("[%s] |%s|(%s)", f.prefix, path, id),
		color:     f.color,
		timestamp: f.timestamp,
	}
}
