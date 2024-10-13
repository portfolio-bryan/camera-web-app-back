package logger

import (
	"encoding/json"

	"github.com/bperezgo/rtsp/shared/platform/handlertypes"
	"github.com/rs/zerolog/log"
)

type LogHttpInput struct {
	Request  handlertypes.Request
	Response handlertypes.Response
}

type Error struct {
	Message string
}

type LogState int64

const (
	SUCCESS LogState = iota
	FAILED
	PENDING
)

func (s LogState) String() string {
	switch s {
	case SUCCESS:
		return "SUCCESS"
	case FAILED:
		return "FAILED"
	case PENDING:
		return "PENDING"
	}
	return "unknown"
}

type LogInput struct {
	Action  string
	State   LogState
	Message string
	Http    *LogHttpInput
	Error   *Error
	Meta    *handlertypes.Meta
}

type Logger struct{}

var logger *Logger

func InitLogger() {
	if logger == nil {
		logger = &Logger{}
	}
}

func GetLogger() *Logger {
	InitLogger()
	return logger
}

func (Logger) Info(input LogInput) {
	if e := log.Info(); e.Enabled() {
		bHttp, err := json.Marshal(input.Http)
		if err != nil {
			e.Str("action", input.Action).Str("state", input.State.String()).Msg(err.Error())
			return
		}

		bMeta, err := json.Marshal(input.Meta)
		if err != nil {
			e.Str("action", input.Action).Str("state", input.State.String()).Msg(err.Error())
			return
		}
		e.Str("action", input.Action).
			Str("action", input.State.String()).
			RawJSON("http", bHttp).
			RawJSON("meta", bMeta).
			Msg(input.Message)
	}
}

func (Logger) Error(input LogInput) {
	if e := log.Error(); e.Enabled() {
		e.Str("action", input.Action).Str("action", input.State.String()).Str("error", input.Error.Message).Msg(input.Message)
		e.Any("http", input.Http)
	}
}

func (Logger) Warn(input LogInput) {
	if e := log.Warn(); e.Enabled() {
		e.Str("action", input.Action).Str("action", input.State.String()).Msg(input.Message)
		e.Any("http", input.Http)
	}
}
