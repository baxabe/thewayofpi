package logs

import (
	"fmt"
	"time"
	"bytes"
)

const (
	i = iota
	w
	e
)

func log(level uint8, msg string) error {
	var severity string
	switch level {
	case i:
		severity = "Info"
	case w:
		severity = "Warning"
	case e:
		severity = "Error"
	default:
		severity = "Unknown"
	}
	now := time.Now().Format(time.RFC3339Nano)
	err := fmt.Errorf("[%v]%v: %v\n", now, severity, msg)
	fmt.Print(err)
	return err
}

//func Info(err ...error) error {
//	return log(i, expandError(err))
//}

func Warning(err ...error) error {
	return log(w, expandError(err))
}

func Error(err ...error) error {
	return log(e, expandError(err))
}

func expandError(err []error) string {
	var result bytes.Buffer
	if len(err) > 0 {
		result.WriteString(fmt.Sprintf("[%v", err[0].Error()))
		for i := 1; i < len(err); i++ {
			result.WriteString(fmt.Sprintf(" | %v", err[i].Error()))
		}
		result.WriteString("]")
	}
	return result.String()
}
