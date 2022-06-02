package utility

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/log"
)

type LogLevel string

const (
	INFO LogLevel = "INFO"
)

type Log struct {
	Level   LogLevel
	Message string
	Data    []any
}

func NewLog(level LogLevel, message string, data []any) *Log {
	return &Log{Level: level, Message: message, Data: data}
}

func InfoLog(str string) Log {
	return Log{Level: INFO, Message: str, Data: nil}
}

func InfoLogf(message string, data []any) Log {
	return Log{Level: INFO, Message: message, Data: data}
}

type LoggingActor struct {
	Logger *log.Logger
}

func (state *LoggingActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Log:
		fields := make([]log.Field, len(msg.Data))
		for i, datum := range msg.Data {
			fields[i] = log.Message(datum)
		}
		switch msg.Level {
		case INFO:
			state.Logger.Info(msg.Message, fields...)
		}

	}
}
