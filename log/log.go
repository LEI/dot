package log

import (
	"fmt"
	"log"
	"os"
)

type Logger log.Logger

func (l *Logger) Debug(msg ...interface{}) {
	// if debug {}
	fmt.Println("Debug:", fmt.Sprint(msg...))
}

func (l *Logger) Info(msg ...interface{}) {
	fmt.Println("Info:", fmt.Sprint(msg...))
}

func (l *Logger) Warn(msg ...interface{}) {
	fmt.Println("Warning:", fmt.Sprint(msg...))
}

func (l *Logger) Error(msg ...interface{}) {
	fmt.Println("Error:", fmt.Sprint(msg...))
}

func (l *Logger) Fatal(msg ...interface{}) {
	fmt.Println("Fatal:", fmt.Sprint(msg...))
	os.Exit(1)
}

func (l *Logger) Panic(msg ...interface{}) {
	panic(fmt.Sprintf("Panic: %s\n", msg...))
}
