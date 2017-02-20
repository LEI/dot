package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

var (
	PanicPrefix = "Panic: "
	FatalPrefix = "Fatal: "
	ErrorPrefix = "Error: " // × ✕ ✖ ✗ ✘
	WarnPrefix  = "Warn: "  // ⚠ !
	// SuccessPrefix = "✓" // ✔
	InfoPrefix  = "" // ›
	DebugPrefix = "Debug: "
)

func (l *Logger) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

func (l *Logger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) { l.Output(2, fmt.Sprintln(v...)) }

func (l *Logger) Print(v ...interface{}) { l.Output(2, fmt.Sprint(v...)) }

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level >= DebugLevel {
		l.Output(2, DebugPrefix+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.level >= DebugLevel {
		l.Output(2, DebugPrefix+fmt.Sprintln(v...))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level >= DebugLevel {
		l.Output(2, DebugPrefix+fmt.Sprint(v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level >= InfoLevel {
		l.Output(2, InfoPrefix+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infoln(v ...interface{}) {
	if l.level >= InfoLevel {
		l.Output(2, InfoPrefix+fmt.Sprintln(v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level >= InfoLevel {
		l.Output(2, InfoPrefix+fmt.Sprint(v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level >= WarnLevel {
		l.Output(2, WarnPrefix+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	if l.level >= WarnLevel {
		l.Output(2, WarnPrefix+fmt.Sprintln(v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level >= WarnLevel {
		l.Output(2, WarnPrefix+fmt.Sprint(v...))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level >= ErrorLevel {
		l.Output(2, ErrorPrefix+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	if l.level >= ErrorLevel {
		l.Output(2, ErrorPrefix+fmt.Sprintln(v...))
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.level >= ErrorLevel {
		l.Output(2, ErrorPrefix+fmt.Sprint(v...))
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level >= FatalLevel {
		l.Output(2, FatalPrefix+fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
	if l.level >= FatalLevel {
		l.Output(2, FatalPrefix+fmt.Sprintln(v...))
	}
	os.Exit(1)
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.level >= FatalLevel {
		l.Output(2, FatalPrefix+fmt.Sprint(v...))
	}
	os.Exit(1)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprint(v...)
	if l.level >= PanicLevel {
		l.Output(2, PanicPrefix+fmt.Sprintf(format, s))
	}
	panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprint(v...)
	if l.level >= PanicLevel {
		l.Output(2, PanicPrefix+fmt.Sprintln(s))
	}
	panic(s)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	if l.level >= PanicLevel {
		l.Output(2, PanicPrefix+fmt.Sprint(s))
	}
	panic(s)
}

func (l *Logger) Output(calldepth int, s string) error {
	now := time.Now()
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.flag&LUTC != 0 {
		t = t.UTC()
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
