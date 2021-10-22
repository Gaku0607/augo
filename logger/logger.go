package logger

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type LoggerOptions func(*Logger)

func SetLoggerConfig(config *LoggerConfig) LoggerOptions {
	return func(l *Logger) {
		l.Config = config
	}
}

type Logger struct {
	Config *LoggerConfig
	mu     *sync.Mutex
}

func NowLogger(opts ...LoggerOptions) *Logger {
	l := &Logger{}
	l.defaultParms()
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Logger) defaultParms() {
	l.mu = &sync.Mutex{}
	l.Config = DefaultLoggerConfig()
}

func (l *Logger) DebugPrint(fprint func(io.Writer, ...interface{}) (int, error), val string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if fprint == nil {
		fprint = fmt.Fprint
	}
	if !strings.HasSuffix(val, "\n") {
		val += "\n"
	}
	fprint(l.Config.Out, val)
}

func (l *Logger) Log(parms *LoggerParms) {
	l.DebugPrint(parms.setTypeColor(), l.Config.Format(parms))
}

type LoggerConfig struct {
	Format LoggerFormatter
	Out    io.Writer
}

func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Format: defaultFormatter,
		Out:    defaultWriter,
	}
}

//關閉資源
func (l *LoggerConfig) Close() error {
	if wc, ok := l.Out.(io.WriteCloser); ok {
		return wc.Close()
	}
	return nil
}
