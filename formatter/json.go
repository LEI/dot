package formatter

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultTimestampFormat ...
	DefaultTimestampFormat = time.RFC3339
	// FieldKeyMsg ...
	FieldKeyMsg = "msg"
	// FieldKeyLevel ...
	FieldKeyLevel = "level"
	// FieldKeyTime ...
	FieldKeyTime = "time"
)

type fieldKey string

// FieldMap ...
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

// JSONFormatter ...
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// FieldMap allows users to customize the names of keys for various fields.
	// As an example:
	// formatter := &JSONFormatter{
	//	FieldMap: FieldMap{
	//   FieldKeyTime: "@timestamp",
	//   FieldKeyLevel: "@level",
	//   FieldKeyMsg: "@message",
	//    },
	// }
	FieldMap FieldMap
}

// Format entry
func (f *JSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields, len(entry.Data)+3)

	// timestampFormat := f.TimestampFormat
	// if timestampFormat == "" {
	// 	timestampFormat = DefaultTimestampFormat
	// }

	// if !f.DisableTimestamp {
	// 	data[f.FieldMap.resolve(FieldKeyTime)] = entry.Time.Format(timestampFormat)
	// }
	data[f.FieldMap.resolve(FieldKeyMsg)] = entry.Message
	data[f.FieldMap.resolve(FieldKeyLevel)] = entry.Level.String()

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
