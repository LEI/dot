package formatter

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
)

type CLIFormatter struct {
	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool
}

func (f *CLIFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// f.appendKeyValue(b, "level", entry.Level.String())
	// if entry.Message != "" {
	// 	f.appendKeyValue(b, "msg", entry.Message)
	// }
	// for _, key := range keys {
	// 	f.appendKeyValue(b, key, entry.Data[key])
	// }


	f.printCommand(b, entry, keys)

	b.WriteByte('\n')
	return b.Bytes(), nil
}


func (f *CLIFormatter) printCommand(b *bytes.Buffer, entry *log.Entry, keys []string) {
	lvl := strings.Title(entry.Level.String())
	b.WriteString(fmt.Sprintf("%s: ", lvl))
	b.WriteString(entry.Message)
	for _, key := range keys {
		b.WriteString(fmt.Sprintf(" %s=%v", key, entry.Data[key]))
	}

}

func (f *CLIFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}


func (f *CLIFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
	b.WriteByte(' ')
}

func (f *CLIFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
