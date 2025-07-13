package helper

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const resetColor = "\033[0m"

func redText(str string) string {
	return "\033[31m" + str + resetColor
}

func greenText(str string) string {
	return "\033[32m" + str + resetColor
}

// NewLogger creates new logger with prefix.
func NewLogger(prefix string) *log.Logger {
	return log.New(log.Writer(), fmt.Sprintf("[%s] ", prefix), log.Ldate|log.Ltime|log.Lmsgprefix)
}

func (helper *ServerHelper) LogRequest(statusCode int, r *http.Request, start time.Time, errorMsg string) {
	logText := fmt.Sprintf(
		"%d %s %s %s in %v ",
		statusCode,
		r.Method,
		r.URL.Path,
		r.URL.Query().Encode(),
		time.Since(start),
	)

	if errorMsg != "" {
		helper.MainLogger.Printf(redText("%s ERROR: %s"), logText, errorMsg)
	} else {
		helper.MainLogger.Println(greenText(logText))
	}
}
