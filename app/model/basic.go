package model

type LogLevel string
type SimpleResponseCode string

type SimpleResponse struct {
	Code    SimpleResponseCode `json:"code"`
	Message string             `json:"message"`
}

const (
	Version = "0.0"
)

const (
	LogLevelTrace LogLevel           = "trace"
	LogLevelDebug LogLevel           = "debug"
	LogLevelInfo  LogLevel           = "info"
	LogLevelWarn  LogLevel           = "warn"
	LogLevelErr   LogLevel           = "error"
	ErrCode       SimpleResponseCode = "ERROR"
	WarnCode      SimpleResponseCode = "WARN"
	EmptyCode     SimpleResponseCode = "EMPTY"
	SuccessCode   SimpleResponseCode = "SUCCESS"
)
