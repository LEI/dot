package logger

// https://github.com/sirupsen/logrus/blob/master/exported.go

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

// var AllLevels = []Level{
// 	PanicLevel,
// 	FatalLevel,
// 	ErrorLevel,
// 	WarnLevel,
// 	InfoLevel,
// 	DebugLevel,
// }

// var std = log.New(os.Stderr, "", log.LstdFlags)

type Level uint8

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

type Logger struct {
	flag   int
	level  Level
	mu     sync.Mutex
	out    io.Writer
	prefix string
	// buf []byte
}

func New(out io.Writer, prefix string, flag int) *Logger {
	logger := &Logger{out: out, prefix: prefix, flag: flag}
	logger.SetLevel(InfoLevel)
	return logger
}

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level >= DebugLevel {
		l.log("Debug: "+fmt.Sprint(v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level >= InfoLevel {
		l.log(fmt.Sprint(v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level >= WarnLevel {
		l.log("Warning: "+fmt.Sprint(v...))
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.level >= ErrorLevel {
		l.log("Error: "+fmt.Sprint(v...))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.level >= FatalLevel {
		l.log("Fatal: "+fmt.Sprint(v...))
	}
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	if l.level >= PanicLevel {
		l.log(fmt.Sprint("Panic:", s))
	}
	panic(s)
}

func (l *Logger) Output(calldepth int, s string) error {
	return log.Output(calldepth, s)
}

func (l *Logger) log(s string) {
	fmt.Fprintf(l.out, "%s\n", s)
}
