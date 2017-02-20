package log

// https://golang.org/src/log/log.go
// https://github.com/sirupsen/logrus

import (
	"io"
	"os"
	"sync"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type Logger struct {
	// *log.Logger
	flag   int
	level  Level
	mu     sync.Mutex
	out    io.Writer
	prefix string
	buf    []byte
}

func New(out io.Writer, prefix string, flag int) *Logger {
	logger := &Logger{out: out, prefix: prefix, flag: flag}
	logger.SetLevel(InfoLevel)
	return logger
}

var std = New(os.Stderr, "", LstdFlags)

func StandardLogger() *Logger {
	return std
}

func Flags() int {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.flag
}

func SetFlags(flag int) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.flag = flag
}

func GetLevel() Level {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.level
}

func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.level = level
}

func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}

func Prefix() string {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.prefix
}

func SetPrefix(prefix string) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.prefix = prefix
}

func Printf(format string, v ...interface{}) { std.Printf(format, v...) }
func Println(v ...interface{})               { std.Println(v...) }
func Print(v ...interface{})                 { std.Print(v...) }
func Debugf(format string, v ...interface{}) { std.Debugf(format, v...) }
func Debugln(v ...interface{})               { std.Debugln(v...) }
func Debug(v ...interface{})                 { std.Debug(v...) }
func Infof(format string, v ...interface{})  { std.Infof(format, v...) }
func Infoln(v ...interface{})                { std.Infoln(v...) }
func Info(v ...interface{})                  { std.Info(v...) }
func Warnf(format string, v ...interface{})  { std.Warnf(format, v...) }
func Warnln(v ...interface{})                { std.Warnln(v...) }
func Warn(v ...interface{})                  { std.Warn(v...) }
func Errorf(format string, v ...interface{}) { std.Errorf(format, v...) }
func Errorln(v ...interface{})               { std.Errorln(v...) }
func Error(v ...interface{})                 { std.Error(v...) }
func Panicf(format string, v ...interface{}) { std.Panicf(format, v...) }
func Panicln(v ...interface{})               { std.Panicln(v...) }
func Panic(v ...interface{})                 { std.Panic(v...) }

func Output(calldepth int, s string) error {
	return std.Output(calldepth+1, s)
}
