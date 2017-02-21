package log

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
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
