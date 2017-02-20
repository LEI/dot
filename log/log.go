package log

import (
	"fmt"
	"log"
	"os"
)

type Logger log.Logger

func (l *Logger) Debug(v ...interface{}) {
	// if debug {}
	log.Output(2, "Debug: " + fmt.Sprint(v...))
}

func (l *Logger) Info(v ...interface{}) {
	log.Output(2, "Info: " + fmt.Sprint(v...))
}

func (l *Logger) Warn(v ...interface{}) {
	log.Output(2, "Warning: " + fmt.Sprint(v...))
}

func (l *Logger) Error(v ...interface{}) {
	log.Output(2, "Error: " + fmt.Sprint(v...))
}

func (l *Logger) Fatal(v ...interface{}) {
	log.Output(2, "Fatal: " + fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprintf("Panic:", v...)
	log.Output(2, s)
	panic(s)
}
